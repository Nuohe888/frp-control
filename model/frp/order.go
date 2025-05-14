package frp

type Order struct {
	Model
	Id            uint   `gorm:"primarykey"`
	Username      string `json:"username"`
	ExpTime       int64  `json:"expTime"` //到期时间
	NodeId        uint   `json:"nodeId"`
	NodeName      string `json:"nodeName"`
	NodeIp        string `json:"nodeIp"`
	UserUuid      string `json:"userUuid"`
	Type          string `json:"type"`
	Port          string `json:"port"`
	RunId         string `json:"runId"`
	Status        int    `json:"status"`
	Speed         int    `json:"speed"`         //带宽限速
	BandwidthUp   int64  `json:"bandwidthUp"`   //带宽上行
	BandwidthDown int64  `json:"bandwidthDown"` //带宽下行
}

func (Order) TableName() string {
	return "frp_order"
}
