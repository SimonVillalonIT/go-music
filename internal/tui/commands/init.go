package cmds

import (
	"log"

	"github.com/DexterLB/mpvipc"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	NAME = iota
	CURRENT_FRAME
	DURATION
	POSITION
	LENGTH
)

func SetupConnection(p *tea.Program) {
	conn := mpvipc.NewConnection("/tmp/mpvsocket")

	err := conn.Open()
	for err != nil {
		err = conn.Open()
	}

	_, err = conn.Call("observe_property", NAME, "filename/no-ext")
	if err != nil {
		log.Print(err)
	}

	_, err = conn.Call("observe_property", CURRENT_FRAME, "audio-pts")

	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Call("observe_property", DURATION, "duration")

	if err != nil {
		log.Print(err)
	}
	_, err = conn.Call("observe_property", POSITION, "playlist-pos-1")

	if err != nil {
		log.Print(err)
	}
	_, err = conn.Call("observe_property", LENGTH, "playlist-count")

	if err != nil {
		log.Print(err)
	}

	p.Send(ConnMsg(conn))

	events, stopListening := conn.NewEventListener()

	go func() {
		conn.WaitUntilClosed()
	}()

	for event := range events {
		if event.ID == NAME && event.Data != nil {
			p.Send(TrackNameMsg(event.Data.(string)))
		}
		if event.ID == CURRENT_FRAME && event.Data != nil {
			p.Send(TrackCurrentFrameMsg(event.Data.(float64)))
		}
		if event.ID == DURATION && event.Data != nil {
			p.Send(TrackDurationMsg(event.Data.(float64)))
		}
		if event.ID == POSITION && event.Data != nil {
			p.Send(PlaylistPositionMsg(event.Data.(float64)))
		}
		if event.ID == LENGTH && event.Data != nil {
			p.Send(PlaylistLengthMsg(event.Data.(float64)))
		}
	}
	stopListening <- struct{}{}
}
