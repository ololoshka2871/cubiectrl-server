package main 

import (
	"bitbucket.org/Olololshka/cubiectrl"
	"log"
	"os"
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
	port, ok1 := settings.Value("port", "/dev/ttyS0").(string)
	baudRate, ok2 := settings.Value("boudRate", 57600).(int)
	rtsPin, ok3 := settings.Value("RtsPin", "gpio3_pg8").(string)
	if ok1 && ok2 && ok3 {
		dataToServerChan = make(chan cubiectrl.CellData)
		
		resChan, err := StartModbusClient(port, baudRate, rtsPin)
		if err != nil {
			log.Fatal(err)
		} else {
			go func(c <-chan Cell) {
				for v := range c {
					val, err := v.valueAsFloat()
					data4server := cubiectrl.CellData{Name : v.Name}
					if err != nil {
						data4server.Error = err
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
	srv := cubiectrl.NewServer(srvPort, dataToServerChan)
	if err := srv.ListenAndServe(); err != nil {
		panic("server failed to start")
	}
}

