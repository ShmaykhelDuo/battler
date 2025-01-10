package game

import "errors"

var ErrHasConnection = errors.New("user is already connected to a match")
var ErrChanClosed = errors.New("channel closed")
