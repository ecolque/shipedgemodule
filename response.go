package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SendData(w http.ResponseWriter, data interface{}, err interface{}) {
	if err != nil {
		SendError(w, err, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	output, _ := json.Marshal(&data)
	fmt.Fprintln(w, string(output))
}

func SendError(w http.ResponseWriter, msg interface{}, status int) {
	w.WriteHeader(status)
	fmt.Fprintln(w, msg)
}
