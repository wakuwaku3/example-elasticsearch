package main

import (
	"os"
	"strconv"

	"github.com/wakuwaku3/example-elasticsearch/golang/elasticsearch"
	"github.com/wakuwaku3/example-elasticsearch/golang/user"
)

func main() {
	opt, err := createOptions()
	if err != nil {
		panic(err)
	}

	client := elasticsearch.NewClient(opt)
	if err := client.CreateIndices(); err != nil {
		panic(err)
	}

	u, err := user.NewModel("name", 20)
	if err != nil {
		panic(err)
	}

	repo := user.NewRepository(client)
	if err := repo.Upsert(u); err != nil {
		panic(err)
	}

	queryService := user.NewQueryService(client)
	query := &user.GetByNameQuery{Name: "name"}
	if _, err := queryService.GetByName(query); err != nil {
		panic(err)
	}

	if err := repo.Refresh(); err != nil {
		panic(err)
	}

	if _, err := queryService.GetByName(query); err != nil {
		panic(err)
	}

	u.Rename("new name")
	if err := repo.Upsert(u); err != nil {
		panic(err)
	}

	query2 := &user.GetByNameQuery{Name: "new name"}
	if _, err := queryService.GetByName(query2); err != nil {
		panic(err)
	}

	if err := repo.Refresh(); err != nil {
		panic(err)
	}

	if _, err := queryService.GetByName(query2); err != nil {
		panic(err)
	}

	if err := repo.Delete(u); err != nil {
		panic(err)
	}

	if err := repo.Refresh(); err != nil {
		panic(err)
	}

	if _, err := queryService.GetByName(query2); err != nil {
		panic(err)
	}
}

func createOptions() (*elasticsearch.ClientOption, error) {
	o := &elasticsearch.ClientOption{
		IndexDefinitionDir: "../shared/elasticsearch/indices",
	}

	o.UseSSL = os.Getenv("ELASTICSEARCH_USE_SSL") == "true"
	o.Host = os.Getenv("ELASTICSEARCH_HOST")
	port, err := strconv.Atoi(os.Getenv("ELASTICSEARCH_PORT"))
	if err != nil {
		return nil, err
	}
	o.Port = port
	o.UserName = os.Getenv("ELASTICSEARCH_ELASTIC_USERNAME")
	o.Password = os.Getenv("ELASTICSEARCH_ELASTIC_PASSWORD")
	o.IndexNameSuffix = os.Getenv("ELASTICSEARCH_INDEX_SUFFIX")

	indexDefinitionDir := os.Getenv("ELASTICSEARCH_INDEX_DEFINITION_DIR")
	if indexDefinitionDir != "" {
		o.IndexDefinitionDir = indexDefinitionDir
	}

	return o, nil
}
