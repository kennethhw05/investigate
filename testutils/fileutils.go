package testutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/pkg/errors"
	"gopkg.in/h2non/gock.v1"
)

func ReadFileToStruct(file string, v interface{}) error {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	jsonFile, err := os.Open(fmt.Sprintf("%s/http_mocks/%s", basepath, file))
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error opening file %s", file))
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error reading from file %s", file))
	}

	err = json.Unmarshal(byteValue, v)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error reading file %s into type %s", file, reflect.TypeOf(v).String()))
	}
	return nil
}

// CreateMockRequest Attaches a mock request onto gock for the duration of the test
// IMPORTANT: Make sure you have defer gock.Off() if using gock or it will hang tests
// that don't use it.
func CreateMockRequest(host string, path string, respCode int, fixtureFile string, method string) error {
	if method == "GET" {

		mockResp := make(map[string]interface{})
		err := ReadFileToStruct(fixtureFile, &mockResp)
		if err == nil {
			gock.New(host).
				Get(path).
				Persist().
				Reply(respCode).
				JSON(mockResp)
			return nil
		}
		arrayMockResp := make([]map[string]interface{}, 0)
		err = ReadFileToStruct(fixtureFile, &arrayMockResp)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("error reading file %s to struct", fixtureFile))
		}
		gock.New(host).
			Get(path).
			Persist().
			Reply(respCode).
			JSON(arrayMockResp)
		return nil
	}
	if method == "POST" {
		//gock.Observe(gock.DumpRequest)
		arrayMockResp := make([]map[string]interface{}, 0)
		err := ReadFileToStruct(fixtureFile, &arrayMockResp)
		if err == nil {
			gock.New(host).
				Post(path).
				Persist().
				Reply(respCode).
				JSON(arrayMockResp)
			return nil
		}
	}

	return nil
}
