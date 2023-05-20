package util

import (
	"encoding/json"
	"fmt"
	"monitoring/service/logging"
	"net/http"
	"reflect"
)

func WriteResult(rw http.ResponseWriter, result interface{}, statusCode int) {
	if result == nil || reflect.ValueOf(result).IsNil() {
		rw.WriteHeader(statusCode)
		return
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		logging.Errorf("error on Marshal http response: %s", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err, ok := result.(error)
	if ok {
		logging.Errorf("error on http request: %s", err.Error())
		_, _ = rw.Write([]byte(fmt.Sprintf("error: %s", err.Error())))
		rw.WriteHeader(statusCode)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	_, _ = rw.Write(jsonBytes)
	rw.WriteHeader(http.StatusOK)
}
