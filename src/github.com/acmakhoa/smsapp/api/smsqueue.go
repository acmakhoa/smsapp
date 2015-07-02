package api

import(
	// "fmt"  	
   	"log"
   	"time"
   	"net/http"
   	"strings"	 
   	"strconv" 	 
//  	"encoding/json"
   	"github.com/acmakhoa/smsapp/db"
   	"github.com/acmakhoa/smsapp/worker"
   	"github.com/acmakhoa/smsapp/conf"
   	"github.com/robfig/cron"

 //  	"github.com/satori/go.uuid"
 // 	"github.com/gorilla/mux"
  	
)

type SmsQueueAPI struct{}

func (_ *SmsQueueAPI) Push(list db.List,sms db.SMS, no int) error{
	model := db.SmsQueue{Id:list.Id+"_"+ sms.Id,Content:list.Body,Status:0,Retries:0,Created:time.Now(),No:no}
	model.Save();
	log.Println("log mode:",model)
	return nil
}

func (_ *SmsQueueAPI) Pop() (db.SmsQueue, error){
	model := db.SmsQueue{}
	smsQueue,_ := model.FindOne();
	log.Println("Pop message before:",smsQueue)
	smsQueue.Status=1;
	smsQueue.Save()
	return smsQueue,nil	
}


//Init cron tab to send SmsQueue
func (smsQueue *SmsQueueAPI) InitCronTab(messagesFile string){
	c := cron.New()		
	c.AddFunc(""+strconv.Itoa(1)+" * * * * *", func() {
			log.Println("contab start") 
			sms,_ := smsQueue.Pop()
			if(sms.Id != ""){
				Ids:=strings.Split(sms.Id,"_")
				log.Println(Ids[0],Ids[1])				
				//get roudrobin message
				strMessage := conf.GetRandomMessage(messagesFile,sms.No)
				log.Println("get random message from file:",strMessage)
				//worker.AddSMS(Ids[1],sms.Content)
				worker.AddSMS(Ids[1],strMessage)
			}
			log.Println("contab finish") 
		})
	c.Start()				
}

func (smsQueue *SmsQueueAPI) FindHandler(w http.ResponseWriter, r *http.Request){
	sms,_ := smsQueue.Pop()
	if(sms.Id != ""){
		Ids:=strings.Split(sms.Id,"_")
		log.Println(Ids[0],Ids[1])
	}
	
	//worker.AddSMS(Ids[1],sms.Content)	
	Response(w,map[string]interface{}{})
}

//Find all queue SMS
func (_ *SmsQueueAPI) FindAllHandler(w http.ResponseWriter, r *http.Request){
	model := db.SmsQueue{}	
	queues := model.FindAll()	
	data:=map[string]interface{}{"queues":queues}
	Response(w,data)
}


func (smsQueue * SmsQueueAPI) ResetQueueHandler(w http.ResponseWriter, r *http.Request){
	model := db.SmsQueue{}	
	err := model.ResetAll()	
	if err !=nil{
		ErrorResponse(w,"error",1000)
	}else{
		Response(w,map[string]interface{}{})
	}

}