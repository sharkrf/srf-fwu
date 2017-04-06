package srffwu

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jacobsa/go-serial/serial"
)

var port io.ReadWriteCloser

// readLine returns a single line (without the ending \n) from the input buffered reader.
// An error is returned if there is an error with the buffered reader.
func readLine(r *bufio.Reader) (string, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

// SerialPortReader is a goroutine which reads lines from the given buffered reader,
// and sends them to the given channel.
func SerialPortReader(c chan string) {
	r := bufio.NewReaderSize(port, 1)

	for {
		line, err := readLine(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading serial port, exiting\n")
			os.Exit(1)
		}
		c <- strings.TrimSpace(line)
	}
}

// SerialPortWrite writes the given string to the serial port.
func SerialPortWrite(s string) {
	port.Write([]byte(s))
}

// SerialPortOpen opens the serial port device portName, and returns it's handle.
func SerialPortOpen(portName string) error {
	options := serial.OpenOptions{
		PortName:        portName,
		BaudRate:        115200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	fmt.Println("opening serial port " + portName)

	var err error
	port, err = serial.Open(options)
	return err
}

// SerialPortClose closes the serial port device.
func SerialPortClose() {
	fmt.Println("closing serial port")
	port.Close()
}
