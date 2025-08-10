package controllers

import (
	"go-api-find-my-friend/pkg/errors"
	"net/http"
)

func getErrStatusCode(error error) int {
	statusCode, success := error.(*errors.AppError)
	if success {
		return statusCode.Code
	}
	return http.StatusInternalServerError
}
