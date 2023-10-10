package graphql

import (
	context "context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/audit"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/auth"
	models "gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
)

type poolDefaultResolver struct{ *Resolver }

func (r *poolDefaultResolver) ID(ctx context.Context, obj *models.PoolDefault) (string, error) {
	return obj.GetID(), nil
}

func (r *poolDefaultResolver) LegCount(ctx context.Context, obj *models.PoolDefault) (string, error) {
	return obj.LegCount.Decimal.String(), nil
}

func (r *poolDefaultResolver) Guarantee(ctx context.Context, obj *models.PoolDefault) (string, error) {
	return obj.Guarantee.Decimal.String(), nil
}

func (r *poolDefaultResolver) CarryIn(ctx context.Context, obj *models.PoolDefault) (string, error) {
	return obj.CarryIn.Decimal.String(), nil
}

func (r *poolDefaultResolver) Allocation(ctx context.Context, obj *models.PoolDefault) (string, error) {
	return obj.Allocation.Decimal.String(), nil
}

func (r *poolDefaultResolver) MinUnitPerLine(ctx context.Context, obj *models.PoolDefault) (string, error) {
	return obj.MinUnitPerLine.Decimal.String(), nil
}

func (r *poolDefaultResolver) MaxUnitPerLine(ctx context.Context, obj *models.PoolDefault) (string, error) {
	return obj.MaxUnitPerLine.Decimal.String(), nil
}

func (r *poolDefaultResolver) MinUnitPerTicket(ctx context.Context, obj *models.PoolDefault) (string, error) {
	return obj.MinUnitPerTicket.Decimal.String(), nil
}

func (r *poolDefaultResolver) MaxUnitPerTicket(ctx context.Context, obj *models.PoolDefault) (string, error) {
	return obj.MaxUnitPerTicket.Decimal.String(), nil
}

func (r *poolDefaultResolver) UnitValue(ctx context.Context, obj *models.PoolDefault) (string, error) {
	return obj.UnitValue.Decimal.String(), nil
}

func (r *queryResolver) PoolDefault(ctx context.Context, id string) (*models.PoolDefault, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
		models.AccessRoleUser,
		models.AccessRoleGuestAPIOnly,
	})

	if err != nil {
		return &models.PoolDefault{}, err
	}

	poolDefault := models.PoolDefault{ID: repository.NewSQLCompatUUIDFromStr(id)}
	err = r.DB.Select(&poolDefault)

	return &poolDefault, err
}

func (r *queryResolver) AllPoolDefaults(ctx context.Context, filter *models.PoolDefaultFilter, page *int, perPage *int, sortField *string, sortOrder *string) ([]*models.PoolDefault, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
		models.AccessRoleUser,
	})
	if err != nil {
		return []*models.PoolDefault{}, err
	}

	var poolDefaults []*models.PoolDefault

	sqlbuilder := sq.Select("pool_defaults.*").From("pool_defaults")
	sql, args := createAllPoolDefaultsSQL(sqlbuilder, filter, page, perPage, sortField, sortOrder)
	_, err = r.DB.Query(&poolDefaults, sql, args...)

	return poolDefaults, err
}

func (r *queryResolver) _allPoolDefaultsMeta(ctx context.Context, filter *models.PoolDefaultFilter, page *int, perPage *int) (*models.ListMetadata, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
		models.AccessRoleUser,
	})
	if err != nil {
		return &models.ListMetadata{}, err
	}

	var count int
	sqlbuilder := sq.Select("COUNT(id)").From("pool_defaults")
	sql, args := createAllPoolDefaultsSQL(sqlbuilder, filter, nil, nil, nil, nil)
	_, err = r.DB.Query(&count, sql, args...)

	return &models.ListMetadata{Count: count}, err
}

func createAllPoolDefaultsSQL(builder sq.SelectBuilder, filter *models.PoolDefaultFilter, page *int, perPage *int, sortField *string, sortOrder *string) (sql string, args []interface{}) {
	if filter != nil {
		if filter.Type != nil {
			builder = builder.Where("pool_defaults.type = ?", filter.Type.String())
		}

		if filter.Game != nil {
			builder = builder.Where("pool_defaults.game = ?", filter.Game.String())
		}

		builder = addStandardPagination(builder, page, perPage)
	}

	builder = addDefaultSort(builder, sortField, sortOrder)

	sql, args, _ = builder.ToSql()
	return sql, args
}

