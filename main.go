package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	version = "dev"
	commit  = "unknown"
)

const (
	name        = "casecmp"
	defaultPort = 8080
	defaultBind = "0.0.0.0"
)

// Argument parsing setup.
var (
	portFlag       = flag.Int("p", defaultPort, "Port to listen on")
	bindFlag       = flag.String("b", defaultBind, "Bind address")
	forceHTTPSFlag = flag.Bool(
		"f", false, "Use https:// in example curl commands",
	)
	versionFlag = flag.Bool("v", false, "Print version info")
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

func startServer() error {
	if *portFlag == defaultPort {
		envPort := os.Getenv("PORT")
		if envPort != "" {
			v, err := strconv.Atoi(envPort)
			if err != nil {
				return err
			}

			*portFlag = v
		}
	}

	if !*forceHTTPSFlag && os.Getenv("FORCE_HTTPS") != "" {
		*forceHTTPSFlag = true
	}

	address := fmt.Sprintf("%s:%d", *bindFlag, *portFlag)
	fmt.Printf("Listening on %s\n", address)

	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler:      http.HandlerFunc(handler),
		Addr:         address,
	}

	return srv.ListenAndServe()
}

func main() {
	flag.Parse()

	if *versionFlag {
		printVersion()
		return
	}

	err := startServer()
	if err != nil {
		log.Fatal(err)
	}
}
