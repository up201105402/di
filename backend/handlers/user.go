package handlers

import (
	"di/model"
	"di/service"
	"di/util"
	"di/util/errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

type tokensReq struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

func EditUser(services *service.Services, I18n *i18n.Localizer) gin.HandlerFunc {
	return func(context *gin.Context) {

		var req model.UserReq

		if ok := bindData(context, &req); !ok {
			return
		}

		user, err := getUser(context)
		if err != nil {
			context.JSON(err.Status(), gin.H{
				"error": err.Error(),
			})
		}

		oldUser, getErr := services.UserService.Get(user.ID)

		if getErr != nil {
			log.Printf(getErr.Error())
			context.JSON(http.StatusNotFound, gin.H{
				"error": getErr.Error(),
			})
			return
		}

		if req.Username != "" {
			oldUser.Username = req.Username
		}

		if req.Password != "" {
			oldPw, err := util.HashPassword(req.Password)

			if oldPw != oldUser.Password {
				errMessage, _ := I18n.Localize(&i18n.LocalizeConfig{
					MessageID: "user.handler.edit.password.old.failed",
				})

				log.Printf(errMessage)
				context.JSON(http.StatusInternalServerError, gin.H{
					"error": errMessage,
				})
				return
			}

			pw, err := util.HashPassword(req.NewPassword)

			if err != nil {
				errMessage, _ := I18n.Localize(&i18n.LocalizeConfig{
					MessageID: "user.repository.create.user.failed",
					TemplateData: map[string]interface{}{
						"Reason": err.Error(),
					},
				})

				log.Printf(errMessage)
				context.JSON(http.StatusInternalServerError, gin.H{
					"error": errMessage,
				})
				return
			}

			oldUser.Password = pw
		}

		updateErr := services.UserService.UpdateDetails(oldUser)

		if updateErr != nil {
			log.Printf(updateErr.Error())
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": updateErr.Error(),
			})
			return
		}

		token, tokenErr := services.TokenService.NewFirstPairFromUser(context.Request.Context(), user)

		if tokenErr != nil {
			log.Printf(tokenErr.Error())
			context.JSON(http.StatusInternalServerError, gin.H{
				"error": tokenErr.Error(),
			})
			return
		}

		context.JSON(http.StatusCreated, gin.H{
			"tokens": token,
		})
	}
}

func LogIn(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		var req model.UserReq

		if ok := bindData(context, &req); !ok {
			return
		}

		user := &model.User{
			Username: req.Username,
			Password: req.Password,
		}

		err := services.UserService.Signin(user)

		if err != nil {
			log.Printf(err.Error())
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		tokens, err := services.TokenService.NewFirstPairFromUser(context.Request.Context(), user)

		if err != nil {
			log.Printf("Failed to create tokens for user: %v\n", err.Error())

			context.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprint("Failed to create tokens for user: %v\n", err.Error()),
			})

			return
		}

		context.JSON(http.StatusOK, gin.H{
			"tokens": tokens,
		})
	}
}

func SignUp(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {

		var req model.UserReq

		if ok := bindData(context, &req); !ok {
			return
		}

		err := services.UserService.Signup(req.Username, req.Password)

		if err != nil {
			log.Printf(err.Error())
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		user, err := services.UserService.GetByUsername(req.Username)

		if err != nil {
			log.Printf(err.Error())
			context.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		token, err := services.TokenService.NewFirstPairFromUser(context.Request.Context(), user)

		if err != nil {
			log.Printf("Failed to create tokens for user: %v\n", err.Error())

			context.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprint("Failed to create tokens for user: %v\n", err.Error()),
			})
			return
		}

		context.JSON(http.StatusCreated, gin.H{
			"tokens": token,
		})
	}
}

func SignOut(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {
		user := context.MustGet("user")

		ctx := context.Request.Context()
		if err := services.TokenService.Signout(ctx, user.(*model.User).ID); err != nil {
			context.JSON(errors.Status(err), gin.H{
				"error": err,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"message": "User signed out successfully!",
		})
	}
}

func NewAccessToken(services *service.Services) gin.HandlerFunc {
	return func(context *gin.Context) {
		// bind JSON to req of type tokensRew
		var req tokensReq

		if ok := bindData(context, &req); !ok {
			return
		}

		ctx := context.Request.Context()

		// verify refresh JWT
		refreshToken, err := services.TokenService.ValidateRefreshToken(req.RefreshToken)

		if err != nil {
			context.JSON(errors.Status(err), gin.H{
				"error": err,
			})
			return
		}

		// get up-to-date user
		user, err := services.UserService.Get(refreshToken.UID)

		if err != nil {
			context.JSON(errors.Status(err), gin.H{
				"error": err,
			})
			return
		}

		tokens, err := services.TokenService.NewPairFromUser(ctx, user, *refreshToken)

		if err != nil {
			log.Printf("Failed to create tokens for user: %+v. Error: %v\n", user, err.Error())

			context.JSON(errors.Status(err), gin.H{
				"error": err,
			})
			return
		}

		context.JSON(http.StatusOK, gin.H{
			"tokens": tokens,
		})
	}
}

// bindData is helper function, returns false if data is not bound
func bindData(c *gin.Context, req interface{}) bool {
	if c.ContentType() != "application/json" {
		msg := fmt.Sprintf("%s only accepts Content-Type application/json", c.FullPath())

		c.JSON(http.StatusUnsupportedMediaType, gin.H{
			"error": msg,
		})
		return false
	}
	// Bind incoming json to struct and check for validation errors
	if err := c.ShouldBind(req); err != nil {
		log.Printf("Error binding data: %+v\n", err)

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

			err := fmt.Sprintf("Bad request. Reason: Invalid request parameters. See invalidArgs")

			c.JSON(http.StatusBadRequest, gin.H{
				"error":       err,
				"invalidArgs": invalidArgs,
			})
			return false
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error.",
		})
		return false
	}

	return true
}

func getUser(context *gin.Context) (*model.User, *errors.Error) {
	contextUser, exists := context.Get("user")

	if !exists {
		err := errors.NewNotFound("User does not exist!")
		return nil, err
	}

	return contextUser.(*model.User), nil
}
