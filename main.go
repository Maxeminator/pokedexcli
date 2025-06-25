package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	cfg := config{}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanned := scanner.Scan()
		if !scanned {
			break
		}
		line := scanner.Text()
		words := cleanInput(line)
		if cmd, ok := commands[words[0]]; ok {
			err := cmd.callback(&cfg)
			if err != nil {
				fmt.Println("Error:", err)
			}

		} else {
			fmt.Println("Unknown command")
		}
	}
}
