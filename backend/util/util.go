package util

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

func BindData(c *gin.Context, req interface{}) bool {
	if c.ContentType() != "application/json" {
		msg := fmt.Sprintf("%s only accepts Content-Type application/json", c.FullPath())

		c.JSON(http.StatusUnsupportedMediaType, gin.H{
			"error": msg,
		})
		return false
	}

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

func Filter[T any](data []T, f func(T) bool) []T {

	fltd := make([]T, 0, len(data))

	for _, e := range data {
		if f(e) {
			fltd = append(fltd, e)
		}
	}

	return fltd
}

func Map[T, U any](data []T, f func(T) U) []U {

	res := make([]U, 0, len(data))

	for _, e := range data {
		res = append(res, f(e))
	}

	return res
}

func StringArrayContains(array []string, str string) bool {
	for _, a := range array {
		if a == str {
			return true
		}
	}
	return false
}

func ReadCsvFile(filePath string) ([][]string, bool) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, false
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Println(err.Error())
	}

	return records, true
}
