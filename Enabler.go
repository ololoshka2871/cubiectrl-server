package main

import (
	"log"
	"time"
)

const (
	 readErrStr = "Failed to read ctrl pin"
	 Enable_update_interval = 200 * time.Millisecond
)

func Started() {
	log.Println("Start event detected!")
	StartDefault()
	
	PulseGenerate(2 * 8)
}

func Stopped() {
	ControlSmallDisplay(false)
	
	ControlBigDisplay(Diable_bigDisplay)
	log.Println("Shutdown event detected!")
}

func SetupEnabler(enableInput string) error {
	if enablePin, err := NewGpioPin(enableInput); err != nil {
		return err
	} else {
		if err := enablePin.SetDirection(false); err != nil {
			return err
		}
		
		go func() {
			time.Sleep(time.Second) // prestart
			
			var v bool
			prevState, err := enablePin.Value()
			if err != nil {
				log.Println(readErrStr)
				return
			}
			
			for {
				time.Sleep(Enable_update_interval)
				v, err = enablePin.Value()
				if err != nil {
					log.Println(readErrStr)
					return
				} else {
					if v != prevState {
						if v {
							// Start
							go Started()
						} else {
							// stop
							go Stopped()
						}
					}
				}
			}
		}()
		
		return nil
	}
}

