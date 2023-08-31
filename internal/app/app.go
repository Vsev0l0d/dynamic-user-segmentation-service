package app

import (
	"context"
	"dynamic-user-segmentation-service/internal/config"
	"dynamic-user-segmentation-service/internal/db"
	"dynamic-user-segmentation-service/internal/handler"
	"dynamic-user-segmentation-service/internal/periodic"
	"dynamic-user-segmentation-service/internal/s3"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(configPath string) {
	conf, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	listener, err := net.Listen("tcp", conf.HTTP.Addr)
	if err != nil {
		log.Fatalf("Error occurred: %s", err.Error())
	}
	database, err := db.Connect(conf.DB)
	if err != nil {
		log.Fatalf("Could not set up database: %v", err)
	}
	defer database.DB.Close()

	err = periodic.GoPeriodicDeletionOfInactiveSegments(conf.PeriodForDeletingInactiveSegments.CronExpression, database)
	if err != nil {
		log.Fatalf("Error starting periodic deletion of inactive segments task: %v", err)
	}

	s3, err := s3.Connect(conf.Minio)
	if err != nil {
		log.Fatalf("Could not set up s3: %v", err)
	}

	httpHandler := handler.NewHandler(database, s3)
	server := &http.Server{
		Handler: httpHandler,
	}
	go func() {
		server.Serve(listener)
	}()
	defer Stop(server)
	log.Printf("Started server on %s", conf.HTTP.Addr)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(fmt.Sprint(<-ch))
	log.Println("Stopping API server.")
}

func Stop(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Could not shut down server correctly: %v\n", err)
		os.Exit(1)
	}
}
