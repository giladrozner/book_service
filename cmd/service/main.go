package main

import (
	"github.com/gin-gonic/gin"

	"github.com/giladrozner/book_service/pkg/config"
	"github.com/giladrozner/book_service/pkg/consts"
	"github.com/giladrozner/book_service/pkg/service"
)

func main() {
	if err := config.Setup(); err != nil {
		panic(err)
	}

	engine := gin.Default()
	service.Routes(engine)

	if err := engine.Run(consts.ServerPort); err != nil {
		panic(err)
	}
}
