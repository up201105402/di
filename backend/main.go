package main

import (
	"di/model"
	"di/service"
	"di/util"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var db = make(map[string]string)

// signupReq is not exported, hence the lowercase name
// it is used for validation and json marshalling
type signupReq struct {
	Username string `json:"username" binding:"required,gte=5,lte=30"`
	Password string `json:"password" binding:"required,gte=5,lte=30"`
}

// used to help extract validation errors
type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

type Services struct {
	userService     model.UserService
	pipelineService model.PipelineService
	//tokenService model.TokenService
}

func main() {
	dbConnection := util.ConnectToDB()
	error := util.CreateOrUpdateSchema(dbConnection)
	services := &Services{
		userService:     service.NewUserService(dbConnection),
		pipelineService: service.NewPipelineService(dbConnection),
		//tokenService: service.NewTokenService(),
	}

	r := setupRouter(services)

	if error != nil {
		return
	}

	// Listen and Server in 0.0.0.0:8001
	r.Run(":8001")
}

func setupRouter(services *Services) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	router := gin.Default()

	router.Static("/public", "../client/public")
	router.Static("/data-sources", "../client/public/data-sources")
	router.Static("/assets", "../client/src/assets")

	router.LoadHTMLFiles("../client/index.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// Ping test
	router.GET("/ping", func(context *gin.Context) {
		context.String(http.StatusOK, "pong")
	})

	private := router.Group("/api/user")

	private.POST("/login", func(context *gin.Context) {
		LogIn(context, services)
	})

	private.POST("/signup", func(context *gin.Context) {
		SignUp(context, services)
	})

	return router
}

func LogIn(context *gin.Context, services *Services) {
	// define a variable to which we'll bind incoming
	// json body, {email, password}
	var req signupReq

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(context, &req); !ok {
		return
	}

	user := &model.User{
		Username: req.Username,
		Password: req.Password,
	}

	// ctx := context.Request.Context()
	err := services.userService.Signin(user)

	if err != nil {
		log.Printf("Failed to log in user: %v\n", err.Error())
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// create token pair as strings
	//tokens, err := h.TokenService.NewPairFromUser(ctx, u, "")

	// if err != nil {
	// 	log.Printf("Failed to create tokens for user: %v\n", err.Error())

	// 	context.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err,
	// 	})
	// 	return
	// }

	context.JSON(http.StatusOK, gin.H{
		// "tokens": tokens,
	})
}

func SignUp(context *gin.Context, services *Services) {
	// define a variable to which we'll bind incoming
	// json body, {email, password}
	var req signupReq

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(context, &req); !ok {
		return
	}

	// ctx := context.Request.Context()
	err := services.userService.Signup(req.Username, req.Password)

	if err != nil {
		log.Printf("Failed to sign up user: %v\n", err.Error())
		context.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// create token pair as strings
	//tokens, err := h.TokenService.NewPairFromUser(ctx, u, "")

	if err != nil {
		log.Printf("Failed to create tokens for user: %v\n", err.Error())

		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		// "tokens": tokens,
	})
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
			// could probably extract this, it is also in middleware_auth_user
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
