package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sharkrf/srf-fwu/srf-fwu"
)

func main() {
	var settings srffwu.Settings

	fmt.Println("srf-fwu - SharkRF Bootloader USB serial console firmware upgrade tool")

	settings.PortName = "/dev/ttyACM0"
	flag.StringVar(&settings.PortName, "p", settings.PortName, "port (ex. /dev/ttyACM0 or COM12)")
	flag.StringVar(&settings.FwFileName, "f", settings.FwFileName, "firmware file name to load")
	flag.BoolVar(&settings.Verbose, "v", settings.Verbose, "verbose")
	flag.Parse()

	if err := srffwu.FwDataRead(settings.FwFileName); err != nil {
		fmt.Fprintf(os.Stderr, "%s, exiting\n", err.Error())
		os.Exit(1)
	}

	retryOnce := true
	for {
		retry, err := srffwu.Start(settings)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())

			if err.Error() == "timeout" {
				if !retry && retryOnce {
					retryOnce = false
					retry = true

					fmt.Println("rebooting bootloader and retrying")
					srffwu.SerialPortOpen(settings.PortName)
					if settings.Verbose {
						fmt.Println("out: rbb\r")
					}
					srffwu.SerialPortWrite("rbb\r")
					srffwu.SerialPortClose()
				}
			} else {
				os.Exit(1)
			}
		}

		if retry {
			for i := 3; i > 0; i-- {
				fmt.Printf("retrying firmware upgrade in %d seconds...\n", i)
				time.Sleep(time.Second)
			}
		} else {
			break
		}
	}
}
