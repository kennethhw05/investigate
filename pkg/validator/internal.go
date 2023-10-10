package validator

import (
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
)

func InternalValidateOverUnderDefaults(db repository.DataSource, game models.Game, matches *[]models.Match) error {
	formats := make(map[models.MatchFormat]bool)
	for _, match := range *matches {
		formats[match.Format] = true
	}
	keys := make([]string, 0, len(*matches))
	for k := range formats {
		keys = append(keys, k.String())
	}

	count, err := db.Model((*models.OverUnderDefault)(nil)).
		Where("game = ? AND match_format in (?)", game.String(), pg.In(keys)).
		Count()

	if err != nil {
		return err
	}

	if count != len(keys) {
		return errors.Errorf("Unknown match format is present in the list of matches")
	}
	return nil
}
