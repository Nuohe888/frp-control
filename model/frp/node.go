package frp

type Node struct {
	Model
	Id       uint   `json:"id" gorm:"primarykey"`
	Name     string `json:"name"` //节点名(唯一)
	Ip       string `json:"ip"`   //节点ip
	Password string `json:"password"`
	Desc     string `json:"desc"`
	Status   *int   `json:"status"`
}

func (Node) TableName() string {
	return "frp_node"
}
