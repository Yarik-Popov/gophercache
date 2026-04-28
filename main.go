package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	cache := CreateCache[string, string](3, 1)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		if len(parts) < 1 {
			break
		}
		switch parts[0] {
		case "g":
			value, ok := cache.Get(parts[1])
			if ok {
				fmt.Printf("Got %s: %s\n", parts[1], value)
			} else {
				fmt.Println("Not found: ", parts[1])
			}
		case "p":
			fmt.Printf("Putting %s: %s\n", parts[1], parts[2])
			cache.Put(parts[1], parts[2])
		}
	}
}
