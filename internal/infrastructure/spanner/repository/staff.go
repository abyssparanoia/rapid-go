package repository

import (
	"context"

	"cloud.google.com/go/spanner"
	"github.com/abyssparanoia/memeduck"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/spanner/internal/dbmodel"
	"github.com/abyssparanoia/rapid-go/internal/infrastructure/spanner/internal/marshaller"
	"google.golang.org/grpc/codes"
)

type staff struct{}

func NewStaff() repository.Staff {
	return &staff{}
}

func (r *staff) Get(
	ctx context.Context,
	query repository.GetStaffQuery,
) (*model.Staff, error) {
	conds := []memeduck.WhereCond{}
	params := map[string]interface{}{}
	if query.ID.Valid {
		conds = append(conds, memeduck.Eq(memeduck.Ident("StaffID"), memeduck.Param("StaffID")))
		params["StaffID"] = query.ID.String
	}
	if query.AuthUID.Valid {
		conds = append(conds, memeduck.Eq(memeduck.Ident("AuthUID"), memeduck.Param("AuthUID")))
		params["AuthUID"] = query.AuthUID.String
	}
	sql, err := memeduck.Select(
		dbmodel.StaffTableName(),
		dbmodel.StaffColumns(),
	).
		Where(conds...).
		Limit(1).
		SQL()
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	rows, err := dbmodel.GetSpannerTransaction(ctx).QueryContext(ctx, sql, params)
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	defer func() { _ = rows.Close() }()

	if ok, err := rows.Next(); err != nil && spanner.ErrCode(err) != codes.NotFound {
		return nil, errors.InternalErr.Wrap(err)
	} else if !ok {
		if !query.OrFail {
			return nil, nil
		} else {
			return nil, errors.StaffNotFoundErr.New().
				WithDetail("staff is not found").
				WithValue("query", query)
		}
	}

	var dst dbmodel.Staff
	if err := rows.ToStruct(&dst); err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}

	return marshaller.StaffToModel(&dst), nil
}

func (r *staff) Create(
	ctx context.Context,
	staff *model.Staff,
) error {
	if err := marshaller.StaffToDBModel(staff).Insert(ctx); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}
