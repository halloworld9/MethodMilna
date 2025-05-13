package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func Logger() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		start := time.Now()
		u, err := uuid.NewUUID()
		if err != nil {
			log.Err(err).Send()
			return fiber.ErrInternalServerError
		}
		us := u.String()
		log.Info().Time("start", start).
			Str("Path", ctx.Path()).
			Str("uuid", us).
			Str("IP", ctx.IP()).
			Str("METHOD", ctx.Method()).
			Str("User-Agent", ctx.Get("User-Agent")).
			Send()
		defer log.Info().
			Str("uuid", u.String()).
			Int("status", ctx.Response().StatusCode()).
			TimeDiff("completion time", time.Now(), start).
			Send()
		return ctx.Next()
	}
}
