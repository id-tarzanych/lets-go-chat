/*
Package hasher implements a simple library that allows to hash and verify passwords.
Currently supports only SHA256

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

import "testing"

func TestHashPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{ "Simple password", args{"12345678"}, "ef797c8118f02dfb649607dd5d3f8c7623048c9c063d532cc95c5ed7a898a64f", false},
		{ "Strongs password", args{"T3stP@s$w0rd"}, "ee3e03e2b3a920fb5e983ce3efb7cf85c5b0d9ac4f2ebfb98562642da83850db", false},
		{ "Empty value", args{""}, "", true},
	}
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashPassword(tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HashPassword() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	type args struct {
		password string
		hash     string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Empty value", args{"", ""}, false},
		{"Simple password", args{"admin123", "240be518fabd2724ddb6f04eeb1da5967448d7e831c08c8fa822809f74c720a9"}, true},
		{"Complex password", args{"nSQ+vGq@zA%c58aj", "3c87a26a13284ca4a125312c5c870ac0b7fba4a19558d009e4ccb44ced534fab"}, true},
		{"Incorrect value", args{"Ep@MLearning", "67f74d0df98c79f50e3b230d42774d3a0c6715f3a583df6f756e266108165dd8"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPasswordHash(tt.args.password, tt.args.hash); got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

