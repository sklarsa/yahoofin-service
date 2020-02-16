package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sklarsa/yahoofin"
	"net/http"
	"strings"
	"time"
)

var yahooClient *yahoofin.Client

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/{ticker}", priceHandler)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)

}

func getClient() (*yahoofin.Client, error) {
	if yahooClient == nil {
		return yahoofin.NewClient()
	}
	return yahooClient, nil
}

func parseQsDate(val string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", val)
	if err != nil {
		return time.Time{}, fmt.Errorf("Invalid or missing date: '%v'", val)
	}
	return parsed, nil
}

func priceHandler(w http.ResponseWriter, r *http.Request) {
	client, err := getClient()
	if err != nil {
		http.Error(w, "Error retrieving yahoofin client", 500)
		return
	}

	query := r.URL.Query()

	errors := make([]string, 0)
	startDate, err := parseQsDate(query.Get("startDate"))
	if err != nil {
		errors = append(errors, fmt.Sprintf("startDate: %v", err.Error()))
	}

	endDate, err := parseQsDate(query.Get("endDate"))
	if err != nil {
		errors = append(errors, fmt.Sprintf("endDate: %v", err.Error()))
	}

	if len(errors) > 0 {
		http.Error(w, strings.Join(errors, "\n"), 400)
		return
	}

	vars := mux.Vars(r)
	data, err := client.GetSecurityDataString(vars["ticker"], startDate, endDate, yahoofin.History)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte(data))

}
