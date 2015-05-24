package definitions

import (
	"fmt"
	"strings"

	"github.com/cihangir/schema"
)

// DefineExtensions creates definition for extensions
func DefineExtensions(settings schema.Generator, s *schema.Schema) (res string) {
	exts := make([]string, 0)

	for _, val := range s.Properties {
		if val.Default == nil {
			continue
		}

		def := fmt.Sprintf("%v", val.Default)
		switch def {
		case "uuid_generate_v1()", "uuid_generate_v1mc()", "uuid_generate_v4()":
			exts = append(exts, "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
		}
	}

	if len(exts) > 0 {
		res = `
-- ----------------------------
--  Required extensions
-- ----------------------------
` + strings.Join(exts, "\n")
	}

	return res
}
