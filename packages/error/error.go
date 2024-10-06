package cerr

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type ServerError struct {
	HttpStats int
	Msg string
}

func NewServerError(httpStatus int, msg string) *ServerError {
	return &ServerError{
		HttpStats: httpStatus,
		Msg: msg,
	}
}

func (se *ServerError) Error() string {
	return fmt.Sprintln(se.Msg)
}

func (se *ServerError) JSON(c echo.Context) error {
	return c.JSON(se.HttpStats, map[string]string{"message":se.Msg})
}
