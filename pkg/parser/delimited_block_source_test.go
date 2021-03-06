package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("source blocks", func() {

	Context("draft documents", func() {

		Context("delimited blocks", func() {
			sourceCode := [][]interface{}{
				{
					types.StringElement{
						Content: "package foo",
					},
				},
				{},
				{
					types.StringElement{
						Content: "// Foo",
					},
				},
				{
					types.StringElement{
						Content: "type Foo struct{",
					},
				},
				{
					types.StringElement{
						Content: "    Bar string",
					},
				},
				{
					types.StringElement{
						Content: "}",
					},
				},
			}

			It("with source attribute only", func() {
				source := `[source]
----
package foo

// Foo
type Foo struct{
    Bar string
}
----`
				expected := types.DraftDocument{
					Elements: []interface{}{
						types.ListingBlock{
							Attributes: types.Attributes{
								types.AttrBlockKind: types.Source,
							},
							Lines: sourceCode,
						},
					},
				}
				Expect(ParseDraftDocument(source)).To(MatchDraftDocument(expected))
			})

			It("with source attribute and comma", func() {
				source := `[source,]
----
package foo

// Foo
type Foo struct{
    Bar string
}
----`
				expected := types.DraftDocument{
					Elements: []interface{}{
						types.ListingBlock{
							Attributes: types.Attributes{
								types.AttrBlockKind: types.Source,
							},
							Lines: sourceCode,
						},
					},
				}
				Expect(ParseDraftDocument(source)).To(MatchDraftDocument(expected))
			})

			It("with title, source and language attributes", func() {
				source := `[source,go]
.foo.go
----
package foo

// Foo
type Foo struct{
    Bar string
}
----`
				expected := types.DraftDocument{
					Elements: []interface{}{
						types.ListingBlock{
							Attributes: types.Attributes{
								types.AttrBlockKind: types.Source,
								types.AttrLanguage:  "go",
								types.AttrTitle:     "foo.go",
							},
							Lines: sourceCode,
						},
					},
				}
				Expect(ParseDraftDocument(source)).To(MatchDraftDocument(expected))
			})

			It("with id, title, source and language and other attributes", func() {
				source := `[#id-for-source-block]
[source,go,linenums]
.foo.go
----
package foo

// Foo
type Foo struct{
    Bar string
}
----`
				expected := types.DraftDocument{
					Elements: []interface{}{
						types.ListingBlock{
							Attributes: types.Attributes{
								types.AttrBlockKind: types.Source,
								types.AttrLanguage:  "go",
								types.AttrID:        "id-for-source-block",
								types.AttrCustomID:  true,
								types.AttrTitle:     "foo.go",
								types.AttrLineNums:  nil,
							},
							Lines: sourceCode,
						},
					},
				}
				Expect(ParseDraftDocument(source)).To(MatchDraftDocument(expected))
			})
		})
	})

	Context("final documents", func() {

		Context("source block", func() {

			It("with source attribute only", func() {
				source := `[source]
----
require 'sinatra'

get '/hi' do
  "Hello World!"
end
----`
				expected := types.Document{
					Elements: []interface{}{
						types.ListingBlock{
							Attributes: types.Attributes{
								types.AttrBlockKind: types.Source,
							},
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "require 'sinatra'",
									},
								},
								{},
								{
									types.StringElement{
										Content: "get '/hi' do",
									},
								},
								{
									types.StringElement{
										Content: "  \"Hello World!\"",
									},
								},
								{
									types.StringElement{
										Content: "end",
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with title, source and languages attributes", func() {
				source := `[source,ruby]
.Source block title
----
require 'sinatra'

get '/hi' do
  "Hello World!"
end
----`
				expected := types.Document{
					Elements: []interface{}{
						types.ListingBlock{
							Attributes: types.Attributes{
								types.AttrBlockKind: types.Source,
								types.AttrLanguage:  "ruby",
								types.AttrTitle:     "Source block title",
							},
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "require 'sinatra'",
									},
								},
								{},
								{
									types.StringElement{
										Content: "get '/hi' do",
									},
								},
								{
									types.StringElement{
										Content: "  \"Hello World!\"",
									},
								},
								{
									types.StringElement{
										Content: "end",
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("with id, title, source and languages attributes", func() {
				source := `[#id-for-source-block]
[source,ruby]
.app.rb
----
require 'sinatra'

get '/hi' do
  "Hello World!"
end
----`
				expected := types.Document{
					Elements: []interface{}{
						types.ListingBlock{
							Attributes: types.Attributes{
								types.AttrBlockKind: types.Source,
								types.AttrLanguage:  "ruby",
								types.AttrID:        "id-for-source-block",
								types.AttrCustomID:  true,
								types.AttrTitle:     "app.rb",
							},
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "require 'sinatra'",
									},
								},
								{},
								{
									types.StringElement{
										Content: "get '/hi' do",
									},
								},
								{
									types.StringElement{
										Content: "  \"Hello World!\"",
									},
								},
								{
									types.StringElement{
										Content: "end",
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})
		})
	})
})
