package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/git-starter-go/utils"
)

// Usage: your_program.sh <command> <arg1> <arg2> ...
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintf(os.Stderr, "Logs from your program will appear here!\n")

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		// Uncomment this block to pass the first stage!

		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")
	case "hash-object":
		option := os.Args[2]

		if option == "-w" {
			file := os.Args[3]

			if !utils.FileExists(file) {
				panic("File does not exist")
			}

			data, err := os.ReadFile(file)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
			}

			data = []byte(fmt.Sprintf("blob %d\x00%s", len(data), data))
			fileData, err := utils.CompressData(data)

			hash, err := utils.GenerateHash(data)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
			}

			objectFolder := hash[:2]

			dir := fmt.Sprintf(".git/objects/%s", objectFolder)
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}

			of := hash[2:]

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
			}
			if err := os.WriteFile(fmt.Sprintf("%s/%s", dir, of), []byte(fileData), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
			}

			fmt.Fprintf(os.Stdout, hash)
		}
	case "cat-file":
		option := os.Args[2]

		if option == "-p" {
			hash := os.Args[3]
			dir := ".git/objects/%s/%s"

			objectFolder := hash[:2]

			of := hash[2:]

			data, err := os.ReadFile(fmt.Sprintf(dir, objectFolder, of))

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
			}

			decodedData, err := utils.UnCompressData(data)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
			}

			dataParts := strings.Split(decodedData, "\x00")

			resp := dataParts[len(dataParts)-1]

			fmt.Fprintf(os.Stdout, resp)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
