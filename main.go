package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sklarsa/yahoofin"
	"net/http"
	"strings"
)

var yahooClient *yahoofin.Client

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/{ticker}/{field}", priceHandler)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)

}

func getClient() (*yahoofin.Client, error) {
	if yahooClient == nil {
		return yahoofin.NewClient()
	}
	return yahooClient, nil
}

func priceHandler(w http.ResponseWriter, r *http.Request) {
	_, err := getClient()
	if err != nil {
		http.Error(w, "Error retrieving yahoofin client", 500)
		return
	}

	query := r.URL.Query()

	errors := make([]string, 0)
	startDate := query.Get("startDate")
	if startDate == "" {
		errors = append(errors, "Missing querystring param: startDate")
	}

	endDate := query.Get("endDate")
	if endDate == "" {
		errors = append(errors, "Missing querystring param: endDate")
	}

	if len(errors) > 0 {
		http.Error(w, strings.Join(errors, "\n"), 500)
	}

	vars := mux.Vars(r)
	w.Write([]byte(fmt.Sprintf("%v", vars)))

	//client.GetSecurityData(vars["ticker"])

}
