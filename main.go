package main

import (
	"fmt"
	"os/exec"
	"strings"
	"gopkg.in/yaml.v2"
	"github.com/spf13/cobra"
)

// todo: read from file
var data = `
screen:
  lark:
    fw-lib: ~/job/repo/fw-lib
  dmr:
    arch: ~/code/task/arch-setup/ansibook
    env: '~'
    hclpy: ~/code/python/terrapy
monitor:
  lark: 2x-w-lg.sh
  work: 2x-w-lg.sh
net:
  work: lark
  print_work: omg
`

type ScreenList map[string]string

type Config struct {
	Screen map[string]ScreenList
	Monitor map[string]string
	Net map[string]string
}

var rootCmd = &cobra.Command{
	Use: "env-knife-go",
	Short: "hello cobra",
	Long: "hello hello hello cobra",
}

func Parse(data string) Config {
	cfg := Config{}
	err := yaml.Unmarshal([]byte(data), &cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

func CurrentScreensList() []string {
	pipeline := "screen -ls | cut -f2 | head -n -1 | tail -n +2 | cut -d'.' -f2"
	c := exec.Command("/bin/sh", "-c", pipeline)
	out, err := c.Output()
	if err != nil {
		panic(err)
	}
	outstr := string(out)
	names := strings.Split(outstr, "\n")
	actualNames := names[:len(names)-1]
	return actualNames
}

func CurrentScreens() map[string]bool {
	ret := make(map[string]bool)
	names := CurrentScreensList()
	for _, n := range names {
		ret[n] = true
	}
	return ret
}

func MakeScreens(screens map[string]string) {
	names := CurrentScreens()
	for name, path := range screens {
		fmt.Printf("%s: %s\n", name, path)
		if names[name] {
			fmt.Printf("have %s\n", name)
		} else {
			fmt.Printf("need %s\n", name)
		}
	}
}

func MakeScreenRunner(screens map[string]string) func(*cobra.Command, []string) {
	return func(c *cobra.Command, args []string) {
		MakeScreens(screens)
	}
}

func MakeMonitorRunner(monitorScript string) func(*cobra.Command, []string) {
	return func(c *cobra.Command, args []string) {
		fmt.Println(monitorScript)
	}
}

func MakeNetRunner(ssid string) func(*cobra.Command, []string) {
	return func(c *cobra.Command, args []string) {
		fmt.Println(ssid)
	}
}

func Build(cfg Config) {
	if cfg.Screen != nil {
		screenCmd := &cobra.Command{
			Use: "screen",
			Short: "hello screen",
			Long: "hello hello hello screen",
		}
		for name, screens := range cfg.Screen {
			runner := MakeScreenRunner(screens)
			cmd := cobra.Command{
				Use: name,
				Aliases: []string{ name[0:1] },
				Short: fmt.Sprintf("%s screens", name),
				Long: fmt.Sprintf("%s screens: %s", name, screens),
				Run: runner,
			}
			screenCmd.AddCommand(&cmd)
		}
		rootCmd.AddCommand(screenCmd)
	}

	if cfg.Monitor != nil {
		monitorCmd := &cobra.Command{
			Use: "monitor",
			Short: "hello monitor",
			Long: "hello hello hello monitor",
		}
		for name, monitorScript := range cfg.Monitor {
			cmd := &cobra.Command{
				Use: name,
				Short: name,
				Long: name,
				Run: MakeMonitorRunner(monitorScript),
			}
			monitorCmd.AddCommand(cmd)
		}
		rootCmd.AddCommand(monitorCmd)
	}

	if cfg.Net != nil {
		netCmd := &cobra.Command{
			Use: "net",
			Short: "hello net",
			Long: "hello hello hello net",
		}
		for name, ssid := range cfg.Net {
			cmd := &cobra.Command{
				Use: name,
				Short: name,
				Long: name,
				Run: MakeNetRunner(ssid),
			}
			netCmd.AddCommand(cmd)
		}
		rootCmd.AddCommand(netCmd)
	}
}


func main() {
	cfg := Parse(data)
	Build(cfg)
	rootCmd.Execute()
}
