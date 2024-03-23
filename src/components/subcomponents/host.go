package subcomponents

import (
	"server/src/utils"
	"time"
)

type Endpoint struct {
	AllowedMethods []string `json:"allowed_methods"`
	AllowedHosts   []string `json:"allowed_hosts"`
}

type HealthCheck struct {
	Delay    int     `json:"delay"`
	Endpoint string  `json:"endpoint"`
	Duration float64 `json:"duration"`
	GotError bool    `json:"got_error"`
}

type HostInfo struct {
	Id          string              `json:"id"`
	Url         string              `json:"url"`
	Alive       bool                `json:"alive"`
	HealthCheck HealthCheck         `json:"healthcheck"`
	SecretKey   string              `json:"secret_key"`
	Timeout     int                 `json:"timeout"`
	Endpoints   map[string]Endpoint `json:"endpoints"`
}

func (hi *HostInfo) GetDelay() time.Duration {
	return utils.GetTimeInSeconds(hi.HealthCheck.Delay)
}

func (hi *HostInfo) GetTimeout() time.Duration {
	return utils.GetTimeInSeconds(hi.Timeout)
}
