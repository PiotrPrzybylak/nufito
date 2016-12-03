package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/piotrprz/nufito/db"
	"github.com/piotrprz/nufito/shared"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
)

type nufitoService struct {
	Trainers []string
}

func (svc nufitoService) GetTrainers() ([]string, error) {
	return svc.Trainers, nil
}

func (svc *nufitoService) AddTrainer(trainer string) error {
	svc.Trainers = append(svc.Trainers, trainer)
	return nil
}

func main() {

	logger := log.NewLogfmtLogger(os.Stderr)

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "nufito",
		Subsystem: "trainers_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "nufito",
		Subsystem: "trainers_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "nufito",
		Subsystem: "trainers_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

	ctx := context.Background()
	// var svc NufitoService = &nufitoService{Trainers: []string{"Marian", "Stefan", "Roman"}}
	var svc shared.NufitoService = db.NewService()

	svc = instrumentingMiddleware{requestCount, requestLatency, countResult, svc}

	trainersEndpoint := makeTrainersEndpoint(svc)
	trainersEndpoint = loggingMiddleware(log.NewContext(logger).With("method", "getTrainers"))(trainersEndpoint)

	trainersHandler := httptransport.NewServer(
		ctx,
		trainersEndpoint,
		decodeGetTrainersRequest,
		encodeResponse,
	)

	addTrainerEndpoint := makeAddTrainerEndpoint(svc)
	addTrainerEndpoint = loggingMiddleware(log.NewContext(logger).With("method", "AddTrainer"))(addTrainerEndpoint)

	addTrainerHandler := httptransport.NewServer(
		ctx,
		addTrainerEndpoint,
		decodeAddTrainerRequest,
		encodeResponse,
	)

	http.Handle("/trainers/add", addTrainerHandler)
	http.Handle("/trainers", trainersHandler)
	http.Handle("/metrics", stdprometheus.Handler())
	logger.Log("msg", "HTTP", "addr", ":8080")
	logger.Log("err", http.ListenAndServe(":8080", nil))
}

func makeTrainersEndpoint(svc shared.NufitoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		trainers, err := svc.GetTrainers()
		if err != nil {
			return shared.GetTrainersResponse{Trainers: trainers, Err: err.Error()}, nil
		}
		return shared.GetTrainersResponse{Trainers: trainers, Err: ""}, nil
	}
}

func makeAddTrainerEndpoint(svc shared.NufitoService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(shared.AddTrainerRequest)
		err := svc.AddTrainer(req.Name)
		if err != nil {
			return shared.AddTrainerResponse{Err: err.Error()}, nil
		}
		return shared.AddTrainerResponse{Err: ""}, nil
	}
}

func decodeGetTrainersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request shared.GetTrainersRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeAddTrainerRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request shared.AddTrainerRequest
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
