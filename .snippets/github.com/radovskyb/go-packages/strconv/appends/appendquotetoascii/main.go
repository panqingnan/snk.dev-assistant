package main

import (
	"fmt"
	"strconv"
)

func main() {
	b := []byte("quote (ascii):")

	// AppendQuoteToASCII appends a double-quoted Go string
	// literal representing s, as generated by QuoteToASCII,
	// to dst and returns the extended buffer.
	b = strconv.AppendQuoteToASCII(b, `"I'm inside of quotes"`)

	fmt.Println(string(b))
}
