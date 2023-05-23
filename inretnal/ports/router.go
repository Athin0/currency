package ports

import (
	"currency/inretnal/service"
	"github.com/gin-gonic/gin"
)

func AppRouter(g *gin.RouterGroup, a service.App) {
	g.GET("/avg", GetAvg(a))        // Метод для получения среднего значение курса рубля за весь период по всем валютам
	g.GET("/max", GetMax(a))        // Метод для получения максимального курса валюты, название этой валюты и дату этого максимального значения
	g.GET("/min", GetMin(a))        // Метод для получения минимального курса валюты, название этой валюты и дату этого минимального значения
	g.GET("/update", UpdateData(a)) // Метод обновления данных последних 90 дней
}
