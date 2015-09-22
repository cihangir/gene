package kit

import (
	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
)

func GenerateTransportHTTP(context *common.Context, s *schema.Schema) ([]common.Output, error) {
	return generate(context, s, TransportHTTPTemplate, "transport_http")
}
