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

type legResolver struct{ *Resolver }

func (r *legResolver) ID(ctx context.Context, obj *models.Leg) (string, error) {
	return obj.GetID(), nil
}

func (r *legResolver) MatchID(ctx context.Context, obj *models.Leg) (string, error) {
	return obj.GetMatchID(), nil
}

func (r *legResolver) PoolID(ctx context.Context, obj *models.Leg) (string, error) {
	return obj.GetPoolID(), nil
}

func (r *legResolver) Threshold(ctx context.Context, obj *models.Leg) (string, error) {
	return obj.Threshold.Decimal.String(), nil
}

func (r *queryResolver) Leg(ctx context.Context, id string) (*models.Leg, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
		models.AccessRoleUser,
		models.AccessRoleGuestAPIOnly,
	})
	if err != nil {
		return &models.Leg{}, err
	}
	leg := models.Leg{}
	sqlmatches := sq.Select("legs.*").From("legs").Where(sq.Eq{"legs.id": id})

	sql, args, _ := sqlmatches.ToSql()
	_, err = r.DB.QueryOne(&leg, sql, args...)

	return &leg, err
}

func (r *queryResolver) AllLegs(ctx context.Context, filter *models.LegFilter, page *int, perPage *int, sortField *string, sortOrder *string) ([]*models.Leg, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
		models.AccessRoleUser,
	})
	if err != nil {
		return []*models.Leg{}, err
	}
	var legs []*models.Leg

	sqlbuilder := sq.Select("legs.*").From("legs")
	sql, args := createAllLegsSQL(sqlbuilder, filter, page, perPage, sortField, sortOrder)
	_, err = r.DB.Query(&legs, sql, args...)
	return legs, err
}

func (r *queryResolver) _allLegsMeta(ctx context.Context, filter *models.LegFilter, page *int, perPage *int) (*models.ListMetadata, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
		models.AccessRoleUser,
	})
	if err != nil {
		return &models.ListMetadata{}, err
	}

	var count int
	sqlbuilder := sq.Select("COUNT(id)").From("legs")
	sql, args := createAllLegsSQL(sqlbuilder, filter, nil, nil, nil, nil)
	_, err = r.DB.Query(&count, sql, args...)
	return &models.ListMetadata{Count: count}, err
}

func createAllLegsSQL(builder sq.SelectBuilder, filter *models.LegFilter, page *int, perPage *int, sortField *string, sortOrder *string) (sql string, args []interface{}) {
	if filter != nil {
		if filter.MatchID != nil {
			builder = builder.Where("legs.match_id = ?", *filter.MatchID)
		}

		if filter.PoolID != nil {
			builder = builder.Where("legs.pool_id = ?", *filter.PoolID)
		}

		if filter.MatchID != nil {
			builder = builder.Where("legs.match_id = ?", *filter.MatchID)
		}

		if filter.ID != nil {
			filter.Ids = append(filter.Ids, *filter.ID)
		}

		builder = filterByIDs(builder, filter.Ids, "legs")
		builder = addDefaultSort(builder, sortField, sortOrder)
		builder = addStandardPagination(builder, page, perPage)
	}
	sql, args, _ = builder.ToSql()

	return sql, args
}

func (r *mutationResolver) CreateLeg(ctx context.Context, input models.CreateLegInput) (*models.Leg, error) {
	userID, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})
	if err != nil {
		return &models.Leg{}, err
	}

	leg := models.Leg{
		MatchID: repository.NewSQLCompatUUIDFromStr(input.MatchID),
		PoolID:  repository.NewSQLCompatUUIDFromStr(input.PoolID),
	}

	if input.Threshold != nil {
		leg.Threshold = parseNullDecimal(*input.Threshold)
	}

	err = r.DB.Insert(&leg)

	if err == nil {
		audit.CreateAudit(r.DB, leg.ID, repository.NewSQLCompatUUIDFromStr(userID), models.EditActionCreate, leg)
	}
	return &leg, err
}

func (r *mutationResolver) UpdateLeg(ctx context.Context, input models.UpdateLegInput) (*models.Leg, error) {
	userID, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})
	if err != nil {
		return &models.Leg{}, err
	}
	var leg models.Leg

	sqlbuilder := sq.Update("legs")
	field := fmt.Sprintf("%s.id", "legs")
	sqlbuilder = sqlbuilder.Where(sq.Eq{field: input.ID})

	if input.MatchID != nil {
		sqlbuilder = sqlbuilder.Set("match_id", repository.NewSQLCompatUUIDFromStr(*input.MatchID))
	}

	if input.PoolID != nil {
		sqlbuilder = sqlbuilder.Set("pool_id", repository.NewSQLCompatUUIDFromStr(*input.PoolID))
	}

	if input.Threshold != nil {
		sqlbuilder = sqlbuilder.Set("threshold", parseNullDecimal(*input.Threshold))
	}

	sqlbuilder = sqlbuilder.Suffix("RETURNING *")
	sql, args, _ := sqlbuilder.ToSql()
	_, err = r.DB.Query(&leg, sql, args...)

	if err == nil {
		audit.CreateAudit(r.DB, leg.ID, repository.NewSQLCompatUUIDFromStr(userID), models.EditActionUpdate, leg)
	}

	return &leg, err
}

func (r *mutationResolver) DeleteLeg(ctx context.Context, id string) (*models.Leg, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})
	if err != nil {
		return &models.Leg{}, err
	}
	panic("not implemented")
}
