package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"
	"weeb_bot/internal/storage/couch"
)

func main() {
	var err error
	couchdb := couch.New()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = couchdb.Init(ctx)
	if err != nil {
		log.Fatalf("Init failed: %v", err)
	}
	m, err := couchdb.GetAverages(ctx, "jealcFNOjVdIZAiDIPcUlEGXeEa4tyK4voa59xJkRZF7GtpgTMC2stABJjFxRHB4dwPyiUcJEqfLYA")
	if err != nil {
		log.Fatalf("GetAverages failed: %v", err)
	}
	log.Debugf("%v", m)
}
