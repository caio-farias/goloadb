package components

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"server/src"
	"server/src/utils"
	"sync"
	"time"
)

type LoadBalancer struct {
	Host            string       `json:"host"`
	Port            string       `json:"port"`
	FilePath        string       `json:"-"`
	Services        []HostInfo   `json:"services"`
	mutexSync       sync.RWMutex `json:"-"`
	enableBalancing bool         `json:"-"`
	client          *http.Client `json:"-"`
	targetHost      string       `json:"-"`
}

func NewLoadBalancer(path string) *LoadBalancer {
	file_bytes, err := os.ReadFile(path)
	if err != nil {
		log.Panicln(utils.FAILED_LOAD_CONFIGFILE, err)
		return nil
	}

	var lb LoadBalancer
	json.Unmarshal(file_bytes, &lb)
	os.Setenv("host", lb.Host)
	os.Setenv("port", lb.Port)

	return &lb
}

func (lb *LoadBalancer) toJson() ([]byte, error) {
	lb_bytes, err := json.MarshalIndent(&lb, "", "	")
	if err != nil {
		log.Println("Could not parse LoadBalancer to JSON.", err)
		return nil, err
	}
	return lb_bytes, nil
}

func (lb *LoadBalancer) SyncFile() error {
	(*lb).mutexSync.Lock()
	defer (*lb).mutexSync.Unlock()

	lb_bytes, _ := lb.toJson()
	err := os.WriteFile(lb.FilePath, lb_bytes, 0644)
	if err != nil {
		log.Println(utils.COULD_NOT_WRITE_CONFIG_FILE, err)
		return err
	}
	if src.VERBOSE {
		log.Println("Saving data on file... ", string(lb_bytes))
	}

	return nil
}

func (lb *LoadBalancer) manageLockReadStringFunction(fn func() string) string {
	(*lb).mutexSync.RLock()
	defer (*lb).mutexSync.RUnlock()
	res := fn()
	return res
}

func (lb *LoadBalancer) DiscoverService() string {
	return lb.manageLockReadStringFunction(func() string {
		return lb.targetHost
	})
}

func (lb *LoadBalancer) GetDiscoveryInfo() string {
	return lb.manageLockReadStringFunction(func() string {
		lb_bytes, _ := lb.toJson()
		return string(lb_bytes)
	})
}

func (lb *LoadBalancer) updateHealthCheck(hostId string, status bool, diffTime float64) {
	(*lb).mutexSync.Lock()
	defer (*lb).mutexSync.Unlock()

	for idx, host := range lb.Services {
		if host.Id == hostId {
			lb.Services[idx].Alive = status
			lb.Services[idx].HealthCheck.GotError = !status
			lb.Services[idx].HealthCheck.Duration = diffTime
			break
		}
	}
}

func (lb *LoadBalancer) EnableLoadBalancing() {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	lb.client = &http.Client{Transport: tr}
	lb.enableServiceDiscovery()
}

func (lb *LoadBalancer) enableServiceDiscovery() {
	for _, host := range lb.Services {
		go func(host HostInfo) {
			for lb.enableBalancing {
				start := time.Now()
				service_url := host.Url
				hc_endpoint := host.HealthCheck.Endpoint

				if _, err := lb.client.Get(service_url + hc_endpoint); err != nil {
					log.Println(utils.SERVICE_COULD_NOT_RESPOND, err)
					lb.updateHealthCheck(host.Id, false, utils.GetTimeDiffInMillisec(start))
					time.Sleep(host.GetDelay())
					continue
				}

				lb.updateHealthCheck(host.Id, true, utils.GetTimeDiffInMillisec(start))
				time.Sleep(host.GetDelay())
			}
		}(host)
	}
}
