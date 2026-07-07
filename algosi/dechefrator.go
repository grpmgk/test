package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func EncodeRLE(input string) string {
	if input == "" {
		return ""
	}
	var str strings.Builder
	count := 1
	for i := 0; i < len(input); i++ {
		if i+1 == len(input) || input[i] != input[i+1] {
			str.WriteByte(input[i])
			str.WriteString(strconv.Itoa(count))
			count = 1
		} else {
			count++
		}
	}
	return str.String()
}

func DecodeRLE(input string) string {
	if input == "" {
		return ""
	}
	var result strings.Builder
	for i := 0; i < len(input); i++ {
		char := input[i]
		i++

		var countStr strings.Builder

		for i < len(input) && unicode.IsDigit(rune(input[i])) {
			countStr.WriteByte(input[i])
			i++
		}
		i--
		count, _ := strconv.Atoi(countStr.String())
		result.WriteString(strings.Repeat(string(char), count))
	}

	return result.String()
}

func main() {
	test := "aaabbbghb"
	enc := EncodeRLE(test)
	dec := DecodeRLE(enc)

	fmt.Printf("Вход: %s\n", test)
	fmt.Printf("Сжато: %s\n", enc)
	fmt.Printf("Восстановлено: %s\n", dec)
	fmt.Printf("OK: %v\n", test == dec)
}
