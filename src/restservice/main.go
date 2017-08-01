package main

// Inspired by this blogpost https://appliedgo.net/rest/

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	// This is httprouter. Ensure to install it first via go get.
	"github.com/julienschmidt/httprouter"
)

type store struct {
	// We need a data store. For our purposes, a simple map from string to string is completely sufficient.
	data map[string]string

	// Handlers run concurrently, and maps are not thread-safe. This mutex is used to ensure that only one goroutine can update data.
	m sync.RWMutex
}

var (
	// We need a flag for setting the listening address. We set the default to port 8000, which is a common HTTP port for servers with local-only access.
	//addr = flag.String("addr", ":8000", "http service address")

	// Now we create the data store.
	s = store{
		data: map[string]string{},
		m:    sync.RWMutex{},
	}
)

func main() {
	// Read the env variable REST_PORT from the OS.
	listentPort := os.Getenv("REST_PORT")

	// If its empty, set a default value of 8000
	if len(listentPort) == 0 {
		listentPort = "8000"
		fmt.Print("Warn: REST_PORT not set. Using default port : " + listentPort + "\n")
	}

	// Printing information to screen
	fmt.Print("REST service started using port: ", listentPort, "...\n")

	// The main function starts by parsing the commandline.
	flag.Parse()

	// Now we can create a new httprouter instance
	r := httprouter.New()

	// and add some routes. httprouter provides functions named after HTTP verbs.
	// So to create a route for HTTP GET, we simply need to call the GET function
	// and pass a route and a handler function. The first route is /entry followed
	// by a key variable denoted by a leading colon. The handler function is set to show.
	r.GET("/entry/:key", show)

	// We do the same for /list. Note that we use the same handler function here;
	// we’ll switch functionality within the show function based on the existence of a key variable.
	r.GET("/list", show)

	// For updating, we need a PUT operation. We want to pass a key and a value to the URL,
	// so we add two variables to the path. The handler function for this PUT operation is update.
	r.PUT("/entry/:key/:value", update)

	// A simple webpage
	r.GET("/", Index)

	// Finally, we just have to start the http Server. We pass the listening address as well as our router instance.
	err := http.ListenAndServe(*flag.String("addr", ":"+listentPort, "http service address"), r)
	//err := http.ListenAndServe(*addr, r)

	// For this demo, let’s keep error handling simple. log.Fatal prints out an error message and exits the process.
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func show(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	// To access these parameters, we call the ByName method, passing the variable name that we chose when defining the route in main
	k := p.ByName("key")

	// The show function serves two purposes. If there is no key in the URL, it lists all entries of the data map.
	if k == "" {
		// Lock the store for reading.
		s.m.RLock()
		fmt.Fprintf(w, "Read list: %v", s.data)
		s.m.RUnlock()
		return
	}

	// If a key is given, the show function returns the corresponding value. It does so by simply printing to the ResponseWriter parameter, which is sufficient for our purposes.
	s.m.RLock()
	fmt.Fprintf(w, "Read entry: s.data[%s] = %s", k, s.data[k])
	s.m.RUnlock()
}

func update(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Fetch key and value from the URL parameters.
	k := p.ByName("key")
	v := p.ByName("value")

	// We just need to either add or update the entry in the data map.
	s.m.Lock()
	s.data[k] = v
	s.m.Unlock()

	// Finally, we print the result to the ResponseWriter.
	fmt.Fprintf(w, "Updated: s.data[%s] = %s", k, v)

}

func Index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprint(w, "<html><body><p>Welcome!</p></body></html>")
}
