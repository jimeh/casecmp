package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"text/template"
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

type IndexData struct {
	Name    string
	Version string
	Scheme  string
	Host    string
}

var indexTpl = template.Must(template.New("index").Parse(`{{.Name}} {{.Version}}

Case-insensitive string comparison, as an API. Because ¯\_(ツ)_/¯

Example usage:
curl -X POST -F "a=Foo Bar" -F "b=FOO BAR" {{.Scheme}}://{{.Host}}/
curl -X GET "{{.Scheme}}://{{.Host}}/?a=Foo+Bar&b=FOO+BAR"
curl -X GET -H "Accept: application/json" "{{.Scheme}}://{{.Host}}/?a=Foo+Bar&b=FOO+BAR"
curl -X POST -H "Content-Type: application/json" -d '{"a":"Foo Bar","b":"FOO BAR"}' {{.Scheme}}://{{.Host}}/
`))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	scheme := "http"
	if r.TLS != nil || *forceHTTPSFlag {
		scheme = "https"
	}

	err := indexTpl.Execute(w, &IndexData{
		Name:    name,
		Version: version,
		Scheme:  scheme,
		Host:    r.Host,
	})
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

type JSONData struct {
	A string `json:"a"`
	B string `json:"b"`
}

func casecmpHandler(w http.ResponseWriter, r *http.Request) error {
	var a, b string

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}

		d := JSONData{}
		err = json.Unmarshal(body, &d)
		if err != nil {
			return err
		}

		a = d.A
		b = d.B
	} else {
		a = r.FormValue("a")
		b = r.FormValue("b")
	}

	resp := "0"
	if strings.EqualFold(string(a), string(b)) {
		resp = "1"
	}

	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, err := fmt.Fprintf(w, `{"result":%s}`, resp)
		return err
	}

	_, err := fmt.Fprint(w, resp)
	return err
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		if r.Method != "GET" || r.URL.RawQuery != "" {
			err := casecmpHandler(w, r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprint(w, err.Error())
			}
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

	buffer.WriteString(fmt.Sprintf(", built with %s", runtime.Version()))

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

	if *bindFlag == defaultBind {
		envBind := os.Getenv("BIND")
		if envBind != "" {
			*bindFlag = envBind
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
