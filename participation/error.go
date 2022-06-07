package participation

import "errors"

var (
	ErrNotHost        = errors.New("only host are authorized for such actions")
	ErrNotMember      = errors.New("only members are authorized for such actions")
	ErrAlreadyJoined  = errors.New("already joined")
	ErrAlreadyApplied = errors.New("already applied")
)
