package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/csj8520/clash-for-flutter-service/constant"
	"github.com/csj8520/clash-for-flutter-service/server"

	"github.com/kardianos/service"
)

var Service service.Service

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	fmt.Println("Start success")
	return nil
}

func (p *program) run() {
	err := server.StartServer()
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), "http: Server closed") {
			return
		}
		if strings.Contains(err.Error(), "bind: address already in use") {
			if os.Args[1] == "service-mode" {
				Service.Stop()
			} else {
				p.Stop(Service)
				os.Exit(101)
			}
		}

	}
}

func (p *program) Stop(s service.Service) error {
	if server.Server != nil {
		server.Server.Shutdown(context.TODO())
	}
	if server.Cmd != nil {
		server.Cmd.Process.Kill()
		server.Cmd.Process.Wait()
	}
	fmt.Println("Stop success")
	return nil
}

func main() {

	if len(os.Args) <= 1 {
		fmt.Println("Please use command: install, uninstall, status, start, stop, restart, version, service-mode, user-mode")
		return
	}

	svcConfig := &service.Config{
		Name:        "clash-for-flutter-service",
		DisplayName: "Clash For Flutter Service",
		Description: "This is a Clash For Flutter Service.",
		Arguments:   []string{"service-mode"},
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	Service = s
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, it := range os.Args[1:] {
		switch it {
		case "install":
			handleInstall(s)
		case "uninstall":
			handleStop(s)
			handleUnInstall(s)
		case "status":
			handleStatus(s)
		case "start":
			handleStart(s)
		case "stop":
			handleStop(s)
		case "restart":
			handleRestart(s)
		case "version":
			handleVersion(s)
		case "service-mode":
			handleRun(s)
		case "user-mode":
			handleRun(s)
		default:
			fmt.Println("Command does not exist")
		}
	}

}

func handleInstall(s service.Service) {
	err := s.Install()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Install Service Success")
	}
}

func handleUnInstall(s service.Service) {
	err := s.Uninstall()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("UnInstall Service Success")
	}
}
func handleStatus(s service.Service) {
	status, err := s.Status()
	if err != nil {
		fmt.Println(err)
	} else {
		switch status {
		case service.StatusRunning:
			fmt.Println("Service Status is runing")
		case service.StatusStopped:
			fmt.Println("Service Status is stoped")
		default:
			fmt.Println("Service Status is unknow")
		}
	}
}
func handleStart(s service.Service) {
	err := s.Start()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Service Start Success ")
	}
}
func handleStop(s service.Service) {
	err := s.Stop()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Service Stop Success")
	}
}
func handleRestart(s service.Service) {
	err := s.Restart()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Service Restart Success")
	}
}
func handleVersion(s service.Service) {
	fmt.Println(constant.Version)
}
func handleRun(s service.Service) {
	err := s.Run()
	if err != nil {
		fmt.Println(err)
	}
}
