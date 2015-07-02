// +build windows

package device

import(
	//"io/ioutil"
	//"strings"
	"strconv"
	"log"	
)


func FindDevices() []string {
	var devices []string
	//Scan port from 0 - 99
	var port string
	for i:=0;i<100;i++ {
		port = "COM"+strconv.Itoa(i)
		if CanOpenPort(port) {
			
			devices = append(devices, port)
		}			
	}
	log.Println("COM Found:",devices)
	return devices
}