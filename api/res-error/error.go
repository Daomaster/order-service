package res_error

import "errors"

type ResponseError struct {
	Error string `json:"res-error"`
}

var ErrOrderAlreadyTaken = errors.New("the order is already taken")