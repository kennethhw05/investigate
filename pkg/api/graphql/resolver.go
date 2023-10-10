//go:generate gorunpkg github.com/99designs/gqlgen

package graphql

import (
	"fmt"
	"regexp"

	sq "github.com/Masterminds/squirrel"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/config"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
)

type Resolver struct {
	DB     repository.DataSource
	CFG    *config.Config
	Logger *logrus.Logger
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

func (r *Resolver) Event() EventResolver {
	return &eventResolver{r}
}
func (r *Resolver) Leg() LegResolver {
	return &legResolver{r}
}
func (r *Resolver) Match() MatchResolver {
	return &matchResolver{r}
}
func (r *Resolver) Player() PlayerResolver {
	return &playerResolver{r}
}
func (r *Resolver) Pool() PoolResolver {
	return &poolResolver{r}
}
func (r *Resolver) ConsolationPrize() ConsolationPrizeResolver {
	return &consolationPrizeResolver{r}
}
func (r *Resolver) PoolDefault() PoolDefaultResolver {
	return &poolDefaultResolver{r}
}
func (r *Resolver) Competitor() CompetitorResolver {
	return &competitorResolver{r}
}
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}
func (r *Resolver) Audit() AuditResolver {
	return &auditResolver{r}
}
func (r *Resolver) OverUnderDefault() OverUnderDefaultResolver {
	return &overUnderDefaultResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func addStandardPagination(query sq.SelectBuilder, page *int, perPage *int) sq.SelectBuilder {
	if page == nil {
		page = new(int)
		*page = 0
	}
	if perPage == nil {
		perPage = new(int)
		*perPage = 10
	}

	query = query.Offset(uint64(*page * *perPage))
	query = query.Limit(uint64(*perPage))

	return query
}

func addDefaultSort(query sq.SelectBuilder, sortField *string, sortOrder *string) sq.SelectBuilder {
	if sortField == nil {
		return query
	}

	//injection prevention per: https://stackoverflow.com/questions/30867337/golang-order-by-issue-with-mysql
	var validColumn = regexp.MustCompile("^[A-Za-z0-9_]+$")
	if validColumn.MatchString(*sortField) {
		if sortOrder == nil {
			sortOrder = new(string)
			*sortOrder = "ASC"
		} else if *sortOrder != "ASC" && *sortOrder != "DESC" {
			*sortOrder = "ASC"
		}

		query = query.OrderBy(fmt.Sprintf("%s %s", *sortField, *sortOrder))
	}
	return query
}

func filterByIDs(query sq.SelectBuilder, inputIds []string, tableName string) sq.SelectBuilder {
	if inputIds != nil {
		field := fmt.Sprintf("%s.id", tableName)
		query = query.Where(sq.Eq{field: inputIds})
	}

	return query
}

func parseNullDecimal(input string) decimal.NullDecimal {
	dec, err := decimal.NewFromString(input)
	if err == nil {
		return decimal.NullDecimal{
			Decimal: dec,
			Valid:   true,
		}
	}

	return decimal.NullDecimal{}
}
