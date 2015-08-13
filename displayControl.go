package main

import (
	"os"
	"errors"
	"log"
	"os/exec"
	"syscall"
)

const (
	Diable_bigDisplay = 0
	ShowVideo_bigDisplay = 1
	ShowQMLForm_bigDisplay = 2
	
	Player = "/usr/bin/mpv"
	//QmlProgramm = "" /*TODO*/
	Big_Display = "DISPLAY=:0.0"
	Small_Display = "DISPLAY=:0.1"
	
	PauseCmd = "cycle pause\n"
	ToggleFSCmd = "cycle fullscreen\n"
	CmdPipeFile = "/tmp/mpvctrl.fifo"
)

var PlayerArgsCommon = []string{Player, "--fs", "--loop=inf", "--input-file=" + CmdPipeFile}

type tCurrentDisplayState struct {
	SmallDisplayMode bool
	BigDisplayMode int
	
	SmallDisplayPlayerProcess 	*os.Process
	BigDisplayPlayerProcess 	*exec.Cmd	
	BigDisplayValuesProcess		*exec.Cmd
}

var CurrentDisplayState tCurrentDisplayState

func prepareBigDisplay() error {
	if media, ok := settings.Value("BigDispFileName", "").(string); ok && media != "" {
		PlayerArgs := append(PlayerArgsCommon, media)
		CurrentDisplayState.BigDisplayPlayerProcess = exec.Command(Player, PlayerArgs...)
		CurrentDisplayState.BigDisplayPlayerProcess.Env = append(
			CurrentDisplayState.BigDisplayPlayerProcess.Env,Big_Display)
		
		/* FIFO */
		if _, err := os.Stat(CmdPipeFile); os.IsNotExist(err) {
			if err := syscall.Mknod(CmdPipeFile, syscall.S_IFIFO|0666, 0); err == nil {
				err := CurrentDisplayState.BigDisplayPlayerProcess.Start()
				if err != nil {
					CurrentDisplayState.BigDisplayPlayerProcess = nil
					return err
				}
			} else {
				CurrentDisplayState.BigDisplayPlayerProcess = nil
				return err
			}
		}
		CurrentDisplayState.BigDisplayMode = ShowVideo_bigDisplay
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

func togglePlayBigDisplay() error {
	
	if f, err := os.OpenFile(CmdPipeFile, os.O_WRONLY, 0664); err == nil {
		defer f.Close()
		f.Write([]byte(PauseCmd))
		f.Write([]byte(ToggleFSCmd))
		return nil
	} else {
		return err
	}
}

func ControlBigDisplay(ctrl int) error {
	if ctrl != CurrentDisplayState.BigDisplayMode {
		switch(ctrl) {
			case Diable_bigDisplay:
				log.Println("Stop player big display")
				if CurrentDisplayState.BigDisplayPlayerProcess != nil {
					if err := togglePlayBigDisplay(); err != nil {
						return err;
					}
				}

			case ShowVideo_bigDisplay:
				log.Println("Start player big display")
				if CurrentDisplayState.BigDisplayPlayerProcess == nil {
					if err := prepareBigDisplay(); err != nil { 
						return err
					}
				}
				
				// TODO hide values form
				
				if err := togglePlayBigDisplay(); err != nil {
					return err;
				}
				
			case ShowQMLForm_bigDisplay:
				log.Println("Show values big display")
				if CurrentDisplayState.BigDisplayPlayerProcess != nil {
					if err := togglePlayBigDisplay(); err != nil {
						return err;
					}
				}
				// TODO bring values form to front
				
			default :
				return errors.New("Incorrect ctrl request")
		}
		
		CurrentDisplayState.BigDisplayMode = ctrl
	}
	return nil
}