package main 

import (
	"bitbucket.org/Olololshka/cubiectrl"
	"log"
	"os"
	"time"
)

func main() {
	var (
		dataToServerChan chan cubiectrl.CellData
	)
	
	log.Println("Reading settings..")
	settings, err := NewSettings(os.Getenv("HOME") + "/.config/cubiectrl/cubiectrl.json")
	if err != nil {
		panic("Failed to create settings")
	}
	log.Println("Starting modbus..")
	port, ok1 := settings.Value("Port", "/dev/ttyS0").(string)
	baudRate, ok2 := settings.Value("BoudRate", 57600).(float64)
	rtsPin, ok3 := settings.Value("RtsPin", "gpio3_pg8").(string)
	if ok1 && ok2 && ok3 {
		dataToServerChan = make(chan cubiectrl.CellData)
		
		resChan, err := StartModbusClient(port, int(baudRate), rtsPin, settings)
		if err != nil {
			log.Fatal(err)
		} else {
			go func(c <-chan Cell) {
				for v := range c {
					val, err := v.valueAsFloat()
					data4server := cubiectrl.CellData{ Name : v.Name, Timestamp : time.Now() }
					if err != nil {
						data4server.Error = true
					} else {
						data4server.Value = val
					}
					dataToServerChan <- data4server
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
	srv := cubiectrl.NewServer(srvPort, dataToServerChan, settings) 
	if err := srv.ListenAndServe(); err != nil {
		panic("server failed to start: " + err.Error())
	}
}

