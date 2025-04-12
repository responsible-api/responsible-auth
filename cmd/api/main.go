package main

import (
	"fmt"
	"responsible-api-go/config"
)

func main() {
	c := config.New()

	fmt.Println("Config loaded:", c)

	fmt.Println("Hello world from Responsible Go API!")
	fmt.Println("This is the first Version of the Responsible API.")
	fmt.Println("This is a personal project please don't use it in your production.")

	fmt.Println("To do: - Add a database connection")
	fmt.Println("To do: - Add a web server")
	fmt.Println("To do: - Add a web server Authorization validation")
}
