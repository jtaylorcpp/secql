package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import (
	"github.com/aws/aws-sdk-go/aws/session"
)

type Resolver struct {
	Session *session.Session
}
