package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/miekg/dns"
)

// getSOA retrieves the SOA record for a given domain
func getSOA(domain string) (string, error) {
	// Set up a DNS client
	client := dns.Client{}

	// Create a message to send to the server
	msg := dns.Msg{}
	msg.SetQuestion(dns.Fqdn(domain), dns.TypeSOA)

	// Send the message and get the response
	r, _, err := client.Exchange(&msg, net.JoinHostPort("8.8.8.8", "53"))
	if err != nil {
		return "", err
	}

	// Check the response to see if it contains an SOA record
	if len(r.Answer) < 1 {
		return "", fmt.Errorf("No SOA record found for domain %s", domain)
	}
	soa, ok := r.Answer[0].(*dns.SOA)
	if !ok {
		return "", fmt.Errorf("No SOA record found for domain %s", domain)
	}

	return soa.String(), nil
}

// watchSOA watches the SOA for a given domain and sends a message
// if the SOA changes
func watchSOA(w http.ResponseWriter, r *http.Request) {
	// Get the domain from the request parameters
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		http.Error(w, "Missing domain parameter", http.StatusBadRequest)
		return
	}

	// Get the initial SOA value
	soa, err := getSOA(domain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set up a timer to check the SOA every 10 minutes
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	// Loop indefinitely, checking the SOA every 10 minutes
	for range ticker.C {
		newSOA, err := getSOA(domain)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// If the SOA has changed, send a message
		if soa != newSOA {
			soa = newSOA
			// Send message indicating that the SOA has changed
			fmt.Fprintf(w, "SOA for domain %s has changed to %s", domain, soa)
		}
	}
}

func main() {
	http.HandleFunc("/watch", watchSOA)
	http.ListenAndServe(":8080", nil)
}
