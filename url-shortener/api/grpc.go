package api

import (
	"context"
	"log"
	"net"
	"sync"

	pb "url-shortener/api/grpc" // Your full module path

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Entry struct {
	ID   string
	URL  string
	Hash string
}

type GrpcServer struct {
	pb.UnimplementedEntryServiceServer
	handler   *Handler
	mu        sync.RWMutex
	streams   []*streamSubscriber
	broadcast chan *pb.GetEntryResponse
}

type streamSubscriber struct {
	stream pb.EntryService_WatchEntriesServer
	done   chan struct{}
}

func NewGrpcServer(handler *Handler) *GrpcServer {
	return &GrpcServer{
		handler: handler,
		//entries:   make([]Entry, 0),
		streams:   make([]*streamSubscriber, 0),
		broadcast: make(chan *pb.GetEntryResponse, 100),
	}
}

func (s *GrpcServer) AddEntry(_ context.Context, req *pb.CreateEntryRequest) (*pb.GetEntryResponse, error) {
	entry, err := s.handler.CreateEntry(req.Url)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Error Processing the URL: %v", err)
	}

	resp := &pb.GetEntryResponse{
		Id:  entry.ID,
		Url: entry.URL,
	}

	select {
	case s.broadcast <- resp:
	default:
		log.Printf("Broadcast channel full; dropping update")
	}

	return resp, nil
}

func (s *GrpcServer) GetEntry(_ context.Context, req *pb.GetEntryRequest) (*pb.GetEntryResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, found := s.handler.database.GetEntry(req.Id)

	if !found {
		return nil, status.Errorf(codes.InvalidArgument, "Entry not found: %s", req.Id)
	}

	return &pb.GetEntryResponse{
		Id:  entry.ID,
		Url: entry.URL,
	}, nil
}

func (s *GrpcServer) WatchEntries(_ *emptypb.Empty, stream pb.EntryService_WatchEntriesServer) error {
	sub := &streamSubscriber{
		stream: stream,
		done:   make(chan struct{}),
	}
	s.mu.Lock()
	s.streams = append(s.streams, sub)
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		for i, st := range s.streams {
			if st == sub {
				s.streams = append(s.streams[:i], s.streams[i+1:]...)
				close(st.done)
				break
			}
		}
		s.mu.Unlock()
	}()

	// No initial entries sent; only listen for new additions
	for {
		select {
		case newEntry, ok := <-s.broadcast:
			if !ok {
				return nil
			}
			if err := stream.Send(newEntry); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}

func registerGrpcServer(handler *Handler) {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	pb.RegisterEntryServiceServer(srv, NewGrpcServer(handler))
	reflection.Register(srv) // Assuming reflection is enabled
	log.Printf("grpcServer listening at %v", lis.Addr())
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
