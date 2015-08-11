package main

import (
	"os"
	"errors"
)

const (
	Diable_bigDisplay = 0
	ShowVideo_bigDisplay = 1
	ShowQMLForm_bigDisplay = 2
	
	Player = "mpv"
	Big_Display = "DISPLAY=:0.1"
	Small_Display = "DISPLAY=:0.0"
)

var PlayerArgsCommon = []string{"--fs", "--loop=inf"}

type tCurrentDisplayState struct {
	SmallDisplayMode bool
	BigDisplayMode int
	
	SmallDisplayPlayerProcess *os.Process
}

var CurrentDisplayState tCurrentDisplayState

func StartDefault() {
	//ControlSmallDisplay(true)
	//ControlBigDisplay(ShowVideo_bigDisplay)
} 

func ControlSmallDisplay(enable bool) error {
	if enable != CurrentDisplayState.SmallDisplayMode {
		if enable {
			env := append(os.Environ(), Big_Display) // set DISPLAY env
			var PlayerArgs []string
			if media, ok := settings.Value("SmallDispFileName", "").(string); !ok || media == "" {
				return errors.New("Playing media error")
			} else {
				PlayerArgs = append(PlayerArgsCommon, media)
			}
			if proc, err :=	os.StartProcess(Player, PlayerArgs, &os.ProcAttr{Env : env}); err != nil {
				CurrentDisplayState.SmallDisplayPlayerProcess = proc
			} else {
				return err
			}
		} else {
			if CurrentDisplayState.SmallDisplayPlayerProcess != nil {
				CurrentDisplayState.SmallDisplayPlayerProcess.Kill()
				CurrentDisplayState.SmallDisplayPlayerProcess = nil
			} else {
				CurrentDisplayState.SmallDisplayMode = false
				return errors.New("Player not running")
			}
		}
		CurrentDisplayState.SmallDisplayMode = enable
	}
	return nil
}

func ControlBigDisplay(ctrl int) error {
	if ctrl != CurrentDisplayState.BigDisplayMode {
		//
		
		CurrentDisplayState.BigDisplayMode = ctrl
	}
	return nil
}