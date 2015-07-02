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
	No int	
}

func (sq *SmsQueue) Save() error{
	Ddb.Update(func(tx *bolt.Tx) error {
	    	log.Println("--- Save SmsQueue")
		    b := tx.Bucket([]byte("SmsQueue"))	
		    if(b==nil){
		    	log.Println("--- Create SmsQueue")
		    	tx.CreateBucket([]byte("SmsQueue"))
		    	b = tx.Bucket([]byte("SmsQueue"))	
		    }	    		    
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
	var itemFirst = SmsQueue{}
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


func (sq *SmsQueue) FindAll() []SmsQueue{
	var queues []SmsQueue
	chanQueues :=make(chan SmsQueue)	
	Ddb.View(func(tx *bolt.Tx) error {		
		b := tx.Bucket([]byte("SmsQueue"))	
		if b != nil{
			c := b.Cursor()		
		    _,vl := c.Last()
		   	xlast := SmsQueue{}
			json.Unmarshal(vl,&xlast)
			count:=0;
		    for k, v := c.First(); k != nil; k, v = c.Next() {	 
		    	count++
		        go func(v []byte){
		        	var s = SmsQueue{}   
		        	_ = json.Unmarshal(v,&s)	        	
					chanQueues<-s	
								
		        }(v)	    			    
		    }
		    for i:=0;i<count;i++ {
		    	select {
		    		case im:=<-chanQueues:
		    			queues=append(queues,im)		       
		    	}
		    }
		}
	  
	    return nil
	})	
	return queues
}

func (sms *SmsQueue) ResetAll() error{		
    Ddb.Update(func(tx *bolt.Tx) error {    	
	    err := tx.DeleteBucket([]byte("SmsQueue"))		    		   		
	    return err
	})
	return nil
}