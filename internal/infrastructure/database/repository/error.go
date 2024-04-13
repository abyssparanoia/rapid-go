package repository

import (
	"context"
	goerrors "errors"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
)

func handleError(err error) error {
	if err == nil {
		return nil
	}

	if goerrors.Is(err, context.Canceled) {
		return errors.CanceledErr.Wrap(err)
	}
	return errors.InternalErr.Wrap(err)
}
