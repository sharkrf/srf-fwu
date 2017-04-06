package srffwu

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math"
)

var fwData []byte
var fwDataPos int

// FwDataRead reads the given firmware file data to a byte array and returns it.
func FwDataRead(fwFileName string) error {
	if fwFileName == "" {
		return fmt.Errorf("no firmware data file specified")
	}

	fmt.Println("reading firmware data from " + fwFileName + "...")
	var err error
	fwData, err = ioutil.ReadFile(fwFileName)
	if err != nil {
		return fmt.Errorf("error reading firmware data from %s", fwFileName)
	}

	fwDataPos = 0
	return nil
}

// FwDataGetHexChunk returns the next 512-byte chunk in hex character pairs.
func FwDataGetHexChunk() string {
	remaining := len(fwData) - fwDataPos
	sendSliceLength := int(math.Min(float64(remaining), 512))
	res := hex.EncodeToString(fwData[fwDataPos : fwDataPos+sendSliceLength])
	fwDataPos += sendSliceLength
	return res
}
