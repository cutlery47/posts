package mem

import "errors"

var (
	ErrBadDump    = errors.New("error when dumping")
	ErrBadRestore = errors.New("error when restoring")
)
