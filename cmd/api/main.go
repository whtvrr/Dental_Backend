package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/whtvrr/Dental_Backend/pkg/logger"
)

func main() {
	l := logger.GerLogger()
	// connectionLink := os.Getenv("SHODY_MONGO_DB")
	// database, err := mongodb.NewClient(context.Background(), "Dental", connectionLink)
	// if err != nil {
	// 	l.Fatalf("not connected to database, err: %s", err)
	// }
	r := gin.Default()
	l.Info("router created")

	r.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"http://localhost:3000"},
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// userStorage := UserDB.NewUserStorage(database, "Users", &l)
	// userHandler := user.NewHangler(userStorage, &l)
	// userHandler.Register(r)
	l.Logger.Info("User Handler added")

	r.Run("localhost:8080")
}
