package db

import (
	"github.com/boltdb/bolt"
	"log"		
)
var Ddb *bolt.DB

func InitDB( dbname string) (*bolt.DB, error) {
	db, err := bolt.Open(dbname, 0600, nil)
    if err != nil {
        log.Fatal(err)
    }    
    Ddb = db
    Ddb.Update(func(tx *bolt.Tx) error {
	    	log.Println("--- Create Bucket SMS")
		    b := tx.Bucket([]byte("SMS"))		
		    if b==nil{
		    	//create bucket `SMS`
		    	_, err1 := tx.CreateBucket([]byte("SMS"))		
		    	if err1!=nil{
		    		return err1
		    	}
		    }   
		    //create List
		    bList := tx.Bucket([]byte("List"))		
		    if bList==nil{
		    	//create bucket `List`
		    	_, err2 := tx.CreateBucket([]byte("List"))		
		    	if err2!=nil{
		    		return err2
		    	}
		    }		   
		    return nil		    
		})
    return Ddb, nil;
}
