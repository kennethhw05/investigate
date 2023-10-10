package graphql

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	sq "github.com/Masterminds/squirrel"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/audit"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/auth"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type userResolver struct{ *Resolver }

func (r *userResolver) ID(ctx context.Context, obj *models.User) (string, error) {
	return obj.GetID(), nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input models.CreateUserInput) (*models.User, error) {
	userID, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
	})

	if err != nil {
		return &models.User{}, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return &models.User{}, errors.New("error signing up with provided password")
	}

	_, err = mail.ParseAddress(input.Email)
	if err != nil {
		return &models.User{}, errors.New("error parsing provided email")
	}

	user := models.User{
		Email:      input.Email,
		Password:   passwordHash,
		AccessRole: input.Role,
	}
	err = r.DB.Insert(&user)

	if err == nil {
		audit.CreateAudit(r.DB, user.ID, repository.NewSQLCompatUUIDFromStr(userID), models.EditActionCreate, user)
	}

	return &user, err
}

func (r *mutationResolver) CreateSession(ctx context.Context, input models.AuthInput) (*models.Session, error) {
	user := models.User{}
	err := r.DB.Model(&user).Where("email = ?", input.Email).Select()
	if err != nil {
		r.Logger.Infof("Error finding user %s, err: %s", input.Email, err.Error())
		return &models.Session{}, errors.New("invalid email or password")
	}

	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(input.Password)); err != nil {
		r.Logger.Infof("Invalid password for user %s, err: %s", input.Email, err.Error())
		return &models.Session{}, errors.Errorf("invalid email or password, err: %s", err.Error())
	}

	token, err := r.generateJWT(&user)
	if err != nil {
		r.Logger.Errorf("Could not generate a valid JWT for user %s, err: %s", input.Email, err.Error())
		return &models.Session{}, errors.New("invalid email or password")
	}
	return &models.Session{Token: token}, nil
}

func (r *mutationResolver) generateJWT(user *models.User) (string, error) {
	expiresAt := time.Now().Add(time.Hour).Unix()
	tokenString, err := r.generateJWTWithExpiry(user, expiresAt)
	return tokenString, err
}

func (r *mutationResolver) generateJWTWithExpiry(user *models.User, expiresAt int64) (string, error) {
	claims := &auth.JWTCustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
		ID:         user.GetID(),
		Email:      user.Email,
		AccessRole: user.AccessRole,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(r.CFG.JWTSignature))
	return tokenString, err
}

func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
		models.AccessRoleUser,
		models.AccessRoleGuestAPIOnly,
	})

	if err != nil {
		return &models.User{}, err
	}

	user := models.User{ID: repository.NewSQLCompatUUIDFromStr(id)}
	err = r.DB.
		Select(&user)
	return &user, err
}

func (r *queryResolver) AllUsers(ctx context.Context, filter *models.UserFilter, page *int, perPage *int, sortField *string, sortOrder *string) ([]*models.User, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})

	if err != nil {
		return []*models.User{}, err
	}

	var users []*models.User

	sqlbuilder := sq.Select("users.*").From("users")
	sql, args := createAllUsersSQL(sqlbuilder, filter, page, perPage, sortField, sortOrder)
	_, err = r.DB.Query(&users, sql, args...)

	return users, err
}

func (r *queryResolver) _allUsersMeta(ctx context.Context, filter *models.UserFilter, page *int, perPage *int) (*models.ListMetadata, error) {
	_, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
		models.AccessRoleAdmin,
	})

	if err != nil {
		return &models.ListMetadata{Count: 0}, err
	}

	var count int

	sqlbuilder := sq.Select("COUNT(*)").From("users")
	sql, args := createAllUsersSQL(sqlbuilder, filter, nil, nil, nil, nil)
	_, err = r.DB.Query(&count, sql, args...)
	return &models.ListMetadata{Count: count}, err
}

func createAllUsersSQL(builder sq.SelectBuilder, filter *models.UserFilter, page *int, perPage *int, sortField *string, sortOrder *string) (sql string, args []interface{}) {
	if filter != nil {
		if filter.Email != nil {
			builder = builder.Where("users.email ILIKE ?", fmt.Sprint("%", *filter.Email, "%"))
		}

		if filter.Role != nil {
			builder = builder.Where("users.access_role = ?", filter.Role.String())
		}

		builder = filterByIDs(builder, filter.Ids, "users")
	}
	builder = addDefaultSort(builder, sortField, sortOrder)
	builder = addStandardPagination(builder, page, perPage)

	sql, args, _ = builder.ToSql()
	return sql, args
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input models.UpdateUserInput) (*models.User, error) {
	userID, err := auth.Authorize(ctx, []models.AccessRole{
		models.AccessRoleSuperAdmin,
	})
	if err != nil {
		return &models.User{}, err
	}
	var user models.User

	sqlbuilder := sq.Update("users")
	field := fmt.Sprintf("%s.id", "users")
	sqlbuilder = sqlbuilder.Where(sq.Eq{field: input.ID})

	if input.Email != nil {
		_, err = mail.ParseAddress(*input.Email)
		if err != nil {
			return &models.User{}, errors.New("error parsing provided email to update")
		}
		sqlbuilder = sqlbuilder.Set("email", *input.Email)
	}

	if input.Password != nil {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return &models.User{}, errors.New("error updating password up with provided password")
		}
		sqlbuilder = sqlbuilder.Set("password", passwordHash)
	}

	if input.AccessRole != nil {
		sqlbuilder = sqlbuilder.Set("access_role", *input.AccessRole)
	}

	sqlbuilder = sqlbuilder.Suffix("RETURNING *")
	sql, args, _ := sqlbuilder.ToSql()
	_, err = r.DB.Query(&user, sql, args...)

	if err == nil {
		audit.CreateAudit(r.DB, user.ID, repository.NewSQLCompatUUIDFromStr(userID), models.EditActionUpdate, user)
	}
	return &user, err
}
