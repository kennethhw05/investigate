package audit

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/models"
	"gitlab.com/siimpl/esp-betting/betting-feed/pkg/repository"
)

func CreateAudit(db repository.DataSource, targetId uuid.NullUUID, userID uuid.NullUUID, action models.EditAction, editedObject interface{}) error {
	jsonObj, err := json.Marshal(editedObject)
	if err != nil {
		return err
	}

	targetType := reflect.TypeOf(editedObject).String()
	audit := models.Audit{
		Time:       time.Now(),
		TargetID:   targetId,
		TargetType: strings.Split(targetType, ".")[1],
		UserID:     userID,
		Content:    string(jsonObj),
		EditAction: action,
	}

	err = db.Insert(&audit)
	if err != nil {
		fmt.Printf("Error making audit insert %s", err.Error())
	}
	return err

}
