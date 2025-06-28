package main

import (
	common "url/db"
	"url/db/in_mem"
)

func main() {
	var database common.Database = in_mem.NewInMemoryDatabase()
	entry1 := database.AddEntry("https://www.github.com")
	database.AddEntry("https://www.sapo.pt")
	database.String()
	database.DeleteEntry(entry1.ID) // delete entry1
	database.String()
}
