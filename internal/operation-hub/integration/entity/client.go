package entity

import "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"

type Client struct {
	Id string `json:"client_id"`
}

func (c Client) ToModel() model.Client {
	return model.Client{
		Id: c.Id,
	}
}
