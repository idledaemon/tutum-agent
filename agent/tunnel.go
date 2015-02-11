package agent

import (
	"encoding/json"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/ActiveState/tail"
	"github.com/tutumcloud/tutum-agent/utils"
)

type TunnelPatchForm struct {
	Tunnel  string `json:"tunnel:"`
	Version string `json:"agent_version"`
}

func NatTunnel(url, ngrokPath, ngrokLogPath string) {
	if !utils.FileExist(ngrokPath) {
		Logger.Printf("Cannot find ngrok binary(%s), skipping NAT tunnel\n", ngrokPath)
		return
	}

	var cmd *exec.Cmd
	if *FlagNgrokToken != "" {
		Logger.Println("About to tunnel to public ngrok service")
		cmd = exec.Command(ngrokPath,
			"-log", "stdout",
			"-authtoken", *FlagNgrokToken,
			"-proto", "tcp",
			DockerHostPort)
	} else {
		Logger.Println("About to tunnel to private ngrok service")
		confPath := path.Join(TutumHome, NgrokConfName)
		if !utils.FileExist(confPath) {
			Logger.Println("Cannot find ngrok conf, skipping NAT tunnel")
			return
		}
		cmd = exec.Command(ngrokPath,
			"-config", confPath,
			"-log", "stdout",
			"-proto", "tcp",
			DockerHostPort)
	}

	os.RemoveAll(ngrokLogPath)
	logFile, err := os.OpenFile(ngrokLogPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Logger.Println(err)
	} else {
		defer logFile.Close()
		cmd.Stdout = logFile
	}

	Logger.Println("Starting montoring tunnel:", cmd.Args)
	go monitorTunnels(url, ngrokLogPath)
	Logger.Println("Starting NAT tunnel:", cmd.Args)

	for {
		runGronk(cmd)
		Logger.Println("Restarting NAT tunnel:", cmd.Args)
		time.Sleep(10 * time.Second)
	}
}

func runGronk(cmd *exec.Cmd) bool {
	if err := cmd.Start(); err != nil {
		return true
	}
	cmd.Wait()
	return true
}

func monitorTunnels(url, ngrokLogPath string) {
	update, _ := tail.TailFile(ngrokLogPath, tail.Config{
		Follow: true,
		ReOpen: true})
	for line := range update.Lines {
		if strings.Contains(line.Text, "[INFO] [client] Tunnel established at") {
			terms := strings.Split(line.Text, " ")
			tunnel := terms[len(terms)-1]
			Logger.Printf("Found new tunnel:%s\n", tunnel)
			patchTunnelToTutum(url, tunnel)
		}
	}
}

func patchTunnelToTutum(url, tunnel string) {
	Logger.Println("Patching tunnel address to Tutum")
	form := TunnelPatchForm{}
	form.Version = VERSION
	form.Tunnel = tunnel
	data, err := json.Marshal(form)
	if err != nil {
		Logger.Printf("Cannot marshal the TunnelPatch form:%s\f", err)
	}

	headers := []string{"Authorization TutumAgentToken " + Conf.TutumToken,
		"Content-Type", "application/json"}
	_, err = SendRequest("PATCH", utils.JoinURL(url, Conf.TutumUUID), data, headers)
	if err != nil {
		Logger.Println("Failed to patch tunnel address to Tutum,", err.Error())
	} else {
		Logger.Println("Successfully Patched tunnel address to Tutum")
	}
}

func DownloadNgrok(url, ngrokBinPath string) {
	if utils.FileExist(ngrokBinPath) {
		Logger.Printf("Found ngrok locally(%s), skip downloading\n", ngrokBinPath)
	} else {
		Logger.Println("No ngrok binary is found locally. Starting to download ngrok...")
		downloadFile(url, ngrokBinPath, "gnrok")
	}
}
