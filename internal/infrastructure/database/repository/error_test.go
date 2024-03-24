package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
)

func Test_handleError(t *testing.T) {
	testcase := map[string]struct {
		err  error
		want error
	}{
		"no error": {
			err:  nil,
			want: nil,
		},
		"cancelled error": {
			err:  context.Canceled,
			want: errors.CanceledErr.Wrap(context.Canceled),
		},
		"internal error": {
			err:  errors.InternalErr.New(),
			want: errors.InternalErr.Wrap(errors.InternalErr.New()),
		},
		"other error": {
			err:  fmt.Errorf("other error"),
			want: errors.InternalErr.Wrap(fmt.Errorf("other error")),
		},
	}

	for name, tc := range testcase {
		t.Run(name, func(t *testing.T) {
			got := handleError(tc.err)
			if got == nil {
				if tc.want != nil {
					t.Fatalf("got = nil, want %v", tc.want)
				}
				return
			}

			if tc.want.Error() != got.Error() {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}
