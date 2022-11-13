package system

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func run(cmd string) string {
	out := exec.Command("bash", "-c", cmd)
	stdout, _ := out.Output()
	return string(stdout)
}

func GetHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}

func GetMacAddress() string {
	ifas, err := net.Interfaces()
	if err != nil {
		return ""
	}
	return ifas[1].HardwareAddr.String()
}

func GetOS() string {
	var os string
	if IsLinux() {
		os = strings.ReplaceAll(run("uname -osri"), "\n", "")
	}
	return os
}

func GetCores() int {
	var cores int
	var _ error
	if IsLinux() {
		cores, _ = strconv.Atoi(strings.ReplaceAll(run("grep \"^processor\" /proc/cpuinfo | wc -l"), "\n", ""))
	}
	return cores
}

func GetMemory() string {
	var memory string
	var _ error
	if IsLinux() {
		temp, _ := strconv.Atoi(strings.ReplaceAll(run("free -m | awk '$1 == \"Mem:\" {print $2}'"), "\n", ""))
		memory = fmt.Sprintf("%dGB", (temp / 1000))
	}
	return memory
}

func GetDisk() string {
	var disk string
	var _ error
	if IsLinux() {
		disk = strings.ReplaceAll(run("df -Bg | grep '^/dev/' | grep -v '/boot$' | awk '{ft += $2} END {print ft}'"), "\n", "") + "GB"
	}
	return disk
}