package models

type GrpcGroup struct {
	Id         int64
	Name       string
	TrainerId  int64
	StudentIds []int64
	Link       string
}

type GroupResponse struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Trainer  User   `json:"trainer"`
	Students []User `json:"students"`
	Link     string `json:"link"`
}
