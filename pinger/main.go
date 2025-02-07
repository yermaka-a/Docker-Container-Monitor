package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var wg sync.WaitGroup

type Container struct {
	Id       int    `json:"id"`
	Ip       string `json:"ip"`
	TimeMs   string `json:"timeMs"`
	PingDate string `json:"pingDate"`
}

const url = "http://web:80/containers"

var containersL []Container

var mx sync.Mutex

func pingContainer(ip string, i int) {
	cmd := exec.Command("ping", "-c", "1", ip)
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Error pinging container %v:  %v", ip, err)
		return
	}
	re := regexp.MustCompile(`time=([\d.]+)`)
	matches := re.FindStringSubmatch(string(out))
	mx.Lock()
	containersL = append(containersL, Container{
		Id:       i,
		Ip:       ip,
		TimeMs:   matches[1],
		PingDate: time.Now().String(),
	})
	mx.Unlock()
	wg.Done()
}

func postData(containersL []Container) {
	jsonData, err := json.Marshal(containersL)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Response Status:", resp.Status)
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

func main() {
	tick := time.Tick(7 * time.Second)
	for {
		select {
		case <-tick:
			containersL = make([]Container, 0, 10)
			getContainersData()
			wg.Wait()
			postData(containersL)
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
