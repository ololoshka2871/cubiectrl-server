package main

import (
	"os"
	"errors"
	"log"
)

const (
	Diable_bigDisplay = 0
	ShowVideo_bigDisplay = 1
	ShowQMLForm_bigDisplay = 2
	
	Player = "/usr/bin/mpv"
	Big_Display = "DISPLAY=:0.1"
	Small_Display = "DISPLAY=:0.0"
)

var PlayerArgsCommon = []string{"--fs", "--loop=inf"}

type tCurrentDisplayState struct {
	SmallDisplayMode bool
	BigDisplayMode int
	
	SmallDisplayPlayerProcess 	*os.Process
	BigDisplayPlayerProcess 	*os.Process
	BigDisplayValuesProcess		*os.Process
}

var CurrentDisplayState tCurrentDisplayState

func StartDefault() {
	//ControlSmallDisplay(true)
	//ControlBigDisplay(ShowVideo_bigDisplay)
} 

func ControlSmallDisplay(enable bool) error {
	if enable != CurrentDisplayState.SmallDisplayMode {
		if enable {
			env := append(os.Environ(), Small_Display) // set DISPLAY env
			var PlayerArgs []string
			if media, ok := settings.Value("SmallDispFileName", "").(string); !ok || media == "" {
				return errors.New("Playing media error")
			} else {
				PlayerArgs = append(PlayerArgsCommon, media)
			}
			if proc, err :=	os.StartProcess(Player, PlayerArgs, &os.ProcAttr{Env : env}); err == nil {
				CurrentDisplayState.SmallDisplayPlayerProcess = proc
			} else {
				log.Println("Failed to start plaing on SMALL display")
				return err
			}
		} else {
			if CurrentDisplayState.SmallDisplayPlayerProcess != nil {
				CurrentDisplayState.SmallDisplayPlayerProcess.Kill()
				CurrentDisplayState.SmallDisplayPlayerProcess = nil
			} else {
				CurrentDisplayState.SmallDisplayMode = false
				S := "Player not running"
				log.Println(S)
				return errors.New(S)
			}
		}
		CurrentDisplayState.SmallDisplayMode = enable
		if enable {
			log.Println("Start playing on small display")
		} else {
			log.Println("Stop playing on small display")
		}
	}
	return nil
}

func ControlBigDisplay(ctrl int) error {
	if ctrl != CurrentDisplayState.BigDisplayMode {
		env := append(os.Environ(), Big_Display) // set DISPLAY env
		switch(ctrl) {
			case Diable_bigDisplay:
				if CurrentDisplayState.BigDisplayPlayerProcess != nil {
					CurrentDisplayState.BigDisplayPlayerProcess.Kill()
					CurrentDisplayState.BigDisplayPlayerProcess = nil
				}
				if CurrentDisplayState.BigDisplayValuesProcess != nil {
					CurrentDisplayState.BigDisplayValuesProcess.Kill()
					CurrentDisplayState.BigDisplayValuesProcess = nil
				}
				log.Println("Big Display disabled")
				return nil
				
			case ShowVideo_bigDisplay:
				if CurrentDisplayState.BigDisplayValuesProcess != nil {
					CurrentDisplayState.BigDisplayValuesProcess.Kill()
					CurrentDisplayState.BigDisplayValuesProcess = nil
				}
				if CurrentDisplayState.BigDisplayPlayerProcess != nil {
					CurrentDisplayState.BigDisplayMode = ShowVideo_bigDisplay
					return errors.New("Allready playing")
				}
				
				var PlayerArgs []string
				if media, ok := settings.Value("BigDispFileName", "").(string); !ok || media == "" {
					return errors.New("Playing media error")
				} else {
					PlayerArgs = append(PlayerArgsCommon, media)
				}
				if proc, err :=	os.StartProcess(Player, PlayerArgs, &os.ProcAttr{Env : env}); err == nil {
					CurrentDisplayState.BigDisplayPlayerProcess = proc
				} else {
					log.Println("Failed to start plaing on BIG display")
					return err
				}
				
				log.Println("Start playing on Big display")
				return nil
				
			case ShowQMLForm_bigDisplay:
				if CurrentDisplayState.BigDisplayPlayerProcess != nil {
					CurrentDisplayState.BigDisplayPlayerProcess.Kill()
					CurrentDisplayState.BigDisplayPlayerProcess = nil
				}
				if CurrentDisplayState.BigDisplayValuesProcess != nil {
					CurrentDisplayState.BigDisplayMode = ShowQMLForm_bigDisplay
					return errors.New("Allready Displaying")
				}
				
				log.Println("Displaing values on big display")
				return nil
				
			default :
				return errors.New("Incorrect ctrl request")
		}
		
		CurrentDisplayState.BigDisplayMode = ctrl
	}
	return nil
}