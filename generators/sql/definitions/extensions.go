package definitions

import (
	"fmt"

	"github.com/cihangir/schema"
)

// DefineExtensions creates definition for extensions
func DefineExtensions(settings schema.Generator, s *schema.Schema) (res string) {
	for _, val := range s.Properties {
		if val.Default == nil {
			continue
		}

		def := fmt.Sprintf("%v", val.Default)
		switch def {
		case "uuid_generate_v1()", "uuid_generate_v1mc()", "uuid_generate_v4()":
			res += "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";\n"
		}
	}

	if res != "" {
		res = `-- ----------------------------
--  Required extensions
-- ----------------------------
` + res
	}

	return res
}
