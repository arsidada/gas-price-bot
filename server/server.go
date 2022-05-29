package server

import (
	"net/http"
	"strings"

	"github.com/arsidada/gas-price-bot/fetcher"
	"github.com/gin-gonic/gin"
)

func StartServer() error {
	r := gin.Default()
	r.GET("/price/:location", priceHandler)

	err := r.Run()
	if err != nil {
		return err
	}

	return nil
}

func priceHandler(c *gin.Context) {
	location := c.Param("location")
	if location == "" {
		c.String(http.StatusBadRequest, "location parameter not set")
	}
	price, err := fetcher.FetchPrice(strings.ToUpper(location))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.JSON(200, gin.H{
		"price": price,
	})
}
