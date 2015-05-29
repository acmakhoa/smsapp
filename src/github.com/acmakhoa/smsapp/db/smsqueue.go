package db

import (
	"github.com/boltdb/bolt"	
	"log"	
	"time"
	"encoding/json"
)

type SmsQueue struct{
	Id string
	Content string
	Status int
	Retries int
	Created time.Time	
}

func (sq *SmsQueue) Save() error{
	Ddb.Update(func(tx *bolt.Tx) error {
	    	log.Println("--- Save SmsQueue")
		    b := tx.Bucket([]byte("SmsQueue"))		    		    
		    encoded, err := json.Marshal(sq)
		    log.Println("--- Save SmsQueue content",sq)
		    if err != nil {
		        return err
		    }		    
			_ = b.Put([]byte(sq.Id), encoded)
		   
		    return err
		})
	return nil
}


func (sq * SmsQueue) FindOne() (SmsQueue,error){
	itemFirst := SmsQueue{}
	Ddb.View(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte("SmsQueue"))
		if b!=nil {
			c :=b.Cursor()
			for k,v := c.First();k!=nil; k,v=c.Next(){
				var x =SmsQueue{}
				_ = json.Unmarshal(v,&x)	
				if x.Status==0{
					itemFirst=x
					return nil
				}
			}						
		}		
		return nil	
	})
	return itemFirst,nil
}