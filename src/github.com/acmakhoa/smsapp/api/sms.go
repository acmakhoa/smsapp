package api

import (
	"fmt"  	
  	"log"
  	"net/http"	
  	"strings"
 	"time"
 	"os"
 	"html/template"
 	"encoding/csv"
 	"encoding/json"
  	"github.com/acmakhoa/smsapp/db"
  	"github.com/satori/go.uuid"
 	"github.com/gorilla/mux"
 	"github.com/acmakhoa/smsapp/worker"
 )

type SMSResponse struct {
	Status   int            `json:"status"`
	Message  string         `json:"message"`	
	Data []db.SMS    `json:"data"`
}

type SMSCreateRequest struct{
	Name string `json:"name"`
	Phone string `json:"phone"`
}

type SmsAPI struct{}

/* end dashboard handlers */

/* API handlers */

// push sms, allowed methods: POST
func (_ *SmsAPI) ImportCsvHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		r.ParseMultipartForm(1000000)
		body:=r.FormValue("body")
		//process to Create List
		name:=r.FormValue("listname")
		log.Println("name:====",name);
		id := uuid.NewV1()
		var t=&db.List{Id:id.String(),Name:name,Body:body,Created:time.Now()}
		t.Save()
		log.Println("log list created",t.Id)
		csvfile,_,err := r.FormFile("csv-file")
		//log.Println("param post:",f)	
		//var filePath="/home/khoa/Desktop/Go-lang/smsapp/assets/files/customer.csv"
		//csvfile, err := os.Open(filePath)

         if err != nil {
                 fmt.Println(err)
                 return
         }

         defer csvfile.Close()

         reader := csv.NewReader(csvfile)

         reader.FieldsPerRecord = -1 // see the Reader struct information below

         rawCSVdata, err := reader.ReadAll()

         if err != nil {
                 fmt.Println(err)
                 os.Exit(1)
         }

         // sanity check, display to standard output
         for _, each := range rawCSVdata {
                 //fmt.Printf("phone : %s and name : %s\n", each[0], each[1])
         	id = uuid.NewV1()
			newsms := &db.SMS{Id:"0"+each[0],Name:each[0],Created:time.Now()}			
			//save message
			newsms.Save()	

			//save sms in list
			db.SaveSmsInList(t.Id,newsms)
         }
     }

    log.Println("--- importCsv1111")	
	t1:=template.New("import_csv.html")
	t1.Delims("{{%", "%}}")
	t, _ := t1.ParseFiles("./templates/import_csv.html")
	t, _ =t.ParseFiles("./templates/header.html","./templates/footer.html")
	t.Execute(w, nil)

}

// push sms, allowed methods: POST
func (_ *SmsAPI) CreateSMSHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--- CreateSMSHandler")
	w.Header().Set("Content-type", "application/json")

	r.ParseForm()
	decoder := json.NewDecoder(r.Body)
	
	var t=SMSCreateRequest{}
	_ = decoder.Decode(&t)	
	mobiles  :=strings.Split(t.Phone, "\n");
	

	for i:=0;i<len(mobiles);i++{		
		newsms := &db.SMS{Id:mobiles[i],Name:t.Name,Created:time.Now()}
		log.Println("save new sms:",newsms)
		//save message
		newsms.Save()	
		//send sms 		
		worker.AddSMS(mobiles[i],t.Name)
	}

	Response(w,map[string]interface{}{})
}

// push sms, allowed methods: POST
func (_ *SmsAPI) FindSMSHandler(w http.ResponseWriter, r *http.Request) {	
	s := db.SMS{}
	smses:=s.FindAll()
	data:=map[string]interface{}{"sms":smses}
	Response(w,data)	
}

// delete sms, allowed methods: POST
func (_ *SmsAPI) DeleteSMSHandler(w http.ResponseWriter, r *http.Request) {	
	params := mux.Vars(r)
	
	sms := db.SMS{Id:params["id"]}	
	error:=sms.Delete()
	if error !=nil {		
		ErrorResponse(w,"error",1000)
	}else {		
		Response(w,map[string]interface{}{})
	}		
}

// delete sms, allowed methods: POST
func (_ *SmsAPI) ResetSMSHandler(w http.ResponseWriter, r *http.Request) {	
	sms := db.SMS{}	
	//reset all SMS
	error:=sms.ResetAll()
	if error !=nil {		
		ErrorResponse(w,"error",1000)
	}else {
		Response(w,map[string]interface{}{})
	}	
}





