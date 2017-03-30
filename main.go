package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Version gets populated with version at build-time.
var Version string
var defaultPort = "8080"

var (
	port = kingpin.Flag("port", "Port to listen to.").Short('p').
		Default(defaultPort).String()
	bind = kingpin.Flag("bind", "Bind address.").Short('b').
		Default("0.0.0.0").String()
	version = kingpin.Flag("version", "Print version info.").
		Short('v').Bool()
)

func indexHandler(c *routing.Context) error {
	c.Write([]byte(
		"Case-insensitive string comparison, as an API. Because ¯\\_(ツ)_/¯\n" +
			"\n" +
			"Example:\n" +
			"curl -X POST -F \"a=Foo Bar\" -F \"b=FOO BAR\" " +
			"http://" + string(c.Host()) + "/",
	))
	return nil
}

func casecmpHandler(c *routing.Context) error {
	a := c.FormValue("a")
	b := c.FormValue("b")

	resp := "0"
	if strings.EqualFold(string(a), string(b)) {
		resp = "1"
	}

	c.Write([]byte(resp))
	return nil
}

func printVersion() {
	fmt.Println("casecmp " + Version)
}

func startServer() {
	r := routing.New()
	r.Get("/", indexHandler)
	r.Post("/", casecmpHandler)

	server := fasthttp.Server{Handler: r.HandleRequest}

	if *port == defaultPort {
		envPort := os.Getenv("PORT")
		if envPort != "" {
			*port = envPort
		}
	}

	address := *bind + ":" + *port
	fmt.Println("Listening on " + address)
	log.Fatal(server.ListenAndServe(address))
}

func main() {
	kingpin.Parse()

	if *version {
		printVersion()
	} else {
		startServer()
	}
}
