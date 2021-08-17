package mongo

import (
	"fmt"
	"github.com/xiaop0817/ftgoutils/c"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

var Session *mgo.Session
var prefix = c.C(fmt.Sprintf("%-10s", "[Mongo]"), c.Yellow)

func InitMongo(host string, db string, usrName string, pwd string) {
	dialInfo := &mgo.DialInfo{
		Addrs:     []string{host},
		Timeout:   60 * time.Second,
		Source:    db,
		Username:  usrName,
		Password:  pwd,
		PoolLimit: 4096,
	}

	defer func() {
		if err := recover(); err != nil {
			log.Printf("%s Mongo连接失败[%s]", prefix, c.C(err, c.LightRed))
			time.Sleep(time.Second * 4)
		}
	}()

	s, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Printf("%s Mongo创建连接失败[%s]", prefix, c.C(err, c.LightRed))
		return
	}
	log.Printf("%s %s", prefix, c.C("Mongo创建连接完成!", c.LightGreen))
	Session = s
}

// 获取文档对象
func connect(db string, collection string) (*mgo.Session, *mgo.Collection) {
	if Session != nil {
		s := Session.Copy()
		c := s.DB(db).C(collection)
		s.SetMode(mgo.Monotonic, true)
		return s, c
	}
	return nil, nil
}

// Insert 插入数据
// db string 操作的数据库
// collection string 操作的文档(表)
// docs 数据
func Insert(db, collection string, docs ...interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Insert(docs...)
}

// FindOne 查询单条数据
// db string 操作的数据库
// collection string 操作的文档(表)
// query:查询条件
// selector:需要过滤的数据(projection)
// result:查询到的结果
func FindOne(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).One(result)
}

// FindAll 查询数据列表
// db string 操作的数据库
// collection string 操作的文档(表)
// query:查询条件
// selector:需要过滤的数据(projection)
// result:查询到的结果
func FindAll(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Find(query).Select(selector).All(result)
}

// FindAllSort 查询数据列表
// db string 操作的数据库
// collection string 操作的文档(表)
// query:查询条件
// selector:需要过滤的数据(projection)
// result:查询到的结果
// sort:排序 exam:(-create_time,+index)
func FindAllSort(db, collection string, query, selector, result interface{}, sort ...string) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Find(query).Sort(sort...).Select(selector).All(result)
}

// UpdateOne 更新一条数据
// db string 操作的数据库
// collection string 操作的文档(表)
// query:要更新的文档id
// docs:需要更新的文档数据
// result:查询到的结果
func UpdateOne(db, collection string, id bson.ObjectId, docs interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	//ObjectIdHex
	return c.Update(bson.M{"_id": id}, docs)
}

func Update(db, collection string, selector, update interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Update(selector, update)
}

func Count(db, collection string, query interface{}) (int, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Count()
}

func GroupBy(db, collection string, query interface{}) (int, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Count()
}

// BulkInsert 批量插入记录
// db string 操作的数据库
// collection string 操作的文档(表)
// docs:需要写入的记录列表
func BulkInsert(db, collection string, docs ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.Insert(docs...)
	return bulk.Run()
}

// RemoveAll 删除符合条件的全部记录
// db string 操作的数据库
// collection string 操作的文档(表)
// selector 记录条件
func RemoveAll(db, collection string, selector interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	_, err := c.RemoveAll(selector)
	return err
}

func BulkRemove(db, collection string, selector ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()

	bulk := c.Bulk()
	bulk.Remove(selector...)
	return bulk.Run()
}

// Upsert 插入或更新
// db string 操作的数据库
// collection string 操作的文档(表)
// update 需要写入的记录
func Upsert(db, collection string, selector, update interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	_, err := c.Upsert(selector, update)
	return err
}

func GetAndUpdate(db, collection string) int {
	ms, c := connect(db, collection)
	defer ms.Close()

	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"seq": 1}},
		Upsert:    true,
		ReturnNew: true,
	}
	doc := struct{ Seq int }{}
	c.Find(bson.M{"_id": "order_seq_" + time.Now().Format("20060102")}).Apply(change, &doc)
	return doc.Seq
}
