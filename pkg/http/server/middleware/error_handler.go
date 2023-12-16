package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/PopescuStefanRadu/ent-demo/pkg/ent"
	"github.com/PopescuStefanRadu/ent-demo/pkg/http/server/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type ErrorHandler struct {
	Logger zerolog.Logger
}

func (eh *ErrorHandler) HandleErrors(c *gin.Context) {
	c.Next()
	errs := c.Errors

	if errs == nil {
		return
	}

	var (
		//nolint:godox
		// TODO do not use ent.NotFoundError, instead create a business error that wraps these cases.
		notFound     *ent.NotFoundError
		constraint   *ent.ConstraintError
		responseErr  *response.Error
		validatorErr validator.ValidationErrors
	)

	r := response.Response[*any]{
		Errors: map[string][]response.Error{},
	}

	for _, err := range errs {
		switch {
		case errors.As(err, &responseErr):
			r.Errors[responseErr.Path] = append(r.Errors[responseErr.Path], *responseErr)
		case errors.As(err, &notFound):
			r.Errors["global"] = append(r.Errors["global"], response.Error{
				Code:    "NotFound",
				Message: "resource not found",
			})
		case errors.As(err, &constraint):
			r.Errors["global"] = append(r.Errors["global"], response.Error{
				Code:    "Constraint",
				Message: err.Error(),
			})
		case errors.As(err, &validatorErr):
			for _, fieldError := range validatorErr {
				r.Errors[fieldError.Namespace()] = append(r.Errors[fieldError.Namespace()], response.Error{
					Code:    fieldError.ActualTag(),
					Message: fmt.Sprintf("Validation for %s failed on the '%s' tag", fieldError.Field(), fieldError.ActualTag()),
				})
			}
		default:
			r.Errors["global"] = append(r.Errors["global"], response.Error{
				Code:    "unknown",
				Message: err.Error(),
			})

			eh.Logger.Error().Msgf("Unhandled error of type %T", err)
		}
	}

	c.JSON(http.StatusNotFound, r)
}
