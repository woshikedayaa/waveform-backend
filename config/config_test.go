package config

import (
	"fmt"
	"os"
	"testing"
)

func TestParseConfig(t *testing.T) {
	err := InitConfig()
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	fmt.Printf("%#v", G())
}
