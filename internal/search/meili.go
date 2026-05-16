package search

import (
	"fmt"

	"github.com/meilisearch/meilisearch-go"
)

type Meili struct {
	Client meilisearch.ServiceManager
}

func NewMeili(host, apiKey string) *Meili {
	client := meilisearch.New(
		host,
		meilisearch.WithAPIKey(apiKey),
	)

	return &Meili{Client: client}
}

func (m *Meili) UpsertDocuments(indexName string, document interface{}) error {
	pk := "id"
	_, err := m.Client.Index(indexName).AddDocuments(document, &meilisearch.DocumentOptions{
		PrimaryKey: &pk,
	})
	return err
}

func (m *Meili) DeleteDocument(indexName string, id int) error {
	_, err := m.Client.Index(indexName).DeleteDocument(fmt.Sprintf("%d", id), nil)
	return err
}
