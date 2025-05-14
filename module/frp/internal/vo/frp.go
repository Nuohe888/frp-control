package vo

type AuthReq struct {
	Node  string `json:"node"`
	Uuid  string `json:"uuid"`
	RunId string `json:"runId"`
}

type Auth2Req struct {
	Name  string `json:"name"`
	RunId string `json:"runId"`
	Node  string `json:"node"`
	Uuid  string `json:"uuid"`
	Port  string `json:"port"`
	Type  string `json:"type"`
}

type SpeedReq struct {
	Uuid  string `json:"uuid"`
	RunId string `json:"runId"`
	Node  string `json:"node"`
}

type SpeedRes struct {
	Limit int `json:"limit"`
}
