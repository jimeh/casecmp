package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Name of application.
var Name = "casecmp"

// Version gets populated with version at build-time.
var Version string

// DefaultPort that service runs on.
var DefaultPort = "8080"

// Argument parsing setup.
var (
	port = kingpin.Flag("port", "Port to listen to.").Short('p').
		Default("").String()
	bind = kingpin.Flag("bind", "Bind address.").Short('b').
		Default("0.0.0.0").String()
	version = kingpin.Flag("version", "Print version info.").
		Short('v').Bool()
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	resp := Name + " " + Version + "\n" +
		"\n" +
		"Case-insensitive string comparison, as an API. Because ¯\\_(ツ)_/¯\n" +
		"\n" +
		"Example usage:\n" +
		"curl -X POST -F \"a=Foo Bar\" -F \"b=FOO BAR\" " +
		"http://" + r.Host + "/\n" +
		"curl -X POST http://" + r.Host + "/?a=Foo%%20Bar&b=FOO%%20BAR"

	io.WriteString(w, resp)
}

func casecmpHandler(w http.ResponseWriter, r *http.Request) {
	a := r.FormValue("a")
	b := r.FormValue("b")

	resp := "0"
	if strings.EqualFold(string(a), string(b)) {
		resp = "1"
	}
	fmt.Fprintf(w, resp)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "GET" {
		indexHandler(w, r)
	} else {
		casecmpHandler(w, r)
	}
}

func printVersion() {
	fmt.Println(Name + " " + Version)
}

func startServer() {
	http.HandleFunc("/", rootHandler)

	if *port == "" {
		envPort := os.Getenv("PORT")
		if envPort != "" {
			*port = envPort
		} else {
			*port = DefaultPort
		}
	}

	address := *bind + ":" + *port
	fmt.Println("Listening on " + address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func main() {
	kingpin.Parse()

	if *version {
		printVersion()
	} else {
		startServer()
	}
}
