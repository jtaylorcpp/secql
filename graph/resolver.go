package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jtaylorcpp/secql/osquery"
)

type Resolver struct {
	Session *session.Session
	Cache   *OSQueryClientCache
}

type OSQueryClientCache struct {
	Cache map[string]osquery.Client
}

func (c *OSQueryClientCache) Exists(key string) bool {
	_, ok := c.Cache[key]
	return ok
}

func (c *OSQueryClientCache) Get(key string) osquery.Client {
	if client, ok := c.Cache[key]; ok {
		return client
	} else {
		return nil
	}
}

func (c *OSQueryClientCache) Put(key string, client osquery.Client) {
	c.Cache[key] = client
}
