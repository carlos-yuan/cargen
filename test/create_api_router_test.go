package test

import (
	"testing"

	"github.com/carlos-yuan/cargen/cmd/gen"
)

func TestCreateApiRouter(t *testing.T) {
	gen.CreateApiRouter("/media/ysgk/DATA/carlos/hc_enterprise_server")
}
