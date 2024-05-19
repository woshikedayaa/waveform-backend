package wave

import (
	"fmt"
	"go.bug.st/serial"
	"testing"
)

func TestWaveformWrite(t *testing.T) {
	wf, err := New("/dev/pts/0", &serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.TwoStopBits,
	})
	if err != nil {
		t.Fatal(err)
	}
	//
	wf.Port.Write([]byte("HELLO,WORLD\r\n"))
	//wf.Port.Write([]byte("HELLO,WORLD\r\n"))
	//wf.Port.Write([]byte("HELLO,WORLD\r\n"))

	bs := make([]byte, 128)
	for err == nil {
		var n int

		n, err = wf.Port.Read(bs)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(string(bs[:n]))
	}
	//
	defer wf.Close()
}
