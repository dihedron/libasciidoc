package parser

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/configuration"
	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ----------------------------------------------------------------------------
// Substitutions
// ----------------------------------------------------------------------------

// ApplySubstitutions applies all the substitutions on delimited blocks, standalone paragraphs and paragraphs
// in continued list items, and then attribute substitutions, and as a result returns a `DraftDocument`.
func ApplySubstitutions(rawDoc types.RawDocument, config configuration.Configuration) (types.DraftDocument, error) {
	attrs := types.AttributesWithOverrides{
		Content:   types.Attributes{},
		Overrides: config.AttributeOverrides,
		Counters:  map[string]interface{}{},
	}
	// also, add all front-matter key/values
	attrs.Add(rawDoc.FrontMatter.Content)
	// also, add all AttributeDeclaration at the top of the document
	attrs.Add(rawDoc.Attributes())

	elements, err := applySubstitutions(rawDoc.Elements, attrs)
	if err != nil {
		return types.DraftDocument{}, err
	}
	if len(elements) == 0 {
		elements = nil // avoid carrying empty slice
	}
	return types.DraftDocument{
		Attributes:  attrs.All(),
		FrontMatter: rawDoc.FrontMatter,
		Elements:    elements,
	}, nil
}

// applySubstitutions applies the substitutions on paragraphs and delimited blocks (including when in continued list elements)
func applySubstitutions(elements []interface{}, attrs types.AttributesWithOverrides) ([]interface{}, error) {
	if len(elements) == 0 {
		return nil, nil
	}
	result := make([]interface{}, 0, len(elements))
	for _, e := range elements {
		var err error
		switch e := e.(type) {
		case types.BlockWithElementSubstitution:
			subs, err := substitutionsFor(e)
			if err != nil {
				return nil, err
			}
			elements, err := applySubstitutionsOnElements(e.ElementsToSubstitute(), subs, attrs)
			if err != nil {
				return nil, err
			}
			result = append(result, e.ReplaceElements(elements))
		case types.BlockWithLineSubstitution:
			subs, err := substitutionsFor(e)
			if err != nil {
				return nil, err
			}
			elements, err := applySubstitutionsOnLines(e.LinesToSubstitute(), subs, attrs)
			if err != nil {
				return nil, err
			}
			result = append(result, e.ReplaceLines(elements))
		case types.MarkdownQuoteBlock: // slightly different since there is an extraction for the author attributions
			e, err := applySubstitutionsOnMarkdownQuoteBlock(e, attrs)
			if err != nil {
				return nil, err
			}
			result = append(result, e)
		case types.ContinuedListItemElement:
			r, err := applySubstitutions([]interface{}{e.Element}, attrs)
			if err != nil {
				return nil, err
			}
			e.Element = r[0]
			result = append(result, e)
		case types.ImageBlock:
			if e, err = applySubstitutionsOnImageBlock(e, attrs); err != nil {
				return nil, err
			}
			result = append(result, e)
		case types.Section:
			if e, err = applySubstitutionsOnSection(e, attrs); err != nil {
				return nil, err
			}
			result = append(result, e)
		default:
			// no support for element substitution here
			// so let's proceed with attribute substitutions
			if e, err = applyAttributeSubstitutionsOnElement(e, attrs); err != nil {
				return nil, err
			}
			result = append(result, e)
		}
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("after all substitutions:")
		spew.Fdump(log.StandardLogger().Out, result)
	}
	return result, nil
}

// ----------------------------------------------------------------------------
// Delimited Block substitutions
// ----------------------------------------------------------------------------

var substitutions = map[string]elementsSubstitution{
	"inline_passthrough": substituteInlinePassthrough,
	"callouts":           substituteCallouts,
	"specialcharacters":  substituteSpecialCharacters,
	"specialchars":       substituteSpecialCharacters,
	"quotes":             substituteQuotedTexts,
	"attributes":         substituteAttributes,
	"replacements":       substituteReplacements,
	"macros":             substituteInlineMacros,
	"post_replacements":  substitutePostReplacements,
	"none":               substituteNone,
}

