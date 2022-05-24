package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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
	forceHTTPSFlag = kingpin.Flag(
		"force-https", "Use https:// in example curl commands",
	).Bool()
	versionFlag = kingpin.Flag("version", "Print version info.").
			Short('v').Bool()
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	scheme := "http"
	if r.TLS != nil || *forceHTTPSFlag {
		scheme = "https"
	}

	_, err := fmt.Fprintf(w, `%s %s

Case-insensitive string comparison, as an API. Because ¯\_(ツ)_/¯

Example usage:
curl -X POST -F "a=Foo Bar" -F "b=FOO BAR" %s://%s/
curl -X GET "%s://%s/?a=Foo+Bar&b=FOO+BAR"
`,
		name, version, scheme, r.Host, scheme, r.Host)
	if err != nil {
		log.Fatal(err)
	}
}

func aboutHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w,
		`%s %s

https://github.com/jimeh/casecmp
`,
		name, version)
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
	_, err := fmt.Fprint(w, resp)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		if r.Method != "GET" || r.URL.RawQuery != "" {
			casecmpHandler(w, r)
			return
		}
		indexHandler(w, r)
	case "/about":
		aboutHandler(w, r)
	default:
		http.NotFound(w, r)
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
	if *portFlag == "" {
		envPort := os.Getenv("PORT")
		if envPort != "" {
			*portFlag = envPort
		} else {
			*portFlag = defaultPort
		}
	}

	if !*forceHTTPSFlag && os.Getenv("FORCE_HTTPS") != "" {
		*forceHTTPSFlag = true
	}

	address := *bindFlag + ":" + *portFlag
	fmt.Printf("Listening on %s\n", address)
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler:      http.HandlerFunc(handler),
		Addr:         address,
	}

	log.Fatal(srv.ListenAndServe())
}

func main() {
	kingpin.Parse()

	if *versionFlag {
		printVersion()
	} else {
		startServer()
	}
}
