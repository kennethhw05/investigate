package validator

import (
	"time"

	"github.com/pkg/errors"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
)

func ColossusValidateNumberOfLegsOrMatches(count int) error {
	switch count {
	case 4, 6, 8, 10:
		return nil
	default:
		return errors.New("only pools with 4/6/8/10 legs can be sent to colossus")
	}
}

func ColossusValidate48HoursAdvancedCreation(matches []models.Match) error {
	if len(matches) == 0 {
		return errors.New("need at least one match to validate we're in the 48 hour window")
	}

	if matches[0].StartTime.After(time.Now().Add(48 * time.Hour)) {
		return errors.New("earliest scheduled match must be within the next 48 hours")
	}

	return nil
}
