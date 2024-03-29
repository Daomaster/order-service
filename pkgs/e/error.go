package e

import (
	"errors"
)

type ResponseError struct {
	Error string `json:"error"`
}

var (
	// Error when google map can't calculate the distance
	ErrDistanceUnknown = errors.New("the distance between origin and destination is unknown")
	// Error for an order already taken
	ErrOrderAlreadyTaken = errors.New("the order is already taken")
	// Error for query string invalid
	ErrQueryStringInvalid = errors.New("the query strings provided are invalid")
	// Error for order quest invalid
	ErrOrderRequestInvalid = errors.New("the order request is invalid")
	// Error for trying to take order which does not exist
	ErrOrderNotExist = errors.New("the order requested does not exist")
	// Error for all internal error should not be exposed
	// Note: it should add error id and check it in the logs
	ErrInternalError = errors.New("the request failed by internal error")
)

// function to create a response error
func CreateErr(err error) *ResponseError {
	res := &ResponseError{Error: err.Error()}
	return res
}
