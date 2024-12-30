package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestStep(t *testing.T) {
	os.Setenv("ELASTICSEARCH_INDEX_DEFINITION_DIR", "../../../shared/elasticsearch/indices")
	defer os.Unsetenv("ELASTICSEARCH_INDEX_DEFINITION_DIR")
	os.Setenv("ELASTICSEARCH_INDEX_SUFFIX", fmt.Sprintf("-test-%s", uuid.NewString()))
	defer os.Unsetenv("ELASTICSEARCH_INDEX_SUFFIX")

	main()

	t.Fatal("output log")
}
