package conf

import(
	"log"
	"os"
	"encoding/json"
)

type SmsConfig struct{
	Messages []string	
}

func GetSmsConfig(jsonfile string) []string{
	file,_ :=os.Open(jsonfile)
	decoder := json.NewDecoder(file)
	configuration := SmsConfig{}
	err :=decoder.Decode(&configuration)
	if err !=nil {
		log.Println("error:",err)		
	}
	return configuration.Messages	
}

func GetRandomMessage(jsonfile string,index int) string{
	listMessages := GetSmsConfig(jsonfile)
	lenListMessages:=len(listMessages)
	if lenListMessages>0{
		number := (index%lenListMessages)
		return listMessages[number]
	}
	return ""
}
