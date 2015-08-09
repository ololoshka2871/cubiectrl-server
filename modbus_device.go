package main

import (
	"github.com/ololoshka2871/go-modbus"
	"github.com/tarm/serial"
)

const (
	DB00 = 0xdb00
)

type DB00Device struct {
	Address uint8
	
	val [4]float32
}

func decodeAnsver(d []byte) ([]uint16, error) {
	var (
		sliceStart, sliceStop 	int
		data                  	int16
		decodeErr             	error
		result					[]uint16
	)
	result = make([]uint16, len(d) - 3 - 2)

	for i := 0; i < int(d[2]); i++ {
		// take the next two bytes, if available
		sliceStart = 3 + i
		sliceStop = 3 + i + 2

					// decode them into integers
		data, decodeErr = modbusclient.DecodeHiLo(d[sliceStart:sliceStop])
		if decodeErr != nil {
			return nil, decodeErr	
		} else {
			result[i] = uint16(data)
		}
	}
	
	return result, nil
}

func (this *DB00Device) Test(ctx *serial.Port) bool {
	readResult, readErr := modbusclient.RTURead(ctx,
	                        byte(0x04),
	                        modbusclient.FUNCTION_READ_HOLDING_REGISTERS,
	                        uint16(0),
	                        uint16(1),
	                        300, true)
	
	if readErr != nil { 
		return false
		}
	
	value, readErr := decodeAnsver(readResult)
	
	return (readErr == nil) && (value[0] == DB00) 
}

