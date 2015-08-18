package main

import (
	"time"
	"log"
)

func PrepareModeSwitcher(pin string) error {
	if switcherPin, err := NewGpioPin(pin); err != nil {
		return err
	} else {
		if err := switcherPin.SetDirection(false); err != nil {
			return err
		}
		
		go func() {
			time.Sleep(time.Second) // prestart
			
			var v bool
			prevState, err := switcherPin.Value()
			if err != nil {
				log.Println(readErrStr)
				return
			}
			
			for {
				time.Sleep(Enable_update_interval)
				v, err = switcherPin.Value()
				if err != nil {
					log.Println(readErrStr)
					return
				} else {
					if v != prevState {
						if !v {
							switch CurrentDisplayState.BigDisplayMode {
								case Diable_bigDisplay:
									ControlBigDisplay(ShowVideo_bigDisplay)
								case ShowVideo_bigDisplay:
									ControlBigDisplay(ShowQMLForm_bigDisplay)
								case ShowQMLForm_bigDisplay:
									ControlBigDisplay(ShowVideo_bigDisplay)
								default : 
									panic("Unknown big display state!")
							}
						}
					}
				}
			}
		}()
		
		return nil
	}
}