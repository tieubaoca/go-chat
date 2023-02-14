package request

type CreateNewGroupReq struct {
	Name    string   `json:"name" bson:"name"`
	Members []string `json:"members" bson:"members"`
}
