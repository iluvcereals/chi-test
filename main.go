package main

import "log"

func main() {
	cfg := config{
		addr: ":5001",
	}
	api := &application{
		config: cfg,
	}

	router := api.mount()
	log.Fatal(api.run(router))
}
