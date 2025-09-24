package transactable

import (
	"context"
	"database/sql"
	std_errors "errors"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/avast/retry-go"
	"github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

var ctxTxKey = struct{}{}

// GetContextExecutor :.
func GetContextExecutor(ctx context.Context) boil.ContextExecutor {
	if tx, ok := ctx.Value(&ctxTxKey).(*sql.Tx); ok {
		return tx
	}
	return boil.GetContextDB()
}

// RunTx :.
var RunTx = func(ctx context.Context, fn func(context.Context) error) error {
	db, ok := GetContextExecutor(ctx).(boil.ContextBeginner)
	if !ok {
		panic("The database in the context does not support boil.ContextBeginner")
	}
	return runTxWithDB(ctx, db, fn)
}

// RunTxWithDB :.
func runTxWithDB(ctx context.Context, db boil.ContextBeginner, fn func(context.Context) error) error {
	txFn := func() error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return errors.InternalErr.Wrap(err)
		}
		defer func() {
			if err := recover(); err != nil {
				if err := tx.Rollback(); err != nil {
					panic(err)
				}
				panic(err)
			}
		}()

		ctxWithTx := context.WithValue(ctx, &ctxTxKey, tx)
		if err := fn(ctxWithTx); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.InternalErr.Wrap(rollbackErr)
			}
			return err
		}
		if err := tx.Commit(); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.InternalErr.Wrap(rollbackErr)
			}
			return errors.InternalErr.Wrap(err)
		}
		return nil
	}

	if err := retry.Do(
		txFn,
		retry.RetryIf(isDeadLock),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			return time.Duration(n) * time.Second
		}),
		retry.Attempts(4),
		retry.LastErrorOnly(true),
	); err != nil {
		return err
	}

	return nil
}

func isDeadLock(err error) bool {
	if err == nil {
		return false
	}

	for err != nil {
		var pqErr *pq.Error
		if std_errors.As(err, &pqErr) {
			switch string(pqErr.Code) {
			case "40P01", "55P03", "40001":
				return true
			}
		}
		err = std_errors.Unwrap(err)
	}

	return false
}
