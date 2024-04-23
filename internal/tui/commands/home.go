package cmds

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/DexterLB/mpvipc"
	"github.com/SimonVillalonIT/music-golang/internal/services"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/viper"
)

type MpvMsg bool
type ConnMsg mpvipc.Connection
type PlayMsg mpvipc.Connection

type ErrMsg error
type ConnErr error

func StartMpv() tea.Msg {
	if !isMpvRunning() {
		err := exec.Command("mpv", "--no-video", "--idle", "--really-quiet", "--input-ipc-server=/tmp/mpvsocket").Run()

		if err != nil {
			log.Printf("Failed to start mpv: %v", err)

			return ErrMsg(err)
		}
	}

	return MpvMsg(true)
}

func isMpvRunning() bool {
	cmd := exec.Command("ps", "-A")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Printf("Error checking for running processes: %v", err)
		return false
	}
	return strings.Contains(out.String(), "mpv")
}

func SearchSongs() ([]services.Item, []list.Item, error) {
	rawFile, err := os.ReadFile(viper.GetString(services.STORE_PATH))
	if err != nil {
		return nil, nil, err
	}

	var jsonFile []services.Item

	err = json.Unmarshal(rawFile, &jsonFile)
	if err != nil {
		return nil, nil, err
	}

	items := make([]list.Item, len(jsonFile))

	for i, song := range jsonFile {
		items[i] = song
	}

	return jsonFile, items, nil
}

func ConnCmd() tea.Msg {
	conn := mpvipc.NewConnection("/tmp/mpvsocket")

	err := conn.Open()

	if conn.IsClosed() || conn == nil {
		return ConnErr(err)
	}

	return ConnMsg(*conn)
}

func PlayCmd(conn mpvipc.Connection, item services.Item) tea.Cmd {
	return func() tea.Msg {
		err := services.Play(conn, item)
		if err != nil {
			return ErrMsg(err)
		}
		return PlayMsg(conn)
	}
}

func Kill() {
	exec.Command("killall", "mpv").Run()
}
