package provider

import (
	fmt "fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	providerMutex sync.RWMutex
)

type Source struct {
	CmdName     string
	Type        ProviderType
	Cmd         *exec.Cmd
	Description string
	SrcDir      string
	BinName     string
	GoFiles     []string
	Options     []*Option
	SubFunc     func(cmd *exec.Cmd)
}

////////////////////////////////////////////////////////////
////////////         Provider Class         ////////////////
///////////////////////////////////////////////////////////

func NewProvider(name string, providerType ProviderType) *Provider {
	p := &Provider{
		Id:   0,
		Name: name,
		Type: providerType,
	}
	return p
}

func NewScenarioProvider(name string, providerType ProviderType, scenario *Scenario) *Provider {
	p := &Provider{
		Id:   0,
		Name: name,
		Type: providerType,
	}
	p.WithScenario(scenario)
	return p
}

func (p *Provider) WithScenario(s *Scenario) *Provider {
	p.Data = &Provider_Scenario{s}
	return p
}

func NewClockProvider(name string, providerType ProviderType, clock *Clock) *Provider {
	p := &Provider{
		Id:   0,
		Name: name,
		Type: providerType,
	}
	p.WithClock(clock)
	return p
}

func (p *Provider) WithClock(c *Clock) *Provider {
	p.Data = &Provider_Clock{c}
	return p
}

func NewVisualizationProvider(name string, providerType ProviderType, vis *Visualization) *Provider {
	p := &Provider{
		Id:   0,
		Name: name,
		Type: providerType,
	}
	p.WithVisualization(vis)
	return p
}

func (p *Provider) WithVisualization(v *Visualization) *Provider {
	p.Data = &Provider_Visualization{v}
	return p
}

func NewAgentProvider(name string, providerType ProviderType, agent *Agent) *Provider {
	p := &Provider{
		Id:   0,
		Name: name,
		Type: providerType,
	}
	p.WithAgent(agent)
	return p
}

func (p *Provider) WithAgent(a *Agent) *Provider {
	p.Data = &Provider_Agent{a}
	return p
}

func (p *Provider) Run(source *Source) error {
	log.Printf("Run '%s'\n", p.Name)

	cmd, err := createCmd(source)
	if err != nil {
		return err
	}
	runMyCmd(cmd, source)
	return nil
}

////////////////////////////////////////////////////////////
////////////         Run Commands           ////////////////
///////////////////////////////////////////////////////////

type Log struct {
	ID          uint64
	Description string
}

type Option struct {
	Key   string
	Value string
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

func runMyCmd(cmd *exec.Cmd, source *Source) {

	err := cmd.Start()
	if err != nil {
		log.Printf("Error for executing %s %v\n", cmd.Args[0], err)
		return
	}
	log.Printf("Starting %s..\n", cmd.Args[0])

	// run SubFuncition
	source.SubFunc(cmd)

	log.Printf("[%s]:Now ending...", source.CmdName)

	cmd.Wait()

	log.Printf("Command [%s] closed\n", source.CmdName)
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
