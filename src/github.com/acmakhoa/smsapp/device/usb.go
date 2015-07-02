package device

import(
	"time"
	"github.com/haxpax/goserial"		
)

func CanOpenPort(port string) bool{		
	c := &serial.Config{Name: port, Baud: 115200, ReadTimeout: time.Second}
	s, err := serial.OpenPort(c)
	if err == nil {
		s.Close()
		return true
	}		
	return false	
}


