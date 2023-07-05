package main

import (
	"fmt"

	api "github.com/Axit88/BiteSpeedTask/src/api"
	"github.com/Axit88/BiteSpeedTask/src/constants"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/identify", api.PostRequest)

	url := fmt.Sprintf(":%v", constants.REST_PORT)
	router.Run(url)
}
