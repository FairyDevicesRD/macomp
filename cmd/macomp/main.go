package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"

	"github.com/FairyDevicesRD/macomp"
	"github.com/fatih/color"
	flags "github.com/jessevdk/go-flags"
	runewidth "github.com/mattn/go-runewidth"
)

func handleWidths(mr *macomp.MaResource) (map[string]string, int) {
	name2fname := map[string]string{}
	max := 0
	for name := range mr.Mecabs {
		name2fname[name] = name
		l := runewidth.StringWidth(name)
		if l > max {
			max = l
		}
	}
	for name := range mr.Jumans {
		name2fname[name] = name
		l := runewidth.StringWidth(name)
		if l > max {
			max = l
		}
	}

	for name := range name2fname {
		name2fname[name] = runewidth.FillRight(name, max)
	}

	return name2fname, max
}

func listup(settings map[string]macomp.MaSetting) {
	outs := make([]string, 0, len(settings))
	for orgname, setting := range settings {
		var alias string
		if len(setting.Aliases) != 0 {
			alias = fmt.Sprintf(" (%s)", strings.Join(setting.Aliases, ", "))
		}
		outs = append(outs, orgname+alias+"\n")
	}
	sort.Strings(outs)
	for _, out := range outs {
		fmt.Print(out)
	}
}

func initConfig(opts *cmdOptions) error {
	if len(opts.MaConfigFile) == 0 {
		return errors.New("Undesignated output path")
	}
	if _, err := os.Stat(opts.MaConfigFile); err == nil {
		return fmt.Errorf("Already exist: %s", opts.MaConfigFile)
	}

	settings := map[string]macomp.MaSetting{}
	if mc, err := exec.LookPath("mecab-config"); err == nil {
		if dicpath, err := exec.Command(mc, "--dicdir").Output(); err == nil {
			dp := strings.Replace(string(dicpath), "\n", "", 1)
			fileInfos, err := ioutil.ReadDir(dp)
			if err == nil {
				for _, fi := range fileInfos {
					setting := macomp.MaSetting{
						MaType:  "mecab",
						Aliases: []string{},
						Options: map[string]string{
							"dicdir": path.Join(dp, fi.Name()),
						},
					}
					settings[fi.Name()] = setting
				}
			}
		}
	}
	for _, name := range []string{"juman", "jumanpp"} {
		if path, err := exec.LookPath(name); err == nil {
			settings[name] = macomp.MaSetting{
				MaType:  name,
				Path:    path,
				Aliases: []string{},
				Options: map[string]string{},
			}
		}
	}

	if len(settings) == 0 {
		return errors.New("Morphological analyzers are not found")
	}

	b, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(opts.MaConfigFile, b, 0644)
}

func operateLine(resource *macomp.MaResource, text string, pos, check bool, name2fname map[string]string, maxw int) {
	text = strings.Replace(text, " ", "　", -1)
	plainText := text
	if check {
		plainText = strings.Replace(plainText, "|", "", -1)
		plainText = strings.Replace(plainText, "?", "", -1)
	}
	if len(plainText) == 0 {
		return
	}
	results := resource.Parse(plainText)

	for _, result := range results {
		if check {
			pl, isValid := macomp.DecorateSurfaces(text, result.Surfaces)
			if isValid {
				fmt.Print(color.GreenString("O "))
			} else {
				fmt.Print(color.RedString("X "))
			}
			fmt.Print(name2fname[result.Name])
			fmt.Print(" ")
			fmt.Print(pl)
		} else {
			fmt.Print(name2fname[result.Name])
			fmt.Print(" ")
			pl := macomp.PrettySurfaces(result.Surfaces)
			fmt.Print(pl)
		}
		fmt.Print("\n")

		if pos {
			len := maxw + 1
			if check {
				len += 2
			}
			for i := 0; i < len; i++ {
				fmt.Print(" ")
			}
			pl2 := macomp.PrettyFeatures(result.Surfaces, result.Features)
			fmt.Print(pl2)
			fmt.Print("\n")
		}
	}
}

func operation(opts *cmdOptions) error {
	var resource *macomp.MaResource

	if opts.Init {
		return initConfig(opts)
	}

	if c, err := ioutil.ReadFile(opts.MaConfigFile); err == nil {
		var settings map[string]macomp.MaSetting
		if err := json.Unmarshal(c, &settings); err != nil {
			return err
		}

		if opts.List { //Just list up
			listup(settings)
			return nil
		}

		aliases := map[string]string{}
		for orgname, setting := range settings {
			for _, alias := range setting.Aliases {
				aliases[alias] = orgname
			}
		}
		if opts.Targets != nil && len(opts.Targets) != 0 {
			nsettings := map[string]macomp.MaSetting{}
			for _, name := range opts.Targets {
				if orgname, ok := aliases[name]; ok {
					name = orgname //replace
				}
				s, ok := settings[name]
				if !ok {
					return fmt.Errorf("Unknown name: %s", name)
				}
				nsettings[name] = s
			}
			settings = nsettings
		}
		if resource, err = macomp.NewMaResource(settings); err != nil {
			return err
		}
	} else {
		return err
	}
	defer resource.Destroy()

	name2fname, maxw := handleWidths(resource)

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		text := s.Text()
		operateLine(resource, text, opts.Pos, opts.Check, name2fname, maxw)
	}
	return s.Err()
}

type cmdOptions struct {
	MaConfigFile string   `short:"m" long:"ma"  description:"MA Config File"`
	Pos          bool     `long:"pos" description:"Show POS"`
	Check        bool     `short:"c" long:"check" description:"Check partial annotation constraint"`
	Targets      []string `short:"t" long:"target"  description:"Target names to execute"`
	List         bool     `long:"list"  description:"List available targets"`
	ForceColor   bool     `short:"C" long:"color" description:"Force color output"`
	Init         bool     `long:"init"  description:"Initialize the setting file"`
}

func main() {
	opts := cmdOptions{}
	optparser := flags.NewParser(&opts, flags.Default)
	optparser.Name = ""
	optparser.Usage = ""
	_, err := optparser.Parse()

	//show help
	if err != nil {
		for _, arg := range os.Args {
			if arg == "-h" || arg == "--help" {
				macomp.PrintDefaultPath()
				os.Exit(0)
			}
		}
		os.Exit(1)
	}

	if opts.ForceColor {
		color.NoColor = false
	}

	//Get config path
	if len(opts.MaConfigFile) == 0 {
		opts.MaConfigFile = macomp.GetConfigPath()
	}

	if err := operation(&opts); err != nil {
		log.Fatal(err)
	}
}
