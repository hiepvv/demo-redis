package models

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"errors"

	"gopkg.in/mgo.v2/bson"
)

// User model declaration
type User struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Email    string        `json:"email" bson:"email"`
	Password string        `json:"-" bson:"password"`
	Salt     string        `json:"-" bson:"salt"`
}

// GenerateRandomString returns a URL-safe, base64 encoded
func GenerateRandomString() (string, error) {
	length := 128

	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), err
}

// GenerateHMACCrypto returns a Std-safe, base64 encoded
func GenerateHMACCrypto(salt string, data ...string) (cryptoString string) {
	crypto := hmac.New(sha512.New, []byte(salt))
	decodeData := salt
	for _, value := range data {
		decodeData += value
	}
	crypto.Write([]byte(decodeData))
	cryptoString = base64.StdEncoding.EncodeToString(crypto.Sum(nil))
	return
}

// FindUser will find user data by params in DB
func FindUser(params map[string]interface{}) (user User, err error) {
	fields := []bson.M{}
	for key, value := range params {
		fields = append(fields, bson.M{key: value})
	}
	err = usersCollection.Find(bson.M{"$and": fields}).One(&user)
	return
}

// AddNewUser will handle insert user data into User Collection
func AddNewUser(email string) (newUser User, err error) {

	salt, errGenerateSalt := GenerateRandomString()
	if errGenerateSalt != nil {
		return User{}, errors.New("generate salt was failed")
	}
	newUser = User{}
	newUser.ID = bson.NewObjectId()
	newUser.Email = email
	newUser.Password = GenerateHMACCrypto(salt, "123456")
	newUser.Salt = salt

	err = usersCollection.Insert(&newUser)
	return
}

// VerifyUser will verify user data and return a token if data is valid.
func VerifyUser(email string, password string) (user User, err error) {
	user, err = FindUser(map[string]interface{}{"email": email, "active": true})
	if err != nil {
		return User{}, errors.New("user is not exist")
	}
	if GenerateHMACCrypto(user.Salt, password) != user.Password {
		return User{}, errors.New("user is not exist")
	}
	return user, nil
}
