package main

import (
	"time"
)

const (
	PulseOnTime = 100 * time.Millisecond
)

var pulseOutput *GpioPin

func PreparePulseGenerator(gpio string) error {
	var err error
	var pin *GpioPin
	if pin, err = NewGpioPin(gpio); err != nil {
		return err
	}
	if err = pin.SetDirection(true); err != nil {
		return err
	}
	if err = pin.SetValue(false); err != nil {
		return err
	}
	pulseOutput = pin
	return nil
}

func PulseGenerate(count int) {
	go func() {
		for i := 0; i < count; i++ {
			if err := pulseOutput.SetValue(true); err != nil {
				break;
			}
			time.Sleep(100 * time.Millisecond)
			if err := pulseOutput.SetValue(false); err != nil {
				break;
			}
			SleepTime := time.Duration(settings.Value("FacePlaneCtrlPeriod", 1.0).(float64) * float64(time.Second))
			time.Sleep(SleepTime)
		}
	}()
}

