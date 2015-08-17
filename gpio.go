package main

import (
	"os"
	"errors"
	"fmt"
)

type GpioPin struct {
	valueFile 	*os.File
	dirFile		*os.File
	direction	bool
}

const (
	sysGpioPath = "/sys/class/gpio/"
)
 

func NewGpioPin(pin string) (*GpioPin, error) {
	// is exists
	if _, err := os.Stat(sysGpioPath + pin); os.IsNotExist(err) {
		return nil, errors.New(fmt.Sprintf("Pin %s not exported!", pin))
	}

	direction := sysGpioPath + pin + "/direction"
	value := sysGpioPath + pin + "/value"

	valueFile, err := os.OpenFile(value, os.O_RDWR, 0664)
	if err != nil {
		return nil, err 
	}

	directionFile, err := os.OpenFile(direction, os.O_RDWR, 0664)
	if err != nil {
		return nil, err
	}
 	
	result := &GpioPin{valueFile, directionFile, false}
	if dir, err := result.Direction(); err != nil {
		valueFile.Close()
		directionFile.Close()
		return nil, err
	} else {
		result.direction = dir
	}

	return result, nil
}

func (this *GpioPin) Direction() (bool, error) {
	var dirStr []byte
	
	if _, err := this.dirFile.Read(dirStr); err != nil {
		return false, nil
	}
	
	switch string(dirStr) {
		case "in" : return false, nil;
		case "out" : return true, nil;
		default: panic("Unknow direction: " + string(dirStr))
	}
}

func (this *GpioPin) SetDirection(out bool) error {
	dir := "in"
	if out {
		dir = "out"
	}
	_, err := this.dirFile.Write([]byte(dir))
	if err != nil {  
		return err
	}
	this.direction = out
	return nil
}

func (this *GpioPin) SetValue(val bool) error {
	
	if !this.direction {
		return errors.New("Pin is input")
	}
	
	v := "0"
	if val {
		v = "1"
	}
	_, err := this.valueFile.Write([]byte(v))
	if err != nil {  
		return err
	}
	return nil
}

func (this *GpioPin) Value() (bool, error) {
	var dirStr []byte
	
	if _, err := this.valueFile.Read(dirStr); err != nil {
		return false, nil
	}
	
	switch string(dirStr) {
		case "0" : return false, nil;
		case "1" : return false, nil;
		default: panic("Unknow value: " + string(dirStr))
	}
}