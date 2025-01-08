package sdblueprint

import (
	"fmt"
)

func defaultSchemaRefs(p *Project, schema Schema) map[string]string {
	return map[string]string{
		"go": fmt.Sprintf("*%s/%s.%s", p.ModName(), goModelsDir, schema.Names().Go()),
	}
}
