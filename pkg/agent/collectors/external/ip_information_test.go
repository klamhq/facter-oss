package external

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var googleGeoApikey = "AIzaSyB9i6aHGUa5Iv0tfDb9sy_HqXUebFmk8wI"
var googleGeoUrl = "https://www.googleapis.com/geolocation/v1/geolocate"
var url = "https://ifconfig.me/"
var badUrl = "https://ifconfig.meaa/"
var timeout = 2

func TestIpInformation(t *testing.T) {
	timeout := 2
	ip, err := GetIpInfo(url, timeout)
	assert.NoError(t, err)
	assert.NotNil(t, ip.Ip)
	assert.NotNil(t, ip.Forwarded)
	assert.NotEqual(t, ip.Ip, "8.8.8.8")
	assert.NotEqual(t, ip.Forwarded, "8.8.8.8, 0.0.0.0")
}

func TestFailedIpInformation(t *testing.T) {
	// test with bad url
	ip, err := GetIpInfo(badUrl, timeout)
	assert.Error(t, err)
	assert.Nil(t, ip)
}

func TestGeoIpInformation(t *testing.T) {
	floatTest := 47.726592
	var i32 int32 = -1
	geoIp, err := GetGeoIpLocalisation(googleGeoApikey, googleGeoUrl, timeout)
	assert.NoError(t, err)
	assert.NotNil(t, geoIp.GeoIpInfoLocationLatitude)
	assert.NotNil(t, geoIp.GeoIpInfoLocationLongitude)
	assert.NotNil(t, geoIp.GeoIpInfoAccuracy)
	assert.IsType(t, geoIp.GeoIpInfoLocationLatitude, floatTest)
	assert.IsType(t, geoIp.GeoIpInfoLocationLongitude, floatTest)
	assert.IsType(t, geoIp.GeoIpInfoAccuracy, i32)
}
