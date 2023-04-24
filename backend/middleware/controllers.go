package middleware

import (
	"github.com/gin-contrib/sessions"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginGetHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		session := sessions.Default(context)
		user := session.Get("user")
		if user != nil {
			context.HTML(http.StatusBadRequest, "login.html",
				gin.H{
					"content": "Please logout first",
					"user":    user,
				})
			return
		}
		context.HTML(http.StatusOK, "login.html", gin.H{
			"content": "",
			"user":    user,
		})
	}
}

func LoginPostHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		session := sessions.Default(context)
		user := session.Get("user")
		if user != nil {
			context.HTML(http.StatusBadRequest, "login.html", gin.H{"content": "Please logout first"})
			return
		}

		username := context.PostForm("username")
		// password := context.PostForm("password")

		// if helpers.EmptyUserPass(username, password) {
		// 	context.HTML(http.StatusBadRequest, "login.html", gin.H{"content": "Parameters can't be empty"})
		// 	return
		// }

		// if !helpers.CheckUserPass(username, password) {
		// 	context.HTML(http.StatusUnauthorized, "login.html", gin.H{"content": "Incorrect username or password"})
		// 	return
		// }

		session.Set("user", username)
		if err := session.Save(); err != nil {
			context.HTML(http.StatusInternalServerError, "login.html", gin.H{"content": "Failed to save session"})
			return
		}

		context.Redirect(http.StatusMovedPermanently, "/dashboard")
	}
}

func LogoutGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		log.Println("logging out user:", user)
		if user == nil {
			log.Println("Invalid session token")
			return
		}
		session.Delete("user")
		if err := session.Save(); err != nil {
			log.Println("Failed to save session:", err)
			return
		}

		c.Redirect(http.StatusMovedPermanently, "/")
	}
}

func IndexGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": "This is an index page...",
			"user":    user,
		})
	}
}

func DashboardGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"content": "This is a dashboard",
			"user":    user,
		})
	}
}
