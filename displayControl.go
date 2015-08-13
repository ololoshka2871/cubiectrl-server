package main

import (
	"os"
	"errors"
	"log"
	"os/exec"
	"io"
)

const (
	Diable_bigDisplay = 0
	ShowVideo_bigDisplay = 1
	ShowQMLForm_bigDisplay = 2
	
	Player = "/usr/bin/mpv"
	//QmlProgramm = "" /*TODO*/
	Big_Display = "DISPLAY=:0.0"
	Small_Display = "DISPLAY=:0.1"
	
	PauseKey = "p"
	ToggleFSkey = "f"
)

var PlayerArgsCommon = []string{Player, "--fs", "--loop=inf"}

type tCurrentDisplayState struct {
	SmallDisplayMode bool
	BigDisplayMode int
	
	SmallDisplayPlayerProcess 	*os.Process
	BigDisplayPlayerProcess 	*exec.Cmd
	BigDisplayplayerStdinPipe	io.WriteCloser				
	
	BigDisplayValuesProcess		*exec.Cmd
}

var CurrentDisplayState tCurrentDisplayState

func prepareBigDisplay() error {
	if media, ok := settings.Value("BigDispFileName", "").(string); ok && media != "" {
		PlayerArgs := append(PlayerArgsCommon, media)
		CurrentDisplayState.BigDisplayPlayerProcess = exec.Command(Player, PlayerArgs...)
		CurrentDisplayState.BigDisplayPlayerProcess.Env = append(
			CurrentDisplayState.BigDisplayPlayerProcess.Env,Big_Display)
		if pipe, err := CurrentDisplayState.BigDisplayPlayerProcess.StdinPipe(); err == nil {
			CurrentDisplayState.BigDisplayplayerStdinPipe = pipe
			err := CurrentDisplayState.BigDisplayPlayerProcess.Start()
			if err != nil {
				CurrentDisplayState.BigDisplayPlayerProcess = nil
				return err
			}
		} else {
			CurrentDisplayState.BigDisplayPlayerProcess = nil
			return err
		}
		log.Println("Start player big display")
		return nil
	} else {
		return errors.New("Big display player config error")
	}
}

func StartDefault() {
	ControlSmallDisplay(true)
	
	prepareBigDisplay()
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

func togglePlayBigDisplay() {
	CurrentDisplayState.BigDisplayplayerStdinPipe.Write([]byte(PauseKey))
	CurrentDisplayState.BigDisplayplayerStdinPipe.Write([]byte(ToggleFSkey))
}

func ControlBigDisplay(ctrl int) error {
	if ctrl != CurrentDisplayState.BigDisplayMode {
		switch(ctrl) {
			case Diable_bigDisplay:
				log.Println("Stop player big display")
				if CurrentDisplayState.BigDisplayPlayerProcess != nil {
					togglePlayBigDisplay()
				}

			case ShowVideo_bigDisplay:
				log.Println("Start player big display")
				if CurrentDisplayState.BigDisplayPlayerProcess == nil {
					if err := prepareBigDisplay(); err != nil { 
						return err
					}
				}
				
				// TODO hide values form
				
				togglePlayBigDisplay()
				
			case ShowQMLForm_bigDisplay:
				log.Println("Show values big display")
				if CurrentDisplayState.BigDisplayPlayerProcess != nil {
					togglePlayBigDisplay()
				}
				// TODO bring values form to front
				
			default :
				return errors.New("Incorrect ctrl request")
		}
		
		CurrentDisplayState.BigDisplayMode = ctrl
	}
	return nil
}