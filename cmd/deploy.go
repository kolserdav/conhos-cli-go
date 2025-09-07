package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
)

type Deploy struct {
	Conn             *websocket.Conn
	CacheWorked      bool
	CacheFilePath    map[string]string
	FileList         []string
	UploadedServices []string
	IsNewUpload      bool
	WaitGitUpload    map[string]bool
	Token            string
	ConnID           string
	UserID           string
	Project          string
	Config           *ConfigFile
}

type ConfigFile struct {
	Services map[string]ServiceConfig
	Volumes  []VolumeConfig
}

type ServiceConfig struct {
	PWD     string
	Image   string
	Volumes []string
}

type VolumeConfig struct {
	Name string
}

type WSMessage struct {
	Type   string          `json:"type"`
	Data   json.RawMessage `json:"data"`
	Status string          `json:"status"`
}

type PrepareDeployData struct {
	Service string `json:"service"`
	Exclude string `json:"exclude"`
	PWD     string `json:"pwd"`
	Active  bool   `json:"active"`
	Cache   string `json:"cache"`
	Git     string `json:"git"`
}

func NewDeploy(conn *websocket.Conn, token, connID, userID string) *Deploy {
	return &Deploy{
		Conn:          conn,
		CacheFilePath: make(map[string]string),
		WaitGitUpload: make(map[string]bool),
		Token:         token,
		ConnID:        connID,
		UserID:        userID,
	}
}

func (d *Deploy) Listener() {
	for {
		_, message, err := d.Conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("JSON unmarshal error:", err)
			continue
		}

		switch msg.Type {
		case "prepareDeployCli":
			var data PrepareDeployData
			if err := json.Unmarshal(msg.Data, &data); err == nil {
				d.PrepareUpload(data)
			}
		case "deployDeleteFilesCli":
			// TODO: Implement
		case "deployGitCli":
			// TODO: Implement
		case "deployPrepareVolumeUploadCli":
			// TODO: Implement
		case "deployProgressCli":
			// TODO: Implement
		case "acceptDeleteCli":
			// TODO: Implement
		default:
			log.Println("Unhandled message type:", msg.Type)
		}
	}
}

func (d *Deploy) PrepareUpload(data PrepareDeployData) {
	fmt.Printf("Prepare upload for service: %s\n", data.Service)

	d.UploadedServices = append(d.UploadedServices, data.Service)
	needUploadServices := d.GetNeedUploadServices()
	last := len(needUploadServices) <= len(d.UploadedServices)

	if data.Git != "" {
		fmt.Println("Starting synchronize git")
		d.SendMessage("deployGitServer", map[string]interface{}{
			"service": data.Service,
			"last":    last,
			"pwd":     data.PWD,
			"git":     data.Git,
			"active":  data.Active,
		})
		d.WaitGitUpload[data.Service] = true
		return
	}

	d.SetCacheFilePath(data.Service, d.Project)

	if !data.Active {
		fmt.Printf("Skipping upload for deleted service: %s\n", data.PWD)
		d.SendMessage("deployEndServer", map[string]interface{}{
			"service": data.Service,
			"skip":    true,
			"last":    last,
			"file":    "",
			"latest":  true,
			"num":     0,
		})
		return
	}

	fmt.Printf("Preparing to upload service files: %s\n", data.PWD)
	// TODO: Implement cache checking and file upload logic
}

func (d *Deploy) GetNeedUploadServices() []string {
	var services []string
	for name, svc := range d.Config.Services {
		if svc.PWD != "" && d.IsCustomService(svc.Image) {
			services = append(services, name)
		}
	}
	return services
}

func (d *Deploy) IsCustomService(image string) bool {
	return !strings.Contains(image, "/") // Simplified check for custom service
}

func (d *Deploy) SetCacheFilePath(service, project string) {
	packagePath := filepath.Join(os.TempDir(), project, service)
	if err := os.MkdirAll(packagePath, 0755); err == nil {
		d.CacheFilePath[service] = filepath.Join(packagePath, "cache.json")
	}
}

func (d *Deploy) SendMessage(msgType string, data map[string]interface{}) {
	message := map[string]interface{}{
		"token":   d.Token,
		"type":    msgType,
		"userId":  d.UserID,
		"connId":  d.ConnID,
		"status":  "info",
		"package": "conhos-cli",
		"data":    data,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshaling message:", err)
		return
	}

	if err := d.Conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
		log.Println("WebSocket write error:", err)
	}
}
