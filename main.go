// Copyright (c) 2018 Popcorn
// Author(s) Calvin Feng

package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"popcorn/model"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	} else {
		port = fmt.Sprintf(":%s", port)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	db, err := SetupDatabase()
	if err != nil {
		logrus.Error("Failed to set up database", err)
		return
	}

	defer db.Close()

	// Client connection map is meant for keeping track of all web socket connection to every client. It is also being
	// used in the recommend engine for notifying clients that their preference vector is ready.
	clientConnMap := make(map[uint]*websocket.Conn)

	// This is the channel for communication between http handlers and a background running engine asynchronously.
	updateUserPreferenceQueue := make(chan *model.User, 100)

	// Set up online learning engine for serving the incoming requests.
	engine := NewOnlineLearningEngine(db, clientConnMap)
	go engine.ListenToInbound(updateUserPreferenceQueue)

	server := &http.Server{
		Handler:      LoadRoutes(db, updateUserPreferenceQueue),
		Addr:         port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.Infof("HTTP server is listening and serving on port %v", port)
	logrus.Fatal(server.ListenAndServe())
}
