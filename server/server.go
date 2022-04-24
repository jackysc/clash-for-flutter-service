package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/csj8520/clash-for-flutter-service/constant"
	"github.com/gorilla/websocket"
)

var (
	StatusRunning string = "running"
	StatusStopped string = "stopped"
)

var (
	Server      *http.Server
	MaxLog      int                        = 1000
	Status      string                     = StatusStopped
	Logs        []string                   = make([]string, 0)
	logsClients map[string]*websocket.Conn = make(map[string]*websocket.Conn)
	Cmd         *exec.Cmd
	upgrader    websocket.Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type startArgs struct {
	Args []string `json:"args"`
}

func StartServer() error {
	Server = &http.Server{Addr: "127.0.0.1:9089"}
	http.HandleFunc("/start", handleReqStart)
	http.HandleFunc("/stop", handleReqStop)
	http.HandleFunc("/info", hanldeReqInfo)
	http.HandleFunc("/logs", hanldeReqLogs)
	err := Server.ListenAndServe()
	return err
}

func checkRequest(w http.ResponseWriter, r *http.Request) (reject bool) {
	reject = !strings.Contains((r.Header.Get("User-Agent")), "clash-for-flutter/")
	if reject {
		w.WriteHeader(403)
	}
	return reject
}

func FileIsExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func hanldeReqInfo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	reject := checkRequest(w, r)

	if reject {
		return
	}

	res := make(map[string]interface{})
	res["code"] = 0
	res["status"] = Status
	res["mode"] = os.Args[1]
	res["version"] = constant.Version
	json, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(json))
}

func hanldeReqLogs(w http.ResponseWriter, r *http.Request) {
	reject := checkRequest(w, r)
	if reject {
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Print("upgrade:", err)
		return
	}

	logsClients[r.RemoteAddr] = c
	fmt.Println("Connect Success:", r.RemoteAddr, "Connect Count:", len(logsClients))

	c.WriteMessage(websocket.TextMessage, []byte("Connect Success"))
	if len(Logs) > 0 {
		c.WriteMessage(websocket.TextMessage, []byte(strings.Join(Logs, "\n")))
	}

	for {
		_, msg, err := c.ReadMessage()

		if string(msg) == "PING" {
			c.WriteMessage(websocket.TextMessage, []byte("PONG"))
		}

		if err != nil {
			break
		}
	}

	c.Close()
	delete(logsClients, r.RemoteAddr)
	fmt.Println("UnConnect:", r.RemoteAddr, "Connect Count:", len(logsClients))
}

func handleReqStop(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	reject := checkRequest(w, r)
	if reject {
		return
	}
	res := make(map[string]interface{})
	if Status == StatusRunning && Cmd != nil {
		Cmd.Process.Kill()
		Cmd.Process.Wait()
		Status = StatusStopped
		res["code"] = 0
	} else {
		res["code"] = 1
		res["msg"] = "clash core is not running"
	}
	json, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func handleReqStart(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	reject := checkRequest(w, r)
	if reject {
		return
	}
	res := make(map[string]interface{})

	if Status == StatusRunning {
		res["code"] = 1
		res["msg"] = "clash core is running"
	} else {
		body, _ := ioutil.ReadAll(r.Body)
		parmas := &startArgs{}
		json.Unmarshal(body, parmas)
		startClashCore(parmas.Args)
		res["code"] = 0
	}
	json, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func startClashCore(args []string) {
	file, _ := os.Executable()
	path := filepath.Dir(file)

	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	name := fmt.Sprintf("clash-%s-%s%s", runtime.GOOS, runtime.GOARCH, ext)
	clashCorePath := filepath.Join(path, name)
	EchoLog("Clash Core Path Is: ", clashCorePath)
	if !FileIsExist(clashCorePath) {
		// for dev
		d, _ := os.Getwd()
		EchoLog("pwd is: ", d)
		clashCorePath = filepath.Join(d, name)
		if !FileIsExist(clashCorePath) {
			EchoLog("Clash Core Is Not Exist")
			return
		}
		EchoLog("Clash Core Path Is: ", clashCorePath)
	}

	EchoLog("Args Is: ", args)

	Cmd = exec.Command(clashCorePath, args...)
	stdout, _ := Cmd.StdoutPipe()
	stderr, _ := Cmd.StderrPipe()
	go listenLog(&stderr)
	go listenLog(&stdout)
	Cmd.Start()
	Status = StatusRunning
	go func() {
		s, _ := Cmd.Process.Wait()
		Status = StatusStopped
		Cmd = nil
		EchoLog("Clash Core exit with code: ", s.ExitCode())
	}()
}

func listenLog(i *io.ReadCloser) {
	reader := bufio.NewReader(*i)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		EchoLog(strings.TrimSpace(line))
		// fmt.Printf("len=%d cap=%d slice=%v\n", len(Logs), cap(Logs), Logs)
	}
}

func EchoLog(logs ...interface{}) {
	log := fmt.Sprint(logs...)
	Logs = append(Logs, log)
	if len(Logs) > MaxLog {
		Logs = Logs[1:]
	}
	for _, c := range logsClients {
		c.WriteMessage(websocket.TextMessage, []byte(log))
	}
	fmt.Println(log)
}
