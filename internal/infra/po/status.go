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

var (
	DeliverStatusList = []Status{
		Create,
		PackageReceived,
		InTransit,
		OutForDelivery,
		DeliveryAttempted,
		Delivered,
		ReturnedToSender,
		Exception,
	}
)
