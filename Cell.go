package main

import (
	"github.com/ololoshka2871/go-modbus"
	"github.com/tarm/serial"
	
	"encoding/binary"
    "math"
    "errors"
    "log"
    "fmt"
)

const (
	InputRegister = 0
	HoldingRegister = 1
	Coil = 2
	DiscreteInput = 3
	
	MBTimeout = 100
	MBDebug = true
)

type Cell struct {
	DevAddr byte
	CellType int
	CellStartAddr uint16
	CellLen_mbCells uint16
	value []byte
}

func (this *Cell) valueAsFloat() (float32, error) {
	
	if len(this.value) != 4 {
		return 0.0, errors.New("sizeof(Value) != sizeof(float)") 
	}
	bits := binary.LittleEndian.Uint32(this.value)
    res := math.Float32frombits(bits)
    return res, nil
}

func (this *Cell) Read(ctx *serial.Port) error {
	var fun byte
	
	switch this.CellType {
		case InputRegister:
			fun = modbusclient.FUNCTION_READ_INPUT_REGISTERS
		case HoldingRegister:
			fun = modbusclient.FUNCTION_READ_HOLDING_REGISTERS
		case Coil:
			fun = modbusclient.FUNCTION_READ_COILS
		case DiscreteInput:
			fun = modbusclient.FUNCTION_READ_DISCRETE_INPUTS
	}
	
	readResult, readErr := modbusclient.RTURead(ctx,
	                        this.DevAddr,
	                        fun,
	                        this.CellStartAddr,
	                        this.CellLen_mbCells,
	                        MBTimeout, MBDebug)
	if MBDebug {
		log.Println(fmt.Sprintf("Rx: %x", readResult))
	}
	
	if readErr != nil {
		if MBDebug {
			log.Printf("Read error: dev=%#x, Cell=%#x, error=%v", this.DevAddr, this.CellStartAddr, readErr)
		}
		return readErr
	}
	this.value = readResult[3:len(readResult) - 2]
	return nil
}