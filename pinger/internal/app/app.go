package app

import (
	"context"
	"log"
	"os/exec"
	"pinger/internal/config"
	"pinger/internal/models"
	"pinger/internal/services"
	"regexp"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var wg sync.WaitGroup

var containersL []models.Container

var mx sync.Mutex

func pingContainer(ip string, i int) {
	cmd := exec.Command("ping", "-c", "1", ip)
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Error pinging container %v:  %v", ip, err)
		wg.Done()
		return
	}
	re := regexp.MustCompile(`time=([\d.]+)`)
	matches := re.FindStringSubmatch(string(out))
	mx.Lock()
	containersL = append(containersL, models.Container{
		Id:       i,
		Ip:       ip,
		TimeMs:   matches[1],
		PingDate: time.Now().Format(time.ANSIC),
	})
	mx.Unlock()
	wg.Done()
}

func getContainersData() {

	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Panic(err)
	}
	defer apiClient.Close()
	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		log.Panic(err)
	}

	for _, ctr := range containers {
		ctrJSON, err := apiClient.ContainerInspect(context.Background(), ctr.ID)

		if err != nil {
			log.Printf("Error inspecting container: %v", err)
			continue
		}
		i := 1
		for _, networkProps := range ctrJSON.NetworkSettings.Networks {
			wg.Add(1)
			go pingContainer(networkProps.IPAddress, i)
			i += 1
		}
	}
}

func Start() {
	pConfig := config.New().PConfig
	tick := time.Tick(time.Duration(pConfig.PINGER_WAIT) * time.Second)
	for {
		select {
		case <-tick:
			containersL = make([]models.Container, 0, 10)
			getContainersData()
			wg.Wait()
			services.PostData(containersL)
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
