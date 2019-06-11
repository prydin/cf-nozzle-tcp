package nozzle

import (
	"github.com/cloudfoundry-community/go-cfclient"
	"net/url"
)

func NewClient(conf *NozzleConfig) (*cfclient.Client, error) {
	apiURL := conf.APIURL
	if !isValidURL(apiURL) {
		apiURL = "https://" + apiURL
	}
	config := &cfclient.Config{
		ApiAddress:        apiURL,
		ClientID:          conf.ClientID,
		ClientSecret:      conf.ClientSecret,
		SkipSslValidation: conf.SkipSSL,
	}

	client, err := cfclient.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// isValidUrl tests a string to determine if it is a url or not.
func isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}
	return true
}
