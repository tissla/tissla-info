package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func main() {

	port := "8099"

	s := &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	http.HandleFunc("/info", infoHandler)

	log.Printf("> Listening on %s", port)
	log.Fatal(s.ListenAndServe())
}

func infoHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	data := getData()
	json.NewEncoder(w).Encode(data)
}

func getData() []DockerContainer {

	cmd := exec.Command("docker", "ps", "--format", "{{json .}}")
	raw, _ := cmd.Output()

	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")

	out := make([]DockerContainer, 0, len(lines))

	for _, line := range lines {
		var c dockerPSRow
		err := json.Unmarshal([]byte(line), &c)
		if err != nil {
			log.Printf("Error unmarshaling json: %s", err)
		}

		out = append(out, DockerContainer{
			Name:    c.Names,
			Status:  c.Status,
			Image:   c.Image,
			Created: c.CreatedAt,
			Uptime:  c.RunningFor,
		})
	}

	return out
}

// types
type DockerContainer struct {
	Name    string
	Status  string
	Image   string
	Created string
	Uptime  string
}

type dockerPSRow struct {
	Names      string `json:"Names"`
	Image      string `json:"Image"`
	Status     string `json:"Status"`
	RunningFor string `json:"RunningFor"`
	CreatedAt  string `json:"CreatedAt"`
}
