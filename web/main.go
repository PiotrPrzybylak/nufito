package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/piotrprz/nufito/shared"
	"golang.org/x/net/context"
)

var templates = template.Must(template.ParseFiles("trainers.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *shared.GetTrainersResponse) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	ctx := context.Background()

	trainersEndpoint := makeTrainersEndpoint(ctx, "http://backend:8080/trainers")

	addTrainerEndpoint := makeAddTrainerEndpoint(ctx, "http://backend:8080/trainers/add")

	trainersHandler := func(w http.ResponseWriter, r *http.Request) {

		response, _ := trainersEndpoint(ctx, shared.GetTrainersRequest{})
		res := response.(shared.GetTrainersResponse)
		renderTemplate(w, "trainers", &res)
	}

	addTrainerHandler := func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		_, _ = addTrainerEndpoint(ctx, shared.AddTrainerRequest{Name: name})
		http.Redirect(w, r, "trainers", 301)
	}

	http.HandleFunc("/trainers", trainersHandler)
	http.HandleFunc("/add-trainer", addTrainerHandler)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Fatal(http.ListenAndServe(":8082", nil))
}

func makeTrainersEndpoint(ctx context.Context, proxyURL string) endpoint.Endpoint {
	u, err := url.Parse(proxyURL)
	if err != nil {
		panic(err)
	}
	return httptransport.NewClient(
		"GET",
		u,
		encodeRequest,
		decodeGetTrainersResponse,
	).Endpoint()
}

func makeAddTrainerEndpoint(ctx context.Context, proxyURL string) endpoint.Endpoint {
	u, err := url.Parse(proxyURL)
	if err != nil {
		panic(err)
	}
	return httptransport.NewClient(
		"GET",
		u,
		encodeRequest,
		decodeAddTrainerResponse,
	).Endpoint()
}

func encodeRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeGetTrainersResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response shared.GetTrainersResponse

	//	buf := new(bytes.Buffer)
	//	buf.ReadFrom(r.Body)
	//	s := buf.String() // Does a complete copy of the bytes in the buffer.
	// fmt.Print("s:")
	// fmt.Println(s)

	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		fmt.Print("err:")
		fmt.Println(err)
		return nil, err
	}

	fmt.Print("response: ")
	fmt.Println(response)

	return response, nil
}

func decodeAddTrainerResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response shared.AddTrainerResponse

	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		fmt.Print("err:")
		fmt.Println(err)
		return nil, err
	}

	fmt.Print("response: ")
	fmt.Println(response)

	return response, nil
}
