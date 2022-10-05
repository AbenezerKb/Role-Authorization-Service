package sqlcerr

import "errors"

var (
	ErrNoRows = errors.New("no rows in result set")
)

func Is(err, target error) bool {
	return err.Error() == target.Error() // FIXME: a better way to do this?
}
