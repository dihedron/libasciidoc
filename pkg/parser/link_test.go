package parser_test

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
	. "github.com/bytesparadise/libasciidoc/testsupport"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("links", func() {

	Context("final documents", func() {

		Context("external links", func() {

			It("external link without text", func() {
				source := "a link to https://foo.bar"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "a link to "},
									types.InlineLink{
										Location: types.Location{
											Scheme: "https://",
											Path: []interface{}{
												types.StringElement{
													Content: "foo.bar",
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("external link with empty text", func() {
				source := "a link to https://foo.bar[]"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "a link to "},
									types.InlineLink{
										Location: types.Location{
											Scheme: "https://",
											Path: []interface{}{
												types.StringElement{
													Content: "foo.bar",
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("external link with text only", func() {
				source := "a link to mailto:foo@bar[the foo@bar email]."
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "a link to "},
									types.InlineLink{
										Location: types.Location{
											Scheme: "mailto:",
											Path: []interface{}{
												types.StringElement{
													Content: "foo@bar",
												},
											},
										},
										Attributes: types.Attributes{
											"positional-1": []interface{}{
												types.StringElement{
													Content: "the foo@bar email",
												},
											},
										},
									},
									types.StringElement{Content: "."},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("external link with text and extra attributes", func() {
				source := "a link to mailto:foo@bar[the foo@bar email, foo=bar]"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "a link to "},
									types.InlineLink{
										Location: types.Location{
											Scheme: "mailto:",
											Path: []interface{}{
												types.StringElement{
													Content: "foo@bar",
												},
											},
										},
										Attributes: types.Attributes{
											"positional-1": []interface{}{
												types.StringElement{
													Content: "the foo@bar email",
												},
											},
											"foo": "bar",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("external link inside a multiline paragraph -  without attributes", func() {
				source := `a http://website.com
and more text on the
next lines`

				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{
									Content: "a ",
								},
									types.InlineLink{
										Location: types.Location{
											Scheme: "http://",
											Path: []interface{}{
												types.StringElement{
													Content: "website.com",
												},
											},
										},
									},
								},
								{
									types.StringElement{
										Content: "and more text on the",
									},
								},
								{
									types.StringElement{
										Content: "next lines",
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("external link inside a multiline paragraph -  with attributes", func() {
				source := `a http://website.com[]
and more text on the
next lines`

				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{
									Content: "a ",
								},
									types.InlineLink{
										Location: types.Location{
											Scheme: "http://",
											Path: []interface{}{
												types.StringElement{
													Content: "website.com",
												},
											},
										},
									},
								},
								{
									types.StringElement{
										Content: "and more text on the",
									},
								},
								{
									types.StringElement{
										Content: "next lines",
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("external link with more text afterwards", func() {
				source := `a link to https://foo.bar and more text`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "a link to "},
									types.InlineLink{
										Location: types.Location{
											Scheme: "https://",
											Path: []interface{}{
												types.StringElement{
													Content: "foo.bar",
												},
											},
										},
									},
									types.StringElement{Content: " and more text"},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			Context("text attribute with comma", func() {

				It("external link only with text having comma", func() {
					source := `a link to http://website.com[A, B, and C]`
					expected := types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "a link to "},
										types.InlineLink{
											Location: types.Location{
												Scheme: "http://",
												Path: []interface{}{
													types.StringElement{
														Content: "website.com",
													},
												},
											},
											Attributes: types.Attributes{
												"positional-1": []interface{}{
													types.StringElement{
														Content: "A",
													},
												},
												"positional-2": []interface{}{
													types.StringElement{
														Content: " B",
													},
												},
												"positional-3": []interface{}{
													types.StringElement{
														Content: " and C",
													},
												},
											},
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("external link only with doublequoted text having comma", func() {
					source := `a link to http://website.com["A, B, and C"]`
					expected := types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "a link to "},
										types.InlineLink{
											Location: types.Location{
												Scheme: "http://",
												Path: []interface{}{
													types.StringElement{
														Content: "website.com",
													},
												},
											},
											Attributes: types.Attributes{
												"positional-1": []interface{}{
													types.StringElement{
														Content: "A, B, and C",
													},
												},
											},
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("external link with doublequoted text having comma and other attrs", func() {
					source := `a link to http://website.com["A, B, and C", role=foo]`
					expected := types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "a link to "},
										types.InlineLink{
											Location: types.Location{
												Scheme: "http://",
												Path: []interface{}{
													types.StringElement{
														Content: "website.com",
													},
												},
											},
											Attributes: types.Attributes{
												"positional-1": []interface{}{
													types.StringElement{
														Content: "A, B, and C",
													},
												},
												"role": "foo",
											},
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("external link with text having comma and other attributes", func() {
					source := `a link to http://website.com[A, B, and C, role=foo]`
					expected := types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "a link to "},
										types.InlineLink{
											Location: types.Location{
												Scheme: "http://",
												Path: []interface{}{
													types.StringElement{
														Content: "website.com",
													},
												},
											},
											Attributes: types.Attributes{
												"positional-1": []interface{}{
													types.StringElement{
														Content: "A",
													},
												},
												"positional-2": []interface{}{
													types.StringElement{
														Content: " B",
													},
												},
												"positional-3": []interface{}{
													types.StringElement{
														Content: " and C",
													},
												},
												"role": "foo",
											},
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})
			})

			It("external link with special characters", func() {
				source := "a link to https://foo*_.com"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "a link to "},
									types.InlineLink{
										Location: types.Location{
											Scheme: "https://",
											Path: []interface{}{
												types.StringElement{
													Content: "foo*_.com",
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("external link with quoted text", func() {
				source := "a link to https://foo.com[_a_ *b* `c`]"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "a link to "},
									types.InlineLink{
										Attributes: types.Attributes{
											"positional-1": []interface{}{
												types.QuotedText{
													Kind: types.Italic,
													Elements: []interface{}{
														types.StringElement{
															Content: "a",
														},
													},
												},
												types.StringElement{
													Content: " ",
												},
												types.QuotedText{
													Kind: types.Bold,
													Elements: []interface{}{
														types.StringElement{
															Content: "b",
														},
													},
												},
												types.StringElement{
													Content: " ",
												},
												types.QuotedText{
													Kind: types.Monospace,
													Elements: []interface{}{
														types.StringElement{
															Content: "c",
														},
													},
												},
											},
										},
										Location: types.Location{
											Scheme: "https://",
											Path: []interface{}{
												types.StringElement{
													Content: "foo.com",
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("external link in bold text", func() {
				source := `a link to *https://foo.com[]*`

				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "a link to "},
									types.QuotedText{
										Kind: types.Bold,
										Elements: []interface{}{
											types.InlineLink{
												Location: types.Location{
													Scheme: "https://",
													Path: []interface{}{
														types.StringElement{
															Content: "foo.com",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("external link with special characters", func() {
				source := "a link to https://foo*_.com"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "a link to "},
									types.InlineLink{
										Location: types.Location{
											Scheme: "https://",
											Path: []interface{}{
												types.StringElement{
													Content: "foo*_.com",
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("external link in bold text", func() {
				source := `a link to *https://foo.com[]*`

				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "a link to "},
									types.QuotedText{
										Kind: types.Bold,
										Elements: []interface{}{
											types.InlineLink{
												Location: types.Location{
													Scheme: "https://",
													Path: []interface{}{
														types.StringElement{
															Content: "foo.com",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("external link in italic text", func() {
				source := `a link to _https://foo.com[]_`

				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "a link to "},
									types.QuotedText{
										Kind: types.Italic,
										Elements: []interface{}{
											types.InlineLink{
												Location: types.Location{
													Scheme: "https://",
													Path: []interface{}{
														types.StringElement{
															Content: "foo.com",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			Context("with document attribute substitutions", func() {

				It("external link with a document attribute substitution for the whole URL", func() {
					source := `
:url: https://foo.bar
:url: https://foo2.bar
	
a link to {url}`

					expected := types.Document{
						Attributes: types.Attributes{
							"url": "https://foo2.bar", // overridden by second declaration
						},
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{
											Content: "a link to ",
										},
										types.InlineLink{
											Location: types.Location{
												Scheme: "https://",
												Path: []interface{}{
													types.StringElement{
														Content: "foo2.bar",
													},
												},
											},
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("external link with two document attribute substitutions only", func() {
					source := `
:scheme: https
:path: foo.bar
	
a link to {scheme}://{path} and https://foo.baz`

					expected := types.Document{
						Attributes: types.Attributes{
							"scheme": "https",
							"path":   "foo.bar",
						},
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "a link to "},
										types.InlineLink{
											Location: types.Location{
												Scheme: "https://",
												Path: []interface{}{
													types.StringElement{
														Content: "foo.bar",
													},
												},
											},
										},
										types.StringElement{Content: " and "},
										types.InlineLink{
											Location: types.Location{
												Scheme: "https://",
												Path: []interface{}{
													types.StringElement{
														Content: "foo.baz",
													},
												},
											},
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("external link with two document attribute substitutions in bold text", func() {
					source := `
:scheme: https
:path: foo.bar
	
a link to *{scheme}://{path}[] and https://foo.baz[]*`

					expected := types.Document{
						Attributes: types.Attributes{
							"scheme": "https",
							"path":   "foo.bar",
						},
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "a link to "},
										types.QuotedText{
											Kind: types.Bold,
											Elements: []interface{}{
												types.InlineLink{
													Location: types.Location{
														Scheme: "https://",
														Path: []interface{}{
															types.StringElement{
																Content: "foo.bar",
															},
														},
													},
												},
												types.StringElement{
													Content: " and ",
												},
												types.InlineLink{
													Location: types.Location{
														Scheme: "https://",
														Path: []interface{}{
															types.StringElement{
																Content: "foo.baz",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("external link with two document attribute substitutions and a reset", func() {
					source := `
:scheme: https
:path: foo.bar
	
:!path:
	
a link to {scheme}://{path} and https://foo.baz`

					expected := types.Document{
						Attributes: types.Attributes{
							"scheme": "https",
							"path":   nil,
						},
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "a link to "},
										types.InlineLink{
											Location: types.Location{
												Scheme: "https://",
												Path: []interface{}{
													types.StringElement{
														Content: "{path}",
													},
												},
											},
										},
										types.StringElement{Content: " and "},
										types.InlineLink{
											Location: types.Location{
												Scheme: "https://",
												Path: []interface{}{
													types.StringElement{
														Content: "foo.baz",
													},
												},
											},
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("external link with document attribute in section 0 title", func() {
					source := `= a title to {scheme}://{path} and https://foo.baz
:scheme: https
:path: foo.bar`

					title := []interface{}{
						types.StringElement{Content: "a title to "},
						types.InlineLink{
							Location: types.Location{
								Scheme: "https://",
								Path: []interface{}{
									types.StringElement{
										Content: "foo.bar",
									},
								},
							},
						},
						types.StringElement{Content: " and "},
						types.InlineLink{
							Location: types.Location{
								Scheme: "https://",
								Path: []interface{}{
									types.StringElement{
										Content: "foo.baz",
									},
								},
							},
						},
					}
					expected := types.Document{
						Attributes: types.Attributes{
							"scheme": "https",
							"path":   "foo.bar",
						},
						ElementReferences: types.ElementReferences{
							"_a_title_to_httpsfoo_bar_and_httpsfoo_baz": title,
						},
						Elements: []interface{}{
							types.Section{
								Level: 0,
								Attributes: types.Attributes{
									types.AttrID: "_a_title_to_httpsfoo_bar_and_httpsfoo_baz",
								},
								Title:    title,
								Elements: []interface{}{},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("external link with document attribute in section 1 title", func() {
					source := `:scheme: https
:path: foo.bar
	
== a title to {scheme}://{path} and https://foo.baz`

					title := []interface{}{
						types.StringElement{Content: "a title to "},
						types.InlineLink{
							Location: types.Location{
								Scheme: "https://",
								Path: []interface{}{
									types.StringElement{
										Content: "foo.bar",
									},
								},
							},
						},
						types.StringElement{Content: " and "},
						types.InlineLink{
							Location: types.Location{
								Scheme: "https://",
								Path: []interface{}{
									types.StringElement{
										Content: "foo.baz",
									},
								},
							},
						},
					}
					expected := types.Document{
						Attributes: types.Attributes{
							"scheme": "https",
							"path":   "foo.bar",
						},
						ElementReferences: types.ElementReferences{
							"_a_title_to_httpsfoo_bar_and_httpsfoo_baz": title,
						},
						Elements: []interface{}{
							types.Section{
								Level: 1,
								Attributes: types.Attributes{
									types.AttrID: "_a_title_to_httpsfoo_bar_and_httpsfoo_baz",
								},
								Title:    title,
								Elements: []interface{}{},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})
			})
		})

		Context("relative links", func() {

			It("relative link to doc without text", func() {
				source := "a link to link:foo.adoc[]"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "a link to "},
									types.InlineLink{
										Location: types.Location{
											Scheme: "",
											Path: []interface{}{
												types.StringElement{
													Content: "foo.adoc",
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("relative link to doc with text", func() {
				source := "a link to link:foo.adoc[foo doc]"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "a link to "},
									types.InlineLink{
										Location: types.Location{
											Scheme: "",
											Path: []interface{}{
												types.StringElement{
													Content: "foo.adoc",
												},
											},
										},
										Attributes: types.Attributes{
											"positional-1": []interface{}{
												types.StringElement{
													Content: "foo doc",
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("relative link to external URL with text only", func() {
				source := "a link to link:https://foo.bar[foo doc]"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "a link to "},
									types.InlineLink{
										Location: types.Location{
											Scheme: "https://",
											Path: []interface{}{
												types.StringElement{
													Content: "foo.bar",
												},
											},
										},
										Attributes: types.Attributes{
											"positional-1": []interface{}{
												types.StringElement{
													Content: "foo doc",
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("relative link to external URL with text and extra attributes", func() {
				source := "a link to link:https://foo.bar[foo doc, foo=bar]"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "a link to "},
									types.InlineLink{
										Location: types.Location{
											Scheme: "https://",
											Path: []interface{}{
												types.StringElement{
													Content: "foo.bar",
												},
											},
										},
										Attributes: types.Attributes{
											"positional-1": []interface{}{
												types.StringElement{
													Content: "foo doc",
												},
											},
											"foo": "bar",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("relative link to external URL with extra attributes only", func() {
				source := "a link to link:https://foo.bar[foo=bar]"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "a link to "},
									types.InlineLink{
										Location: types.Location{
											Scheme: "https://",
											Path: []interface{}{
												types.StringElement{
													Content: "foo.bar",
												},
											},
										},
										Attributes: types.Attributes{
											"foo": "bar",
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("invalid relative link to doc", func() {
				source := "a link to link:foo.adoc"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{
									Content: "a link to link:foo.adoc",
								},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("relative link with quoted text attribute", func() {
				source := "link:/[a _a_ b *b* c `c`]"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.InlineLink{
									Location: types.Location{
										Scheme: "",
										Path: []interface{}{
											types.StringElement{
												Content: "/",
											},
										},
									},
									Attributes: types.Attributes{
										"positional-1": []interface{}{
											types.StringElement{
												Content: "a ",
											},
											types.QuotedText{
												Kind: types.Italic,
												Elements: []interface{}{
													types.StringElement{
														Content: "a",
													},
												},
											},
											types.StringElement{
												Content: " b ",
											},
											types.QuotedText{
												Kind: types.Bold,
												Elements: []interface{}{
													types.StringElement{
														Content: "b",
													},
												},
											},
											types.StringElement{
												Content: " c ",
											},
											types.QuotedText{
												Kind: types.Monospace,
												Elements: []interface{}{
													types.StringElement{
														Content: "c",
													},
												},
											},
										},
									},
								},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("relative link with all valid characters", func() {
				source := `a link to link:ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789~:/?#@!$&;=()*+,-_.%[as expected]`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{Content: "a link to "},
									types.InlineLink{
										Location: types.Location{
											Scheme: "",
											Path: []interface{}{
												types.StringElement{
													Content: "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789~:/?#@!$",
												},
												types.SpecialCharacter{
													Name: "&",
												},
												types.StringElement{
													Content: ";=()*+,-_.%",
												},
											},
										},
										Attributes: types.Attributes{
											"positional-1": []interface{}{
												types.StringElement{
													Content: "as expected",
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("relative link with encoded space", func() {
				source := `Test 1: link:/test/a b[with space]
Test 2: link:/test/a%20b[with encoded space]`
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "Test 1: link:/test/a b[with space]",
									},
								},
								{
									types.StringElement{Content: "Test 2: "},
									types.InlineLink{
										Location: types.Location{
											Scheme: "",
											Path: []interface{}{
												types.StringElement{
													Content: "/test/a%20b",
												},
											},
										},
										Attributes: types.Attributes{
											"positional-1": []interface{}{
												types.StringElement{
													Content: "with encoded space",
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("relative link with two document attribute substitutions and a reset", func() {
				source := `
:scheme: link
:path: foo.bar

:!path:

a link to {scheme}:{path}[] and https://foo.baz`

				expected := types.Document{
					Attributes: types.Attributes{
						"scheme": "link",
						"path":   nil,
					},
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{
									types.StringElement{
										Content: "a link to ",
									},
									types.InlineLink{
										Location: types.Location{
											Path: []interface{}{
												types.StringElement{
													Content: "{path}",
												},
											},
										},
									},
									types.StringElement{
										Content: " and ",
									},
									types.InlineLink{
										Location: types.Location{
											Scheme: "https://",
											Path: []interface{}{
												types.StringElement{
													Content: "foo.baz",
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			It("relative link within quoted text", func() {
				source := "*link:foo[]*"
				expected := types.Document{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.QuotedText{
									Kind: types.Bold,
									Elements: []interface{}{
										types.InlineLink{
											Location: types.Location{
												Path: []interface{}{
													types.StringElement{
														Content: "foo",
													},
												},
											},
										},
									},
								},
								},
							},
						},
					},
				}
				Expect(ParseDocument(source)).To(MatchDocument(expected))
			})

			Context("text attribute with comma", func() {

				It("relative link only with text having comma", func() {
					source := `a link to link:https://foo.bar[A, B, and C]`
					expected := types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "a link to "},
										types.InlineLink{
											Location: types.Location{
												Scheme: "https://",
												Path: []interface{}{
													types.StringElement{
														Content: "foo.bar",
													},
												},
											},
											Attributes: types.Attributes{
												"positional-1": []interface{}{
													types.StringElement{
														Content: "A",
													},
												},
												"positional-2": []interface{}{
													types.StringElement{
														Content: " B",
													},
												},
												"positional-3": []interface{}{
													types.StringElement{
														Content: " and C",
													},
												},
											},
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("relative link only with doublequoted text having comma", func() {
					source := `a link to link:https://foo.bar["A, B, and C"]`
					expected := types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "a link to "},
										types.InlineLink{
											Location: types.Location{
												Scheme: "https://",
												Path: []interface{}{
													types.StringElement{
														Content: "foo.bar",
													},
												},
											},
											Attributes: types.Attributes{
												"positional-1": []interface{}{
													types.StringElement{
														Content: "A, B, and C",
													},
												},
											},
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("relative link with doublequoted text having comma and other attrs", func() {
					source := `a link to link:https://foo.bar["A, B, and C", role=foo]`
					expected := types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "a link to "},
										types.InlineLink{
											Location: types.Location{
												Scheme: "https://",
												Path: []interface{}{
													types.StringElement{
														Content: "foo.bar",
													},
												},
											},
											Attributes: types.Attributes{
												"positional-1": []interface{}{
													types.StringElement{
														Content: "A, B, and C",
													},
												},
												"role": "foo",
											},
										},
									},
								},
							},
						},
					}
					Expect(ParseDocument(source)).To(MatchDocument(expected))
				})

				It("relative link with text having comma and other attributes", func() {
					source := `a link to link:https://foo.bar[A, B, and C, role=foo]`
					expected := types.Document{
						Elements: []interface{}{
							types.Paragraph{
								Lines: [][]interface{}{
									{
										types.StringElement{Content: "a link to "},
										types.InlineLink{
											Location: types.Location{
												Scheme: "https://",
												Path: []interface{}{
													types.StringElement{
														Content: "foo.bar",
													},
												},
											},
											Attributes: types.Attributes{
												"positional-1": []interface{}{
													types.StringElement{
														Content: "A",
													},
												},
												"positional-2": []interface{}{
													types.StringElement{
														Content: " B",
													},
												},
												"positional-3": []interface{}{
													types.StringElement{
														Content: " and C",
													},
												},
												"role": "foo",
											},
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
})
