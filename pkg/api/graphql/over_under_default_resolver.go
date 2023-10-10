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

type overUnderDefaultResolver struct{ *Resolver }

func (r *overUnderDefaultResolver) ID(ctx context.Context, obj *models.OverUnderDefault) (string, error) {
	return obj.GetID(), nil
}

func (r *overUnderDefaultResolver) Game(ctx context.Context, obj *models.OverUnderDefault) (string, error) {
	return obj.Game.String(), nil
}

func (r *overUnderDefaultResolver) MatchFormat(ctx context.Context, obj *models.OverUnderDefault) (string, error) {
	return obj.MatchFormat.String(), nil
}

func (r *overUnderDefaultResolver) EvenThreshold(ctx context.Context, obj *models.OverUnderDefault) (string, error) {
	return obj.EvenThreshold.Decimal.String(), nil
}

func (r *overUnderDefaultResolver) FavoredThreshold(ctx context.Context, obj *models.OverUnderDefault) (string, error) {
	return obj.FavoredThreshold.Decimal.String(), nil
}

func (r *overUnderDefaultResolver) Note(ctx context.Context, obj *models.OverUnderDefault) (string, error) {
	return obj.Note, nil
}

func (r *queryResolver) OverUnderDefault(ctx context.Context, id string) (*models.OverUnderDefault, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
		models.AccessRoleUser,
		models.AccessRoleGuestAPIOnly,
	})

	if err != nil {
		return &models.OverUnderDefault{}, err
	}

	overUnderDefault := models.OverUnderDefault{ID: repository.NewSQLCompatUUIDFromStr(id)}
	err = r.DB.Select(&overUnderDefault)

	return &overUnderDefault, err
}

func (r *queryResolver) AllOverUnderDefaults(ctx context.Context, filter *models.OverUnderDefaultFilter, page *int, perPage *int, sortField *string, sortOrder *string) ([]*models.OverUnderDefault, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
		models.AccessRoleUser,
	})
	if err != nil {
		return []*models.OverUnderDefault{}, err
	}

	var overUnderDefaults []*models.OverUnderDefault

	sqlbuilder := sq.Select("over_under_defaults.*").From("over_under_defaults")
	sql, args := createAllOverUnderDefaultsSQL(sqlbuilder, filter, page, perPage, sortField, sortOrder)
	_, err = r.DB.Query(&overUnderDefaults, sql, args...)

	return overUnderDefaults, err
}

func (r *queryResolver) _allOverUnderDefaultsMeta(ctx context.Context, filter *models.OverUnderDefaultFilter, page *int, perPage *int) (*models.ListMetadata, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
		models.AccessRoleUser,
	})
	if err != nil {
		return &models.ListMetadata{}, err
	}

	var count int
	sqlbuilder := sq.Select("COUNT(id)").From("over_under_defaults")
	sql, args := createAllOverUnderDefaultsSQL(sqlbuilder, filter, nil, nil, nil, nil)
	_, err = r.DB.Query(&count, sql, args...)

	return &models.ListMetadata{Count: count}, err
}

func createAllOverUnderDefaultsSQL(builder sq.SelectBuilder, filter *models.OverUnderDefaultFilter, page *int, perPage *int, sortField *string, sortOrder *string) (sql string, args []interface{}) {
	if filter != nil {
		if filter.MatchFormat != nil {
			builder = builder.Where("over_under_defaults.match_format = ?", filter.MatchFormat.String())
		}

		if filter.Game != nil {
			builder = builder.Where("over_under_defaults.game = ?", filter.Game.String())
		}

		builder = addStandardPagination(builder, page, perPage)
	}

	builder = addDefaultSort(builder, sortField, sortOrder)

	sql, args, _ = builder.ToSql()
	return sql, args
}

func (r *mutationResolver) CreateOverUnderDefault(ctx context.Context, input models.CreateOverUnderDefaultInput) (*models.OverUnderDefault, error) {
	userID, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})
	if err != nil {
		return &models.OverUnderDefault{}, err
	}

	overUnderDefault := models.OverUnderDefault{
		Game:        input.Game,
		MatchFormat: input.MatchFormat,
	}

	if input.Note != nil {
		overUnderDefault.Note = *input.Note
	}

	overUnderDefault.EvenThreshold = parseNullDecimal(input.EvenThreshold)
	overUnderDefault.FavoredThreshold = parseNullDecimal(input.FavoredThreshold)

	err = r.DB.Insert(&overUnderDefault)
	if err == nil {
		audit.CreateAudit(r.DB, overUnderDefault.ID, repository.NewSQLCompatUUIDFromStr(userID), models.EditActionCreate, overUnderDefault)
	}
	return &overUnderDefault, err
}

func (r *mutationResolver) UpdateOverUnderDefault(ctx context.Context, input models.UpdateOverUnderDefaultInput) (*models.OverUnderDefault, error) {
	userID, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})
	if err != nil {
		return &models.OverUnderDefault{}, err
	}

	var overUnderDefault models.OverUnderDefault

	sqlbuilder := sq.Update("over_under_defaults")
	field := fmt.Sprintf("%s.id", "over_under_defaults")
	sqlbuilder = sqlbuilder.Where(sq.Eq{field: input.ID})

	if input.EvenThreshold != nil {
		sqlbuilder = sqlbuilder.Set("even_threshold", parseNullDecimal(*input.EvenThreshold))
	}

	if input.FavoredThreshold != nil {
		sqlbuilder = sqlbuilder.Set("favored_threshold", parseNullDecimal(*input.FavoredThreshold))
	}

	if input.Note != nil {
		sqlbuilder = sqlbuilder.Set("note", *input.Note)
	}

	sqlbuilder = sqlbuilder.Suffix("RETURNING *")

	sql, args, _ := sqlbuilder.ToSql()
	_, err = r.DB.Query(&overUnderDefault, sql, args...)

	if err == nil {
		audit.CreateAudit(r.DB, overUnderDefault.ID, repository.NewSQLCompatUUIDFromStr(userID), models.EditActionUpdate, overUnderDefault)
	}

	return &overUnderDefault, err
}

func (r *mutationResolver) DeleteOverUnderDefault(ctx context.Context, id string) (*models.OverUnderDefault, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})
	if err != nil {
		return &models.OverUnderDefault{}, err
	}

	panic("not implemented")
}
