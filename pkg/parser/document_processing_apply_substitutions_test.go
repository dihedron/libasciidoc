package parser

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"

	. "github.com/onsi/ginkgo" //nolint golint
	. "github.com/onsi/gomega" //nolint golint
)

var _ = Describe("document attribute substitutions", func() {

	It("should replace with new StringElement on first position", func() {
		// given
		elements := []interface{}{
			types.Paragraph{
				Lines: [][]interface{}{
					{
						types.AttributeSubstitution{
							Name: "foo",
						},
						types.StringElement{
							Content: " and more content.",
						},
					},
				},
			},
		}
		// when
		result, err := applyAttributeSubstitutionsOnElements(elements, types.AttributesWithOverrides{
			Content: map[string]interface{}{
				"foo": "bar",
			},
			Overrides: map[string]string{},
		})
		// then
		Expect(err).To(Not(HaveOccurred()))
		Expect(result).To(Equal([]interface{}{
			types.Paragraph{
				Lines: [][]interface{}{
					{
						types.StringElement{
							Content: "bar and more content.",
						},
					},
				},
			},
		}))
	})

	It("should replace with new StringElement on middle position", func() {
		// given
		elements := []interface{}{
			types.Paragraph{
				Lines: [][]interface{}{
					{
						types.StringElement{
							Content: "baz, ",
						},
						types.AttributeSubstitution{
							Name: "foo",
						},
						types.StringElement{
							Content: " and more content.",
						},
					},
				},
			},
		}
		// when
		result, err := applyAttributeSubstitutionsOnElements(elements, types.AttributesWithOverrides{
			Content: map[string]interface{}{
				"foo": "bar",
			},
			Overrides: map[string]string{},
		})
		// then
		Expect(err).To(Not(HaveOccurred()))
		Expect(result).To(Equal([]interface{}{
			types.Paragraph{
				Lines: [][]interface{}{
					{
						types.StringElement{
							Content: "baz, bar and more content.",
						},
					},
				},
			},
		}))
	})

	It("should replace with undefined attribute", func() {
		// given
		elements := []interface{}{
			types.Paragraph{
				Lines: [][]interface{}{
					{
						types.StringElement{
							Content: "baz, ",
						},
						types.AttributeSubstitution{
							Name: "foo",
						},
						types.StringElement{
							Content: " and more content.",
						},
					},
				},
			},
		}
		// when
		result, err := applyAttributeSubstitutionsOnElements(elements, types.AttributesWithOverrides{
			Content:   map[string]interface{}{},
			Overrides: map[string]string{},
		})

		// then
		Expect(err).To(Not(HaveOccurred()))
		Expect(result).To(Equal([]interface{}{
			types.Paragraph{
				Lines: [][]interface{}{
					{
						types.StringElement{
							Content: "baz, {foo} and more content.",
						},
					},
				},
			},
		}))
	})

	It("should merge without substitution", func() {
		// given
		elements := []interface{}{
			types.Paragraph{
				Lines: [][]interface{}{
					{
						types.StringElement{
							Content: "baz, ",
						},
						types.StringElement{
							Content: "foo",
						},
						types.StringElement{
							Content: " and more content.",
						},
					},
				},
			},
		}
		// when
		result, err := applyAttributeSubstitutionsOnElements(elements, types.AttributesWithOverrides{
			Content:   map[string]interface{}{},
			Overrides: map[string]string{},
		})

		// then
		Expect(err).To(Not(HaveOccurred()))
		Expect(result).To(Equal([]interface{}{
			types.Paragraph{
				Lines: [][]interface{}{
					{
						types.StringElement{
							Content: "baz, foo and more content.",
						},
					},
				},
			},
		}))
	})

	It("should replace with new link", func() {
		// given
		elements := []interface{}{
			types.Paragraph{
				Lines: [][]interface{}{
					{
						types.StringElement{
							Content: "a link to ",
						},
						types.AttributeSubstitution{
							Name: "scheme",
						},
						types.StringElement{
							Content: "://",
						},
						types.AttributeSubstitution{
							Name: "host",
						},
						types.StringElement{
							Content: "[].", // explicit use of `[]` to avoid grabbing the `.`
						},
					},
				},
			},
		}
		// when
		result, err := applyAttributeSubstitutionsOnElements(elements, types.AttributesWithOverrides{
			Content: map[string]interface{}{
				"foo":    "bar",
				"scheme": "https",
				"host":   "foo.bar",
			},
			Overrides: map[string]string{},
		})

		// then
		Expect(err).To(Not(HaveOccurred()))
		Expect(result).To(Equal([]interface{}{
			types.Paragraph{
				Lines: [][]interface{}{
					{
						types.StringElement{
							Content: "a link to https://foo.bar[].",
						},
					},
				},
			},
		}))
	})

	Context("list items", func() {

		It("should replace with new StringElement in ordered list item", func() {
			// given
			elements := []interface{}{
				types.OrderedListItem{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.AttributeSubstitution{
									Name: "foo",
								},
									types.StringElement{
										Content: " and more content.",
									},
								},
							},
						},
					},
				},
			}
			// when
			result, err := applyAttributeSubstitutionsOnElements(elements, types.AttributesWithOverrides{
				Content: map[string]interface{}{
					"foo": "bar",
				},
				Overrides: map[string]string{},
			})
			// then
			Expect(err).To(Not(HaveOccurred()))
			Expect(result).To(Equal([]interface{}{
				types.OrderedListItem{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{
									Content: "bar and more content.",
								},
								},
							},
						},
					},
				},
			}))
		})

		It("should replace with new StringElement in unordered list item", func() {
			// given
			elements := []interface{}{
				types.UnorderedListItem{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.AttributeSubstitution{
									Name: "foo",
								},
									types.StringElement{
										Content: " and more content.",
									},
								},
							},
						},
					},
				},
			}
			// when
			result, err := applyAttributeSubstitutionsOnElements(elements, types.AttributesWithOverrides{
				Content: map[string]interface{}{
					"foo": "bar",
				},
				Overrides: map[string]string{},
			})
			// then
			Expect(err).To(Not(HaveOccurred()))
			Expect(result).To(Equal([]interface{}{
				types.UnorderedListItem{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{
									Content: "bar and more content.",
								},
								},
							},
						},
					},
				},
			}))
		})

		It("should replace with new StringElement in labeled list item", func() {
			// given
			elements := []interface{}{
				types.LabeledListItem{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.AttributeSubstitution{
									Name: "foo",
								},
									types.StringElement{
										Content: " and more content.",
									},
								},
							},
						},
					},
				},
			}
			// when
			result, err := applyAttributeSubstitutionsOnElements(elements, types.AttributesWithOverrides{
				Content: map[string]interface{}{
					"foo": "bar",
				},
				Overrides: map[string]string{},
			})
			// then
			Expect(err).To(Not(HaveOccurred()))
			Expect(result).To(Equal([]interface{}{
				types.LabeledListItem{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{
									Content: "bar and more content.",
								},
								},
							},
						},
					},
				},
			}))
		})
	})

	Context("delimited blocks", func() {

		It("should replace with new StringElement in delimited block", func() {
			// given
			elements := []interface{}{
				types.ExampleBlock{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.AttributeSubstitution{
									Name: "foo",
								},
									types.StringElement{
										Content: " and more content.",
									},
								},
							},
						},
					},
				},
			}
			// when
			result, err := applyAttributeSubstitutionsOnElements(elements, types.AttributesWithOverrides{
				Content: map[string]interface{}{
					"foo": "bar",
				},
				Overrides: map[string]string{},
			})
			// then
			Expect(err).To(Not(HaveOccurred()))
			Expect(result).To(Equal([]interface{}{
				types.ExampleBlock{
					Elements: []interface{}{
						types.Paragraph{
							Lines: [][]interface{}{
								{types.StringElement{
									Content: "bar and more content.",
								},
								},
							},
						},
					},
				},
			}))
		})
	})

	Context("quoted texts", func() {

		It("should replace with new StringElement in quoted text", func() {
			// given
			elements := []interface{}{
				types.Paragraph{
					Lines: [][]interface{}{
						{
							types.StringElement{
								Content: "hello ",
							},
							types.QuotedText{
								Elements: []interface{}{
									types.AttributeSubstitution{
										Name: "foo",
									},
									types.StringElement{
										Content: " and more content.",
									},
								},
							},
						},
						{
							types.StringElement{
								Content: "and another line",
							},
						},
					},
				},
			}
			// when
			result, err := applyAttributeSubstitutionsOnElements(elements, types.AttributesWithOverrides{
				Content: map[string]interface{}{
					"foo": "bar",
				},
				Overrides: map[string]string{},
			})
			// then
			Expect(err).To(Not(HaveOccurred()))
			Expect(result).To(Equal([]interface{}{
				types.Paragraph{
					Lines: [][]interface{}{
						{
							types.StringElement{
								Content: "hello ",
							},
							types.QuotedText{
								Elements: []interface{}{
									types.StringElement{
										Content: "bar and more content.",
									},
								},
							},
						},
						{
							types.StringElement{
								Content: "and another line",
							},
						},
					},
				},
			}))
		})
	})

	Context("tables", func() {

		It("should replace with new StringElement in table cell", func() {
			// given
			elements := []interface{}{
				types.ListingBlock{
					Lines: [][]interface{}{
						{
							types.AttributeSubstitution{
								Name: "foo",
							},
							types.StringElement{
								Content: " and more content.",
							},
						},
					},
				},
			}
			// when
			result, err := applyAttributeSubstitutionsOnElements(elements, types.AttributesWithOverrides{
				Content: map[string]interface{}{
					"foo": "bar",
				},
				Overrides: map[string]string{},
			})
			// then
			Expect(err).To(Not(HaveOccurred()))
			Expect(result).To(Equal([]interface{}{
				types.ListingBlock{
					Lines: [][]interface{}{
						{
							types.StringElement{
								Content: "bar and more content.",
							},
						},
					},
				},
			}))
		})
	})

	Context("attribute overrides", func() {

		It("should replace with new StringElement on first position", func() {
			// given
			elements := []interface{}{
				types.AttributeDeclaration{
					Name:  "foo",
					Value: "foo",
				},
				types.AttributeReset{
					Name: "foo",
				},
				types.Paragraph{
					Lines: [][]interface{}{
						{
							types.AttributeSubstitution{
								Name: "foo",
							},
							types.StringElement{
								Content: " and more content.",
							},
						},
					},
				},
			}
			// when
			result, err := applyAttributeSubstitutionsOnElements(elements, types.AttributesWithOverrides{
				Content: map[string]interface{}{
					"foo": "bar",
				},
				Overrides: map[string]string{
					"foo": "BAR",
				},
			})
			// then
			Expect(err).To(Not(HaveOccurred()))
			Expect(result).To(Equal([]interface{}{ // at this stage, AttributeDeclaration and AttributeReset are still present
				types.AttributeDeclaration{
					Name:  "foo",
					Value: "foo",
				},
				types.AttributeReset{
					Name: "foo",
				},
				types.Paragraph{
					Lines: [][]interface{}{
						{
							types.StringElement{
								Content: "BAR and more content.",
							},
						},
					},
				},
			}))
		})
	})

	Context("counters", func() {

		It("should start at one", func() {
			// given
			elements := []interface{}{
				types.CounterSubstitution{
					Name: "foo",
				},
			}
			result, err := applyAttributeSubstitutionsOnElements(elements, types.AttributesWithOverrides{
				Content:   map[string]interface{}{},
				Overrides: map[string]string{},
				Counters:  map[string]interface{}{},
			})
			Expect(err).To(Not(HaveOccurred()))
			Expect(result).To(Equal([]interface{}{ // at this stage, AttributeDeclaration and AttributeReset are still present
				types.StringElement{
					Content: "1",
				},
			}))
		})

		It("should increment correctly", func() {
			// given
			elements := []interface{}{
				types.CounterSubstitution{
					Name: "foo",
				},
				types.CounterSubstitution{
					Name: "bar",
				},
				types.CounterSubstitution{
					Name: "foo",
				},
				types.CounterSubstitution{
					Name:   "alpha",
					Value:  'a',
					Hidden: true,
				},
				types.CounterSubstitution{
					Name: "alpha",
				},
				types.CounterSubstitution{
					Name:   "set",
					Value:  33,
					Hidden: true,
				},
				types.CounterSubstitution{
					Name:   "set",
					Hidden: true,
				},
				types.CounterSubstitution{
					Name: "set",
				},
			}
			result, err := applyAttributeSubstitutionsOnElements(elements, types.AttributesWithOverrides{
				Content:   map[string]interface{}{},
				Overrides: map[string]string{},
				Counters:  map[string]interface{}{},
			})
			Expect(err).To(Not(HaveOccurred()))
			Expect(result).To(Equal([]interface{}{ // at this stage, AttributeDeclaration and AttributeReset are still present
				types.StringElement{
					Content: "112b35", // elements get concatenated
				},
			}))
		})
	})
})

