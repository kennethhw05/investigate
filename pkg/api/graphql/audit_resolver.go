package graphql

import (
	context "context"

	sq "github.com/Masterminds/squirrel"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/auth"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
)

type auditResolver struct{ *Resolver }

func (r *auditResolver) ID(ctx context.Context, obj *models.Audit) (string, error) {
	return obj.GetID(), nil
}

func (r *auditResolver) User(ctx context.Context, obj *models.Audit) (*models.User, error) {
	user := models.User{ID: obj.UserID}
	err := r.DB.
		Select(&user)

	if err != nil {
		r.Logger.Info("Audit Resolver Get User Error")
		r.Logger.Info(err.Error())
	}

	return &user, nil
}

func (r *auditResolver) TargetID(ctx context.Context, obj *models.Audit) (string, error) {
	return obj.GetTargetID(), nil
}

func (r *queryResolver) Audit(ctx context.Context, id string) (*models.Audit, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})

	if err != nil {
		return &models.Audit{}, err
	}

	audit := models.Audit{ID: repository.NewSQLCompatUUIDFromStr(id)}
	err = r.DB.
		Select(&audit)

	return &audit, err
}

func (r *queryResolver) AllAudits(ctx context.Context, filter *models.AuditFilter, page *int, perPage *int, sortField *string, sortOrder *string) ([]*models.Audit, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})

	if err != nil {
		return []*models.Audit{}, err
	}

	var audits []*models.Audit

	sqlbuilder := sq.Select("audits.*").From("audits")
	sql, args := createAllAuditsSQL(sqlbuilder, filter, page, perPage, sortField, sortOrder)
	_, err = r.DB.Query(&audits, sql, args...)

	return audits, err
}

func (r *queryResolver) _allAuditsMeta(ctx context.Context, filter *models.AuditFilter, page *int, perPage *int) (*models.ListMetadata, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})

	if err != nil {
		return &models.ListMetadata{Count: 0}, err
	}

	var count int

	sqlbuilder := sq.Select("COUNT(audits.*)").From("audits")
	sql, args := createAllAuditsSQL(sqlbuilder, filter, nil, nil, nil, nil)
	_, err = r.DB.Query(&count, sql, args...)
	return &models.ListMetadata{Count: count}, err
}

func createAllAuditsSQL(builder sq.SelectBuilder, filter *models.AuditFilter, page *int, perPage *int, sortField *string, sortOrder *string) (sql string, args []interface{}) {
	if filter != nil {

		if filter.TargetID != nil {
			builder = builder.Where("audits.target_id = ?", *filter.TargetID)
		}

		if filter.TargetType != nil {
			builder = builder.Where("audits.target_type = ?", *filter.TargetType)
		}

		if filter.UserID != nil {
			builder = builder.Where("audits.user_id = ?", *filter.UserID)
		}

		if filter.EditAction != nil {
			builder = builder.Where("audits.edit_action = ?", *filter.EditAction)
		}

		builder = filterByIDs(builder, filter.Ids, "audits")
		builder = addDefaultSort(builder, sortField, sortOrder)
		builder = addStandardPagination(builder, page, perPage)
	}

	sql, args, _ = builder.ToSql()
	return sql, args
}
