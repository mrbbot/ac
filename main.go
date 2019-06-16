package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getDeviceSerial() string {
	cmd := exec.Command("adb", "devices")
	out, _ := cmd.Output()
	devices := strings.Split(string(out), "\r\n")[1:]
	for _, device := range devices {
		if device != "" {
			serial := strings.Split(device, "\t")[0]
			if !strings.Contains(serial, "emulator") && !strings.Contains(serial, ":5555") {
				return serial
			}
		}
	}
	return ""
}

func getDeviceIP(serial string) string {
	cmd := exec.Command("adb", "-s", serial, "shell", "ifconfig", "wlan0")
	out, _ := cmd.Output()
	conf := string(out)
	index := strings.Index(conf, "inet addr:")
	if index == -1 {
		return ""
	}
	conf = conf[(index + len("inet addr:")):]
	index = strings.Index(conf, " ")
	if index == -1 {
		return ""
	}
	conf = conf[0:index]
	return conf
}

func enableTCPIP(serial string) {
	cmd := exec.Command("adb", "-s", serial, "tcpip", "5555")
	_ = cmd.Run()
	return
}

func connect(ip string) {
	cmd := exec.Command("adb", "connect", ip+":5555")
	_ = cmd.Run()
	return
}

func main() {
	serial := getDeviceSerial()
	if serial == "" {
		fmt.Println("Unable to find device to connect to!")
		os.Exit(1)
	}
	fmt.Printf("Found %s ", serial)
	ip := getDeviceIP(serial)
	fmt.Printf("at %s\nEnabling TCP/IP ", ip)
	enableTCPIP(serial)
	fmt.Println("and connecting...")
	connect(ip)
	fmt.Println("Done!")
}
