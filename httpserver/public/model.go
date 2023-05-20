package public

type serverStatusDto struct {
	Domain       string `json:"domain"`
	Available    bool   `json:"available"`
	ResponseTime string `json:"responseTime,omitempty"`
	ServerError  string `json:"error,omitempty"`
	LastUpdate   string `json:"lastUpdate"`
}
