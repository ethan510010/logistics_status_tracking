package api

type QueryResponse struct {
	Status string      `json:"status"`
	Data   *QueryData  `json:"data"`
	Error  *QueryError `json:"error"`
}

type QueryData struct {
	Sno               uint32               `json:"sno"`
	TrackingStatus    string               `json:"tracking_status"`
	EstimatedDelivery string               `json:"estimated_delivery"`
	Details           []QueryDetail        `json:"details"`
	Recipient         QueryRecipient       `json:"recipient"`
	CurrentLocation   QueryCurrentLocation `json:"current_location"`
}

type QueryDetail struct {
	ID            uint32 `json:"id"`
	Date          string `json:"date"`
	Hr            string `json:"time"`
	Status        string `json:"status"`
	LocationID    uint32 `json:"location_id"`
	LocationTitle string `json:"location_title"`
}

type QueryRecipient struct {
	ID      uint32 `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type QueryCurrentLocation struct {
	LocationID uint32 `json:"location_id"`
	Title      string `json:"title"`
	City       string `json:"city"`
	Address    string `json:"address"`
}

type QueryError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type FakeList struct {
	ErrorMsg string         `json:"error,omitempty"`
	Data     []FakeResponse `json:"data,omitempty"`
}

type FakeResponse struct {
	Sno            uint32 `json:"sno"`
	TrackingStatus int8   `json:"tracking_status"`
}
