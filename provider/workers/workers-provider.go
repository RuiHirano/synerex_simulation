package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"

	"github.com/synerex/synerex_alpha/util"
)

var (
	logger  *util.Logger
	manager *Manager
)

func init() {
	manager = NewManager()
	logger = util.NewLogger()
}

type Manager struct {
	Cmds map[int]*exec.Cmd
}

func NewManager() *Manager {
	mn := &Manager{
		Cmds: make(map[int]*exec.Cmd),
	}
	return mn
}

func (mn *Manager) run(cmd *exec.Cmd, name string) {

	logger.Info("[%s] Cmd start!\n", name)
	pipe, err := cmd.StderrPipe()
	if err != nil {
		logger.Info("Error for getting stdout pipe %v\n", err)
		return
	}
	err = cmd.Start()
	if err != nil {
		logger.Info("Error for executing %v\n", err)
		return
	}
	s := bufio.NewScanner(pipe)
	for s.Scan() {
		line := s.Text() // output line 出力一行
		logger.Info("[%s] %s\n", name, line)
	}

	mn.Cmds[cmd.Process.Pid] = cmd
	cmd.Wait()

	logger.Info("[%s] Cmd closed\n", name)

}

func runWorkerProviders() {
	// test
	cmd := exec.Command("ls")
	go manager.run(cmd, "test")

	// run nodeserver
	nodecmd := exec.Command("go", "run", "./../../nodeserv/nodeid-server.go")
	go manager.run(nodecmd, "node")
	// run synerexserver
	sycmd := exec.Command("go", "run", "./../../server/synerex-server.go", "./../../server/message-store.go")
	go manager.run(sycmd, "synerex")
	// run worker provider
	wocmd := exec.Command("go", "run", "./../worker/worker-provider.go")
	go manager.run(wocmd, "worker")
	// run agent provider
	agcmd := exec.Command("go", "run", "./../agent/agent-provider.go", "./../agent/simulator.go")
	go manager.run(agcmd, "agent")
}

func main() {
	fmt.Printf("starting workers provider...\n")
	runWorkerProviders()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
