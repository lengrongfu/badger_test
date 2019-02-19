package main

import (
	"github.com/dgraph-io/badger"
	"log"
	"time"
)

var (
	key  = []byte("answer")
	value = []byte("42")
	profix = []byte("a")
	newValue = []byte("100")
	ttl  = time.Duration(5)
)

func insert(txn *badger.Txn) error {
	err := txn.Set(key,value)
	return err
}

func selecT(txn *badger.Txn) error {
	item, err := txn.Get(key)
	if err != nil {
		return err
	}
	val, err := item.Value()
	if err != nil {
		return err
	}
	log.Println("val:",string(val))
	return nil
}

func iterate(txn *badger.Txn) error {
	opt := badger.DefaultIteratorOptions
	opt.PrefetchSize = 10
	iterator := txn.NewIterator(opt)
	defer iterator.Close()
	for iterator.Rewind(); iterator.Valid() ; iterator.Next() {
		item := iterator.Item()
		value, e := item.Value()
		if e != nil {
			return e
		}
		log.Printf("iterate:: item key is:%s,item values is:%s",string(item.Key()),string(value))
	}
	return nil
}

//scans
func scans(txn *badger.Txn) error {
	iterator := txn.NewIterator(badger.DefaultIteratorOptions)
	defer iterator.Close()
	for iterator.Seek(profix) ; iterator.ValidForPrefix(profix) ; iterator.Next() {
		item := iterator.Item()
		value, e := item.Value()
		if e != nil {
			return e
		}
		log.Printf("scans:: item key is:%s,item values is:%s",string(item.Key()),string(value))
	}
	return nil
}

//delete
func delete(txn *badger.Txn)  error {
	return txn.Delete(key)
}

//update
func update(txn *badger.Txn) error {
	return txn.Set(key,newValue)
}

//ttl update

func keySetTTL(txn *badger.Txn) error {
	return txn.SetWithTTL(key,newValue,ttl)
}

func main() {
	opts := badger.DefaultOptions
	opts.Dir = "./db"
	opts.ValueDir = "./db"
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer db.Close()
	//insert
	err = db.Update(insert)
	if err != nil {
		log.Fatal(err)
	}
	//select
	err = db.View(selecT)
	if err != nil {
		log.Fatal(err)
	}
	//iterate
	err = db.View(iterate)
	if err != nil {
		log.Fatal(err)
	}
	//scans
	err = db.View(scans)
	if err != nil {
		log.Fatal(err)
	}
	//delete
	err = db.Update(delete)
	if err != nil {
		log.Fatal(err)
	}
	//update
	err = db.Update(update)
	if err != nil {
		log.Fatal(err)
	}
	db.View(selecT)
	
	//ttl
	db.Update(keySetTTL)

}