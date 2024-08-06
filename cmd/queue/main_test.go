package main

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	var (
		maxCountQueueMsgValue         = 0
		maxCountQueueValue            = 0
		opGetTimeoutValue             = 0
		opGetDisableTimeoutValue      = false
		portValue                uint = 80
		grpcPortValue            uint = 10000
	)

	maxCountQueueMsg = &maxCountQueueMsgValue
	maxCountQueue = &maxCountQueueValue
	opGetTimeout = &opGetTimeoutValue
	opGetDisableTimeout = &opGetDisableTimeoutValue
	port = &portValue
	grpcPort = &grpcPortValue

	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Minute)
	defer cancel()

	go run(ctx)

	time.Sleep(time.Second)
	os.Exit(m.Run())
}
