package handlers

import (
	"di/model"
	"di/service"
	"di/util/errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
