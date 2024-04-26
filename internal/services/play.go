package services

import (
	"fmt"
	"log"
	"os"

	"github.com/DexterLB/mpvipc"
)

func Play(conn *mpvipc.Connection, items ...Item) error {
	playlistFile, err := os.CreateTemp("", "playlist*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(playlistFile.Name()) // Remove the temporary file when done

	for _, item := range items {
		itemPath := getItemPath(item)
		_, err := fmt.Fprintf(playlistFile, "%s\n", itemPath)
		if err != nil {
			fmt.Println("Error writing to playlist file:", err)
			return err
		}
	}
	playlistFile.Close()

	_, err = conn.Call("loadlist", playlistFile.Name(), "append-play")

	if err != nil {
		log.Print(err.Error())
		return err
	}

	return nil
}
