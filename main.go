package main

import (
	"fmt"
	"os"

	"github.com/csj8520/clash-for-flutter-service/constant"
	"github.com/csj8520/clash-for-flutter-service/server"

	"github.com/kardianos/service"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	go server.StartServer()
	fmt.Println("Start success")
	return nil
}

func (p *program) Stop(s service.Service) error {
	if server.Cmd != nil {
		server.Cmd.Process.Kill()
	}
	fmt.Println("Stop success")
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "clash-for-flutter-service",
		DisplayName: "Clash For Flutter Service",
		Description: "This is a Clash For Flutter Service.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(os.Args) > 1 {
		for _, it := range os.Args[1:] {
			switch it {
			case "install":
				handleInstall(s)
			case "uninstall":
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
			default:
				fmt.Println("Command does not exist")
			}
		}

	} else {
		handleRun(s)
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
