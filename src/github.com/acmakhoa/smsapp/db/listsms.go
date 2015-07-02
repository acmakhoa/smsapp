package db

import (
"github.com/boltdb/bolt"	
"log"	
"time"
"encoding/json"
)

type List struct{
	Id string
	Name string
	Body string
	Created time.Time
}

func (list *List) Save() error{		
	Ddb.Update(func(tx *bolt.Tx) error {
		log.Println("--- Save List")
		b := tx.Bucket([]byte("List"))		    		    
		encoded, err := json.Marshal(list)
		log.Println("--- Save List content",list)
		if err != nil {
			return err
		}
		log.Println("put key uuid,",list.Id)
		_ = b.Put([]byte(list.Id), encoded)

		return err
		})
	return nil
}

func (list *List)FindById(Id string) List{		
	 	var model List
	    Ddb.View(func(tx *bolt.Tx) error {
		    b := tx.Bucket([]byte("List"))		  
		    err := json.Unmarshal(b.Get([]byte(Id)),&model)	  		    
		    if err != nil {			   
			    return err
			}
		    return nil
		})
		return model
}

func (list *List)FindAll() []List{	
 	var listSMS []List
 	chanLists := make(chan List)
    Ddb.View(func(tx *bolt.Tx) error {
	    b := tx.Bucket([]byte("List"))	
	    c := b.Cursor()		
	    _,vl := c.Last()
	   	xlast := List{}
		json.Unmarshal(vl,&xlast)
		count:=0;
	    for k, v := c.First(); k != nil; k, v = c.Next() {	 
	    	count++
	        go func(v []byte){
	        	var s = List{}   
	        	_ = json.Unmarshal(v,&s)	        	
				chanLists<-s	
							
	        }(v)	    			    
	    }
	    for i:=0;i<count;i++ {
	    	select {
	    		case im:=<-chanLists:
	    			listSMS=append(listSMS,im)		       
	    	}
	    }
	    return nil
	})
	return listSMS
}

func (list *List) Delete() error{		
    Ddb.Update(func(tx *bolt.Tx) error {
    	log.Println("--- Delete SMS")
	    b := tx.Bucket([]byte("List"))		    
	    //Delete item in List		   
		err := b.Delete([]byte(list.Id))

		//Delete ListSMS Bucket
		smsb := tx.Bucket([]byte("ListSMS-"+list.Id))
		if(smsb!=nil){
			_ = smsb.DeleteBucket([]byte("ListSMS-"+list.Id))
		}
	    return err
	})
	return nil
}
/*
** Delete sms in ListSMS
*/
func (list *List) DeleteSms(smsId string) error{		
    Ddb.Update(func(tx *bolt.Tx) error {
    	log.Println("--- Delete SMS")
	    b := tx.Bucket([]byte("ListSMS-"+list.Id))		    
	    err := b.Delete([]byte(smsId))		
	    return err
	})
	return nil
}


func SaveSmsInList(listId string, sms *SMS) error{		
	Ddb.Update(func(tx *bolt.Tx) error {
		log.Println("--- SaveSmsInList:","ListSMS-"+listId)
		b := tx.Bucket([]byte("ListSMS-"+listId))		    		    	    	   
		if(b==nil){
			_, err3 := tx.CreateBucket([]byte("ListSMS-"+listId))		
			if err3!=nil{
				return err3
			}
			b = tx.Bucket([]byte("ListSMS-"+listId))		    
		}
		//encode
		encoded, err := json.Marshal(sms)
		log.Println("--- Save List content",encoded)
		if err != nil {
			return err
		}	    
		log.Println("--- Save sms:", sms.Id)
		_ = b.Put([]byte(sms.Id), encoded)

		return err
		})
	return nil
}


