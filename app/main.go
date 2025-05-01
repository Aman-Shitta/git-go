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
	case "ls-tree":

		option := os.Args[2]

		hash := os.Args[3]
		path := fmt.Sprintf(".git/objects/%v/%v", hash[:2], hash[2:])

		data, err := os.ReadFile(path)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		}

		strData, err := utils.UnCompressData(data)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error Uncompressing file: %s\n", err)
		}

		strData = strings.SplitAfterN(strData, "\x00", 2)[1]

		treeObjects := ParseTreeObject([]byte(strData))
		res := ""
		if option == "--name-only" {
			for _, item := range treeObjects {
				res += fmt.Sprintln(item.Name)
			}
		} else {
			for _, item := range treeObjects {
				res += fmt.Sprintf("%s %s %s %s\n", item.Mode, ObjectType[item.Mode], item.Hash, item.Name)
				//  fmt.Sprintf("%s\n", item.Name)
			}
		}

		fmt.Fprint(os.Stdout, strings.TrimPrefix(res, "\n"))

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}

var ObjectType = map[string]string{
	"100644": "blob", // normal file
	"100755": "file", // executable file
	"120000": "link", // symbolic link
	"40000":  "dir",  // directory (note: Git drops the leading zero)
}

type TreeObject struct {
	Mode string
	Hash string
	Name string
}

func ParseTreeObject(data []byte) []TreeObject {
	res := make([]TreeObject, 1)

	i := 0
	for i < len(data) {
		// 1. Read mode (up to space)
		start := i
		for data[i] != ' ' {
			i++
		}
		mode := string(data[start:i])
		i++ // skip space

		// 2. Read filename (up to null byte)
		start = i
		for data[i] != 0 {
			i++
		}
		filename := string(data[start:i])
		i++ // skip null byte

		// 3. Read 20-byte SHA-1 (binary), convert to hex
		if i+20 > len(data) {
			fmt.Println("Invalid tree format: not enough bytes for SHA-1")
			return res
		}
		sha1bin := data[i : i+20]
		sha1hex := fmt.Sprintf("%x", sha1bin)
		i += 20

		res = append(res, TreeObject{mode, sha1hex, filename})
	}

	return res
}
