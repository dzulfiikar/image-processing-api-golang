package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dzulfiikar/image-processing-api-golang/cmd/routes"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gocv.io/x/gocv"
)

type ENV struct {
	AppEnv  string `mapstructure:"APP_ENV"`
	AppPort string `mapstructure:"APP_PORT"`
}

func loadEnv() ENV {
	env := ENV{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Error .env file not found : ", err)
		panic(err)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Error loading .env file : ", err)
		panic(err)
	}

	log.Println("Environment loaded successfully")

	log.Printf("App is running in %s mode", env.AppEnv)

	return env
}

func initializeWebServer(port string) {
	router := gin.Default()

	routes.SetupRoutes(&router.RouterGroup)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("Server is running on port", port)

	shutdownWebServer(server)

}

// shutdown gracefully from docs: https://gin-gonic.com/docs/examples/graceful-restart-or-stop/
func shutdownWebServer(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func main() {
	fmt.Printf("gocv version: %s\n", gocv.Version())
	fmt.Printf("opencv lib version: %s\n", gocv.OpenCVVersion())

	env := loadEnv()
	initializeWebServer(env.AppPort)

}
