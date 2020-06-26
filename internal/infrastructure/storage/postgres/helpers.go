package postgres

import (
	"database/sql"

	"github.com/dnozdrin/detask/internal/app"
)

func deferred(log app.Logger, f func() error) {
	if err := f(); err != nil && err != sql.ErrTxDone {
		log.Errorf("%v", err)
	}
}
