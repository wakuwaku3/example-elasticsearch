package user

import (
	"bytes"
	"encoding/json"

	"github.com/wakuwaku3/example-elasticsearch/golang/elasticsearch"
)

type repository struct {
	client elasticsearch.Client
}

func NewRepository(client elasticsearch.Client) *repository {
	return &repository{client: client}
}

const indexName = "users"

func (r *repository) Upsert(model *model) error {
	jsonBody := model.jsonBody()

	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(jsonBody); err != nil {
		return err
	}
	return r.client.Put(indexName, model.id, buffer)
}

func (r *repository) Delete(model *model) error {
	return r.client.Delete(indexName, model.id)
}

func (r *repository) Refresh() error {
	return r.client.Refresh(indexName)
}
