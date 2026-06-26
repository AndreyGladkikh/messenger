package http

import "net/http"

var apiErrorsMap = map[string]int{
	"ChatNotFound": http.StatusNotFound,
	"MessageNotFound": http.StatusNotFound,
	"UserNotFound": http.StatusNotFound,
	"LoginAlreadyExists": http.StatusBadRequest,
}
