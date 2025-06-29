package main

import (
	"backend3/models"
	"backend3/routers"
	"backend3/utils"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	db, _ := utils.DBConnect()
	godotenv.Load()
	go func() {
		for {
			models.CleanBlacklistTokens()
			time.Sleep(15 * time.Minute)
		}
	}()
	defer db.Close()
	r := gin.Default()
	routers.CombineRouter(r)
	r.Run(fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
}
