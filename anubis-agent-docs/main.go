package main

import (
	"fmt"
	"os"
)

func main() {
	readme, _ := os.ReadFile("README.md")
	fmt.Println(string(readme))
}
