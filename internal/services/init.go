package services

import (
	"os/exec"
)

func StartMpv() {
	exec.Command("mpv", "--no-video", "--idle", "--really-quiet", "--input-ipc-server=/tmp/mpvsocket").Run()
}
