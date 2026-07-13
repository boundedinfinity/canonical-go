package universal_serial_bus

type Connector string

var Connectors = connectors{
	A:                "electronics.usb.connector.a",
	MiniA:            "electronics.usb.connector.mini-a",
	MicroA:           "electronics.usb.connector.micro-a",
	B:                "electronics.usb.connector.b",
	BSuperSpeed:      "electronics.usb.connector.b-super-speed",
	MiniB:            "electronics.usb.connector.mini-b",
	MicroB:           "electronics.usb.connector.micro-b",
	MicroBSuperSpeed: "electronics.usb.connector.micro-b-super-speed",
	C:                "electronics.usb.connector.c",
}

type connectors struct {
	A                Connector
	MiniA            Connector
	MicroA           Connector
	B                Connector
	BSuperSpeed      Connector
	MiniB            Connector
	MicroB           Connector
	MicroBSuperSpeed Connector
	C                Connector
}
