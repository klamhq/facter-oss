package external

import (
	"encoding/json"
	"time"

	"github.com/klamhq/facter-oss/pkg/models"
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
)

// GetIpInfo retrieves the public IP and forwarded IP information from a remote web service.
func GetIpInfo(url string, timeout int) (*models.IpInfo, error) {
	request := gorequest.New().Timeout(time.Duration(timeout) * time.Second)
	_, publicIp, err := request.Get(url + "ip").End()

	if err != nil {
		logrus.Errorf("Error with remote webservice when we want get ip %s", err)
		return nil, err[0]
	}

	_, forwardedHeader, err := request.Get(url + "forwarded").End()
	if err != nil {
		logrus.Errorf("Error with remote webservice when we want get forwarded ip %s", err)
		return nil, err[0]

	}

	return &models.IpInfo{
		Ip:        publicIp,
		Forwarded: forwardedHeader,
	}, nil
}

// GetGeoIpLocalisation retrieves geographical location information based on the public IP address using Google Geo API.
func GetGeoIpLocalisation(googleApiKey, googleGeoUrl string, timeout int) (*models.GeoIpInfo, error) {
	geoRequest := gorequest.New().Timeout(time.Duration(timeout) * time.Second)
	googleFulUrlWithKey := googleGeoUrl + "?key=" + googleApiKey
	_, geoLocateInfo, err := geoRequest.Post(googleFulUrlWithKey).Send(`{"considerIp": true}`).End()
	if err != nil {
		logrus.Errorf("Error with google remote geolocate webservice when we want get latitude and longitude ip %s", err)
		return nil, err[0]

	}
	var geoIp models.GeoIp
	unmarshallErr := json.Unmarshal([]byte(geoLocateInfo), &geoIp)
	if unmarshallErr != nil {
		logrus.Errorf("Unable to unmarshall geoLocateInfo: %s", unmarshallErr)
	}
	return &models.GeoIpInfo{
		GeoIpInfoLocationLatitude:  geoIp.Location.Lat,
		GeoIpInfoLocationLongitude: geoIp.Location.Lng,
		GeoIpInfoAccuracy:          geoIp.Accuracy,
	}, nil

}
