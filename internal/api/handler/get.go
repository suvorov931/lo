package handler

import (
	"fmt"
	"net/http"
)

func GetTask() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GetTask")
	}
}
