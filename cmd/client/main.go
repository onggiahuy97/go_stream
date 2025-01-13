package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get("http://localhost:8080/stream")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		fmt.Print(string(line))
	}
}
