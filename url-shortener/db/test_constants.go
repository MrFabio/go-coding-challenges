package common

import "testing"

// Test constants shared across all test files
const TestURL = "https://codingchallenges.fyi/challenges/challenge-url-shortener"

func RandomURL() string {
	return "https://" + GenerateId() + ".com"
}

// DatabaseTestSuite contains all the generic tests that should work with any Database implementation
type DatabaseTestSuite struct {
	Name     string             // Name of the database implementation
	DB       Database           // Database to run tests against
	Validate func(t *testing.T) // Optional validation to run before all tests
	Cleanup  func() error       // Optional cleanup to run after all tests
	Close    func() error       // Optional close, runs after Cleanup
}
