package db

import (
	"github.com/boltdb/bolt"	
	"log"	
	"time"
	"encoding/json"
)

type SMS struct{
	Id string
	Name string
	Status int
	Retries int
	Created time.Time	
}

func (sms *SMS) Save() error{		
	    Ddb.Update(func(tx *bolt.Tx) error {
	    	log.Println("--- Save SMS")
		    b := tx.Bucket([]byte("SMS"))		    		    
		    encoded, err := json.Marshal(sms)
		    log.Println("--- Save SMS content",sms)
		    if err != nil {
		        return err
		    }		    
			_ = b.Put([]byte(sms.Id), encoded)
		   
		    return err
		})
		return nil
}

func (sms *SMS)FindById(Id string) SMS{		
	 	var s SMS
	    Ddb.View(func(tx *bolt.Tx) error {
		    b := tx.Bucket([]byte("SMS"))		  
		    err := json.Unmarshal(b.Get([]byte(Id)),&s)	  		    
		    if err != nil {			   
			    return err
			}
		    return nil
		})
		return s
}

func (sms *SMS)FindByList(listId string) []SMS{	

 	var listSMS []SMS;
 	messages := make(chan SMS)
 	//done := make(chan bool)
    Ddb.View(func(tx *bolt.Tx) error {
    	log.Println("find in list bucket:","ListSMS-"+listId)
	    b := tx.Bucket([]byte("ListSMS-"+listId))		    
	    if(b!=nil){
	    	c := b.Cursor()		
    	
		    _,vl := c.Last()
		   	xlast := SMS{}
			json.Unmarshal(vl,&xlast)
			count:=0;
		    for k, v := c.First(); k != nil; k, v = c.Next() {	 
		    	count++
		    	//log.Println("Done:",count)  	
		        go func(v []byte){
		        	var s = SMS{}   
		        	_ = json.Unmarshal(v,&s)		          
					messages<-s	
								
		        }(v)	    			    
		    }
		    for i:=0;i<count;i++ {
		    	select {
		    		case im:=<-messages:
		    			listSMS=append(listSMS,im)		       
		    	}
		    }
	    }
	    
	    return nil
	})
	return listSMS
}

func (sms *SMS)FindAll() []SMS{	

 	var listSMS []SMS;
 	messages := make(chan SMS)
 	//done := make(chan bool)
    Ddb.View(func(tx *bolt.Tx) error {
	    b := tx.Bucket([]byte("SMS"))	
	    // Iterate over items in sorted key order.
	    /*
	    b.ForEach(func(k, v []byte) error {
	    	var s = SMS{}    	
	    	err := json.Unmarshal(b.Get([]byte(k)),&s)	

		    if err != nil {				   
			    return err
			}	  	
			
			//messages<-s
	        listSMS=append(listSMS,s)
	        return nil
	    })
	    */
	    c := b.Cursor()		
    	
	    _,vl := c.Last()
	   	xlast := SMS{}
		json.Unmarshal(vl,&xlast)
		count:=0;
	    for k, v := c.First(); k != nil; k, v = c.Next() {	 
	    	count++
	    	//log.Println("Done:",count)  	
	        go func(v []byte){
	        	var s = SMS{}   
	        	_ = json.Unmarshal(v,&s)	        	
				messages<-s	
							
	        }(v)	    			    
	    }
	    for i:=0;i<count;i++ {
	    	select {
	    		case im:=<-messages:
	    			listSMS=append(listSMS,im)		       
	    	}
	    }
	    return nil
	})
	return listSMS
}

func (sms *SMS) Delete() error{		
    Ddb.Update(func(tx *bolt.Tx) error {
    	log.Println("--- Delete SMS")
	    b := tx.Bucket([]byte("SMS"))		    		   
		err := b.Delete([]byte(sms.Id))
	    return err
	})
	return nil
}
func (sms *SMS) ResetAll() error{		
    Ddb.Update(func(tx *bolt.Tx) error {    	
	    err := tx.DeleteBucket([]byte("SMS"))		    		   		
	    return err
	})
	return nil
}

