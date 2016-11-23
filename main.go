package main

import (
	"bitbucket.org/piotrp/nufito-prototype/shared"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"net/http"
	"os"
)

type NufitoService interface {
	GetTrainers() ([]string, error)
}

type nufitoService struct{}

func (nufitoService) GetTrainers() ([]string, error) {
	return []string{"Marian", "Stefan", "Roman"}, nil
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
	var svc NufitoService
	svc = nufitoService{}
	svc = instrumentingMiddleware{requestCount, requestLatency, countResult, svc}

	trainersEndpoint := makeTrainersEndpoint(svc)
	trainersEndpoint = loggingMiddleware(log.NewContext(logger).With("method", "getTrainers"))(trainersEndpoint)

	trainersHandler := httptransport.NewServer(
		ctx,
		trainersEndpoint,
		decodeGetTrainersRequest,
		encodeResponse,
	)

	http.Handle("/trainers", trainersHandler)
	http.Handle("/metrics", stdprometheus.Handler())
	logger.Log("msg", "HTTP", "addr", ":8080")
	logger.Log("err", http.ListenAndServe(":8080", nil))
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
