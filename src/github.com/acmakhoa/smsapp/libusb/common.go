package libusb

import "fmt"
import "os"


func Printlist(){
	xyz,zz:=Init()
	fmt.Println("xyz:::::",xyz)
	fmt.Println("xyz:::::",zz)


	for i, info := range Enum() {
		fmt.Printf("======================================================\n")
		fmt.Printf(" %10d : BUS:%s DEVICE:%s VID:%04x PID:%04x\n", i, info.Bus, info.Device, info.Vid, info.Pid)
		dev := Open(info.Vid, info.Pid)
		if dev != nil {
			fmt.Printf(" Vendor     : %s\n", dev.Vendor())
			fmt.Printf(" Product    : %s\n", dev.Product())
			fmt.Printf(" Last Error : %s\n", dev.LastError())
			dev.Close()
		} else {
			os.Exit(1)
		}
	}
}
