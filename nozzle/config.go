package nozzle

type NozzleConfig struct {
	APIURL                 string `required:"true" envconfig:"api_url"`
	ClientID               string `required:"true" envconfig:"client_id"`
	ClientSecret           string `required:"true" envconfig:"client_secret"`
	FirehoseSubscriptionID string `required:"true" envconfig:"firehose_subscription_id"`
	SkipSSL                bool   `default:"false" envconfig:"skip_ssl"`
	Target                 string `required:"true" envconfig:"target"`
}
