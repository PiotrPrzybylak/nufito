package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"errors"
	"log"
	"net/http"
	//"strings"
	"fmt"

	"net/url"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type NufitoService interface {
	GetTrainers() ([]string, error)
}

type nufitoService struct{}

func (nufitoService) GetTrainers() ([]string, error) {
	return []string{"Marian", "Stefan", "Roman"}, nil
}

func main() {
	ctx := context.Background()

	trainersHandler := httptransport.NewServer(
		ctx,
		makeTrainersEndpoint(ctx, "http://localhost:8080/trainers"),
		decodeGetTrainersRequest,
		encodeResponse,
	)

	http.Handle("/trainers", trainersHandler)
	log.Fatal(http.ListenAndServe(":8081", nil))
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
		decodeUppercaseResponse,
	).Endpoint()
}


func decodeGetTrainersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request getTrainersRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}


func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type getTrainersRequest struct {
}

type getTrainersResponse struct {
	V   []string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}


// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")


func encodeRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeGetTrainersResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response getTrainersResponse

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
