package srffwu

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// BootloaderStatus holds the status of the bootloader.
type BootloaderStatus struct {
	dataproc     string
	app          string
	configarea   string
	flash        string
	flashErrAddr uint64
	fwSize       uint64
	fwProcessed  uint64
}

// BootloaderStatusLineParse parses the given status line (in split tokens)
// and returns it a parsed struct.
func BootloaderStatusLineParse(toks []string) (BootloaderStatus, error) {
	var bts BootloaderStatus

	if len(toks) < 2 || toks[0] != "sercon:" || toks[1] != "status:" {
		return bts, errors.New("invalid status line given")
	}

	bts.dataproc = ""
	bts.app = ""
	bts.configarea = ""
	bts.flash = ""
	bts.flashErrAddr = 0
	bts.fwSize = 0
	bts.fwProcessed = 0

	for i := 2; i < len(toks); i++ {
		switch toks[i] {
		case "dataproc:":
			i++
			bts.dataproc = strings.Trim(toks[i], ",")
		case "app:":
			i++
			bts.app = strings.Trim(toks[i], ",")
		case "configarea:":
			i++
			bts.configarea = strings.Trim(toks[i], ",")
		case "flash:":
			i++
			bts.flash = strings.Trim(toks[i], ",")
		case "erraddr:":
			i++
			bts.flashErrAddr, _ = strconv.ParseUint(strings.Trim(toks[i], ","), 10, 32)
		case "fwsize:":
			i++
			bts.fwSize, _ = strconv.ParseUint(strings.Trim(toks[i], ","), 10, 32)
		case "processed:":
			i++
			bts.fwProcessed, _ = strconv.ParseUint(strings.Trim(toks[i], ","), 10, 32)
		}
	}
	return bts, nil
}

// BootloaderStatusPrint prints the contents of the given status struct.
func BootloaderStatusPrint(bts BootloaderStatus) {
	fmt.Printf("device status:\n"+
		"  dataproc: %s\n"+
		"  app: %s\n"+
		"  configarea: %s\n"+
		"  flash: %s\n"+
		"  flash error address: %d\n"+
		"  fw size from header: %d\n"+
		"  processed fw bytes: %d\n", bts.dataproc, bts.app, bts.configarea, bts.flash,
		bts.flashErrAddr, bts.fwSize, bts.fwProcessed)
}
