package main

import (
	"github.com/ololoshka2871/go-modbus"
	"github.com/tarm/serial"
	"os"
	"io"
	"time"
)

type ModbusReader interface {
	ReadAll(*serial.Port) error
	Test(*serial.Port) bool
} 

type RWControlPin struct {
	valueFile *os.File
}

func NewRWControlPin(pin string) (*RWControlPin, error) {

	direction := "/sys/class/gpio/" + pin + "/direction"
	value := "/sys/class/gpio/" + pin + "/value"

	valueFile, err := os.OpenFile(value, os.O_WRONLY, 0664)
	if err != nil {
		return nil, err
	}

	result := &RWControlPin{valueFile}

	directionFile, err := os.OpenFile(direction, os.O_WRONLY, 0664)
	if err != nil {
		return nil, err
	}
	defer directionFile.Close()

	_, err = directionFile.Write([]byte("out"))
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (this *RWControlPin) WriteHook(_ io.ReadWriteCloser, newval bool) {

        var txt []byte
        if newval {
                txt = []byte("1")
        } else {
                txt = []byte("0")
        }

        _, err := this.valueFile.Write(txt)
        if err != nil {
                panic(err)
        }
}

func StartModbusClient(serialPort string, baudRate int, RTS_Pin string) (<-chan Cell, error) {
	ctx, cerr := modbusclient.ConnectRTU(serialPort, baudRate)
	if cerr != nil {
		return nil, cerr
	} else {
		hook, err := NewRWControlPin(RTS_Pin)
		if err != nil { 
			panic(err.Error()) 
		} else {
			modbusclient.SendHook = hook
		}
		
		cells := BuildCellsTable()
		resultChan := make(chan Cell)
		
		go func(cells []Cell, updateInterval time.Duration) {
			for {
				// update thread
				time.Sleep(updateInterval)
				for _, cell := range cells {
					err := cell.Read(ctx)
					if err == nil {
						resultChan <- cell
					}
				}
			}
		}(cells, 100 * time.Millisecond)
		
		return resultChan, nil
	}
}

func BuildCellsTable() []Cell {
	result := make([]Cell, 4)
	result[0] = Cell{DevAddr : 4, CellType : InputRegister, CellLen_mbCells : 2, CellStartAddr : 0}
	result[1] = Cell{DevAddr : 4, CellType : InputRegister, CellLen_mbCells : 2, CellStartAddr : 2}
	result[2] = Cell{DevAddr : 4, CellType : InputRegister, CellLen_mbCells : 2, CellStartAddr : 4}
	result[3] = Cell{DevAddr : 4, CellType : InputRegister, CellLen_mbCells : 2, CellStartAddr : 6}
	return result
}
