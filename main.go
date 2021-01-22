package main

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/stevelacy/imessage-viewer/core"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("Starting... 🚀")

	usr, _ := user.Current()
	dir := path.Join(usr.HomeDir, "/Library/Messages/chat.db")

	if len(os.Args) == 1 {
		fmt.Println("Please provide the command \nValid options:\n - iphone \n - imessage")
		os.Exit(1)
	}

	if len(os.Args) == 2 {
		fmt.Println("Please provide the recipient phone number or email: '+14151231234'")
		os.Exit(1)
	}

	command := os.Args[1]
	recipient := os.Args[2]

	if command == "imessage" {
		fmt.Printf("Using DB %s \n", dir)
		core.Process(dir, recipient)
		return
	}

	if command == "iphone" {
		core.HandleiOSBackups(dir, recipient)
		return
	}

	fmt.Printf("Unknown command: %s", command)
	os.Exit(1)
}
