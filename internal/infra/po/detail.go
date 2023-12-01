package po

type Status int

const (
	Create Status = iota
	PackageReceived
	InTransit
	OutForDelivery
	DeliveryAttempted
	Delivered
	ReturnedToSender
	Exception
)

type Detail struct {
	ID            int64  `gorm:"column:id;primaryKey"`
	Date          string `gorm:"column:date"`
	TimeHour      string `gorm:"column:time"`
	Status        Status `gorm:"column:status"`
	LocationID    int64  `gorm:"column:location_id"`
	LocationTitle string `gorm:"column:location_title"`
}
