package components

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"server/src"
	"server/src/utils"
	"sync"
	"time"
)

const (
	DEFAULT_SLEEP_TIME  = 10
	REQUEST_SLEEP_TIME  = 20
	RESPONSE_SLEEP_TIME = 0
)

type LoadBalancer struct {
	FilePath        string       `json:"-"`
	Host            string       `json:"host"`
	Port            string       `json:"port"`
	Services        []HostInfo   `json:"services"`
	mutexSync       sync.RWMutex `json:"-"`
	targetHost      string
	enableBalancing bool
}

func NewLoadBalancer(path string) *LoadBalancer {
	file_bytes, err := os.ReadFile(path)
	if err != nil {
		log.Panicln("Error when opening file: ", err)
		return nil
	}
	var lb LoadBalancer
	json.Unmarshal(file_bytes, &lb)
	lb.FilePath = path
	os.Setenv("host", lb.Host)
	os.Setenv("host", lb.Port)
	return &lb
}

func (lb *LoadBalancer) EnableLoadBalancing() *LoadBalancer {
	lb.enableBalancing = true
	lb.enableServiceDiscovery()
	return lb
}

func (lb *LoadBalancer) SyncFile() error {
	lb_bytes, _ := lb.to_json()
	(*lb).mutexSync.Lock()
	err := os.WriteFile(lb.FilePath, lb_bytes, 0644)
	defer (*lb).mutexSync.Unlock()
	if err != nil {
		log.Println("Could not write in config.json: ", err)
		return err
	}
	if src.VERBOSE {
		log.Println("Saving data on file... ", string(lb_bytes))
	}

	return nil
}

func (lb *LoadBalancer) FindService() {
	(*lb).mutexSync.RLock()
	defer (*lb).mutexSync.RUnlock()
	idx := rand.Intn(2)
	lb.targetHost = lb.Services[idx].Url
}

func (lb *LoadBalancer) GetDiscoveryInfo() string {
	(*lb).mutexSync.RLock()
	defer (*lb).mutexSync.RUnlock()
	lb_bytes, _ := lb.to_json()
	return string(lb_bytes)
}

func (lb *LoadBalancer) to_json() ([]byte, error) {
	(*lb).mutexSync.RLock()
	defer (*lb).mutexSync.RUnlock()
	lb_bytes, err := json.MarshalIndent(&lb, "", "	")
	if err != nil {
		log.Println("Could not parse LoadBalancer to JSON.", err)
		return nil, err
	}
	return lb_bytes, nil
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

func (lb *LoadBalancer) enableServiceDiscovery() {
	header_content := map[string]string{
		utils.UserAgenteKey: utils.UserAgenteKey,
		utils.DateKey:       utils.GetTimeHere(),
		utils.ConnectionKey: utils.KeepAlive,
	}

	for _, host := range lb.Services {
		delay := host.GetDelay()
		service_url := host.Url
		hc_endpoint := host.HealthCheck.Endpoint
		header := &Header{
			Method:        "GET",
			Path:          hc_endpoint,
			HeaderContent: header_content,
		}

		go func(host HostInfo) {
			for {
				start := time.Now()
				ctx := &Context{
					Url:      service_url,
					Endpoint: hc_endpoint,
					Header:   header,
				}
				_, err := Ping(ctx)
				end := float64(time.Since(start).Microseconds()) * .001
				if err != nil {
					lb.updateHealthCheck(host.Id, false, utils.ToFixed(end))
					time.Sleep(delay)
					continue
				}

				lb.updateHealthCheck(host.Id, true, utils.ToFixed(end))
				time.Sleep(delay)
			}
		}(host)
	}
}
