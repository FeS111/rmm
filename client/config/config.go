package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"festech.de/rmm/client/http"
	"festech.de/rmm/client/models"
	"festech.de/rmm/client/system"
	"github.com/google/uuid"
)

var Device models.Device

var configPath = system.CreateConfigurationPath()

func CreateConfiguration() models.Device {
	id := uuid.New().String()

	systemInfo := models.SystemInfo{
		HostName:   system.GetHostName(),
		Os:         system.GetOS(),
		MacAddress: system.GetMacAddress(),
		Cores:      system.GetCores(),
		Memory:     system.GetMemory(),
		Disk:       system.GetDisk(),
	}
	device := models.Device{
		DeviceID:   id,
		Name:       system.GetHostName(),
		SystemInfo: systemInfo,
	}

	status, textBody, _ := http.Post("http://localhost:8080/devices", device, &device)
	if status == 400 && textBody == "Device with this mac address is already registered" {
		http.Get(fmt.Sprintf("http://localhost:8080/devices/%s?mac=1", device.SystemInfo.MacAddress), &device)
	}
	WriteConfiguration(device)
	return device
}

func WriteConfiguration(device models.Device) {
	file, _ := json.MarshalIndent(device, "", " ")
	ioutil.WriteFile(configPath, file, 0644)
}

func ReadConfiguration() {
	if !system.FileOrFolderExists(configPath) {
		CreateConfiguration()
	}
	file, _ := ioutil.ReadFile(configPath)
	device := models.Device{}
	json.Unmarshal(file, &device)

	Device = device
}
