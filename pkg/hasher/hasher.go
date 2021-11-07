/*
Package hasher implements a simple library that allows to hash and verify passwords.
Currently supports only SHA256

TODO: Rework solution in v2, replace hardcoded hasher with a dynamic value, use struct for functions.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hasher

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
)

var hasher hash.Hash

// Initializes SHA256 hash.Hash object.
func init()  {
	hasher = sha256.New()
}

// HashPassword returns SHA256 hash for password string.
// Returns error on empty string.
func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", errors.New("no input supplied")
	}

	hasher.Write([]byte(password))
	sum := fmt.Sprintf("%x", hasher.Sum(nil))
	hasher.Reset()

	return sum, nil
}

// CheckPasswordHash verifies if password matches provided hash.
func CheckPasswordHash(password, hash string) bool {
	calculatedHash, err := HashPassword(password)
	
	if err != nil {
		return false
	}
	
	return calculatedHash == hash
}