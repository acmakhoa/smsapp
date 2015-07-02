package main

import (
	"fmt"	
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"github.com/gorilla/mux"
	"github.com/acmakhoa/smsapp/api"
)

//reposne structure to /sms
type SMSResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

//response structure to /smsdata/
type SMSDataResponse struct {
	Status   int            `json:"status"`
	Message  string         `json:"message"`
	Summary  []int          `json:"summary"`
	DayCount map[string]int `json:"daycount"`
	//Messages []SMS          `json:"messages"`
}

// Cache templates
//var templates = template.Must(template.ParseFiles("./templates/index.html"))

/* dashboard handlers */

// dashboard
func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--- indexHandler")
	// templates.ExecuteTemplate(w, "index.html", nil)
	// Use during development to avoid having to restart server
	// after every change in HTML
	t1:=template.New("index.html")
	t1.Delims("{{%", "%}}")
	t, _ := t1.ParseFiles("./templates/index.html")
	t, _ =t.ParseFiles("./templates/header.html","./templates/footer.html")
	t.Execute(w, nil)
}

func listHandler(w http.ResponseWriter, r *http.Request) {	
	params := mux.Vars(r)
	if params["id"]!=""{

		//Load detail template
		t2:=template.New("list_detail.html")
		t2.Delims("{{%", "%}}")
		t3, _ := t2.ParseFiles("./templates/list_detail.html")
		t3, _ =t3.ParseFiles("./templates/header.html","./templates/footer.html")
		data:=map[string]interface{}{"listId":params["id"]}
		t3.Execute(w, data)	
	}else{
		t1:=template.New("list.html")
		t1.Delims("{{%", "%}}")
		t, _ := t1.ParseFiles("./templates/list.html")
		t, _ =t.ParseFiles("./templates/header.html","./templates/footer.html")
		t.Execute(w, nil)	
	}
	
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--- sendHandler")
	t1:=template.New("send.html")
	t1.Delims("{{%", "%}}")
	t, _ := t1.ParseFiles("./templates/send.html")
	t, _ =t.ParseFiles("./templates/header.html","./templates/footer.html")
	t.Execute(w, nil)
}

func settingHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("--- settingHandler")
	t1:=template.New("setting.html")
	t1.Delims("{{%", "%}}")
	t, _ := t1.ParseFiles("./templates/setting.html")
	t, _ =t.ParseFiles("./templates/header.html","./templates/footer.html")
	t.Execute(w, nil)
}


// handle all static files based on specified path
// for now its /assets
func handleStatic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	static := vars["path"]
	http.ServeFile(w, r, filepath.Join("./assets", static))
}

func InitServer(host string, port string) error {
	log.Println("--- InitServer ", host, port)

	r := mux.NewRouter()
	r.StrictSlash(true)

	r.HandleFunc("/", listHandler)	
	r.HandleFunc("/sms", indexHandler)		
	r.HandleFunc(`/list/{id:[a-zA-Z0-9=\-\/\.\_]+}`, listHandler)	
	r.Methods("GET").Path("/send").HandlerFunc(sendHandler)	
	r.Methods("GET").Path("/setting").HandlerFunc(settingHandler)	
	

	// handle static files
	r.HandleFunc(`/assets/{path:[a-zA-Z0-9=\-\/\.\_]+}`, handleStatic)

	// all API handlers
	Api := r.PathPrefix("/api").Subrouter()	

	//API List SMS
	apilist:=api.ListSmsAPI{}
	Api.Methods("GET").Path("/list").HandlerFunc(apilist.FindAllHandler)
	Api.Methods("GET").Path(`/list/{id:[a-zA-Z0-9=\-\/\.\_]+}`).HandlerFunc(apilist.FindByIdHandler)
	Api.Methods("POST").Path(`/list/delete/{id:[a-zA-Z0-9=\-\/\.\_]+}`).HandlerFunc(apilist.DeleteHandler)
	Api.Methods("POST").Path(`/list/{id:[a-zA-Z0-9=\-\/\.\_]+}/delete/{sms_id:[a-zA-Z0-9=\-\/\.\_]+}`).HandlerFunc(apilist.DeleteSmsHandler)
	//send all sms in list
	Api.Methods("POST").Path(`/list/send/{id:[a-zA-Z0-9=\-\/\.\_]+}`).HandlerFunc(apilist.SendHandler)
	//send 1 sms in list
	Api.Methods("POST").Path(`/list/{id:[a-zA-Z0-9=\-\/\.\_]+}/send/{sms_id:[a-zA-Z0-9=\-\/\.\_]+}`).HandlerFunc(apilist.SendSMSHandler)

	//API SMS
	apisms:=api.SmsAPI{}
	r.HandleFunc("/importcsv",apisms.ImportCsvHandler)
	Api.Methods("POST").Path("/sms/create").HandlerFunc(apisms.CreateSMSHandler)
	Api.Methods("POST").Path(`/sms/delete/{id:[a-zA-Z0-9=\-\/\.\_]+}`).HandlerFunc(apisms.DeleteSMSHandler)
	Api.Methods("GET").Path("/sms/reset").HandlerFunc(apisms.ResetSMSHandler)
	Api.Methods("GET").Path("/sms/find").HandlerFunc(apisms.FindSMSHandler)

	apismsqueue:=  api.SmsQueueAPI{}
	Api.Methods("GET").Path("/smsqueue/find").HandlerFunc(apismsqueue.FindAllHandler)
	Api.Methods("GET").Path("/smsqueue/reset").HandlerFunc(apismsqueue.ResetQueueHandler)
	 //api.Methods("GET").Path("/sms/delete/{id:[0-9]+}").HandlerFunc(web.deleteSMSHandler)	
	// api.Methods("GET").Path("/sms/send").HandlerFunc(web.cronJobSMSHandler)
	//api.Methods("GET").Path("/sms/resend/{id:[0-9]+}").HandlerFunc(web.resendSMSHandler)
	//Api.Methods("POST").Path("/sms/").HandlerFunc(api.SendSMSHandler)

	apidevice :=api.DeviceAPI{}
	Api.Methods("GET").Path("/device/rescan").HandlerFunc(apidevice.ReScanHandler)
	Api.Methods("GET").Path("/device/list").HandlerFunc(apidevice.ListHandler)
	Api.Methods("POST").Path("/device/select").HandlerFunc(apidevice.SelectHandler)
	

	//Init crontab with message file
	//apismsqueue.InitCronTab(MessagesFile)

	http.Handle("/", r)


	bind := fmt.Sprintf("%s:%s", host, port)
	log.Println("listening on: ", bind)
	return http.ListenAndServe(bind, nil)

}
