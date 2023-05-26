package errors_test

import (
	"XrayHelper/main/common"
	"XrayHelper/main/log"
	"os"
	"testing"
	"time"
)

func TestError(t *testing.T) {
	external := common.NewExternal(2*time.Second, os.Stdout, os.Stderr, "ping", "-n", "6", "127.0.0.1")
	external.Start()
	if err := external.Wait(); err != nil {
		log.HandleError(err)
	}
	if external.Err() != nil {
		log.HandleError(external.Err())
	}
}
