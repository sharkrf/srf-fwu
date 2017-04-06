package srffwu

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

type state int

const (
	stateInit state = iota
	stateWaitingForFirstStatus
	stateSending
)

// Settings stores parameters for the firmware upgrade process.
type Settings struct {
	PortName   string
	FwFileName string
	Verbose    bool
}

func printProgress(bts BootloaderStatus) {
	p := (float64(bts.fwStored) / float64(bts.fwSize)) * 100.0
	if math.IsNaN(p) {
		p = 0
	}
	if p >= 100 {
		fmt.Printf("progress: 100%% (waiting for result...)\n")
	} else {
		fmt.Printf("progress: %.2f%%\n", p)
	}
}

func sendHexChunk(settings Settings) {
	// Sending the next hex chunk (if needed).
	hexChunk := FwDataGetHexChunk()
	if hexChunk != "" {
		if settings.Verbose {
			fmt.Println("out: dta " + hexChunk + "\r")
		}
		SerialPortWrite("dta " + hexChunk + "\r")
	}
}

// Start starts the firmware upgrade process. Returns true if a retry
// (function recall) is needed.
func Start(settings Settings) (bool, error) {
	if err := SerialPortOpen(settings.PortName); err != nil {
		return false, fmt.Errorf("error opening serial port %s (%v), exiting\n", settings.PortName, err.Error())
	}

	defer SerialPortClose()

	c := make(chan string)
	go SerialPortReader(c)

	fmt.Println("identifying bootloader...")

	var fwuState state
	var bts BootloaderStatus
	var err error
	var deviceIdentifier string

	if settings.Verbose {
		fmt.Println("out: \r\r\r")
	}
	SerialPortWrite("\r\r\r")

	for {
		select {
		case line := <-c:
			if settings.Verbose {
				fmt.Println("in: " + line)
			}

			toks := strings.Split(line, " ")

			switch fwuState {
			default:
				if len(toks) >= 4 && toks[0] == "sercon:" && toks[1] == "inf:" && toks[2] == "SharkRF" && toks[3] == "Bootloader" {
					fmt.Println("found bootloader: " + strings.Join(toks[2:5], " "))
					deviceIdentifier = toks[5]
					fmt.Println("device identifier: " + deviceIdentifier)
					fwuState = stateWaitingForFirstStatus
				}
			case stateWaitingForFirstStatus:
				bts, err = BootloaderStatusLineParse(toks)
				if err == nil {
					if bts.dataproc != "ready" {
						BootloaderStatusPrint(bts)
						fmt.Println("bootloader is not in ready state, rebooting device")
						if settings.Verbose {
							fmt.Println("out: rbb\r")
						}
						SerialPortWrite("rbb\r")
						return true, nil
					}
					fmt.Println("bootloader is ready, starting firmware upgrade")
					fwuState = stateSending
					sendHexChunk(settings)
				}
			case stateSending:
				bts, err = BootloaderStatusLineParse(toks)
				if err == nil {
					printProgress(bts)

					// Checking results.
					if bts.flash != "ok" || bts.configarea != "ok" ||
						(bts.dataproc != "working" && bts.dataproc != "ready") {
						BootloaderStatusPrint(bts)
						if bts.flash == "ok" && bts.configarea == "ok" && bts.dataproc == "success" {
							fmt.Println("firmware upgraded successfully!")
						} else {
							fmt.Println("firmware upgrade failed!")
						}
						return false, nil
					}

					sendHexChunk(settings)
				}
			}
		case <-time.After(time.Second * 5):
			return false, errors.New("timeout")
		}
	}
}
