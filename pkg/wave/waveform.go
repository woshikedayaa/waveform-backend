package wave

import "go.bug.st/serial"

type Waveform struct {
	portName   string
	mode       *serial.Mode
	port       serial.Port
	latestRead []byte
}

func New(portname string, mode *serial.Mode) (*Waveform, error) {
	var (
		wf   = new(Waveform)
		err  error
		port serial.Port
	)

	port, err = serial.Open(portname, mode)
	if err != nil {
		return
	}
}