// blocks of blocks
var defaultSubstitutionsForBlockElements = []string{
	"inline_passthrough",
	"specialcharacters",
	"quotes",
	"attributes",
	"replacements",
	"macros",
	"post_replacements",
}
var defaultExampleBlockSubstitutions = defaultSubstitutionsForBlockElements
var defaultQuoteBlockSubstitutions = defaultSubstitutionsForBlockElements
var defaultSidebarBlockSubstitutions = defaultSubstitutionsForBlockElements
var defaultVerseBlockSubstitutions = defaultSubstitutionsForBlockElements // even though it's a block of lines, not a block of blocks
var defaultParagraphSubstitutions = defaultSubstitutionsForBlockElements  // even though it's a block of lines, not a block of blocks

// blocks of lines
var defaultSubstitutionsForBlockLines = []string{
	"callouts", // must be executed before "specialcharacters"
	"specialcharacters",
}
var defaultFencedBlockSubstitutions = defaultSubstitutionsForBlockLines
var defaultListingBlockSubstitutions = defaultSubstitutionsForBlockLines
var defaultLiteralBlockSubstitutions = defaultSubstitutionsForBlockLines

// other blocks
var defaultPassthroughBlockSubstitutions = []string{}
var defaultCommentBlockSubstitutions = []string{"none"}

func applySubstitutionsOnMarkdownQuoteBlock(b types.MarkdownQuoteBlock, attrs types.AttributesWithOverrides) (types.MarkdownQuoteBlock, error) {
	funcs := []elementsSubstitution{
		substituteInlinePassthrough,
		substituteSpecialCharacters,
		substituteQuotedTexts,
		substituteAttributes,
		substituteReplacements,
		substituteInlineMacros,
		substitutePostReplacements,
		splitLines}
	// attempt to extract the block attributions
	var author string
	if b.Lines, author = extractMarkdownQuoteAttribution(b.Lines); author != "" {
		if b.Attributes == nil {
			b.Attributes = types.Attributes{}
		}
		b.Attributes.Set(types.AttrQuoteAuthor, author)
	}
	if len(b.Lines) == 0 { // no more line to parse after extracting the author
		b.Lines = nil
		return b, nil
	}
	// apply all the substitutions
	var err error
	for _, sub := range funcs {
		if b.Lines, err = sub(b.Lines, attrs); err != nil {
			return types.MarkdownQuoteBlock{}, err
		}
	}
	return b, nil
}

func extractMarkdownQuoteAttribution(lines [][]interface{}) ([][]interface{}, string) {
	log.Debug("extracting attribution on markdown block quote")
	// first, check if last line is an attribution (author)
	if len(lines) == 0 {
		return lines, ""
	}
	if l, ok := lines[len(lines)-1][0].(types.StringElement); ok {
		a, err := ParseReader("", strings.NewReader(l.Content), Entrypoint("MarkdownQuoteAttribution"))
		// assume that the last line is not an author attribution if an error occurred
		if err != nil {
			return lines, ""
		}
		if a, ok := a.(string); ok {
			log.Debugf("found attribution in markdown block: '%s'", a)
			return lines[:len(lines)-1], a
		}
	}
	return lines, ""
}

type funcs []string

func (f funcs) append(others ...string) funcs {
	return append(f, others...)
}

func (f funcs) prepend(other string) funcs {
	return append(funcs{other}, f...)
}

func (f funcs) remove(other string) funcs {
	for i, s := range f {
		if s == other {
			return append(f[:i], f[i+1:]...)
		}
	}
	// unchanged
	return f
}

