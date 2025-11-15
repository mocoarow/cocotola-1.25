package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	uuidString := uuid.NewString()
	fmt.Printf("UUID: %s\n", uuidString)
}
