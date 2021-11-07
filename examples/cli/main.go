package main

import (
	"fmt"

	"github.com/id-tarzanych/lets-go-chat/pkg/generators"
	"github.com/id-tarzanych/lets-go-chat/pkg/hasher"
)

func main() {
	fmt.Println("Generating random password...")
	pwd := generators.RandomString(16)
	fmt.Println("Password:", pwd)

	hash, _ := hasher.HashPassword(pwd)
	fmt.Println()
	fmt.Printf("hasher.HashPassword(\"%s\") = \"%s\"\n", pwd, hash)
	fmt.Println()

	fmt.Println("Generating random password...")
	pwd = generators.RandomString(16)
	hash, _ = hasher.HashPassword(pwd)

	fmt.Println("Checking correct and incorrect hashes")
	fmt.Printf("hasher.CheckPasswordHash(\"%s\", \"%s\") = %v\n", pwd, hash, hasher.CheckPasswordHash(pwd, hash))

	// Modify hash and try once more.
	hash = hash[16:32]
	fmt.Printf("hasher.CheckPasswordHash(\"%s\", \"%s\") = %v\n", pwd, hash, hasher.CheckPasswordHash(pwd, hash))
}