func substitutionsFor(block types.BlockWithSubstitution) ([]elementsSubstitution, error) {
	subs := funcs{}
	for _, s := range strings.Split(block.SubstitutionsToApply(), ",") {
		switch s {
		case "":
			subs = subs.append(defaultSubstitutionsFor(block)...)
		case "normal":
			subs = subs.append(
				"inline_passthrough",
				"specialcharacters",
				"quotes",
				"attributes",
				"replacements",
				"macros",
				"post_replacements",
			)
		case "callouts", "specialcharacters", "specialchars", "quotes", "attributes", "macros", "replacements", "post_replacements", "none":
			subs = subs.append(s)
		case "+callouts", "+specialcharacters", "+specialchars", "+quotes", "+attributes", "+macros", "+replacements", "+post_replacements", "+none":
			if len(subs) == 0 {
				subs = subs.append(defaultSubstitutionsFor(block)...)
			}
			subs = subs.append(strings.ReplaceAll(s, "+", ""))
		case "callouts+", "specialcharacters+", "specialchars+", "quotes+", "attributes+", "macros+", "replacements+", "post_replacements+", "none+":
			if len(subs) == 0 {
				subs = subs.append(defaultSubstitutionsFor(block)...)
			}
			subs = subs.prepend(strings.ReplaceAll(s, "+", ""))
		case "-callouts", "-specialcharacters", "-specialchars", "-quotes", "-attributes", "-macros", "-replacements", "-post_replacements", "-none":
			if len(subs) == 0 {
				subs = subs.append(defaultSubstitutionsFor(block)...)
			}
			subs = subs.remove(strings.ReplaceAll(s, "-", ""))
		default:
			return nil, fmt.Errorf("unsupported substitution: '%s", s)
		}
	}
	result := make([]elementsSubstitution, 0, len(subs)+1)
	// result = append(result, substituteInlinePassthrough)
	for _, s := range subs {
		if f, exists := substitutions[s]; exists {
			result = append(result, f)
		}
	}
	result = append(result, splitLines)
	return result, nil
}

func defaultSubstitutionsFor(block interface{}) []string {
	switch b := block.(type) {
	case types.ExampleBlock:
		return defaultExampleBlockSubstitutions
	case types.QuoteBlock:
		return defaultQuoteBlockSubstitutions
	case types.SidebarBlock:
		return defaultSidebarBlockSubstitutions
	case types.FencedBlock:
		return defaultFencedBlockSubstitutions
	case types.ListingBlock:
		return defaultListingBlockSubstitutions
	case types.VerseBlock:
		return defaultVerseBlockSubstitutions
	case types.LiteralBlock:
		return defaultLiteralBlockSubstitutions
	case types.PassthroughBlock:
		return defaultPassthroughBlockSubstitutions
	case types.CommentBlock:
		return defaultCommentBlockSubstitutions
	case types.Paragraph:
		// support for masquerading
		// treat 'Listing' paragraphs as verbatim blocks
		if k, exists := b.Attributes[types.AttrBlockKind]; exists {
			switch k {
			case types.Listing:
				return defaultListingBlockSubstitutions
			}
		}
		return defaultParagraphSubstitutions
	default:
		log.Warnf("unsupported substitutions on block of type: '%T'", block)
		return nil
	}
}

func applySubstitutionsOnElements(elements []interface{}, subs []elementsSubstitution, attrs types.AttributesWithOverrides) ([]interface{}, error) {
	// var err error
	// apply all the substitutions on blocks that need to be processed
	for i, element := range elements {
		log.Debugf("applying substitution on element of type '%T'", element)
		switch e := element.(type) {
		// if the block contains a block...
		case types.BlockWithElementSubstitution:
			lines, err := applySubstitutionsOnElements(e.ElementsToSubstitute(), subs, attrs)
			if err != nil {
				return nil, err
			}
			elements[i] = e.ReplaceElements(lines)
		case types.BlockWithLineSubstitution:
			lines, err := applySubstitutionsOnLines(e.LinesToSubstitute(), subs, attrs)
			if err != nil {
				return nil, err
			}
			elements[i] = e.ReplaceLines(lines)
		default:
			log.Debugf("nothing to substitute on element of type '%T'", element)
			// do nothing
		}
	}
	return elements, nil
}

func applySubstitutionsOnLines(lines [][]interface{}, subs []elementsSubstitution, attrs types.AttributesWithOverrides) ([][]interface{}, error) {
	var err error
	for _, sub := range subs {
		if lines, err = sub(lines, attrs); err != nil {
			return nil, err
		}
	}
	return lines, nil
}

// ----------------------------------------------------------------------------
// Section substitutions
// ----------------------------------------------------------------------------

