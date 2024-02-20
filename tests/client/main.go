package main

import (
	"encoding/json"
	"net"
	"net/http"

	"log"
)

type ResponseWebhook struct {
	Iteration uint64 `json:"iteration"`
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {

		rwebhook := &ResponseWebhook{}

		if err := json.NewDecoder(r.Body).Decode(rwebhook); err != nil {
			log.Println(err)
			return
		}

		log.Printf("%v\n", *rwebhook)
	})

	listner, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatalln(err)
	}

	srv := &http.Server{
		Addr:    listner.Addr().String(),
		Handler: mux,
	}
	log.Println("lisnet serving on ", listner.Addr().String())
	if err := srv.Serve(listner); err != nil {
		log.Fatalln(err)
	}
}
