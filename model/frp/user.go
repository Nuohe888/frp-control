package frp

type User struct {
	Model
	Uuid     string `gorm:"primarykey" json:"uuid"`
	Nickname string `json:"nickname"`
	Flow     int64  `json:"flow"`  //单位字节
	Speed    int    `json:"speed"` //单位M
	Note     string `json:"note"`
	Status   int    `json:"status"`
}

func (User) TableName() string {
	return "frp_user"
}
