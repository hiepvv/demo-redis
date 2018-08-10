package models

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	mgo "gopkg.in/mgo.v2"
)

const (
	mongoDBHosts = "ds119052.mlab.com:19052"
	authDatabase = "demo-redis"
	authUserName = "groot"
	authPassword = "2XGFXShup89fwNHnQsxnzs4GP78tm"
)

// Session declaration
var session *mgo.Session
var usersCollection *mgo.Collection
var loginsCollection *mgo.Collection

// Client ...
var Client *redis.Client

// InitialDBSession will connect to DB Server
func InitialDBSession() {
	var err error
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{mongoDBHosts},
		Timeout:  60 * time.Second,
		Database: authDatabase,
		Username: authUserName,
		Password: authPassword,
	}
	session, err = mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	usersCollection = session.DB("demo-redis").C("users")
	loginsCollection = session.DB("cafeteria").C("logins")

	Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := Client.Ping().Result()
	fmt.Println(pong, err)
}

// CloseDBSession will turn off session connect to DB
func CloseDBSession() {
	session.Close()
}
