package cmds

import (
	"encoding/json"
	"os"
	"os/exec"

	"github.com/SimonVillalonIT/music-golang/internal/services"
	"github.com/charmbracelet/bubbles/list"
	"github.com/spf13/viper"
)

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

func Kill() {
	exec.Command("killall", "mpv").Run()
}
