package api

type FakeList struct {
	ErrorMsg string         `json:"error,omitempty"`
	Data     []FakeResponse `json:"data,omitempty"`
}

type FakeResponse struct {
	Sno            uint32 `json:"sno"`
	TrackingStatus int8   `json:"tracking_status"`
}
