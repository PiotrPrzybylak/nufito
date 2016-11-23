package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	//"strings"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"

	"bitbucket.org/piotrp/nufito-prototype/shared"
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
	svc := nufitoService{}

	trainersHandler := httptransport.NewServer(
		ctx,
		makeTrainersEndpoint(svc),
		decodeGetTrainersRequest,
		encodeResponse,
	)

	http.Handle("/trainers", trainersHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func makeTrainersEndpoint(svc NufitoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//req := request.(getTrainersRequest)
		v, err := svc.GetTrainers()
		if err != nil {
			return shared.GetTrainersResponse{v, err.Error()}, nil
		}
		return shared.GetTrainersResponse{v, ""}, nil
	}
}

func decodeGetTrainersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request shared.GetTrainersRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")
