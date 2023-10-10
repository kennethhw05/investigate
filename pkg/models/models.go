package models

import (
	"fmt"
	"reflect"
	"strings"
)

// GenerateInternalXID Create internal XID for model and extra metadata
func GenerateInternalXID(model interface{}, metadata string) string {
	modelType := ""
	if t := reflect.TypeOf(model); t.Kind() == reflect.Ptr {
		modelType = t.Elem().Name()
	} else {
		modelType = t.Name()
	}
	modelType = strings.ToLower(modelType)

	return fmt.Sprintf("internal:%s:%s", modelType, metadata)
}
