package dbmodel

import (
	"context"
	"fmt"
	"os"
	"testing"

	"cloud.google.com/go/spanner"
	adminapi "cloud.google.com/go/spanner/admin/database/apiv1"
	"cloud.google.com/go/spanner/admin/database/apiv1/databasepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_RWTx(t *testing.T) {
	ctx := context.Background()

	client := setSpannerClient(ctx, t)
	testcase := map[string]struct {
		query       string
		txErr       error
		execCount   int
		wantErrCode codes.Code
	}{
		"success": {
			query:       `INSERT INTO TestTable (ID, Value) VALUES (2, "test")`,
			execCount:   1,
			wantErrCode: codes.OK,
		},
		"failed: backoff count reaches limit when DeadlineExceeded": {
			txErr:       status.Error(codes.DeadlineExceeded, "deadline exceeded"),
			execCount:   6,
			wantErrCode: codes.DeadlineExceeded,
		},
		"failed: backoff count reaches limit when Unknown": {
			txErr:       status.Error(codes.Unknown, "unknown"),
			execCount:   6,
			wantErrCode: codes.Unknown,
		},
		"failed: internal": {
			txErr:       status.Error(codes.Internal, "internal"),
			execCount:   1,
			wantErrCode: codes.Internal,
		},
		// Cloud Spanner has a built-in Retry process by default, so It will be retried twice.
		// https://github.com/googleapis/google-cloud-go/blob/spanner/v1.61.0/spanner/retry.go#L87
		"failed: already exists": {
			query:       `INSERT INTO TestTable (ID, Value) VALUES (1, "test")`,
			execCount:   2,
			wantErrCode: codes.AlreadyExists,
		},
		"failed: Syntax errors in queries (InvalidArgument)": {
			query:       `INSERT INTO TestTable (ID, Value)`,
			execCount:   2,
			wantErrCode: codes.InvalidArgument,
		},
	}

	rwTx := NewTransactable(client)
	for name, tc := range testcase {
		t.Run(name, func(t *testing.T) {
			execCount := 0
			err := rwTx.RWTx(ctx, func(ctx context.Context) error {
				execCount++
				txn := GetSpannerTransaction(ctx)
				if tc.txErr != nil {
					return tc.txErr
				}
				return txn.ExecContext(ctx, tc.query, nil)
			})

			if spanner.ErrCode(err) != tc.wantErrCode {
				t.Errorf("RWTx() = %v, want %v", spanner.ErrCode(err), tc.wantErrCode)
			}

			if tc.wantErrCode == codes.OK {
				if err != nil {
					t.Errorf("RWTx() = %v, want nil", err)
				}
				return
			}

			if execCount != tc.execCount {
				t.Errorf("execCount = %d, want %d", execCount, tc.execCount)
			}
		})
	}
}

func setSpannerClient(ctx context.Context, t *testing.T) *spanner.Client {
	t.Helper()

	project := os.Getenv("SPANNER_PROJECT_ID")
	instance := os.Getenv("SPANNER_INSTANCE_ID")
	database := os.Getenv("SPANNER_DATABASE_ID")
	databases := fmt.Sprintf("projects/%s/instances/%s/databases/%s", project, instance, database)

	adminClient, err := adminapi.NewDatabaseAdminClient(ctx)
	if err != nil {
		t.Fatalf("Failed to create spanner admin client: %v", err)
	}

	ddl := []string{
		`CREATE TABLE TestTable (
			ID INT64 NOT NULL,
			Value STRING(MAX),
		) PRIMARY KEY (ID)`,
	}

	op, err := adminClient.UpdateDatabaseDdl(ctx, &databasepb.UpdateDatabaseDdlRequest{
		Database:   databases,
		Statements: ddl,
	})
	if err != nil {
		t.Fatalf("Failed to initiate DDL update: %v", err)
	}

	if err := op.Wait(ctx); err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	client, err := spanner.NewClient(ctx, databases)
	if err != nil {
		t.Fatalf("Failed to create spanner client: %v", err)
	}

	_, err = client.Apply(ctx, []*spanner.Mutation{
		spanner.Insert("TestTable", []string{"ID", "Value"}, []interface{}{1, "test"}),
	})
	if err != nil {
		t.Fatalf("Failed to insert default data: %v", err)
	}

	t.Cleanup(func() {
		ddl := []string{
			`DROP TABLE TestTable`,
		}

		op, err := adminClient.UpdateDatabaseDdl(ctx, &databasepb.UpdateDatabaseDdlRequest{
			Database:   databases,
			Statements: ddl,
		})
		if err != nil {
			t.Fatalf("Failed to initiate DDL update: %v", err)
		}
		if err := op.Wait(ctx); err != nil {
			t.Fatalf("Failed to drop table: %v", err)
		}

		defer adminClient.Close()
		defer client.Close()
	})
	return client
}
