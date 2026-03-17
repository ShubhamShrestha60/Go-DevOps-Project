package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	fmt.Printf("HASH_START|%s|HASH_END\n", string(hash))
}