// applies the elements and attributes substitutions on the given section title.
func applySubstitutionsOnSection(s types.Section, attrs types.AttributesWithOverrides) (types.Section, error) {
	elements := [][]interface{}{s.Title} // wrap to match the `elementsSubstitution` arg type
	subs := []elementsSubstitution{
		substituteInlinePassthrough,
		substituteSpecialCharacters,
		substituteQuotedTexts,
		substituteAttributes,
		substituteReplacements,
		substituteInlineMacros,
		substitutePostReplacements,
	}
	var err error
	for _, sub := range subs {
		if elements, err = sub(elements, attrs); err != nil {
			return types.Section{}, err
		}
	}
	s.Title = elements[0]
	if s, err = s.ResolveID(attrs); err != nil {
		return types.Section{}, err
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("section title after substitution:")
		spew.Fdump(log.StandardLogger().Out, s.Title)
	}
	return s, nil
}

// ----------------------------------------------------------------------------
// Image Block substitutions
// ----------------------------------------------------------------------------

// applies the elements and attributes substitutions on the given image block.
func applySubstitutionsOnImageBlock(b types.ImageBlock, attrs types.AttributesWithOverrides) (types.ImageBlock, error) {
	elements := [][]interface{}{b.Location.Path} // wrap to match the `elementsSubstitution` arg type
	subs := []elementsSubstitution{substituteAttributes}
	var err error
	for _, sub := range subs {
		if elements, err = sub(elements, attrs); err != nil {
			return types.ImageBlock{}, err
		}
	}
	b.Location.Path = elements[0]
	b.Location = b.Location.WithPathPrefix(attrs.GetAsStringWithDefault("imagesdir", ""))
	if !b.Attributes.Has(types.AttrImageAlt) {
		alt := filepath.Base(b.Location.Stringify())
		ext := filepath.Ext(alt)
		alt = alt[0 : len(alt)-len(ext)]
		b.Attributes = b.Attributes.Set(types.AttrImageAlt, alt)
	}

	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("image block after substitution:")
		spew.Fdump(log.StandardLogger().Out, b)
	}
	return b, nil
}

// ----------------------------------------------------------------------------
// Individual substitution funcs
// ----------------------------------------------------------------------------

// includes a call to `elementsSubstitution` with some post-processing on the result
var substituteAttributes = func(lines [][]interface{}, attrs types.AttributesWithOverrides) ([][]interface{}, error) {
	lines, err := newElementsSubstitution("AttributeSubs", "AttributeSubs")(lines, attrs)
	if err != nil {
		return nil, err
	}
	for i, line := range lines {
		line, err := applyAttributeSubstitutionsOnElements(line, attrs)
		if err != nil {
			return nil, err
		}
		lines[i] = types.Merge(line)
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("applied the 'attributes' substitution")
		spew.Fdump(log.StandardLogger().Out, lines)
	}
	return lines, nil
}

var (
	substituteInlinePassthrough = newElementsSubstitution("InlinePassthroughSubs", "InlinePassthroughSubs")
	substituteSpecialCharacters = newElementsSubstitution("SpecialCharacterSubs", "SpecialCharacterSubs")
	substituteQuotedTexts       = newElementsSubstitution("QuotedTextSubs", "QuotedTextSubs")
	substituteReplacements      = newElementsSubstitution("ReplacementSubs", "ReplacementSubs")
	substituteInlineMacros      = newElementsSubstitution("InlineMacroSubs", "InlineMacroSubs")
	substitutePostReplacements  = newElementsSubstitution("PostReplacementSubs", "PostReplacementSubs")
	substituteNone              = newElementsSubstitution("NoneSubs", "NoneSubs") // TODO: no need for placeholder support here?
	substituteCallouts          = newElementsSubstitution("CalloutSubs", "CalloutSubs")
)

type elementsSubstitution func(lines [][]interface{}, attrs types.AttributesWithOverrides) ([][]interface{}, error)

