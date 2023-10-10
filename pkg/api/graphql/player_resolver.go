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

type playerResolver struct{ *Resolver }

func (r *playerResolver) ID(ctx context.Context, obj *models.Player) (string, error) {
	return obj.GetID(), nil
}

func (r *playerResolver) TeamID(ctx context.Context, obj *models.Player) (string, error) {
	return obj.GetTeamID(), nil
}

func (r *queryResolver) Player(ctx context.Context, id string) (*models.Player, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
		models.AccessRoleUser,
		models.AccessRoleGuestAPIOnly,
	})
	if err != nil {
		return &models.Player{}, err
	}
	player := models.Player{}
	sqlmatches := sq.Select("players.*").From("players").Where(sq.Eq{"players.id": id})

	sql, args, _ := sqlmatches.ToSql()
	_, err = r.DB.QueryOne(&player, sql, args...)

	return &player, err
}

func (r *queryResolver) AllPlayers(ctx context.Context, filter *models.PlayerFilter, page *int, perPage *int, sortField *string, sortOrder *string) ([]*models.Player, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
		models.AccessRoleUser,
	})
	if err != nil {
		return []*models.Player{}, err
	}

	var players []*models.Player

	sqlbuilder := sq.Select("players.*").From("players")
	sql, args := createAllPlayersSQL(sqlbuilder, filter, page, perPage, sortField, sortOrder)
	_, err = r.DB.Query(&players, sql, args...)

	return players, err
}

func (r *queryResolver) _allPlayersMeta(ctx context.Context, filter *models.PlayerFilter, page *int, perPage *int) (*models.ListMetadata, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
		models.AccessRoleUser,
	})
	if err != nil {
		return &models.ListMetadata{}, err
	}

	var count int
	sqlbuilder := sq.Select("COUNT(id)").From("players")
	sql, args := createAllPlayersSQL(sqlbuilder, filter, nil, nil, nil, nil)
	_, err = r.DB.Query(&count, sql, args...)
	return &models.ListMetadata{Count: count}, err
}

func createAllPlayersSQL(builder sq.SelectBuilder, filter *models.PlayerFilter, page *int, perPage *int, sortField *string, sortOrder *string) (sql string, args []interface{}) {

	if filter != nil {
		if filter.ExternalID != nil {
			builder = builder.Where("players.external_id ILIKE ?", fmt.Sprint("%", *filter.ExternalID, "%"))
		}

		if filter.Name != nil {
			builder = builder.Where("players.name ILIKE ?", fmt.Sprint("%", *filter.Name, "%"))
		}

		if filter.Nickname != nil {
			builder = builder.Where("players.nickname ILIKE ?", fmt.Sprint("%", *filter.Nickname, "%"))
		}

		if filter.TeamID != nil {
			builder = builder.Where("players.team_id = ?", *filter.TeamID)
		}

		if filter.ID != nil {
			filter.Ids = append(filter.Ids, *filter.ID)
		}

		builder = filterByIDs(builder, filter.Ids, "players")
		builder = addStandardPagination(builder, page, perPage)
	}

	builder = addDefaultSort(builder, sortField, sortOrder)

	sql, args, _ = builder.ToSql()
	return sql, args
}

func (r *mutationResolver) CreatePlayer(ctx context.Context, input models.CreatePlayerInput) (*models.Player, error) {
	userID, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})
	if err != nil {
		return &models.Player{}, err
	}

	player := models.Player{
		Name:   input.Name,
		TeamID: repository.NewSQLCompatUUIDFromStr(input.TeamID),
	}

	if input.Nickname != nil {
		player.Nickname = *input.Nickname
	}

	err = r.DB.Insert(&player)

	if err != nil {
		return &models.Player{}, err
	}

	player.ExternalID = models.GenerateInternalXID(&models.Player{}, player.GetID())
	_, err = r.DB.Model(&player).WherePK().Update()

	if err == nil {
		audit.CreateAudit(r.DB, player.ID, repository.NewSQLCompatUUIDFromStr(userID), models.EditActionCreate, player)
	}

	return &player, err
}

func (r *mutationResolver) UpdatePlayer(ctx context.Context, input models.UpdatePlayerInput) (*models.Player, error) {
	userID, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})
	if err != nil {
		return &models.Player{}, err
	}

	var player models.Player

	sqlbuilder := sq.Update("players")
	field := fmt.Sprintf("%s.id", "players")
	sqlbuilder = sqlbuilder.Where(sq.Eq{field: input.ID})

	if input.Name != nil {
		sqlbuilder = sqlbuilder.Set("name", *input.Name)
	}

	if input.Nickname != nil {
		sqlbuilder = sqlbuilder.Set("nickname", *input.Nickname)
	}

	if input.TeamID != nil {
		sqlbuilder = sqlbuilder.Set("team_id", repository.NewSQLCompatUUIDFromStr(*input.TeamID))
	}

	sqlbuilder = sqlbuilder.Suffix("RETURNING *")
	sql, args, _ := sqlbuilder.ToSql()
	_, err = r.DB.Query(&player, sql, args...)

	if err == nil {
		audit.CreateAudit(r.DB, player.ID, repository.NewSQLCompatUUIDFromStr(userID), models.EditActionUpdate, player)
	}

	return &player, err
}

func (r *mutationResolver) DeletePlayer(ctx context.Context, id string) (*models.Player, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})
	if err != nil {
		return &models.Player{}, err
	}
	panic("not implemented")
}
