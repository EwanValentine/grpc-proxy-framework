# gRPC API Gateway Framework

A small library for creating API Gateways that proxy gRPC calls to external services from HTTP/JSON inputs.

## Configuration

```yaml
routes:
  /api/v1/hello/:name:
    name: "service.Greeter"
    method: "SayHello"
    http_method: "GET"
  /api/v1/users:
    name: "service.Greeter"
    method: "CreateUser"
    http_method: "POST"

connections:
  service.Greeter: "localhost:8888"
```

## Implementation

```golang 
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/EwanValentine/grpc-proxy-framework/config"
	"github.com/EwanValentine/grpc-proxy-framework/connman"
	"github.com/gorilla/mux"
)

func helloHandler(conf *config.Config, connManager *connman.ConnectionManager) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		service, ok := conf.Routes.Find(r.URL.Path)
		if !ok {
			http.Error(rw, "no downstream route found for this path", http.StatusBadGateway)
			return
		}

		proxy, err := connManager.GetByName(service.Name)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := struct {
			Name string `json:"name"`
		}{
			Name: vars["name"],
		}

		response, err := proxy.Call(context.Background(), service.Name, service.Method, data)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = rw.Write(response)
	}
}

func userHandler(conf *config.Config, connMan *connman.ConnectionManager) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		service, ok := conf.Routes.Find(r.URL.Path)
		if !ok {
			http.Error(rw, "no downstream route found for this path", http.StatusBadGateway)
			return
		}

		proxy, err := connMan.GetByName(service.Name)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		type User struct {
			Name string `json:"name"`
		}

		data := struct {
			User User `json:"user"`
		}{
			User: User{Name: vars["name"]},
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		response, err := proxy.Call(context.Background(), service.Name, service.Method, data)
		if err != nil {
			log.Println(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = rw.Write(response)
	}
}

func main() {
	router := mux.NewRouter()

	conf, err := config.Parse("config.yaml")
	if err != nil {
		log.Panic(err)
	}

    // Creates a proxy connection for each connection
	connMan := connman.New(conf)
	if err := connMan.Start(); err != nil {
		log.Panic(err)
	}
	defer connMan.Close()

	router.HandleFunc("/api/v1/hello/{name}", helloHandler(conf, connMan)).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/users", userHandler(conf, connMan)).Methods(http.MethodPost)

	log.Println("running on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Panic(err)
	}
}
```
