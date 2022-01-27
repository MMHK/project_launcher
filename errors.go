package main

import "errors"

type LauncherError struct {
	Code int 	`json:"code"`
	TagMessage string `json:"message"`
}

const (
	ERROR_DOCKER_DESKTOP_NOT_INSTALLED = 1000
	ERROR_DOCKER_DESKTOP_NOT_RUNNING = 1001
	ERROR_WINDOW_TERMINAL_NOT_INSTALLED = 1002
)

func (e *LauncherError) Error() string {
	return e.TagMessage
}

func NewLauncherError(code int, message string) *LauncherError {
	return &LauncherError{
		Code:    code,
		TagMessage: message,
	}
}

func IsLauncherError(err error) bool {
	target := new(LauncherError)
	return errors.Is(err, target)
}

func MatchLauncherError(err error, errorCode int) bool {
	target := new(LauncherError)
	return  errors.As(err, &target) && target.Code == errorCode
}