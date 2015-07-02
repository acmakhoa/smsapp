package api

import (
	"fmt"
  	"log"
  	//"time"
  	"net/http"	
  	//"strconv"
  	"github.com/gorilla/mux"
  	//"github.com/robfig/cron"
  	"github.com/acmakhoa/smsapp/db"
  	"github.com/acmakhoa/smsapp/worker"
 )

type ListSmsAPI struct{}

func (_ *ListSmsAPI) FindAllHandler(w http.ResponseWriter, r *http.Request){
	l := db.List{}
	lists:=l.FindAll()
	data:=map[string]interface{}{"lists":lists}
	Response(w,data)
}

func (_ *ListSmsAPI) FindByIdHandler(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	fmt.Println("id:",params["id"])
	l := db.List{}
	list:=l.FindById(params["id"])
	s :=db.SMS{}
	smses :=s.FindByList(params["id"])
	data:=map[string]interface{}{"list":list,"sms":smses}
	Response(w,data)
}

// delete sms, allowed methods: POST
func (_ *ListSmsAPI) DeleteHandler(w http.ResponseWriter, r *http.Request) {	
	params := mux.Vars(r)
	
	list := db.List{Id:params["id"]}	
	error:=list.Delete()
	if error !=nil {		
		ErrorResponse(w,"error",1000)
	}else {
		Response(w,map[string]interface{}{})
	}		
}
// delete sms, allowed methods: POST
func (_ *ListSmsAPI) DeleteSmsHandler(w http.ResponseWriter, r *http.Request) {	
	params := mux.Vars(r)
	log.Println("params;;;;;;:",params)
	list := db.List{Id:params["id"]}	
	error:=list.DeleteSms(params["sms_id"])
	if error !=nil {		
		ErrorResponse(w,"error",1000)
	}else {
		Response(w,map[string]interface{}{})
	}		
}

// send a sms in ListSMS
func (_ *ListSmsAPI) SendSMSHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)	
	model:=db.List{}
	l:=model.FindById(params["id"])
	log.Println("send message in list:",l)
	//Find sms

	worker.AddSMS(params["sms_id"],l.Body)
	log.Println("send message to phone:",params["sms_id"])
	Response(w,map[string]interface{}{})

}

// send a sms in ListSMS
func (_ *ListSmsAPI) SendHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	model:=db.List{}
	l:=model.FindById(params["id"])
	
	modelSms:=db.SMS{}	
	smses:=modelSms.FindByList(params["id"])
	queue:=SmsQueueAPI{}
	for i:=0;i<len(smses);i++{
		//log.Println("send message to phone:",smses[i].Id,l.Body)
		go queue.Push(l,smses[i],i)
	}
	Response(w,map[string]interface{}{})
}