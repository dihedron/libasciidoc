package html5

import (
	"github.com/bytesparadise/libasciidoc/pkg/renderer"
	"github.com/bytesparadise/libasciidoc/pkg/types"
)

func processAttributeDeclaration(ctx *renderer.Context, attr types.DocumentAttributeDeclaration) []byte {
	ctx.Document.Attributes.AddDeclaration(attr)
	return []byte{}
}

func processAttributeReset(ctx *renderer.Context, attr types.DocumentAttributeReset) []byte {
	ctx.Document.Attributes.Reset(attr)
	return []byte{}
}
