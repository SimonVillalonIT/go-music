package cmds

import (
	"github.com/DexterLB/mpvipc"
	"github.com/SimonVillalonIT/music-golang/internal/services"
	tea "github.com/charmbracelet/bubbletea"
)

type Playlist struct {
	CurrentTrack Track
	Position     float64
	Length       float64
}

type Track struct {
	Name         string
	Duration     float64
	CurrentFrame float64
}

func PlayCmd(conn *mpvipc.Connection, item services.Item) tea.Cmd {
	return func() tea.Msg {
		err := services.Play(conn, item)
		if err != nil {
			return ErrMsg(err)
		}
		return nil
	}
}

func PauseCmd(conn *mpvipc.Connection) tea.Cmd {
	return func() tea.Msg {
		val, err := conn.Get("pause")

		if err != nil {
			return ErrMsg(err)
		}
		err = conn.Set("pause", !val.(bool))
		if err != nil {
			return ErrMsg(err)
		}
		return nil
	}
}

func DecreaseCmd(conn *mpvipc.Connection) tea.Cmd {
	return func() tea.Msg {
		_, err := conn.Call("osd-auto", "add", "volume", -2)
		if err != nil {
			return ErrMsg(err)
		}
		return nil
	}
}

func IncreaseCmd(conn *mpvipc.Connection) tea.Cmd {
	return func() tea.Msg {
		_, err := conn.Call("osd-auto", "add", "volume", 2)
		if err != nil {
			return ErrMsg(err)
		}
		return nil
	}
}
