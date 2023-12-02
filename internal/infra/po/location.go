package po

type Location struct {
	ID         uint32 `gorm:"column:id;primaryKey;autoIncrement"`
	LocationID uint32 `gorm:"column:location_id" json:"location_id"`
	Title      string `gorm:"column:title" json:"title"`
	City       string `gorm:"column:city" json:"city"`
	Address    string `gorm:"column:address" json:"address"`
}
