package elasticsearch

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

type client struct {
	options *ClientOption
}

type Client interface {
	CreateIndices() error
	Put(indexName string, id uuid.UUID, data io.Reader) error
	Delete(indexName string, id uuid.UUID) error
	Refresh(indexName string) error
	Search(indexName string, query io.Reader) (io.ReadCloser, error)
}

var _ Client = (*client)(nil)

func NewClient(options *ClientOption) *client {
	return &client{options: options}
}

func (c *client) CreateIndices() error {
	// ../shared/elasticsearch/indices から index 定義の json ファイルを全て取得する
	files, err := os.ReadDir(c.options.IndexDefinitionDir)
	if err != nil {
		return err
	}

	client := http.Client{}
	for _, file := range files {
		name := file.Name()
		if !strings.HasSuffix(name, ".json") {
			continue
		}
		body, err := os.Open(filepath.Join(c.options.IndexDefinitionDir, name))
		if err != nil {
			return err
		}
		defer body.Close()

		indexName := strings.TrimSuffix(name, ".json")
		req, err := http.NewRequest("PUT", c.options.indexURL(indexName), body)
		if err != nil {
			return err
		}
		req.SetBasicAuth(c.options.UserName, c.options.Password)
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if res.StatusCode >= 300 {
			return fmt.Errorf("Index creating is failure. index name: %s, status code: %d, message: %s", indexName, res.StatusCode, string(resBody))
		}

		log.Printf("Index %s created", file.Name())
	}

	return nil
}

func (c *client) Put(indexName string, id uuid.UUID, data io.Reader) error {
	client := http.Client{}
	req, err := http.NewRequest("PUT", c.options.documentURL(indexName, id), data)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.options.UserName, c.options.Password)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= 300 {
		return fmt.Errorf("Put is failure. index name: %s, id: %s, status code: %d, message: %s", indexName, id.String(), res.StatusCode, string(resBody))
	}

	log.Printf("Put is success. index name: %s, id: %s", indexName, id.String())

	return nil
}

func (c *client) Delete(indexName string, id uuid.UUID) error {
	client := http.Client{}
	req, err := http.NewRequest("DELETE", c.options.documentURL(indexName, id), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.options.UserName, c.options.Password)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= 300 {
		return fmt.Errorf("Delete is failure. index name: %s, id: %s, status code: %d, message: %s", indexName, id.String(), res.StatusCode, string(resBody))
	}

	log.Printf("Delete is success. index name: %s, id: %s", indexName, id.String())

	return nil
}

func (c *client) Refresh(indexName string) error {
	client := http.Client{}
	req, err := http.NewRequest("POST", c.options.refreshURL(indexName), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.options.UserName, c.options.Password)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= 300 {
		return fmt.Errorf("Refresh is failure. index name: %s, status code: %d, message: %s", indexName, res.StatusCode, string(resBody))
	}

	log.Printf("Refresh is success. index name: %s", indexName)

	return nil
}

func (c *client) Search(indexName string, query io.Reader) (io.ReadCloser, error) {
	client := http.Client{}
	req, err := http.NewRequest("POST", c.options.indexURL(indexName)+"/_search", query)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.options.UserName, c.options.Password)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 300 {
		defer res.Body.Close()
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("Search is failure. index name: %s, status code: %d, message: %s", indexName, res.StatusCode, string(resBody))
	}

	log.Printf("Search is success. index name: %s", indexName)

	return res.Body, nil
}
