package main

import (
	"fmt"
	"log"	
	"os"	
	"github.com/acmakhoa/smsapp/db"
	"github.com/acmakhoa/smsapp/conf"
	"github.com/acmakhoa/smsapp/worker"
	//"github.com/acmakhoa/smsapp/libusb"
	"strconv"	
	"runtime"	
)

func main() {
	runtime.GOMAXPROCS(4)
	appConfig, err := conf.GetConfig("conf.ini")
	
	if err != nil {
		log.Println("main: ", "Invalid config: ", err.Error(), " Aborting")
		os.Exit(1)
	}

	//libusb.Printlist()

	 dbname, _ := appConfig.Get("SETTINGS", "DBNAME")
	 _, err = db.InitDB(dbname)
     if err != nil {
         log.Fatal(err)
     }
     defer db.Ddb.Close()

     //Init server
	_numDevices, _ := appConfig.Get("SETTINGS", "DEVICES")
	numDevices, _ := strconv.Atoi(_numDevices)
	log.Println("main: number of devices: ", numDevices)

	var modems []*worker.GSMModem
	for i := 0; i < numDevices; i++ {
		dev := fmt.Sprintf("DEVICE%v", i)
		_port, _ := appConfig.Get(dev, "COMPORT")
		_baud := 115200 //appConfig.Get(dev, "BAUDRATE")
		_devid, _ := appConfig.Get(dev, "DEVID")
		m := &worker.GSMModem{Port: _port, Baud: _baud, Devid: _devid}
		modems = append(modems, m)
	}

	_bufferSize, _ := appConfig.Get("SETTINGS", "BUFFERSIZE")
	bufferSize, _ := strconv.Atoi(_bufferSize)

	_bufferLow, _ := appConfig.Get("SETTINGS", "BUFFERLOW")
	bufferLow, _ := strconv.Atoi(_bufferLow)

	_loaderTimeout, _ := appConfig.Get("SETTINGS", "MSGTIMEOUT")
	loaderTimeout, _ := strconv.Atoi(_loaderTimeout)

	_loaderCountout, _ := appConfig.Get("SETTINGS", "MSGCOUNTOUT")
	loaderCountout, _ := strconv.Atoi(_loaderCountout)

	_loaderTimeoutLong, _ := appConfig.Get("SETTINGS", "MSGTIMEOUTLONG")
	loaderTimeoutLong, _ := strconv.Atoi(_loaderTimeoutLong)

	log.Println("main: Initializing worker")
	worker.InitWorker(modems, bufferSize, bufferLow, loaderTimeout, loaderCountout, loaderTimeoutLong)
	

	//Init server
	log.Println("main: Initializing server")
	serverhost, _ := appConfig.Get("SETTINGS", "SERVERHOST")
	serverport, _ := appConfig.Get("SETTINGS", "SERVERPORT")
	err = InitServer(serverhost, serverport)
	if err != nil {
		log.Println("main: ", "Error starting server: ", err.Error(), " Aborting")
		os.Exit(1)
	}	
	
}
