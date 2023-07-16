package middleware

import (
	"di/service"
	"di/util/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
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

func Auth(tokenService service.TokenService, I18n *i18n.Localizer) gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := authHeader{}

		if err := context.ShouldBindHeader(&authHeader); err != nil {
			if errs, ok := err.(validator.ValidationErrors); ok {
				var invalidArgs []invalidArgument

				for _, err := range errs {
					invalidArgs = append(invalidArgs, invalidArgument{
						err.Field(),
						err.Value().(string),
						err.Tag(),
						err.Param(),
					})
				}

				errorMessage, _ := I18n.Localize(&i18n.LocalizeConfig{
					MessageID: "sys.binding.req",
					TemplateData: map[string]interface{}{
						"Reason": err.Error(),
					},
				})

				err := errors.NewBadRequest(errorMessage)

				context.JSON(err.Status(), gin.H{
					"error":       err,
					"invalidArgs": invalidArgs,
				})
				context.Abort()
				return
			}

			err := errors.NewInternal("")
			context.JSON(err.Status(), gin.H{
				"error": err.Message,
			})
			context.Abort()
			return
		}

		if authHeader.IDToken == "" {
			errorMessage, _ := I18n.Localize(&i18n.LocalizeConfig{
				MessageID: "auth.header.authorization.empty",
			})
			err := errors.NewAuthorization(errorMessage)

			context.JSON(err.Status(), gin.H{
				"error": err,
			})
			context.Abort()
			return
		}

		user, err := tokenService.ValidateIDToken(authHeader.IDToken)

		if err != nil {
			err := errors.NewAuthorization(err.Error())
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
