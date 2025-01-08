package sdblueprint

import (
	"fmt"
	"github.com/gaorx/stardust6/sdcodegen"
	"github.com/gaorx/stardust6/sdcodegen/sdgogen"
	"github.com/gaorx/stardust6/sderr"
	"github.com/gaorx/stardust6/sdregexp"
	"regexp"
	"strings"
)

type GoRef struct {
	Pkgs []string
	Code string
}

var refPatt = regexp.MustCompile(`^(?P<prefix>([\[\]*])*)(?P<base>(\w|\.|/)*)?$`)

func GoParseRef(s string) (GoRef, bool) {
	s = strings.Replace(s, " ", "", -1)
	if s == "" {
		return GoRef{}, true
	}
	m := sdregexp.FindStringSubmatchGroup(refPatt, s)
	prefix, base := m["prefix"], m["base"]
	if base == "" {
		return GoRef{}, false
	}
	if !strings.Contains(base, ".") {
		return GoRef{Code: prefix + base}, true
	} else {
		ss := strings.SplitAfter(base, ".")
		ssLen := len(ss)
		pkgPath := strings.TrimSuffix(strings.Join(ss[:ssLen-1], ""), ".")
		pkgName := lastStringOf(strings.Split(pkgPath, "/"))
		return GoRef{
			Pkgs: []string{pkgPath},
			Code: prefix + pkgName + "." + ss[ssLen-1],
		}, true
	}
}

func goNewCodeRef(code string) GoRef {
	return GoRef{Code: code}
}

func (ref GoRef) IsZero() bool {
	return len(ref.Pkgs) <= 0 && ref.Code == ""
}

func (ref GoRef) AddPkgs(pkgs ...string) GoRef {
	ref.Pkgs = append(ref.Pkgs, pkgs...)
	return ref
}

func (ref GoRef) AddArray() GoRef {
	return GoRef{Pkgs: ref.Pkgs, Code: "[]" + ref.Code}
}

func goGenRef(c *sdgogen.Context, typ Type) GoRef {
	goRef := typ.Refs().Go()
	ref, ok := GoParseRef(goRef)
	if !ok {
		panic(sderr.With("ref", goRef).Newf("parse ref failed"))
	}
	if !ref.IsZero() {
		return ref
	}

	switch typ.Kind() {
	case KString:
		return goNewCodeRef("string")
	case KBool:
		return goNewCodeRef("bool")
	case KBytes:
		return goNewCodeRef("[]byte")
	case KEnum:
		return goNewCodeRef("uint8")
	case KInt:
		return goNewCodeRef("int")
	case KInt64:
		return goNewCodeRef("int64")
	case KUint:
		return goNewCodeRef("uint")
	case KUint64:
		return goNewCodeRef("uint64")
	case KFloat64:
		return goNewCodeRef("float64")
	case KSchema:
		schema := typ.Schema()
		if schema == nil {
			return goNewCodeRef("any")
		}
		if schema.Id() != "" && !ref.IsZero() {
			return ref
		} else {
			ref = goGenNamelessStructRef(c, schema)
			return ref
		}
	case KArray:
		elem := typ.Elem()
		if elem == nil {
			return goNewCodeRef("[]any")
		}
		elemRef := goGenRef(c, elem)
		if elemRef.IsZero() {
			return elemRef
		}
		return elemRef.AddArray()
	default:
		return goNewCodeRef("any")
	}
}

func goGenNamelessStructRef(cx *sdgogen.Context, schema Schema) GoRef {
	var r GoRef
	r.Code = cx.GenerateText(func(c1 *sdcodegen.Context) {
		cx1 := sdgogen.C(c1)
		cx1.AnonymousStructPtr(func() {
			props := schema.Properties()
			for _, prop := range props {
				colName := prop.Names().Go()
				colRef := goGenRef(cx1, prop.Type())
				r.AddPkgs(colRef.Pkgs...)
				fieldTags := []string{goMakeJsonTag(prop)}
				fieldTags = append(fieldTags, prop.Anns().Go()...)
				cx1.Tab().Field(colName, colRef.Code, fieldTags, prop.Comment())
			}
		})
	})
	return r
}

func goMakeJsonTag(col *Property) string {
	return fmt.Sprintf(`json:"%s"`, col.Names().Json())
}
