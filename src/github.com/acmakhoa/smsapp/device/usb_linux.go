package device

import(
	"io/ioutil"
	"strings"
	"log"	
)


func FindDevices() []string {
	var devices []string
	contents, _ := ioutil.ReadDir("/dev")

	// Look for what is mostly likely the Arduino device
	var port string
	for _, f := range contents {		

		if strings.Contains(f.Name(), "tty.usbserial") ||
			strings.Contains(f.Name(), "ttyUSB") {
			
			port ="/dev/" + f.Name()
			//Check can open port
			if CanOpenPort(port){
				
				devices = append(devices, port)
			}			
		}
	}
	log.Println("USB Found::::", devices)
	return devices
}