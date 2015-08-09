package main 

import (
	"bitbucket.org/Olololshka/cubiectrl"
	"log"
	"os"
)

func main() {
	log.Println("Reading settings..")
	settings, err := NewSettings(os.Getenv("HOME") + "/.config/cubiectrl/cubiectrl.json")
	if err != nil {
		panic("Failed to create settings")
	}
	log.Println("Starting modbus..")
	port, ok1 := settings.Value("port", "/dev/ttyS0").(string)
	baudRate, ok2 := settings.Value("boudRate", 57600).(int)
	if ok1 && ok2 {
		resChan, err := StartModbusClient(port, baudRate)
		if err != nil {
			log.Fatal(err)
		} else {
			go func(c <-chan Cell) {
				for v := range c {
					val, err := v.valueAsFloat()
					if err != nil {
						log.Print(err)
					} else {
						log.Printf("Cell updated: %0.2f", val)
					}
				}
			}(resChan)
		} 
	} else {
		panic("Settings read error")
	}
	
	log.Println("Starting server..")
	srv := cubiectrl.NewServer(8081)
	if err := srv.ListenAndServe(); err != nil {
		panic("server failed to start")
	}
}

