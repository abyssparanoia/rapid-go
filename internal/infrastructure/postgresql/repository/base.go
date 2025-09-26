package repository

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func addForUpdateFromBaseGetOptions(mods []qm.QueryMod, query repository.BaseGetOptions) []qm.QueryMod {
	if query.ForUpdate {
		if query.SkipLocked {
			mods = append(mods, qm.For("UPDATE SKIP LOCKED"))
		} else {
			mods = append(mods, qm.For("UPDATE"))
		}
	}
	return mods
}

func addForUpdateFromBaseListOptions(mods []qm.QueryMod, query repository.BaseListOptions) []qm.QueryMod {
	if query.ForUpdate {
		if query.SkipLocked {
			mods = append(mods, qm.For("UPDATE SKIP LOCKED"))
		} else {
			mods = append(mods, qm.For("UPDATE"))
		}
	}
	return mods
}
