package models

//用户简略信息，登陆时用
type MyUser struct {
	Userid   int64  `json:"userId"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//用户详细信息
type User struct {
	Userid         int64  `json:"userId"`
	Gender         int    `json:"gender"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Operation      int    `json:"operation"`
	Picture        string `json:"picture"`
	Birthday       string `json:"birthday"`
	Hobby          string `json:"hobby"`
	Friends        string `json:"friends"`
	LastChangeTime int64  `json:"lastChangeTime"`
}

//记事
type Note struct {
	NoteId         int64   `json:"noteId"`
	Time           string  `json:"time"`
	Permission     int     `json:"permission"`
	Weather        int     `json:"weather"`
	Text           string  `json:"text"`
	Title          string  `json:"title"`
	Pictures       string  `json:"pictures"`
	Voice          string  `json:"voice"`
	Video          string  `json:"video"`
	Locationx      float64 `json:"locationx"`
	Locationy      float64 `json:"locationy"`
	LastChangeTime int64   `json:"lastChangeTime"`
	Operation      int     `json:"operation"` //0：已同步 1：新增 2：修改 3：删除
}

//配置
type Configuration struct {
	ConfigurationId int    `json:"configurationId"`
	LoginUser       string `json:"loginUser"`
	SyncByWifi      int    `json:"syncByWifi"`
	TrackOrNot      int    `json:"trackOrNot"`
	AutoPush        int    `json:"autoPush"`
	Info            string `json:"info"`
	Changed         int    `json:"changed"`
}

//路线
type Route struct {
	Id         int64   `json:"id"`
	Latitude   float64 `json:"latitude"`
	Longtitude float64 `json:"longtitude"`
	Date       string  `json:"date"`
}

//积分
type Coin struct {
	CoinId  int64  `json:"coinId"`
	Time    string `json:"time"`
	Grade   int    `json:"grade"`
	Content string `json:"content"`
}

//宝藏
type Treasure struct {
	TreasureId int64  `json:"treasureId"`
	Time       string `json:"time"`
	Content    string `json:"content"`
}

type GetDataResponse struct {
	GetResult     string        `json:"getResult"`
	User          User          `json:"user"`
	Configuration Configuration `json:"configuration"`
	Notes         []Note        `json:"notes"`
	Routes        []Route       `json:"routes"`
	Coins         []Coin        `json:"coins"`
	Treasures     []Treasure    `json:"treasures"`
}
