package cmds

import (
	"github.com/DexterLB/mpvipc"
	"github.com/SimonVillalonIT/music-golang/internal/services"
)

type ConnMsg *mpvipc.Connection
type PlayMsg mpvipc.Connection
type TrackNameMsg string
type TrackCurrentFrameMsg float64
type TrackDurationMsg float64
type PlaylistPositionMsg float64
type PlaylistLengthMsg float64
type PlaylistInfoMsg []interface{}
type SearchResultMsg []services.Item
type DownloadSuccessMsg bool
type DownloadErrorMsg string
type SearchSuccessMsg bool
type QuitSearchMsg bool

type StateMsg *uint

type ErrMsg error
type ClearDownloadError bool
