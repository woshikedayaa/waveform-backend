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
		return nil, &OpError{
			raw:  err,
			op:   "open",
			port: portname,
		}
	}

	wf.port = port
	wf.mode = mode
	wf.portName = portname
	wf.latestRead = nil
	return wf, nil
}

func (wf *Waveform) Close() error {
	return wf.port.Close()
}

func (wf *Waveform) Latest() []byte {
	return wf.latestRead
}

func (wf *Waveform) PortName() string {
	return wf.portName
}
