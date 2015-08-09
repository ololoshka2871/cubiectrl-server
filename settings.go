package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"os"
	"fmt"
)

type Settings struct {
	data map[string]*json.RawMessage
	filename string
}

func NewSettings(jsonFile string) (*Settings, error) {
	jsonData, e := ioutil.ReadFile(jsonFile)
	result := &Settings{make(map[string]*json.RawMessage), jsonFile}
	if e == nil {
		e = json.Unmarshal(jsonData, &result.data)
		if e != nil {
			return nil, e
		}
		return result, nil
	} else {
		err := result.SetDefault()
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func createPath(path string) {
	if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
		createPath(filepath.Dir(path))
	} else {
		err = os.Mkdir(path, 0775)
		if err != nil {
			panic(fmt.Sprintf("Failed to create %s : %s", path, err.Error()))
		}
	}
}

func (this *Settings) Sync() error {
	jsonData, err := json.Marshal(this.data)
	if err != nil {
		return err
	}
	createPath(filepath.Dir(this.filename))
	err = ioutil.WriteFile(this.filename, jsonData, 0660)
	if err != nil {
		return err
	}
	return nil
}

func (this *Settings) SetDefault() error {
	//
	//
	return this.Sync()
}

func (this *Settings) Value(key string, defaultVal interface{}) interface{} {
	var result interface{}
	err := json.Unmarshal(*this.data[key], &result)
	if err != nil {
		return defaultVal
	} 
	return result
}