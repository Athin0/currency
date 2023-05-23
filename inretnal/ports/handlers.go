package ports

import (
	"currency/inretnal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAvg(a service.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		ans, err := a.GetAvg(c)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ans)
	}
}
func GetMax(a service.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		ans, err := a.GetMax(c)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ans)
	}
}
func GetMin(a service.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		ans, err := a.GetMin(c)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, ans)
	}
}

func UpdateData(a service.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := a.LastDaysAddToDB(c, 90)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.Status(http.StatusOK)
	}
}