func newElementsSubstitution(contentRuleName, placeholderRuleName string) elementsSubstitution {
	return func(lines [][]interface{}, attrs types.AttributesWithOverrides) ([][]interface{}, error) {
		log.Debugf("applying the '%s' rule on elements", contentRuleName)
		placeholders := newPlaceHolders()
		s := serializeLines(lines, placeholders)
		imagesdirOption := GlobalStore("imagesdir", attrs.GetAsStringWithDefault("imagesdir", ""))
		// process placeholder content (eg: quoted text may contain an inline link)
		for ref, placeholder := range placeholders.elements {
			switch placeholder := placeholder.(type) { // TODO: create `PlaceHolder` interface?
			case types.QuotedString:
				var err error
				if placeholder.Elements, err = parserPlaceHolderElements(placeholder.Elements, imagesdirOption, Entrypoint(placeholderRuleName)); err != nil {
					return nil, err
				}
				placeholders.elements[ref] = placeholder
			case types.QuotedText:
				var err error
				if placeholder.Elements, err = parserPlaceHolderElements(placeholder.Elements, imagesdirOption, Entrypoint(placeholderRuleName)); err != nil {
					return nil, err
				}
				placeholders.elements[ref] = placeholder
			}
		}
		result := make([][]interface{}, 0, len(lines))
		elmts, err := parseContent("", s, imagesdirOption, Entrypoint(contentRuleName))
		if err != nil {
			return nil, err
		}
		elmts = restoreElements(elmts, placeholders)
		result = append(result, elmts)
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debugf("applied the '%s' rule:", contentRuleName)
			spew.Fdump(log.StandardLogger().Out, result)
		}
		return result, nil
	}
}

func parserPlaceHolderElements(elements []interface{}, options ...Option) ([]interface{}, error) {
	result := make([]interface{}, 0, len(elements)) // default capacity (but may not be enough)
	for _, element := range elements {
		switch element := element.(type) {
		case types.StringElement:
			elmts, err := parseContent("", element.Content, options...)
			if err != nil {
				return nil, err
			}
			result = append(result, elmts...)
		default:
			result = append(result, element)
		}
	}
	return result, nil
}

func parseContent(filename string, content string, options ...Option) ([]interface{}, error) {
	// log.Debugf("parsing content '%s'", content)
	result, err := ParseReader(filename, strings.NewReader(content), options...)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse '%s'", content)
	}
	if result, ok := result.([]interface{}); ok {
		return types.Merge(result), nil
	}
	return []interface{}{result}, nil
}

// replace the placeholders with their original element in the given elements
func restoreElements(elmts []interface{}, placeholders *placeholders) []interface{} {
	// skip if there's nothing to restore
	if len(placeholders.elements) == 0 {
		return elmts
	}
	for i, elmt := range elmts {
		switch elmt := elmt.(type) {
		case types.ElementPlaceHolder:
			elmts[i] = placeholders.elements[elmt.Ref]
		case types.Paragraph:
			for i, line := range elmt.Lines {
				elmt.Lines[i] = restoreElements(line, placeholders)
			}
			elmts[i] = elmt
		case types.InlineLink: // TODO: use an interface and implement the `restoreElements` func on these types, instead
			elmt.Location.Path = restoreElements(elmt.Location.Path, placeholders)
			elmt.Attributes = restoreAttributes(elmt.Attributes, placeholders)
			elmts[i] = elmt
		case types.QuotedText:
			elmt.Elements = restoreElements(elmt.Elements, placeholders)
			elmt.Attributes = restoreAttributes(elmt.Attributes, placeholders)
			elmts[i] = elmt
		case types.QuotedString:
			elmt.Elements = restoreElements(elmt.Elements, placeholders)
			elmts[i] = elmt
		case types.IndexTerm:
			elmt.Term = restoreElements(elmt.Term, placeholders)
			elmts[i] = elmt
		case types.ExternalCrossReference:
			elmt.Label = restoreElements(elmt.Label, placeholders)
			elmts[i] = elmt
		case types.Footnote:
			elmt.Elements = restoreElements(elmt.Elements, placeholders)
			elmts[i] = elmt
		case types.ElementRole:
			elmts[i] = types.ElementRole(restoreElements(elmt, placeholders))
		case []interface{}:
			elmts[i] = restoreElements(elmt, placeholders)
		default:
			// do nothing, keep elmt as-is
		}
	}
	return elmts
}

// replace the placeholders with their original element in the given attributes
func restoreAttributes(attrs types.Attributes, placeholders *placeholders) types.Attributes {
	for key, value := range attrs {
		switch value := value.(type) {
		case types.ElementPlaceHolder:
			attrs[key] = placeholders.elements[value.Ref]
		case types.ElementRole:
			attrs[key] = types.ElementRole(restoreElements(value, placeholders))
		case []interface{}:
			attrs[key] = restoreElements(value, placeholders)
		}
	}
	return attrs
}

type placeholders struct {
	seq      int
	elements map[string]interface{}
}

