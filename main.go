package main

import (
	"github.com/SimonVillalonIT/music-golang/cmd"
	_ "github.com/SimonVillalonIT/music-golang/cmd/init"
	_ "github.com/SimonVillalonIT/music-golang/cmd/playlist"
	_ "github.com/SimonVillalonIT/music-golang/cmd/song"
)

func main() {
	cmd.Execute()
}
