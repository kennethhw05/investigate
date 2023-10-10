package graphql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/audit"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/auth"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
)

type competitorResolver struct{ *Resolver }

func (r *competitorResolver) ID(ctx context.Context, obj *models.Competitor) (string, error) {
	return obj.GetID(), nil
}

func (r *queryResolver) Competitor(ctx context.Context, id string) (*models.Competitor, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
		models.AccessRoleUser,
		models.AccessRoleGuestAPIOnly,
	})
	if err != nil {
		return &models.Competitor{}, err
	}

	competitor := models.Competitor{}
	sqlbuilder := sq.Select("competitors.*, array_to_json(array_agg(matches)) as matches").From("competitors").
		Join("competitor_match on competitor_match.competitor_id = competitors.id").
		Join("matches on matches.id = competitor_match.match_id").Where(sq.Eq{fmt.Sprintf("%s.id", "competitors"): id}).GroupBy("competitors.id")

	sql, args, _ := sqlbuilder.ToSql()
	_, err = r.DB.QueryOne(&competitor, sql, args...)

	return &competitor, err
}

func (r *queryResolver) AllCompetitors(ctx context.Context, filter *models.CompetitorFilter, page *int, perPage *int, sortField *string, sortOrder *string) ([]*models.Competitor, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
		models.AccessRoleUser,
	})
	if err != nil {
		return []*models.Competitor{}, err
	}

	var competitors []*models.Competitor
	sqlbuilder := sq.Select("competitors.*, array_to_json(array_agg(matches)) as matches").From("competitors").
		Join("competitor_match on competitor_match.competitor_id = competitors.id").
		Join("matches on matches.id = competitor_match.match_id").GroupBy("competitors.id")
	sql, args := createAllCompetitorsSQL(sqlbuilder, filter, page, perPage, sortField, sortOrder)
	_, err = r.DB.Query(&competitors, sql, args...)
	return competitors, err
}

func (r *queryResolver) _allCompetitorsMeta(ctx context.Context, filter *models.CompetitorFilter, page *int, perPage *int) (*models.ListMetadata, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
		models.AccessRoleUser,
	})
	if err != nil {
		return &models.ListMetadata{}, err
	}

	var count int

	sqlbuilder := sq.Select("Count(id)").From("competitors")
	sql, args := createAllCompetitorsSQL(sqlbuilder, filter, nil, nil, nil, nil)
	_, err = r.DB.Query(&count, sql, args...)
	return &models.ListMetadata{Count: count}, err
}

func createAllCompetitorsSQL(builder sq.SelectBuilder, filter *models.CompetitorFilter, page *int, perPage *int, sortField *string, sortOrder *string) (sql string, args []interface{}) {
	if filter != nil {
		if filter.ExternalID != nil {
			builder = builder.Where("competitors.external_id ILIKE ?", fmt.Sprint("%", *filter.ExternalID, "%"))
		}

		if filter.Name != nil {
			builder = builder.Where("competitors.name ILIKE ?", fmt.Sprint("%", *filter.Name, "%"))
		}

		// if filter.MatchID != nil {
		// 	builder = builder.Where("competitor.match_id = ?", *filter.MatchID)
		// }

		if filter.ID != nil {
			filter.Ids = append(filter.Ids, *filter.ID)
		}

		builder = filterByIDs(builder, filter.Ids, "competitors")
	}
	builder = addDefaultSort(builder, sortField, sortOrder)
	builder = addStandardPagination(builder, page, perPage)

	sql, args, _ = builder.ToSql()
	return sql, args
}

func (r *mutationResolver) CreateCompetitor(ctx context.Context, input models.CreateCompetitorInput) (*models.Competitor, error) {
	userID, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})
	if err != nil {
		return &models.Competitor{}, err
	}

	competitor := models.Competitor{
		Name: input.Name,
		//TODO MatchID: repository.NewSQLCompatUUIDFromStr(input.MatchID),
	}
	if input.Logo != nil {
		competitor.Logo = *input.Logo
	}
	err = r.DB.Insert(&competitor)

	if err != nil {
		return &models.Competitor{}, err
	}

	competitor.ExternalID = models.GenerateInternalXID(&models.Competitor{}, competitor.GetID())
	_, err = r.DB.Model(&competitor).WherePK().Update()

	if err != nil {
		r.Logger.Errorf("Error while updating custom competitor's external id: %s", err.Error())
	}

	//TODO
	// match := models.Match{ID: team.MatchID}

	// err = match.AddTeam(team, r.DB)

	// if err != nil {
	// 	r.Logger.Errorf("Error while updating custom team's match statistics: %s", err.Error())
	// }

	if err == nil {
		audit.CreateAudit(r.DB, competitor.ID, repository.NewSQLCompatUUIDFromStr(userID), models.EditActionCreate, competitor)
	}
	return &competitor, err
}

func (r *mutationResolver) UpdateCompetitor(ctx context.Context, input models.UpdateCompetitorInput) (*models.Competitor, error) {
	userID, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})
	if err != nil {
		return &models.Competitor{}, err
	}

	var competitor models.Competitor

	sqlbuilder := sq.Update("competitors")
	field := fmt.Sprintf("%s.id", "competitors")
	sqlbuilder = sqlbuilder.Where(sq.Eq{field: input.ID})

	if input.Name != nil {
		sqlbuilder = sqlbuilder.Set("name", input.Name)
	}

	if input.Logo != nil {
		sqlbuilder = sqlbuilder.Set("logo", input.Logo)
	}

	sqlbuilder = sqlbuilder.Suffix("RETURNING *")
	sql, args, _ := sqlbuilder.ToSql()
	_, err = r.DB.Query(&competitor, sql, args...)

	if err == nil {
		audit.CreateAudit(r.DB, competitor.ID, repository.NewSQLCompatUUIDFromStr(userID), models.EditActionUpdate, competitor)
	}

	return &competitor, err
}

func (r *mutationResolver) DeleteCompetitor(ctx context.Context, id string) (*models.Competitor, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})
	if err != nil {
		return &models.Competitor{}, err
	}

	panic("not implemented")
}
