package main 

import (
    "net/http"
    "io"
    "sort"
    "strings"
    "log"
    "fmt"
    "time"
)

type CubieCtrlHttpServer struct {
	http.Server
	mux map[string]func(http.ResponseWriter, *http.Request) // карта страница - обработчик
	surtedURLList []string
	}

// Это пустой класс, соответствующий интерфейсу myHandler
type myHandler struct { 
	server *CubieCtrlHttpServer
} 

type CellData struct {
	Name string
	Value float32
	Error bool	
	Timestamp time.Time
}


var results map[string]CellData
var settingsValues SettingsHolder

func (this *myHandler) findBestMatchHandler(url string) func(http.ResponseWriter, *http.Request) {
	    
    var result func(http.ResponseWriter, *http.Request)
    
    for _, template := range this.server.surtedURLList {
    	if strings.HasPrefix(url, template) {
    		result = this.server.mux[template]
    	}
    }
    
    return result
}

func GetPathsList(m *map[string]func(http.ResponseWriter, *http.Request)) []string {
	keys := make([]string, len(*m))
	i := 0
    for k,_ := range ( *m ) {
    	keys[i] = k
    	i++
    }
    sort.Strings(keys)
    
    return keys
}

// метод ServeHTTP
func (this *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	
	h := this.findBestMatchHandler(url)

	if h != nil {
		log.Printf("Processing url: %s -> %#x", url, h)
		h(w, r)
		return
	}
	
	// default ansver
	io.WriteString(w, "My server: "+url)
}

func NewServer(port int, hewData <-chan CellData, settings SettingsHolder) *CubieCtrlHttpServer {
	settingsValues = settings
	results = make(map[string]CellData)
	go func() {
		for r := range hewData {
			results[r.Name] = r
		}
	}()
	
	res := new(CubieCtrlHttpServer) 
	
	res.Addr = fmt.Sprintf(":%d", port)
	res.Handler = &myHandler{res} // типо конструктор
	
	res.mux = make(map[string]func(http.ResponseWriter, *http.Request))
	FillPagesMap(&res.mux)
	res.surtedURLList = GetPathsList(&res.mux)
	
	return res
}


