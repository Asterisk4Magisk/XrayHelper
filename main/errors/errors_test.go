package errors_test

import (
	"XrayHelper/main/log"
	"XrayHelper/main/utils"
	"os"
	"testing"
	"time"
)

func TestError(t *testing.T) {
	external := utils.NewExternal(2*time.Second, os.Stdout, os.Stderr, "ping", "-n", "6", "127.0.0.1")
	external.Start()
	err := external.Wait()
	if err != nil {
		log.HandleError(err)
	}
	if external.Err() != nil {
		log.HandleError(external.Err())
	}
}
