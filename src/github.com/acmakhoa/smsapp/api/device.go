package api

import(
	"log"
	"net/http"
	"encoding/json"
	//"github.com/gorilla/mux"
	"github.com/acmakhoa/smsapp/worker"
	"github.com/acmakhoa/smsapp/device"
)

type DeviceAPI struct{}
type DeviceSelectRequest struct{
	Name string `json:"name"`
}

func (_ *DeviceAPI) ListHandler(w http.ResponseWriter, r *http.Request){
	listDevices := worker.GetAvaliableNameModem()
	data := map[string]interface{}{"devices":listDevices}
	Response(w,data)
}
func (_ *DeviceAPI) ReScanHandler(w http.ResponseWriter, r *http.Request){
	//listDevices := worker.GetAvaliableModem()
	var modems []*worker.GSMModem
	//Find Usb Device
    devices:=device.FindDevices()
    //log.Println("list device",jj)

	for i := 0; i < len(devices); i++ {
		//dev := fmt.Sprintf("DEVICE%v", i)
		_port := devices[i]
		_baud := 115200 //appConfig.Get(dev, "BAUDRATE")
		//_devid, _ := appConfig.Get(dev, "DEVID")
		m := &worker.GSMModem{Port: _port, Baud: _baud, Devid: ""}
		//log.Println("gsm modem :",m)
		err := m.Connect()
		if err!=nil{
			log.Println("rescan: connet error :", err)
			continue
		}
		modems = append(modems, m)
	}
	data:=map[string]interface{}{"modems":modems}
	Response(w,data)
}
func (_ *DeviceAPI) SelectHandler(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-type", "application/json")
	
	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	
	t:=DeviceSelectRequest{}
	_ = decoder.Decode(&t)	
	log.Println("decoder:",t)
	modemSelected:=t.Name
	if modemSelected!=""{
		log.Println("Select Handler ::",modemSelected)
		modems:=worker.GetAvaliableModem()
		for i:=0;i<len(modems);i++{
			log.Println("modem::",modems[i])
			if modems[i].Port==modemSelected{
				var aModem []*worker.GSMModem
				aModem =append(aModem,modems[i])
				worker.InitWorker(aModem,10,1,1,1,1)
			}
		}	
	}
	Response(w,map[string]interface{}{})
	
}
