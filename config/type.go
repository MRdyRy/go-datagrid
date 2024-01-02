package config

type CacheTemplate struct {
	DistributedCache struct {
		Mode     string `json:"mode"`
		Encoding struct {
			MediaType string `json:"media-type"`
		} `json:"encoding"`
		Statistics bool `json:"statistics"`
	} `json:"distributed-cache"`
}

func GenerateTemplate() string {
	template := `{
		"distributed-cache": {
		  "mode": "SYNC",
		  "encoding": {
			"media-type": "application/x-protostream"
		  },
		  "statistics": true
		}
	  }`
	return template
}
