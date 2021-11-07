# Let's Go Chat

## Introduction
This repository is used for learning purposes.  
Currently it contains only **hasher** package

## Packages
### hasher
Provides possibility to calculate and verify password hashes

#### Examples

##### Hash password
```go
package main

import (
	"fmt"

	"github.com/id-tarzanych/lets-go-chat/pkg/hasher"
)

func main() {
	pwd := "S0m3Pas$w0rd"
	hash, _ := hasher.HashPassword(pwd)
	
	fmt.Println(hash)
}
```

##### Check Password Hash
```go
package main

import (
	"fmt"

	"github.com/id-tarzanych/lets-go-chat/pkg/hasher"
)

func main() {
	pwd := "S0m3Pas$w0rd"
	hash := "9d2924208aac19fe770d9271fc221b28340a56f12c7f3c7d3d35b3944db907b8"
	
	if result := hasher.CheckPasswordHash(pwd, hash); result {
		fmt.Println("Password is valid")
    } else {
		fmt.Println("Invalid password")
    }
}
```