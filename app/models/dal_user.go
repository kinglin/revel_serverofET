package models

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

func (d *Dal) HandleRegister(user User) string {
	collection := d.session.DB("alluser").C("userlist")

	var registerResult string

	i, _ := collection.Find(bson.M{"username": user.Username}).Count()
	if i != 0 {
		registerResult = "username existed"
		return registerResult
	}

	var newUser User
	newUser.Userid = time.Now().UnixNano() / 1000000
	newUser.Gender = 0
	newUser.Username = user.Username
	newUser.Password = user.Password
	newUser.Operation = 0
	newUser.Picture = ""
	newUser.Hobby = ""
	newUser.Birthday = ""
	newUser.Friends = ""
	newUser.LastChangeTime = 0

	collection.Insert(&MyUser{
		newUser.Userid, newUser.Username, newUser.Password})

	//创建这个用户的专属数据库，并将注册信息加入进去
	collection = d.session.DB(strconv.FormatInt(newUser.Userid, 10)).C("user")
	collection.Insert(newUser)

	var conf Configuration
	conf.ConfigurationId = 1
	conf.LoginUser = strconv.FormatInt(newUser.Userid, 10)
	conf.SyncByWifi = 0
	conf.TrackOrNot = 1
	conf.Info = ""
	conf.Changed = 0

	collection = d.session.DB(strconv.FormatInt(newUser.Userid, 10)).C("configuration")
	collection.Insert(conf)

	//var note Note
	//note.NoteId = 0
	//collection = d.session.DB(strconv.FormatInt(newUser.Userid, 10)).C("note")
	//collection.Insert(note)
	//note.NoteId = 1
	//collection.Insert(note)

	registerResult = "success"
	return registerResult
}

func (d *Dal) HandleLogin(user User) (string, int64) {
	collection := d.session.DB("alluser").C("userlist")

	var searchUser User

	err := collection.Find(bson.M{"username": user.Username}).One(&searchUser)
	if err != nil {
		return "user not exist", 0
	} else if searchUser.Password != user.Password {
		return "password is wrong", 0
	} else {
		return "yes", searchUser.Userid
	}
	return "no", 0
}

func (d *Dal) HandleGetAllUserData(struserId string) GetDataResponse {

	var response GetDataResponse

	db := d.session.DB(struserId)

	response.GetResult = "yes"

	//进入user表，拿到所有信息，放到user对象中
	collection := db.C("user")
	if i, _ := collection.Find(nil).Count(); i <= 1 {
		collection.Find(nil).One(&response.User)
	} else {
		collection.Find(nil).All(&response.User)
	}

	//进入configuration表，拿到所有信息，放到configuration对象中
	collection = db.C("configuration")
	if i, _ := collection.Find(nil).Count(); i <= 1 {
		collection.Find(nil).One(&response.Configuration)
	} else {
		collection.Find(nil).All(&response.Configuration)
	}

	//进入note表，循环拿到每条note的所有信息放到note对象中，然后将这些对象放到notes数组中
	collection = db.C("note")
	if i, _ := collection.Find(nil).Count(); i <= 1 {
		collection.Find(nil).One(&response.Notes)
	} else {
		//collection.Find(nil).All(&response.Notes)
		var ids []int
		var note Note
		collection.Find(nil).Distinct("noteid", &ids)
		for _, id := range ids {
			collection.Find(bson.M{"noteid": id}).One(&note)
			response.Notes = append(response.Notes, note)
		}
	}

	//进入coin表，循环拿到每条coin的所有信息放到coin对象中，然后将这些对象放到coins数组中
	collection = db.C("coin")
	if i, _ := collection.Find(nil).Count(); i <= 1 {
		collection.Find(nil).One(&response.Coins)
	} else {
		collection.Find(nil).All(&response.Coins)
	}

	//进入route表，循环拿到每条route的所有信息放到route对象中，然后将这些对象放到routes数组中
	collection = db.C("route")
	if i, _ := collection.Find(nil).Count(); i <= 1 {
		collection.Find(nil).One(&response.Routes)
	} else {
		//collection.Find(nil).All(&response.Routes)
		var ids []int
		var route Route
		collection.Find(nil).Distinct("id", &ids)
		for _, id := range ids {
			collection.Find(bson.M{"id": id}).One(&route)
			response.Routes = append(response.Routes, route)
		}
	}

	//进入treasure表，循环拿到每条treasure的所有信息放到treasure对象中，然后将这些对象放到treasures数组中
	collection = db.C("treasure")
	if i, _ := collection.Find(nil).Count(); i <= 1 {
		collection.Find(nil).One(&response.Treasures)
	} else {
		collection.Find(nil).All(&response.Treasures)
	}

	return response
}

func (d *Dal) HandleUpdateConfiguration(struserId string, conf Configuration) string {
	db := d.session.DB(struserId)
	collection := db.C("configuration")

	err := collection.Update(bson.M{"configurationid": conf.ConfigurationId},
		bson.M{"$set": bson.M{
			"loginuser":  conf.LoginUser,
			"syncbywifi": conf.SyncByWifi,
			"trackornot": conf.TrackOrNot,
			"autopush":   conf.AutoPush,
			"changed":    0,
		}})

	if err != nil {
		return err.Error()
	} else {
		return "yes"
	}
}

func (d *Dal) HandleUpdateUser(struserId string, user User) string {
	db := d.session.DB(struserId)
	collection := db.C("user")

	//userId,_:=strconv.ParseInt(struserId, 0, 64)

	err := collection.Update(bson.M{"userid": user.Userid},
		bson.M{"$set": bson.M{"picture": user.Picture,
			"operation":      0,
			"lastchangetime": user.LastChangeTime,
			"hobby":          user.Hobby,
			"gender":         user.Gender,
			"friends":        user.Friends,
			"birthday":       user.Birthday,
		}})

	if err != nil {
		return err.Error()
	} else {
		return "yes"
	}
}

func (d *Dal) HandleUpdateNotes(struserId string, notes []Note) string {
	db := d.session.DB(struserId)
	collection := db.C("note")

	var err error

	for _, note := range notes {
		switch note.Operation {
		//增加
		case 1:
			note.Operation = 0
			err = collection.Insert(note)
			if err != nil {
				return err.Error()
			}
		//修改
		case 2:
			err = collection.Update(bson.M{"noteid": note.NoteId},
				bson.M{"$set": bson.M{"text": note.Text,
					"operation":      0,
					"lastchangetime": note.LastChangeTime,
					"permission":     note.Permission,
					"pictures":       note.Pictures,
					"weather":        note.Weather,
					"video":          note.Video,
					"voice":          note.Voice,
					"title":          note.Title,
					"time":           note.Time,
					"locationx":      note.Locationx,
					"locationy":      note.Locationy}})
			if err != nil {
				return err.Error()
			}

		//删除
		case 3:
			err = collection.Remove(bson.M{"noteid": note.NoteId})
			if err != nil {
				return err.Error()
			}
		}
	}

	return "yes"

}

func (d *Dal) HandleUpdateRoutes(struserId string, routes []Route) string {
	db := d.session.DB(struserId)
	collection := db.C("route")

	var err error
	for _, route := range routes {
		err = collection.Insert(route)
		if err != nil {
			return err.Error()
		}
	}
	return "yes"

}
