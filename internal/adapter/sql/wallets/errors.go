package wallets

import (
	"errors"
	"fmt"
)

var ErrWalletNotFound = errors.New("wallet was not found")

type DatabaseError struct {
	Query string
	Err   error
}

func (e *DatabaseError) Error() string { return fmt.Sprintf("%s %s", e.Query, e.Err.Error()) }
