package main

import (
	"Yarik-Popov/gophercache/src"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	keyValueStore := cache.CreateCache[string, string](3, 5*time.Second)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		if len(parts) < 1 {
			break
		}
		switch parts[0] {
		case "g":
			value, ok := keyValueStore.Get(parts[1])
			if ok {
				fmt.Printf("Got %s: %s\n", parts[1], value)
			} else {
				fmt.Println("Not found: ", parts[1])
			}
		case "p":
			fmt.Printf("Putting %s: %s\n", parts[1], parts[2])
			keyValueStore.Put(parts[1], parts[2])
		}
	}
}
