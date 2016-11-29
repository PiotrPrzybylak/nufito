package shared

type GetTrainersRequest struct {
}

type GetTrainersResponse struct {
	V   []string `json:"v"`
	Err string   `json:"err,omitempty"` // errors don't define JSON marshaling
}

type AddTrainerRequest struct {
	Name string `json:"name"`
}

type AddTrainerResponse struct {
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

type NufitoService interface {
	GetTrainers() ([]string, error)
	AddTrainer(string) error
}
