package worker

import (
	"log"
	"github.com/haxpax/goserial"
	"strings"
	"time"
	"errors"
)

const (
	SMSPending   = iota // 0
	SMSProcessed        // 1
	SMSError            // 2	
)
type SMSWorker struct{
	Phone string
	Body string
}

var messages chan SMSWorker

var avaliableModem []*GSMModem

type GSMModem struct {
	Port   string
	Baud   int
	Status bool
	Conn   serial.ReadWriteFlushCloser
	Devid  string
}

func (m *GSMModem) SendCommand(command string, waitForOk bool) string {
	log.Println("--- SendCommand: ", command)
	var status string = ""
	m.Conn.Flush()
	_, err := m.Conn.Write([]byte(command))
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 100)
	var loop int = 1
	if waitForOk {
		loop = 10
	}
	for i := 0; i < loop; i++ {		
		// ignoring error as EOF raises error on Linux
		n, _ := m.Conn.Read(buf)	
		if n > 0 {
			status = string(buf[:n])
			log.Printf("SendCommand: rcvd %d bytes: %s\n", n, status)
			if strings.HasSuffix(status, "OK\r\n") || strings.HasSuffix(status, "ERROR\r\n") {
				break
			}
		}
	}
	return status
}


func (m *GSMModem) IsModem() int {	
	status := m.SendCommand("AT\r", false)
	log.Println("--- Is Modem Status ==== ", status)
	if strings.HasSuffix(status, "OK\r\n") {
		return 1
	} else if strings.HasSuffix(status, "ERROR\r\n") {
		return 0
	} else{
		return -1
	}

}

func (m *GSMModem) Connect() error {	
	c := &serial.Config{Name: m.Port, Baud: m.Baud, ReadTimeout: time.Second}
	s, err := serial.OpenPort(c)

	if err == nil {
		log.Println("Open modem: ",m.Port)
		m.Status = true
		m.Conn = s	
		//only filter device is GSM Modem
		intIsModem := m.IsModem()
		if intIsModem==1{
			//log
			log.Println("Open sucessful modem: ",m.Port)
			return nil
		}	
		return errors.New("Device is not GSM modem")	
	}
	log.Println("Open modem err: ",err)
	return err
}



func (m *GSMModem) SendSMS(mobile string, message string) int {
	log.Println("--- SendSMS ", mobile, message, m.Port)

	// Put Modem in SMS Text Mode
	m.SendCommand("AT+CMGF=1\r", false)

	m.SendCommand("AT+CMGS=\""+mobile+"\"\r", false)

	// EOM CTRL-Z = 26
	status := m.SendCommand(message+string(26), true)
	//log.Println("--- SendSMS Status ==== ", status)
	if strings.HasSuffix(status, "OK\r\n") {
		return SMSProcessed
	} else if strings.HasSuffix(status, "ERROR\r\n") {
		return SMSError
	} else {
		return SMSPending
	}

}


func (m *GSMModem) ReadSMS() int {
	// Put Modem in SMS Text Mode
	m.SendCommand("AT+CMGF=1\r", false)
	//status := m.SendCommand("AT+CMGL=\"ALL\"\r", true)
	status := m.SendCommand("AT+CMGR=0\r", true)
	
	log.Println("--- ReadSMS Status ==== ", status)
	if strings.HasSuffix(status, "OK\r\n") {
		return SMSProcessed
	} else if strings.HasSuffix(status, "ERROR\r\n") {
		return SMSError
	} else {
		return SMSPending
	}

}

func (m *GSMModem) ProcessMessages() {
	defer func() {
		log.Println("--- deferring ProcessMessage")
		
	}()	
	for {
		//waitng message comming
		message := <-messages
		_ = m.SendSMS(message.Phone, message.Body)
	}
}

func InitWorker(modems []*GSMModem, bufferSize, bufferLow, loaderTimeout, countOut, loaderLongTimeout int) {
	messages = make(chan SMSWorker, bufferSize)
	for i := 0; i < len(modems); i++ {
		modem := modems[i]
		//log.Println("==== try to connect device",modem)
		err := modem.Connect()
		if err != nil {
			continue
		}
		//Add to avaliableModem
		avaliableModem=append(avaliableModem,modem)
		//run go routine in background to recieve sms		
		go modem.ProcessMessages()		
	}	
}

func GetAvaliableNameModem() []string{
	var modemNames []string
	for i:=0 ; i<len(avaliableModem); i++ {
		log.Println("avaliableModem::",avaliableModem[i])
		modemNames =append(modemNames, avaliableModem[i].Port)
	}
	return modemNames
}
func GetAvaliableModem() []*GSMModem{
	return avaliableModem
}

func AddSMS(phone string,body string ) {		
	//wakeup message loader and exit
	go func() {
		//notify the message loader only if its been to too long
		//or too many messages since last notification
		s:=SMSWorker{Phone:phone,Body:body}
		messages <- s
		log.Println("AddMessage :Add message to channel")
	}()
}	
