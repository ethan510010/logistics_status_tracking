package po

type Recipient struct {
	ID      uint32 `gorm:"column:id;primaryKey"`
	Name    string `gorm:"column:name"`
	Address string `gorm:"column:address"`
	Phone   string `gorm:"column:phone"`
}
