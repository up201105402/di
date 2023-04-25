package middleware

import (
	"di/service"
	"di/util/errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type authHeader struct {
	IDToken string `header:"Authorization"`
}

type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

// Auth extracts a user from the Authorization header
// which is of the form "Bearer token"
// It sets the user to the context if the user exists
func Auth(tokenService service.TokenService) gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := authHeader{}

		// bind Authorization Header to h and check for validation errors
		if err := context.ShouldBindHeader(&authHeader); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				// we used this type in bind_data to extract desired fields from errs
				// you might consider extracting it
				var invalidArgs []invalidArgument

				for _, err := range errs {
					invalidArgs = append(invalidArgs, invalidArgument{
						err.Field(),
						err.Value().(string),
						err.Tag(),
						err.Param(),
					})
				}

				err := errors.NewBadRequest("Invalid request parameters. See invalidArgs")

				context.JSON(err.Status(), gin.H{
					"error":       err,
					"invalidArgs": invalidArgs,
				})
				context.Abort()
				return
			}

			// otherwise error type is unknown
			err := errors.NewInternal()
			context.JSON(err.Status(), gin.H{
				"error": err,
			})
			context.Abort()
			return
		}

		idTokenHeader := strings.Split(authHeader.IDToken, "Bearer ")

		if len(idTokenHeader) < 2 {
			err := errors.NewAuthorization("Must provide Authorization header with format `Bearer {token}`")

			context.JSON(err.Status(), gin.H{
				"error": err,
			})
			context.Abort()
			return
		}

		// validate ID token here
		user, err := tokenService.ValidateIDToken(idTokenHeader[1])

		if err != nil {
			err := errors.NewAuthorization("Provided token is invalid")
			context.JSON(err.Status(), gin.H{
				"error": err,
			})
			context.Abort()
			return
		}

		context.Set("user", user)

		context.Next()
	}
}
