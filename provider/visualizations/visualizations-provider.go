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

func runVisProviders() {

	// run nodeserver
	nodecmd := exec.Command("./../../nodeserv/nodeid-server")
	go manager.run(nodecmd, "node")
	// run synerexserver
	sycmd := exec.Command("./../../server/synerex-server")
	go manager.run(sycmd, "synerex")
	// run visualization provider
	wocmd := exec.Command("./../visualization/visualization-provider")
	go manager.run(wocmd, "visualization")
}

func main() {
	fmt.Printf("starting visualizetions provider...\n")
	runVisProviders()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
