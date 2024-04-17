build:
	go build -o ~/.local/bin/music-golang main.go

dev:
	go build -o ~/.local/bin/music-golang main.go & music-golang play
