package wave

import (
	"go.bug.st/serial"
)

func ScanAvailablePort() ([]string, error) {
	ports, err := serial.GetPortsList()

	if err != nil {
		return nil, &OpError{
			raw: err,
			op:  "scan",
		}
	}

	return ports, nil
}
