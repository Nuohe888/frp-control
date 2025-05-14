package client

type FlowRes struct {
	Tcp    []*ProxyStatsInfo `json:"tcp"`
	Udp    []*ProxyStatsInfo `json:"udp"`
	Https  []*ProxyStatsInfo `json:"https"`
	Tcpmux []*ProxyStatsInfo `json:"tcpmux"`
	Stcp   []*ProxyStatsInfo `json:"stcp"`
	Sudp   []*ProxyStatsInfo `json:"sudp"`
}

type ProxyStatsInfo struct {
	Name            string `json:"name"`
	Conf            any    `json:"conf"`
	ClientVersion   string `json:"clientVersion,omitempty"`
	TodayTrafficIn  int64  `json:"todayTrafficIn"`
	TodayTrafficOut int64  `json:"todayTrafficOut"`
	CurConns        int64  `json:"curConns"`
	LastStartTime   string `json:"lastStartTime"`
	LastCloseTime   string `json:"lastCloseTime"`
	Status          string `json:"status"`
}
