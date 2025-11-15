package main

import (
	"log"

	"github.com/google/uuid"
)

func main() {
	uuidString := uuid.NewString()
	log.Printf("UUID: %s\n", uuidString)
}
