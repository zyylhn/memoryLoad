package memoryLoad

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"testing"
	"time"
)

//go:embed cmd/app/zscan_mac_x64
var expapp []byte

func TestLoadExecute(t *testing.T) {
	ctx, c := context.WithTimeout(context.Background(), time.Second*50)
	re, err := Load("", ctx, expapp, "ps", "-H", "172.16.95.1/24")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(re))
	c()
}
