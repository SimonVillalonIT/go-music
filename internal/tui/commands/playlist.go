package cmds

import (
	"github.com/DexterLB/mpvipc"
	tea "github.com/charmbracelet/bubbletea"
)

func PrevCmd(conn *mpvipc.Connection) tea.Cmd {
	return func() tea.Msg {
		_, err := conn.Call("osd-auto", "playlist-prev")
		if err != nil {
			return ErrMsg(err)
		}
		return nil
	}
}

func NextCmd(conn *mpvipc.Connection) tea.Cmd {
	return func() tea.Msg {
		_, err := conn.Call("osd-auto", "playlist-next")
		if err != nil {
			return ErrMsg(err)
		}
		return nil
	}
}

//
// func getPlaylistCmd(m *playlist.Model) tea.Msg {
// 	result, err := m.Conn.Get("playlist")
//
// 	if err != nil {
// 		return ErrMsg(err)
// 	}
//
//     for _, result = m.Items
//
// 	if err != nil {
// 		return ErrMsg(err)
// 	}
//
// 	return nil
// }
