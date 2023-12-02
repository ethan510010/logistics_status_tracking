package po

type TrackingStatus struct {
	ID                    uint32 `gorm:"column:id;primaryKey;autoIncrement"`
	Sno                   uint32 `gorm:"column:sno;index:idx_sno"`
	Status                int8   `gorm:"column:status"`
	EstimatedDeliveryTime string `gorm:"column:estimated_delivery"`
	RecipientID           uint32 `gorm:"column:recipient;index:idx_recipient_id"`
	CurrentLocationID     uint32 `gorm:"current_location_id"`
}
