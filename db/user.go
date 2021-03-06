package db

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

// Stats - user statistics
type Stats struct {
	Kills  int `bson:"kills" json:"kills"`
	Fucks  int `bson:"fucks" json:"fucks"`
	Marrys int `bson:"marrys" json:"marrys"`
}

// User - basically stores user information
type User struct {
	VKID     int    `bson:"vkid" json:"vkid"` // used to identify users in database
	Name     string `bson:"name" json:"name"`
	Sex      bool   `bson:"sex" json:"sex"` // 0 - female, 1 - male
	PhotoURL string `bson:"photo_url" json:"photo_url"`
	Stats    Stats  `bson:"stats" json:"stats"`
}

type notEnoughPicsError struct {
	sliceLen int
	yourLen  int
}

func (e *notEnoughPicsError) Error() string {
	return fmt.Sprintf("%d bigger than %d", e.yourLen, e.sliceLen)
}

// AddUser adds user to database
func (dbw *Wrapper) AddUser(user User) (err error) {
	if n, _ := dbw.Users.Find(bson.M{"vkid": user.VKID}).Count(); n == 0 {
		err = dbw.Users.Insert(user)
	}
	return err
}

// DeleteUser deletes a user with given vkid
func (dbw *Wrapper) DeleteUser(vkid int) error {
	return dbw.Users.Remove(bson.M{"vkid": vkid})
}

// GetRandomUsers samples {size} users with given sex from database
func (dbw *Wrapper) GetRandomUsers(size int, sex bool) (users []User, err error) {
	pipe := dbw.Users.Pipe([]bson.M{{"$match": bson.M{"sex": sex}}, {"$sample": bson.M{"size": size}}})
	err = pipe.All(&users)
	fmt.Println()

	return users, err
}

func (dbw *Wrapper) UpdateUserStats(vkids []int) {
	dbw.Users.Update(bson.M{"vkid": vkids[0]}, bson.M{"$inc": bson.M{"stats.fucks": 1}})
	dbw.Users.Update(bson.M{"vkid": vkids[1]}, bson.M{"$inc": bson.M{"stats.marrys": 1}})
	dbw.Users.Update(bson.M{"vkid": vkids[2]}, bson.M{"$inc": bson.M{"stats.kills": 1}})
}

func (dbw *Wrapper) GetUserByName(name string) (user User, err error) {
	err = dbw.Users.Find(bson.M{"name": name}).One(&user)
	return user, err
}

func (dbw *Wrapper) GetUserByVKID(vkid int) (user User, err error) {
	err = dbw.Users.Find(bson.M{"vkid": vkid}).One(&user)
	return
}

func (dbw *Wrapper) UpdateUserInfo(vkid int, name string, sex bool, photo_url string) {
	err := dbw.Users.Update(bson.M{"vkid": vkid}, bson.M{"$set": bson.M{"name": name}})
	if err != nil {
		panic(err)
	}
	err = dbw.Users.Update(bson.M{"vkid": vkid}, bson.M{"$set": bson.M{"sex": sex}})
	if err != nil {
		panic(err)
	}
	err = dbw.Users.Update(bson.M{"vkid": vkid}, bson.M{"$set": bson.M{"photo_url": photo_url}})
	if err != nil {
		panic(err)
	}
}

func (dbw *Wrapper) GetTopRatedUsers(sex bool, stat string) (users []User, err error) {
	var sortBy = []string{"-stats." + stat}
	if stat == "fuck" {
		sortBy = append(sortBy, "-stats."+"marrys")
	} else if stat == "marry" {
		sortBy = append(sortBy, "-stats."+"fucks")
	} else {
		sortBy = []string{"-stats.kills"}
	}
	err = dbw.Users.Find(bson.M{"sex": sex}).Sort(sortBy...).Limit(10).All(&users)
	return
}
