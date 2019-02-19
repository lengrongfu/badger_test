# Badger 使用demo

> Badger是基于最新`LSM`的设计的`K/V`数据库，在`SSD`上有更好的性能，相对于`leveldb`等这类的`K/V`数据库。

## 打开Badger
```go
    opts := badger.DefaultOptions
	opts.Dir = "./db" //存储data
	opts.ValueDir = "./db" //存储value log
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
```

## 只读事务

```go
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
//在view操作中只能进行查看操作
err = db.View(selecT)

```

## 读写事务
```go
//insert
func insert(txn *badger.Txn) error {
	err := txn.Set(key,value)
	return err
}
//delete
func delete(txn *badger.Txn)  error {
	return txn.Delete(key)
}

//update
func update(txn *badger.Txn) error {
	return txn.Set(key,newValue)
}
// 在update中可进行读写操作
db.Update(insert)
db.Update(delete)
db.Update(update)
```

## 参考
- https://colobu.com/2017/10/11/badger-a-performant-k-v-store/
- https://github.com/dgraph-io/badger