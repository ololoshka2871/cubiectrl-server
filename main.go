package main 

import (
	"log"
	"os"
	"time"
)

var settings *Settings

func main() {
	var (
		dataToServerChan chan CellData
	)
	
	log.Println("Reading settings..")
	s, err := NewSettings(os.Getenv("HOME") + "/.config/cubiectrl/cubiectrl.json")
	if err != nil {
		panic("Failed to create settings " + err.Error())
	}
	settings = s
	
	StartDefault()
	
	log.Println("Starting modbus..")
	port, ok1 := settings.Value("Port", "/dev/ttyS0").(string)
	baudRate, ok2 := settings.Value("BoudRate", 57600.0).(float64)
	rtsPin, ok3 := settings.Value("RtsPin", "gpio3_pg8").(string)
	if ok1 && ok2 && ok3 {
		dataToServerChan = make(chan CellData)
		dataToDisplayChan := make(chan CellData)
		
		if err := ValuesFormCtrlInit(dataToDisplayChan); err != nil {
			panic("Failed to start display form ctrl: " + err.Error())
		}
		
		resChan, err := StartModbusClient(port, int(baudRate), rtsPin, settings)
		if err != nil {
			log.Fatal(err)
		} else {
			go func(c <-chan Cell) {
				for v := range c {
					val, err := v.valueAsFloat()
					data4server := CellData{ Name : v.Name, Timestamp : time.Now() }
					if err != nil {
						data4server.Error = true
					} else {
						data4server.Value = val
					}
					dataToServerChan <- data4server
					dataToDisplayChan <- data4server
				}
			}(resChan)
		} 
	} else {
		panic("Settings read error")
	}
	
	log.Println("Starting server..")
	srvPort, ok1 := settings.Value("serverPort", 8081).(int)
	if !ok1 {
		panic("Settings port error")
	} 
	srv := NewServer(srvPort, dataToServerChan, settings) 
	if err := srv.ListenAndServe(); err != nil {
		panic("server failed to start: " + err.Error())
	}
}

