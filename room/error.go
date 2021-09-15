package room

import "errors"

var (
	ErrMaxCountExceed        = errors.New("max count exceed")
	ErrNotHost               = errors.New("only host are authorized for such actions")
	ErrNotMember             = errors.New("only members are authorized for such actions")
	ErrNotAuthorized         = errors.New("user is not authorized for such actions")
	ErrRoomFull              = errors.New("room is full")
	ErrInvalidInvitationCode = errors.New("invalid invitation code")
	ErrAlreadyJoined         = errors.New("already joined")
	ErrRoundAlreadyCreated   = errors.New("already created round")
	ErrNoRound               = errors.New("no round created")
	ErrNotPublic             = errors.New("room not public")
	ErrNotStarted            = errors.New("room not started yet")
	ErrAlreadyStarted        = errors.New("room already started")
)
