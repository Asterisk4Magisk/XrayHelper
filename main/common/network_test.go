package common

import (
	"fmt"
	"testing"
	"time"
)

func TestPing(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(Ping("tcp", "google.com", "443"))
		time.Sleep(1 * time.Second)
	}
}