func newPlaceHolders() *placeholders {
	return &placeholders{
		seq:      0,
		elements: map[string]interface{}{},
	}
}
func (p *placeholders) add(element interface{}) types.ElementPlaceHolder {
	p.seq++
	p.elements[strconv.Itoa(p.seq)] = element
	return types.ElementPlaceHolder{
		Ref: strconv.Itoa(p.seq),
	}

}

func serializeLines(lines [][]interface{}, placeholders *placeholders) string {
	result := strings.Builder{}
	for i, line := range lines {
		for _, e := range line {
			switch e := e.(type) {
			case types.StringElement:
				result.WriteString(e.Content)
			case types.SingleLineComment:
				// replace with placeholder
				p := placeholders.add(e)
				result.WriteString(p.String())
			default:
				// replace with placeholder
				p := placeholders.add(e)
				result.WriteString(p.String())
			}
		}
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		spew.Fdump(log.StandardLogger().Out, result.String())
	}
	return result.String()
}

func splitLines(lines [][]interface{}, _ types.AttributesWithOverrides) ([][]interface{}, error) {
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("splitting lines on")
		spew.Fdump(log.StandardLogger().Out, lines)
	}
	result := [][]interface{}{}
	for _, line := range lines {
		pendingLine := []interface{}{}
		for _, element := range line {
			switch element := element.(type) {
			case types.StringElement:
				// if content has line feeds, then split in multiple lines
				split := strings.Split(element.Content, "\n")
				for i, s := range split {
					if len(s) > 0 { // no need to append an empty StringElement
						pendingLine = append(pendingLine, types.StringElement{Content: s})
					}
					if i < len(split)-1 {
						result = append(result, pendingLine)
						pendingLine = []interface{}{} // reset for the next line
					}
				}
			default:
				pendingLine = append(pendingLine, element)
			}
		}
		// don't forget the last line (if applicable)
		result = append(result, pendingLine)
	}
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("splitted lines")
		spew.Fdump(log.StandardLogger().Out, result)
	}
	return result, nil
}

// ----------------------------------------------------------------------------
// Attribute substitutions
// ----------------------------------------------------------------------------

