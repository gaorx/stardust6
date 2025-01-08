package sdblueprint

import (
	"github.com/gaorx/stardust6/sderr"
	"slices"
)

type APIDecl struct {
	id         string
	httpMethod string
	httpPath   string
	category   string
	inTyp      Type
	outTyp     Type
	protocol   APIProtocol
	guardKind  APIGuard
	guardVals  []string
	description
}

type APIProtocol string

const (
	APISimple APIProtocol = "simple"
)

type APIDeclBuilder APIDecl

var _ builder[*APIDecl] = (*APIDeclBuilder)(nil)

func SimpleAPI(id string, path string, in, out Type) *APIDeclBuilder {
	return &APIDeclBuilder{
		id:         id,
		httpMethod: "POST",
		httpPath:   trimAPIPath(path),
		category:   "",
		inTyp:      in,
		outTyp:     out,
		protocol:   APISimple,
		guardKind:  APIPermitAll,
	}
}

func (a *APIDecl) Id() string {
	return a.id
}

func (a *APIDecl) HttpMethod() string {
	return a.httpMethod
}

func (a *APIDecl) HttpPath() string {
	return a.httpPath
}

func (a *APIDecl) Category() string {
	return a.category
}

func (a *APIDecl) In() Type {
	return a.inTyp
}

func (a *APIDecl) Out() Type {
	return a.outTyp
}

func (a *APIDecl) Protocol() APIProtocol {
	return a.protocol
}

func (a *APIDecl) Guard() APIGuard {
	return a.guardKind
}

func (a *APIDecl) GuardVals() []string {
	return a.guardVals
}

func (a *APIDecl) GuardVal() string {
	if len(a.guardVals) > 0 {
		return a.guardVals[0]
	}
	return ""
}

func (a *APIDecl) Doc() string {
	return a.doc
}

func (a *APIDecl) Comment() string {
	return a.comment
}

func (a *APIDecl) asBuilder() *APIDeclBuilder {
	return (*APIDeclBuilder)(a)
}

func (a *APIDecl) postBuild(p *Project) error {
	if err := resolveType(p, a.inTyp); err != nil {
		return sderr.Wrap(err)
	}
	if err := resolveType(p, a.outTyp); err != nil {
		return sderr.Wrap(err)
	}
	return nil
}

func (b *APIDeclBuilder) Category(category string) *APIDeclBuilder {
	b.category = category
	return b
}

func (b *APIDeclBuilder) PermitAll() *APIDeclBuilder {
	b.guardKind = APIPermitAll
	return b
}

func (b *APIDeclBuilder) RejectAll() *APIDeclBuilder {
	b.guardKind = APIRejectAll
	return b
}

func (b *APIDeclBuilder) Authenticated() *APIDeclBuilder {
	b.guardKind = APIAuthenticated
	return b
}

func (b *APIDeclBuilder) HasAuthority(authorities ...string) *APIDeclBuilder {
	b.guardKind = APIHasAuthority
	b.guardVals = slices.Clone(authorities)
	return b
}

func (b *APIDeclBuilder) IsMatched(expression string) *APIDeclBuilder {
	b.guardKind = APIIsMatched
	b.guardVals = []string{expression}
	return b
}

func (b *APIDeclBuilder) asAPI() *APIDecl {
	return (*APIDecl)(b)
}

func (b *APIDeclBuilder) prepare(_ *buildContext) error {
	return nil
}

func (b *APIDeclBuilder) build(_ *buildContext) (*APIDecl, error) {
	return &APIDecl{
		id:          b.id,
		httpMethod:  b.httpMethod,
		httpPath:    b.httpPath,
		category:    b.category,
		inTyp:       b.inTyp,
		outTyp:      b.outTyp,
		protocol:    b.protocol,
		guardKind:   b.guardKind,
		guardVals:   slices.Clone(b.guardVals),
		description: b.description,
	}, nil
}
