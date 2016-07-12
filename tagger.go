package hello

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/rpc/tag", tagger)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "This is the service which tags MLab pipeline data")
	fmt.Fprint(w, "You almost certainly mean to use the /rpc/tag url")
}

func errorMessage(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, err)
}

func tagger(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// The key for the switch discard bug lookup table will be a string,
	// constructed from other strings, so we don't need to parse and then
	// reformat the user input.  We do, however, need to validate the user
	// input before we use it.
	timestamp, timestamp_present := r.Form["timestamp"]
	client_ip, client_ip_present := r.Form["client_ip"]
	client_port, client_port_present := r.Form["client_port"]
	server_ip, server_ip_present := r.Form["server_ip"]
	server_port, server_port_present := r.Form["server_port"]

	if !timestamp_present || !client_ip_present || !client_port_present || !server_ip_present || !server_port_present || len(timestamp) != 1 || len(client_ip) != 1 || len(client_port) != 1 || len(server_ip) != 1 || len(server_port) != 1 {
		errorMessage(w, "Bad request - did you include exactly one of each of timestamp, client_ip, client_port, server_ip, server_port?")
		return
	}
	// To validate user input, try to parse and see if you got any errors
	if _, err := strconv.ParseInt(timestamp[0], 10, 64); err != nil {
		errorMessage(w, "bad timestamp")
		return
	}
	for _, port := range [...]string{client_port[0], server_port[0]} {
		parsed_port, err := strconv.ParseInt(port, 10, 16)
		if err != nil || parsed_port <= 0 || parsed_port > 32767 {
			errorMessage(w, "bad port")
			return
		}
	}
	for _, ip := range [...]string{client_ip[0], server_ip[0]} {
		parsed_ip := net.ParseIP(ip)
		if parsed_ip == nil {
			errorMessage(w, "bad ip")
			return
		}
	}
        // The input is now validated, and may be used safely
	w.WriteHeader(http.StatusOK)
        key := strings.Join([]string{client_ip[0], client_port[0], server_ip[0], server_port[0], timestamp[0]}, "#")
        fmt.Fprintf(w, "The key is: %s", key)
}