var _ = Describe("substitution funcs", func() {

	It("should append sub", func() {
		// given"
		f := funcs{"attributes", "quotes"}
		// when
		f = f.append("macros")
		// then
		Expect(f).To(Equal(funcs{"attributes", "quotes", "macros"}))
	})

	It("should append subs", func() {
		// given"
		f := funcs{"attributes"}
		// when
		f = f.append("quotes", "macros")
		// then
		Expect(f).To(Equal(funcs{"attributes", "quotes", "macros"}))
	})

	It("should prepend sub", func() {
		// given"
		f := funcs{"attributes", "quotes"}
		// when
		f = f.prepend("macros")
		// then
		Expect(f).To(Equal(funcs{"macros", "attributes", "quotes"}))
	})

	It("should remove first sub", func() {
		// given"
		f := funcs{"attributes", "quotes", "macros"}
		// when
		f = f.remove("attributes")
		// then
		Expect(f).To(Equal(funcs{"quotes", "macros"}))
	})

	It("should remove middle sub", func() {
		// given"
		f := funcs{"attributes", "quotes", "macros"}
		// when
		f = f.remove("quotes")
		// then
		Expect(f).To(Equal(funcs{"attributes", "macros"}))
	})

	It("should remove last sub", func() {
		// given"
		f := funcs{"attributes", "quotes", "macros"}
		// when
		f = f.remove("macros")
		// then
		Expect(f).To(Equal(funcs{"attributes", "quotes"}))
	})

	It("should remove non existinge", func() {
		// given"
		f := funcs{"attributes", "quotes", "macros"}
		// when
		f = f.remove("other")
		// then
		Expect(f).To(Equal(funcs{"attributes", "quotes", "macros"}))
	})
})
