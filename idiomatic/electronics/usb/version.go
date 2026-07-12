package usb

type Version string

var Versions = versions{
	V1_0: "electronics.usb.version.1.0",
	V1_1: "electronics.usb.version.1.1",
	V2_0: "electronics.usb.version.2.0",
	V2_1: "electronics.usb.version.2.1",
	V3_0: "electronics.usb.version.3.0",
	V3_1: "electronics.usb.version.3.1",
	V3_2: "electronics.usb.version.3.2",
	V4_0: "electronics.usb.version.4.0",
}

type versions struct {
	V1_0 Version
	V1_1 Version
	V2_0 Version
	V2_1 Version
	V3_0 Version
	V3_1 Version
	V3_2 Version
	V4_0 Version
}
