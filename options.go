package httprouter

import (
	"github.com/caasmo/restinpieces/core"
)

func WithRouterHttprouter() core.Option {
	r := New()
	return core.WithRouter(r)
}

