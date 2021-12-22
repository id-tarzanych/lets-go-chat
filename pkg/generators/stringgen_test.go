package generators

import (
	"strings"
	"testing"
)

var result string

func TestRandomString(t *testing.T) {
	const generateCount = 5

	tests := []struct {
		name   string
		length int
	}{
		{"16 chars", 16},
		{"24 chars", 24},
		{"32 chars", 32},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generatedStrings := make(map[string]bool, generateCount)

			for i := 1; i <= generateCount; i++ {
				generatedStrings[RandomString(tt.length)] = true
			}

			if len(generatedStrings) != generateCount {
				t.Errorf("Non-unique strings generated for length %d", tt.length)
			}

			for str, _ := range generatedStrings {
				for _, char := range str {
					if !strings.Contains(CHARSET, string(char)) {
						t.Errorf("Unexpected character generated: %c", char)
					}
				}
			}
		})
	}
}

func Test_stringWithCharset(t *testing.T) {
	const generateCount = 5

	type args struct {
		length  int
		charset string
	}
	tests := []struct {
		name string
		args args
	}{
		{"16 numeric chars", args{16, "0123456789"}},
		{"24 alphanumeric chars", args{24, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"}},
		{"32 alpha chars", args{32, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"}},
	}
	for _, tt := range tests {
		generatedStrings := make(map[string]bool, generateCount)

		for i := 1; i <= generateCount; i++ {
			generatedStrings[stringWithCharset(tt.args.length, tt.args.charset)] = true
		}

		if len(generatedStrings) != generateCount {
			t.Errorf("Non-unique strings generated for length %d and scenario %s", tt.args.length, tt.name)
		}

		for str, _ := range generatedStrings {
			for _, char := range str {
				if !strings.Contains(tt.args.charset, string(char)) {
					t.Errorf("Unexpected character generated for scenario %s: %c", tt.name, char)
				}
			}
		}
	}
}

func BenchmarkRandomString8(b *testing.B) {
	var r string

	for n := 0; n < b.N; n++ {
		r = RandomString(8)
	}

	result = r
}

func BenchmarkRandomString16(b *testing.B) {
	var r string

	for n := 0; n < b.N; n++ {
		r = RandomString(16)
	}

	result = r
}

func BenchmarkRandomString32(b *testing.B) {
	var r string

	for n := 0; n < b.N; n++ {
		r = RandomString(32)
	}

	result = r
}
