package cmds

import (
	"log"

	"github.com/DexterLB/mpvipc"
	tea "github.com/charmbracelet/bubbletea"
)

func ObserveTrack(conn *mpvipc.Connection) tea.Cmd {
	return func() tea.Msg {
		if conn == nil {
			return RetryMsg("retry")
		}
		result, err := conn.Call("get_property_string", "filename")
		if err != nil {
			return RetryMsg("retry")
		}
		if result == nil {
			return RetryMsg("retry")
		}
		log.Print(result)
		return EventMsg(result.(string))
	}
}

type EventMsg string
type RetryMsg string
