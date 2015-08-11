package main

import (
	"net/http"
	"html/template"
	"net/url"
	"fmt"
	"path"
	"runtime"
	"encoding/json"
	"log"
	"strings"
	"strconv"
)

var settingsmap = map[string]interface{} { "SmallDispFileName" : "",
			"BigDispFileName" : "",
			"Port" : "/dev/ttyS0",
			"BoudRate" : 57600,
			"RtsPin" : "gpio3_pg8",
			"UpdateDelay" : 100,
			};

type SettingsHolder interface{
	Value(key string, defaultVal interface{}) interface{}
	SetDefault() error
	Sync() error
	SetValue(key string, Val interface{}) error
}

type IndexPageParams struct {
	SmallDispFileName string
	BigDispFileName string
	Port string
	BoudRate int
	RtsPin string
	UpdateDelay int
}

type __Settings struct {
	IndexPageParams
	OK bool
}

var currentPath string
var fileserverHandler http.Handler

func FillPagesMap(m *map[string]func(http.ResponseWriter, *http.Request)) {
	(*m)["/"] = indexHandlr
	(*m)["/asserts"] = assertsServer
	(*m)["/data.api"] = varsJsonHandlr
	
	_, filename, _, _ := runtime.Caller(1)
	currentPath = path.Dir(filename)
	
	fileserverHandler = http.FileServer(http.Dir(currentPath))
}

func patchPath(name1 string, names ...string) []string {
	
	res := make([]string, len(names) + 1)
	
	res[0] = currentPath + "/" + name1
	for i, na := range names {
		res[i + 1] = currentPath + "/" + na
	}
	
	return res
}

func varsJsonHandlr(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Query().Get("req") {
		case "mesurment":
			log.Print("Sending mesurment")
			mtype := r.URL.Query().Get("type")
			if mtype == "" {
				err := json.NewEncoder(w).Encode(results)
				if err != nil {
					log.Print(err.Error())
					return
				}
				
			} else {
				r := make(map[string]CellData)
				for _, cell := range results {
					if strings.Contains(cell.Name, mtype) {
						r[cell.Name] = cell
					}	
				}
				err := json.NewEncoder(w).Encode(r)
				if err != nil {
					log.Print(err.Error())
					return
				}
			}
		case "getSettings":
			params := __Settings{ IndexPageParams : *getParamsMap(), OK : true}
			err := json.NewEncoder(w).Encode(params)
				if err != nil {
					log.Print(err.Error())
					return
				}
		case "setSettings":
			if err := ApplySettings(r.URL.Query()); err != nil {
				fmt.Fprintf(w, err.Error())
			} else {
				fmt.Fprintf(w, "OK")
			}
		case "resetSettings":
			err := settingsValues.SetDefault()
			if err != nil {
				fmt.Fprintf(w, err.Error())
			} else {
				fmt.Fprintf(w, "OK")
			}
		default :
			fmt.Fprintf(w, "No requestParameters")
			log.Printf("Unknown api request: %s", r.Form)
	}
}

func ApplySettings(values url.Values) error {
	for key, _ := range settingsmap {
		val := values.Get(key)
		var v interface{}
		var e error
		
		switch key {
			case "BoudRate":
				if v, e = strconv.Atoi(val); e != nil {
					v = settingsmap["BoudRate"]
				}
			case "UpdateDelay":
				if v, e = strconv.Atoi(val); e != nil {
					v = settingsmap["UpdateDelay"]
				}
			default:
				v = val
		}
		
		if err := settingsValues.SetValue(key, v); err != nil {
			return err
		}
	}
	return nil
}

func getParamsMap() *IndexPageParams {
	return &IndexPageParams{
    	SmallDispFileName : settingsValues.Value("SmallDispFileName", settingsmap["SmallDispFileName"]).(string),
		BigDispFileName : settingsValues.Value("BigDispFileName", settingsmap["BigDispFileName"]).(string),
		Port : settingsValues.Value("Port", settingsmap["Port"]).(string),
		BoudRate : int(settingsValues.Value("BoudRate", settingsmap["BoudRate"]).(float64)),
		RtsPin : settingsValues.Value("RtsPin", settingsmap["RtsPin"]).(string),
		UpdateDelay : int(settingsValues.Value("UpdateDelay", settingsmap["UpdateDelay"]).(float64)),
    	}
}

func indexHandlr(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles(patchPath(
    	"templates/index.html", 
    	"templates/header.html", 
    	"templates/footer.html")...)
        
    params := getParamsMap()
        
    if err != nil {
    	fmt.Fprintf(w, "Error %s", err.Error())
    	return
    }
    
    if err := t.ExecuteTemplate(w, "index", params); err != nil {
    	log.Println(err.Error())
    }
}

func assertsServer(w http.ResponseWriter, r *http.Request) {
	fileserverHandler.ServeHTTP(w, r)
}