func (r *mutationResolver) CreatePoolDefault(ctx context.Context, input models.CreatePoolDefaultInput) (*models.PoolDefault, error) {
	userID, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})
	if err != nil {
		return &models.PoolDefault{}, err
	}

	poolDefault := models.PoolDefault{
		Type:     input.Type,
		Game:     input.Game,
		Currency: input.Currency,
	}

	poolDefault.LegCount = parseNullDecimal(input.LegCount)
	poolDefault.Guarantee = parseNullDecimal(input.Guarantee)
	poolDefault.CarryIn = parseNullDecimal(input.CarryIn)
	poolDefault.Allocation = parseNullDecimal(input.Allocation)
	poolDefault.UnitValue = parseNullDecimal(input.UnitValue)
	poolDefault.MinUnitPerLine = parseNullDecimal(input.MinUnitPerLine)
	poolDefault.MaxUnitPerLine = parseNullDecimal(input.MaxUnitPerLine)
	poolDefault.MinUnitPerTicket = parseNullDecimal(input.MinUnitPerTicket)
	poolDefault.MaxUnitPerTicket = parseNullDecimal(input.MaxUnitPerTicket)

	err = r.DB.Insert(&poolDefault)
	if err == nil {
		audit.CreateAudit(r.DB, poolDefault.ID, repository.NewSQLCompatUUIDFromStr(userID), models.EditActionCreate, poolDefault)
	}
	return &poolDefault, err
}

func (r *mutationResolver) UpdatePoolDefault(ctx context.Context, input models.UpdatePoolDefaultInput) (*models.PoolDefault, error) {
	userID, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})
	if err != nil {
		return &models.PoolDefault{}, err
	}

	var poolDefault models.PoolDefault

	sqlbuilder := sq.Update("pool_defaults")
	field := fmt.Sprintf("%s.id", "pool_defaults")
	sqlbuilder = sqlbuilder.Where(sq.Eq{field: input.ID})

	if input.Guarantee != nil {
		sqlbuilder = sqlbuilder.Set("guarantee", parseNullDecimal(*input.Guarantee))
	}

	if input.CarryIn != nil {
		sqlbuilder = sqlbuilder.Set("carry_in", parseNullDecimal(*input.CarryIn))
	}

	if input.Allocation != nil {
		sqlbuilder = sqlbuilder.Set("allocation", parseNullDecimal(*input.Allocation))
	}

	if input.UnitValue != nil {
		sqlbuilder = sqlbuilder.Set("unit_value", parseNullDecimal(*input.UnitValue))
	}

	if input.MinUnitPerLine != nil {
		sqlbuilder = sqlbuilder.Set("min_unit_per_line", parseNullDecimal(*input.MinUnitPerLine))
	}

	if input.MaxUnitPerLine != nil {
		sqlbuilder = sqlbuilder.Set("max_unit_per_line", parseNullDecimal(*input.MaxUnitPerLine))
	}

	if input.MinUnitPerTicket != nil {
		sqlbuilder = sqlbuilder.Set("min_unit_per_ticket", parseNullDecimal(*input.MinUnitPerTicket))
	}

	if input.MaxUnitPerTicket != nil {
		sqlbuilder = sqlbuilder.Set("max_unit_per_ticket", parseNullDecimal(*input.MaxUnitPerTicket))
	}

	if input.Currency != nil {
		sqlbuilder = sqlbuilder.Set("currency", (*input.Currency).String())
	}

	if input.Note != nil {
		sqlbuilder = sqlbuilder.Set("note", *input.Note)
	}

	sqlbuilder = sqlbuilder.Suffix("RETURNING *")

	sql, args, _ := sqlbuilder.ToSql()
	_, err = r.DB.Query(&poolDefault, sql, args...)

	if err == nil {
		audit.CreateAudit(r.DB, poolDefault.ID, repository.NewSQLCompatUUIDFromStr(userID), models.EditActionUpdate, poolDefault)
	}

	return &poolDefault, err
}

func (r *mutationResolver) DeletePoolDefault(ctx context.Context, id string) (*models.PoolDefault, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})
	if err != nil {
		return &models.PoolDefault{}, err
	}

	panic("not implemented")
}
