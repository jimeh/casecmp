package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	name        = "casecmp"
	version     = "dev"
	commit      = "unknown"
	date        = "unknown"
	defaultPort = "8080"
)

// Argument parsing setup.
var (
	portFlag = kingpin.Flag("port", "Port to listen to.").Short('p').
			Default("").String()
	bindFlag = kingpin.Flag("bind", "Bind address.").Short('b').
			Default("0.0.0.0").String()
	versionFlag = kingpin.Flag("version", "Print version info.").
			Short('v').Bool()
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	_, err := fmt.Fprintf(w, `%s %s

Case-insensitive string comparison, as an API. Because ¯\_(ツ)_/¯

Example usage:
curl -X POST -F "a=Foo Bar" -F "b=FOO BAR" %s://%s/
curl -X POST "%s://%s/?a=Foo+Bar&b=FOO+BAR"`,
		name, version, scheme, r.Host, scheme, r.Host)

	if err != nil {
		log.Fatal(err)
	}
}

func casecmpHandler(w http.ResponseWriter, r *http.Request) {
	a := r.FormValue("a")
	b := r.FormValue("b")
	resp := "0"

	if strings.EqualFold(string(a), string(b)) {
		resp = "1"
	}
	_, err := fmt.Fprintf(w, resp)

	if err != nil {
		log.Fatal(err)
	}
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
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s %s", name, version))

	if commit != "unknown" {
		buffer.WriteString(fmt.Sprintf(" (%s)", commit))
	}

	fmt.Println(buffer.String())
}

func startServer() {
	http.HandleFunc("/", rootHandler)

	if *portFlag == "" {
		envPort := os.Getenv("PORT")
		if envPort != "" {
			*portFlag = envPort
		} else {
			*portFlag = defaultPort
		}
	}

	address := *bindFlag + ":" + *portFlag
	fmt.Println("Listening on " + address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func main() {
	kingpin.Parse()

	if *versionFlag {
		printVersion()
	} else {
		startServer()
	}
}
