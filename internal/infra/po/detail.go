package po

type Detail struct {
	ID            uint32 `gorm:"column:id;primaryKey;autoIncrement"`
	Date          string `gorm:"column:date"`
	TimeHour      string `gorm:"column:time"`
	Status        int8   `gorm:"column:status"`
	LocationID    uint32 `gorm:"column:location_id;index:idx_location_id"`
	LocationTitle string `gorm:"column:location_title"`
	Sno           uint32 `gorm:"column:sno;index:idx_sno"`
}
