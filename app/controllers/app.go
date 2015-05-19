package controllers

import (
	"encoding/json"
	"github.com/robfig/revel"
	"myapp/app/models"
	//"strconv"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	greeting := "welcome"
	return c.Render(greeting)
}

func (c App) Register(username, password string) revel.Result {

	//定义回复结构体
	type Response struct {
		RegisterResult string `json:"registerResult"`
	}

	/*
		这里是到数据库检验新用户有效性
	*/
	registerResult := ""

	var user models.User
	user.Username = username
	user.Password = password

	dal, _ := models.NewDal()
	registerResult = dal.HandleRegister(user)

	defer dal.Close()

	//回复给客户端
	response := &Response{
		RegisterResult: registerResult}
	jsonResponse, _ := json.Marshal(response)

	return c.RenderJson(string(jsonResponse))
}

func (c App) Login(username, password string) revel.Result {

	//定义回复结构体
	type Response struct {
		LoginResult string `json:"loginResult"`
		UserId      int64  `json:"userId"`
	}

	/*
		这里是到数据库检验新用户有效性
	*/
	var user models.User
	user.Username = username
	user.Password = password

	dal, _ := models.NewDal()
	loginResult, userId := dal.HandleLogin(user)

	defer dal.Close()

	response := &Response{
		LoginResult: loginResult,
		UserId:      userId}
	jsonResponse, _ := json.Marshal(response)

	return c.RenderJson(string(jsonResponse))
}

func (c App) GetAllUserData(struserId string) revel.Result {

	dal, _ := models.NewDal()
	defer dal.Close()

	response := dal.HandleGetAllUserData(struserId)
	jsonResponse, _ := json.Marshal(response)

	return c.RenderJson(string(jsonResponse))
}

func (c App) SyncData(struserId string, json_allChangedData string) revel.Result {
	/*
		这里解析json
		然后将各个数据存到相应的数据结构中
		然后返回同步成功
	*/

	type Response struct {
		SyncResult string `json:"syncResult"`
	}

	syncResult := "ok"
	syncResult1 := "ok"
	syncResult2 := "ok"
	syncResult3 := "ok"
	syncResult4 := "ok"

	type Request struct {
		Configuration models.Configuration `json:"configuration"`
		User          models.User          `json:"user"`
		Notes         []models.Note        `json:"notes"`
		Routes        []models.Route       `json:"routes"`
	}

	request := new(Request)

	json.Unmarshal([]byte(json_allChangedData), &request)

	dal, _ := models.NewDal()
	defer dal.Close()

	if request.Configuration.ConfigurationId != 0 {
		syncResult1 = dal.HandleUpdateConfiguration(struserId, request.Configuration)
	}

	if request.User.Userid != 0 {
		syncResult2 = dal.HandleUpdateUser(struserId, request.User)
	}

	if len(request.Notes) != 0 {
		syncResult3 = dal.HandleUpdateNotes(struserId, request.Notes)
	}

	if len(request.Routes) != 0 {
		syncResult4 = dal.HandleUpdateRoutes(struserId, request.Routes)
	}

	if (syncResult1 == "yes" || syncResult1 == "ok") &&
		(syncResult2 == "yes" || syncResult2 == "ok") &&
		(syncResult3 == "yes" || syncResult3 == "ok") &&
		(syncResult4 == "yes" || syncResult4 == "ok") {
		syncResult = "yes"
	}

	syncResult = syncResult2

	response := &Response{
		SyncResult: syncResult}
	jsonResponse, _ := json.Marshal(response)
	return c.RenderJson(string(jsonResponse))

	//return c.RenderText(json_allChangedData)
}
