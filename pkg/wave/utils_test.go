package wave

import (
	"fmt"
	"testing"
)

func TestScanAvailablePort(t *testing.T) {
	port, err := ScanAvailablePort()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(port); i++ {
		fmt.Println(port[i])
	}
}
