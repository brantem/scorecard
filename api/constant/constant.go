package constant

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

const AppID = "scorecard"

var (
	ErrInternalServerError = errors.New("INTERNAL_SERVER_ERROR")
)

var (
	RespInternalServerError = fiber.Map{"code": "INTERNAL_SERVER_ERROR"}
)
