package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/FairyDevicesRD/macomp"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	flags "github.com/jessevdk/go-flags"
)

var resource *macomp.MaResource

func operation(opts *cmdOptions) error {
	if _, err := os.Stat(opts.StaticRoot); err != nil {
		return err
	}

	if c, err := ioutil.ReadFile(opts.MaConfigFile); err == nil {
		var settings map[string]macomp.MaSetting
		if err := json.Unmarshal(c, &settings); err != nil {
			return err
		}
		if resource, err = macomp.NewMaResource(settings); err != nil {
			return err
		}
	} else {
		return err
	}
	defer resource.Destroy()

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(opts.StaticRoot))))
	showTop := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(opts.StaticRoot, "index.html"))
	}
	r.HandleFunc("/", showTop)
	r.HandleFunc("/{query}", showTop)
	r.HandleFunc("/api/v1/ma", DoMA)
	r.HandleFunc("/api/v1/ma/{text}", DoMA)

	log.Printf("Server started at [%s] from pid %d", opts.Bind, os.Getpid())

	loggedRouter := handlers.CombinedLoggingHandler(os.Stdout, r)

	// Wrap our server with our gzip handler to gzip compress ALL responses.
	// TODO: exclude some binary files like image files
	gh := handlers.CompressHandler(loggedRouter)

	gracehttp.Serve(&http.Server{Addr: opts.Bind, Handler: gh})
	return nil
}

type cmdOptions struct {
	Bind         string `short:"b" long:"bind" default:":5000" description:"String to bind"`
	MaConfigFile string `short:"m" long:"ma"  description:"MA Config File"`
	StaticRoot   string `long:"static" default:"" description:"The path to the static directory"`
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
			if arg == "-h" {
				macomp.PrintDefaultPath()
				os.Exit(0)
			}
		}
		os.Exit(1)
	}

	//Get config path
	if len(opts.MaConfigFile) == 0 {
		opts.MaConfigFile = macomp.GetConfigPath()
	}

	if len(opts.StaticRoot) == 0 {
		opts.StaticRoot = macomp.GetStaticRootPath()
	}

	if err := operation(&opts); err != nil {
		log.Fatal(err)
	}
}
