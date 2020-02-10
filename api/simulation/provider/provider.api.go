package provider

import (
	fmt "fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	agent "github.com/synerex/synerex_alpha/api/simulation/agent"
	area "github.com/synerex/synerex_alpha/api/simulation/area"
)

var (
	providerMutex sync.RWMutex
)

func init() {
}

type Source struct {
	CmdName     string
	Type        ProviderType
	Cmd         *exec.Cmd
	Description string
	SrcDir      string
	BinName     string
	GoFiles     []string
	Options     []*Option
	SubFunc     func(pipe io.ReadCloser, name string)
}

////////////////////////////////////////////////////////////
////////////         Provider Class         ////////////////
///////////////////////////////////////////////////////////

func NewProvider(name string, providerType ProviderType) *Provider {
	uid, _ := uuid.NewRandom()
	p := &Provider{
		Id:   uint64(uid.ID()),
		Name: name,
		Type: providerType,
	}
	return p
}

func NewScenarioProvider(name string, scenario *ScenarioStatus) *Provider {
	uid, _ := uuid.NewRandom()
	p := &Provider{
		Id:   uint64(uid.ID()),
		Name: name,
		Type: ProviderType_SCENARIO,
	}
	p.WithScenarioStatus(scenario)
	return p
}

func (p *Provider) WithScenarioStatus(s *ScenarioStatus) *Provider {
	p.Data = &Provider_ScenarioStatus{s}
	return p
}

func NewClockProvider(name string, clock *ClockStatus) *Provider {
	uid, _ := uuid.NewRandom()
	p := &Provider{
		Id:   uint64(uid.ID()),
		Name: name,
		Type: ProviderType_CLOCK,
	}
	p.WithClockStatus(clock)
	return p
}

func (p *Provider) WithClockStatus(c *ClockStatus) *Provider {
	p.Data = &Provider_ClockStatus{c}
	return p
}

func NewVisualizationProvider(name string, vis *VisualizationStatus) *Provider {
	uid, _ := uuid.NewRandom()
	p := &Provider{
		Id:   uint64(uid.ID()),
		Name: name,
		Type: ProviderType_VISUALIZATION,
	}
	p.WithVisualizationStatus(vis)
	return p
}

func (p *Provider) WithVisualizationStatus(v *VisualizationStatus) *Provider {
	p.Data = &Provider_VisualizationStatus{v}
	return p
}

func NewAgentProvider(name string, agentType agent.AgentType, agent *AgentStatus) *Provider {
	uid, _ := uuid.NewRandom()
	p := &Provider{
		Id:   uint64(uid.ID()),
		Name: name,
		Type: ProviderType_AGENT,
	}
	p.WithAgentStatus(agent)
	return p
}

func (p *Provider) WithAgentStatus(a *AgentStatus) *Provider {
	p.Data = &Provider_AgentStatus{a}
	return p
}

func (p *Provider) Run(source *Source) error {
	log.Printf("Run '%s'\n", p.Name)

	cmd, err := createCmd(source)
	if err != nil {
		return err
	}
	go runMyCmd(cmd, source, p.Name)
	return nil
}

////////////////////////////////////////////////////////////
////////////            Status             ////////////////
///////////////////////////////////////////////////////////

func NewAgentStatus(areaInfo *area.Area, agentType agent.AgentType, agentNum uint64) *AgentStatus {
	as := &AgentStatus{
		Area:      areaInfo,
		AgentType: agentType,
		AgentNum:  agentNum,
	}
	return as
}

////////////////////////////////////////////////////////////
////////////         create Option Class          /////////
///////////////////////////////////////////////////////////

type Option struct {
	Key   string
	Value string
}

func NewProviderOptions(serverAddr string, nodeIdAddr string, providerJson string, scenarioProviderJson string) []*Option {
	o := []*Option{
		&Option{
			Key:   "server_addr",
			Value: serverAddr,
		},
		&Option{
			Key:   "nodeid_addr",
			Value: nodeIdAddr,
		},
		&Option{
			Key:   "provider_json",
			Value: providerJson,
		},
		&Option{
			Key:   "scenario_provider_json",
			Value: scenarioProviderJson,
		},
	}
	return o
}

////////////////////////////////////////////////////////////
////////////         Run Commands           ////////////////
///////////////////////////////////////////////////////////

type Log struct {
	ID          uint64
	Description string
}

func createCmd(source *Source) (*exec.Cmd, error) {

	d, err := os.Getwd()
	if err != nil {
		log.Printf("%s", err.Error())
		return nil, fmt.Errorf("cannot get dir: %s", err.Error())
	}

	// get src dir
	srcpath := filepath.FromSlash(filepath.ToSlash(d) + "/../../../" + source.SrcDir)
	binpath := filepath.FromSlash(filepath.ToSlash(d) + "/../../../" + source.SrcDir + "/" + source.BinName)
	//fi, err := os.Stat(binpath)
	_, err = os.Stat(binpath)

	// バイナリが最新かどうか
	modTime := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
	for _, fn := range source.GoFiles {
		sp := filepath.FromSlash(filepath.ToSlash(srcpath) + "/" + fn)
		ss, _ := os.Stat(sp)
		if ss.ModTime().After(modTime) {
			modTime = ss.ModTime()
		}
	}

	// 最新でない場合、run
	var cmd *exec.Cmd
	if err == nil { //&& fi.ModTime().After(modTime) { // check binary time
		cmdArgs := make([]string, 0)
		for _, option := range source.Options {
			cmdArgs = append(cmdArgs, "-"+option.Key)
			cmdArgs = append(cmdArgs, option.Value)
		}
		cmd = exec.Command("./"+source.BinName, cmdArgs...) // run binary
	} else {
		log.Printf("Error: [provider].go file isn't done build command\n")
		return nil, fmt.Errorf("Error: [provider].go file isn't done build command")
	}

	cmd.Dir = srcpath
	cmd.Env = getGoEnv()

	return cmd, nil
}

func runMyCmd(cmd *exec.Cmd, source *Source, name string) {

	pipe, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		log.Printf("Error for executing %s %v\n", cmd.Args[0], err)
		return
	}
	log.Printf("Starting %s..\n", cmd.Args[0])

	// run SubFuncition
	if source.Type != ProviderType_SYNEREX && source.Type != ProviderType_NODE_ID && source.Type != ProviderType_MONITOR {
		source.SubFunc(pipe, name)
	} else {
		for {

		}
	}

	log.Printf("[%s]:Now ending...", name)

	cmd.Wait()

	log.Printf("Command [%s] closed\n", name)
}

func getGoPath() string {
	env := os.Environ()
	for _, ev := range env {
		if strings.Contains(ev, "GOPATH=") {
			return ev
		}
	}
	return ""
}

func getGoEnv() []string { // we need to get/set gopath
	d, _ := os.Getwd() // may obtain dir of se-daemon
	gopath := filepath.FromSlash(filepath.ToSlash(d) + "/../../../../")
	absGopath, _ := filepath.Abs(gopath)
	env := os.Environ()
	newenv := make([]string, 0, 1)
	foundPath := false
	for _, ev := range env {
		if strings.Contains(ev, "GOPATH=") {
			// this might depends on each OS
			newenv = append(newenv, ev+string(os.PathListSeparator)+filepath.FromSlash(filepath.ToSlash(absGopath)+"/"))
			foundPath = true
		} else {
			newenv = append(newenv, ev)
		}
	}
	if !foundPath { // this might happen at in the daemon..
		gp := getGoPath()
		newenv = append(newenv, gp)
	}
	return newenv
}
