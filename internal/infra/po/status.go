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
	StatusMsgMapping = map[Status]string{
		Create:            "Created",
		PackageReceived:   "Package Received",
		InTransit:         "In Transit",
		OutForDelivery:    "Out for Delivery",
		DeliveryAttempted: "Delivery Attempted",
		Delivered:         "Delivered",
		ReturnedToSender:  "Returned to Sender",
		Exception:         "Exception",
	}
)