func applyAttributeSubstitutionsOnElements(elements []interface{}, attrs types.AttributesWithOverrides) ([]interface{}, error) {
	result := make([]interface{}, 0, len(elements)) // maximum capacity should exceed initial input
	for _, element := range elements {
		e, err := applyAttributeSubstitutionsOnElement(element, attrs)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	result = types.Merge(result)
	return result, nil
}

func applyAttributeSubstitutionsOnLines(lines [][]interface{}, attrs types.AttributesWithOverrides) ([][]interface{}, error) {
	for i, line := range lines {
		line, err := applyAttributeSubstitutionsOnElements(line, attrs)
		if err != nil {
			return nil, err
		}
		lines[i] = line
	}
	return lines, nil
}

func applyAttributeSubstitutionsOnElement(element interface{}, attrs types.AttributesWithOverrides) (interface{}, error) {
	var err error
	switch e := element.(type) {
	case types.Paragraph:
		e.Lines, err = applyAttributeSubstitutionsOnLines(e.Lines, attrs)
		return e, err
	// case []interface{}:
	// 	return applyAttributeSubstitutionsOnElements(e, attrs)
	case types.AttributeDeclaration:
		attrs.Set(e.Name, e.Value)
		return e, nil
	case types.AttributeReset:
		attrs.Set(e.Name, nil) // This allows us to test for a reset vs. undefined.
		return e, nil
	case types.ImageBlock:
		e.Location.Path, err = applyAttributeSubstitutionsOnElements(e.Location.Path, attrs)
		return e, err
	case types.Section:
		e.Title, err = applyAttributeSubstitutionsOnElements(e.Title, attrs)
		return e, err
	case types.OrderedListItem:
		e.Elements, err = applyAttributeSubstitutionsOnElements(e.Elements, attrs)
		return e, err
	case types.UnorderedListItem:
		e.Elements, err = applyAttributeSubstitutionsOnElements(e.Elements, attrs)
		return e, err
	case types.LabeledListItem:
		e.Elements, err = applyAttributeSubstitutionsOnElements(e.Elements, attrs)
		return e, err
	case types.ContinuedListItemElement:
		e.Element, err = applyAttributeSubstitutionsOnElement(e.Element, attrs)
		return e, err

	case types.ExampleBlock:
		e.Elements, err = applyAttributeSubstitutionsOnElements(e.Elements, attrs)
		return e, err
	case types.QuoteBlock:
		e.Elements, err = applyAttributeSubstitutionsOnElements(e.Elements, attrs)
		return e, err
	case types.SidebarBlock:
		e.Elements, err = applyAttributeSubstitutionsOnElements(e.Elements, attrs)
		return e, err
	case types.FencedBlock:
		e.Lines, err = applyAttributeSubstitutionsOnLines(e.Lines, attrs)
		return e, err
	case types.ListingBlock:
		e.Lines, err = applyAttributeSubstitutionsOnLines(e.Lines, attrs)
		return e, err
	case types.VerseBlock:
		e.Lines, err = applyAttributeSubstitutionsOnLines(e.Lines, attrs)
		return e, err
	case types.MarkdownQuoteBlock:
		e.Lines, err = applyAttributeSubstitutionsOnLines(e.Lines, attrs)
		return e, err
	case types.PassthroughBlock:
		e.Lines, err = applyAttributeSubstitutionsOnLines(e.Lines, attrs)
		return e, err
	case types.CommentBlock:
		e.Lines, err = applyAttributeSubstitutionsOnLines(e.Lines, attrs)
		return e, err
	case types.LiteralBlock:
		e.Lines, err = applyAttributeSubstitutionsOnLines(e.Lines, attrs)
		return e, err
	case types.AttributeSubstitution:
		if value, ok := attrs.GetAsString(e.Name); ok {
			return types.StringElement{
				Content: value,
			}, nil
		}
		log.Warnf("unable to find attribute '%s'", e.Name)
		return types.StringElement{
			Content: "{" + e.Name + "}",
		}, nil
	case types.CounterSubstitution:
		return applyCounterSubstitution(e, attrs)
	case types.InlineImage:
		e.Location.Path, err = applyAttributeSubstitutionsOnElements(e.Location.Path, attrs)
		return e, err
	case types.InlineLink:
		e.Location.Path, err = applyAttributeSubstitutionsOnElements(e.Location.Path, attrs)
		return e, err
	case types.ExternalCrossReference:
		e.Location.Path, err = applyAttributeSubstitutionsOnElements(e.Location.Path, attrs)
		return e, err
	case types.QuotedText:
		e.Elements, err = applyAttributeSubstitutionsOnElements(e.Elements, attrs)
		return e, err
	default:
		return e, nil
	}
}

// applyCounterSubstitutions is called by applyAttributeSubstitutionsOnElement.  Unless there is an error with
// the element (the counter is the wrong type, which should never occur), it will return a `StringElement, true`
// (because we always either find the element, or allocate one), and `nil`.  On an error it will return `nil, false`,
// and the error.  The extra boolean here is to fit the calling expectations of our caller.  This function was
// factored out of a case from applyAttributeSubstitutionsOnElement in order to reduce the complexity of that
// function, but otherwise it should have no callers.
func applyCounterSubstitution(c types.CounterSubstitution, attrs types.AttributesWithOverrides) (interface{}, error) {
	log.Debugf("applying counter substitution for '%s'", c.Name)
	counter := attrs.Counters[c.Name]
	if counter == nil {
		counter = 0
	}
	increment := true
	if c.Value != nil {
		attrs.Counters[c.Name] = c.Value
		counter = c.Value
		increment = false
	}
	switch counter := counter.(type) {
	case int:
		if increment {
			counter++
		}
		attrs.Counters[c.Name] = counter
		if c.Hidden {
			// return empty string facilitates merging
			return types.StringElement{Content: ""}, nil
		}
		return types.StringElement{
			Content: strconv.Itoa(counter),
		}, nil
	case rune:
		if increment {
			counter++
		}
		attrs.Counters[c.Name] = counter
		if c.Hidden {
			// return empty string facilitates merging
			return types.StringElement{Content: ""}, nil
		}
		return types.StringElement{
			Content: string(counter),
		}, nil

	default:
		return nil, fmt.Errorf("invalid counter type %T", counter)
	}
}
