package user

import (
	"bytes"
	"encoding/json"
	"io"
	"log"

	"github.com/wakuwaku3/example-elasticsearch/golang/elasticsearch"
)

type queryService struct {
	client elasticsearch.Client
}

func NewQueryService(client elasticsearch.Client) *queryService {
	return &queryService{client: client}
}

type GetByNameQuery struct {
	Name string
}

func (q *GetByNameQuery) Reader() (io.Reader, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": q.Name,
			},
		},
	}
	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(query); err != nil {
		return nil, err
	}
	return buffer, nil
}

func (s *queryService) GetByName(query *GetByNameQuery) ([]*model, error) {
	q, err := query.Reader()
	if err != nil {
		return nil, err
	}
	reader, err := s.client.Search(indexName, q)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	res, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	log.Printf("res: %s", res)
	return nil, nil
}
