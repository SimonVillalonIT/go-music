package cmds

import "github.com/DexterLB/mpvipc"

type ConnMsg *mpvipc.Connection
type PlayMsg mpvipc.Connection
type TrackNameMsg string
type TrackCurrentFrameMsg float64
type TrackDurationMsg float64
type PlaylistPositionMsg float64
type PlaylistLengthMsg float64

type ErrMsg error
