package models

type IpInfo struct {
	Ip        string
	Forwarded string
}

type GeoIpInfo struct {
	GeoIpInfoLocationLatitude  float64
	GeoIpInfoLocationLongitude float64
	GeoIpInfoAccuracy          int32
}

type GeoIp struct {
	Location Location
	Accuracy int32
}

type Location struct {
	Lat float64
	Lng float64
}
