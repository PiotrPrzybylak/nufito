package shared

type GetTrainersRequest struct {
}

type GetTrainersResponse struct {
	Trainers []string `json:"trainers"`
	Err      string   `json:"err,omitempty"` // errors don't define JSON marshaling
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
