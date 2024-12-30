package elasticsearch

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type ClientOption struct {
	UseSSL             bool
	Host               string
	Port               int
	UserName           string
	Password           string
	IndexDefinitionDir string
	IndexNameSuffix    string
}

func (o *ClientOption) url() string {
	schema := ""
	if o.UseSSL {
		schema = "https://"
	} else {
		schema = "http://"
	}

	return schema + o.Host + ":" + strconv.Itoa(o.Port)
}

func (o *ClientOption) indexURL(indexName string) string {
	name := indexName
	if !strings.HasSuffix(indexName, o.IndexNameSuffix) {
		name = o.indexName(indexName)
	}

	return o.url() + "/" + name
}

func (o *ClientOption) documentURL(indexName string, id uuid.UUID) string {
	return o.indexURL(indexName) + "/_doc/" + id.String()
}

func (o *ClientOption) refreshURL(indexName string) string {
	return o.indexURL(indexName) + "/_refresh"
}

func (o *ClientOption) indexName(indexName string) string {
	return indexName + o.IndexNameSuffix
}
