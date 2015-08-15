package main

import (
	"net"
	"os"
	"log"
	"fmt"
)

const (
	ValuesExecutable = "/bin/sh" //TODO
	SocketName = "/tmp/ValuesCtrlSock.unix"
)

var ValuesProcArgs = []string{ValuesExecutable, SocketName}

var commandChans []chan string

func SendValuesCmd(cmd string) {
	for _, c:= range commandChans {
		c <- cmd
	}
}

func ValuesFormCtrlInit(d <-chan CellData) error {
	// remove socket file if allready exists
	if _, err := os.Stat(SocketName); !os.IsNotExist(err) {
		os.Remove(SocketName)
	}
	
	if l, err := net.Listen("unix", SocketName); err != nil {
        return err
    } else {
    	commandChans := make([]chan string, 0)
    	
		go func() {
			for {
		        fd, err := l.Accept()
		        if err != nil {
		            log.Fatal("accept error:", err)
		        }
		        cmdChan := make(chan string)
		        commandChans = append(commandChans, cmdChan)
		
		        go func() {
		        	for msg := range cmdChan {
		        		fd.Write([]byte(msg))
		        	}
		        }()
	    	}
		}()
		
		go func() {
			for data := range d {
				if !data.Error {
					res := fmt.Sprintf("%s=%f", data.Name, data.Value)
					for _, c:= range commandChans {
						c <- res
					}
				}
			}
		}()
		
		env := append(os.Environ(), Big_Display) // set DISPLAY env
		if _, err := os.StartProcess(ValuesExecutable, ValuesProcArgs, &os.ProcAttr{Env : env}); err != nil {
			return err
		}
	}
	
	return nil
}
