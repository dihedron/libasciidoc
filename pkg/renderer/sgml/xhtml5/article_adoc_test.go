package xhtml5_test

import (
	"bufio"
	"bytes"
	"github.com/bytesparadise/libasciidoc/pkg/renderer/sgml/xhtml5"
	"os"

	"github.com/bytesparadise/libasciidoc/pkg/configuration"
	"github.com/bytesparadise/libasciidoc/pkg/parser"
	"github.com/bytesparadise/libasciidoc/pkg/renderer"
	"github.com/davecgh/go-spew/spew"
	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("article.adoc", func() {

	It("should render without failure", func() {
		f, err := os.Open("article.adoc")
		Expect(err).ToNot(HaveOccurred())
		reader := bufio.NewReader(f)
		config := configuration.NewConfiguration()
		doc, err := parser.ParseDocument(reader, config)
		Expect(err).ToNot(HaveOccurred())
		GinkgoT().Logf("actual document: `%s`", spew.Sdump(doc))
		buff := &bytes.Buffer{}
		ctx := renderer.NewContext(doc, config)
		_, err = xhtml5.Render(ctx, doc, buff)
		Expect(err).ToNot(HaveOccurred())
	})
})
