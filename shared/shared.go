package shared

type GetTrainersRequest struct {
}

type GetTrainersResponse struct {
	Trainers []Trainer `json:"trainers"`
	Err      string    `json:"err,omitempty"` // errors don't define JSON marshaling
}

type AddTrainerRequest struct {
	Name string `json:"name"`
}

type AddTrainerResponse struct {
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

type Trainer struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type NufitoService interface {
	GetTrainers() ([]Trainer, error)
	AddTrainer(string) error
}
