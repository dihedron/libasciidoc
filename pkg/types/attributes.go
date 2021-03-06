package types

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// ------------------------------------------
// Attributes
// ------------------------------------------

const (
	// AttrDocType the "doctype" attribute
	AttrDocType = "doctype"
	// AttrSyntaxHighlighter the attribute to define the syntax highlighter on code source blocks
	AttrSyntaxHighlighter = "source-highlighter"
	// AttrID the key to retrieve the ID
	AttrID = "id"
	// AttrIDPrefix the key to retrieve the ID Prefix
	AttrIDPrefix = "idprefix"
	// DefaultIDPrefix the default ID Prefix
	DefaultIDPrefix = "_"
	// AttrIDSeparator the key to retrieve the ID Separator
	AttrIDSeparator = "idseparator"
	// DefaultIDSeparator the default ID Separator
	DefaultIDSeparator = "_"
	// AttrTableOfContents the `toc` attribute at document level
	AttrTableOfContents = "toc"
	// AttrTableOfContentsLevels the document attribute which specifies the number of levels to display in the ToC
	AttrTableOfContentsLevels = "toclevels"
	// AttrNoHeader attribute to disable the rendering of document footer
	AttrNoHeader = "noheader"
	// AttrNoFooter attribute to disable the rendering of document footer
	AttrNoFooter = "nofooter"
	// AttrCustomID the key to retrieve the flag that indicates if the element ID is custom or generated
	AttrCustomID = "@customID"
	// AttrTitle the key to retrieve the title
	AttrTitle = "title"
	// AttrAuthors the key to the authors declared after the section level 0 (at the beginning of the doc)
	AttrAuthors = "authors"
	// AttrRevision the key to the revision declared after the section level 0 (at the beginning of the doc)
	AttrRevision = "revision"
	// AttrRole the key to retrieve the role
	AttrRole = "role"
	// AttrInlineLink the key to retrieve the link
	AttrInlineLink = "link"
	// AttrAdmonitionKind the key to retrieve the kind of admonition , if a "masquerade" is used
	AttrAdmonitionKind = "admonitionKind"
	// AttrQuoteAuthor attribute for the author of a verse
	AttrQuoteAuthor = "quoteAuthor"
	// AttrQuoteTitle attribute for the title of a verse
	AttrQuoteTitle = "quoteTitle"
	// AttrSource the `source` attribute for a source block or a source paragraph (this is a placeholder, ie, it does not expect any value for this attribute)
	AttrSource = "source"
	// AttrLanguage the `language` attribute for a source block or a source paragraph
	AttrLanguage = "language"
	// AttrLineNums the `linenums` attribute for a source block or a source paragraph
	AttrLineNums = "linenums"
	// AttrCheckStyle the attribute to mark the first element of an unordered list itemd as a checked or not
	AttrCheckStyle = "checkstyle"
	// AttrStart the `start` attribute in an ordered list
	AttrStart = "start"
	// AttrLevelOffset the `leveloffset` attribute used in file inclusions
	AttrLevelOffset = "leveloffset"
	// AttrLineRanges the `lines` attribute used in file inclusions
	AttrLineRanges = "lines"
	// AttrTagRanges the `tag`/`tags` attribute used in file inclusions
	AttrTagRanges = "tags"
	// AttrLastUpdated the "last updated" data in the document, i.e., the output/generation time
	AttrLastUpdated = "LastUpdated"
	// AttrImageAlt the image `alt` attribute
	AttrImageAlt = "alt"
	// AttrImageHeight the image `height` attribute
	AttrImageHeight = "height"
	// AttrImageWindow the `window` attribute, which becomes the target for the link
	AttrImageWindow = "window"
	// AttrImageAlign is for image alignment
	AttrImageAlign = "align"
	// AttrIconSize the icon `size`, and can be one of 1x, 2x, 3x, 4x, 5x, lg, fw
	AttrIconSize = "size"
	// AttrIconRotate the icon `rotate` attribute, and can be one of 90, 180, or 270
	AttrIconRotate = "rotate"
	// AttrIconFlip the icon `flip` attribute, and if set can be "horizontal" or "vertical"
	AttrIconFlip = "flip"
	// AttrUnicode local libasciidoc attribute to encode output as UTF-8 instead of ASCII.
	AttrUnicode = "unicode"
	// AttrOptions element options (boolean, comma separated)
	AttrOptions = "options"
	// AttrOpts alias for AttrOptions
	AttrOpts = "opts"
	// AttrCaption is the caption for block images, tables, and so forth
	AttrCaption = "caption"
	// AttrStyle block or list style
	AttrStyle = "style"
	// AttrWidth the `width` attribute used ior images, tables, and so forth
	AttrWidth = "width"
	// AttrFrame the frame used mostly for tables (all, topbot, sides, none)
	AttrFrame = "frame"
	// AttrGrid the grid (none, all, cols, rows) in tables
	AttrGrid = "grid"
	// AttrStripes controls table row background (even, odd, all, none, hover)
	AttrStripes = "stripes"
	// AttrFloat is for image or table float (text flows around)
	AttrFloat = "float"
	// AttrCols the table columns attribute
	AttrCols = "cols"
	// AttrPositional2 positional parameter 2
	AttrPositional2 = "@2"
	// AttrPositional3 positional parameter 3
	AttrPositional3 = "@3"
	// AttrVersionLabel labels the version number in the document
	AttrVersionLabel = "version-label"
	// AttrExampleCaption is the example caption
	AttrExampleCaption = "example-caption"
	// AttrFigureCaption is the figure (image) caption
	AttrFigureCaption = "figure-caption"
	// AttrTableCaption is the table caption
	AttrTableCaption = "table-caption"
	// AttrCautionCaption is the CAUTION caption
	AttrCautionCaption = "caution-caption"
	// AttrImportantCaption is the IMPORTANT caption
	AttrImportantCaption = "important-caption"
	// AttrNoteCaption is the NOTE caption
	AttrNoteCaption = "note-caption"
	// AttrTipCaption is the TIP caption
	AttrTipCaption = "tip-caption"
	// AttrWarningCaption is the TIP caption
	AttrWarningCaption = "warning-caption"
	// AttrSubstitutions the "subs" attribute to configure substitutions on delimited blocks and paragraphs
	AttrSubstitutions = "subs"
)

// NewElementID initializes a new attribute map with a single entry for the ID using the given value
func NewElementID(id interface{}) (Attributes, error) {
	// log.Debugf("initializing a new ElementID with ID=%s", id)
	id, err := Reduce(id)
	if err != nil {
		return nil, err
	}
	return Attributes{
		AttrID:       id,
		AttrCustomID: true,
	}, nil
}

// NewElementOption sets a boolean option.
func NewElementOption(options interface{}) (Attributes, error) {
	options, err := Reduce(options)
	if err != nil {
		return nil, err
	}
	return Attributes{
		AttrOptions: options,
	}, nil
}

// NewElementNamedAttr a named (or positional) element
func NewElementNamedAttr(key string, value interface{}) (Attributes, error) {
	value, err := Reduce(value)
	if err != nil {
		return nil, err
	}
	if key == AttrOpts { // Handle the alias
		key = AttrOptions
	}
	switch key {
	case AttrRole:
		return Attributes{key: ElementRole{value}}, nil
	default:
		return Attributes{key: value}, nil
	}
}

// NewInlineElementID initializes a new attribute map with a single entry for the ID using the given value
func NewInlineElementID(id string) (Attributes, error) {
	log.Debugf("initializing a new inline ElementID with ID=%s", id)
	return Attributes{AttrID: id}, nil
}

// NewElementTitle initializes a new attribute map with a single entry for the title using the given value
func NewElementTitle(title string) (Attributes, error) {
	// log.Debugf("initializing a new ElementTitle with content=%s", title)
	return Attributes{
		AttrTitle: title,
	}, nil
}

// ElementRole the attribute value of an element role
type ElementRole []interface{}

// NewElementRole initializes a new attribute map with a single entry for the title using the given value
func NewElementRole(role interface{}) (Attributes, error) {
	// log.Debugf("initializing a new ElementRole with content=%s", role)
	role, err := Reduce(role)
	if err != nil {
		return nil, err
	}
	switch role := role.(type) {
	case []interface{}:
		return Attributes{
			AttrRole: ElementRole(role), // convert
		}, nil
	default:
		return Attributes{
			AttrRole: ElementRole{role}, // wrap
		}, nil
	}
}

// NewElementStyle initializes a new attribute map with a single entry for the style
func NewElementStyle(style interface{}) (Attributes, error) {
	style, err := Reduce(style)
	if err != nil {
		return nil, err
	}
	return Attributes{AttrStyle: style}, nil
}

// NewAdmonitionAttribute initializes a new attribute map with a single entry for the admonition kind using the given value
func NewAdmonitionAttribute(k AdmonitionKind) (Attributes, error) {
	return Attributes{AttrAdmonitionKind: k}, nil
}

// NewAttributeGroup initializes a group of attributes from the given generic attributes.
func NewAttributeGroup(attributes []interface{}) (Attributes, error) {
	log.Debugf("initializing a new AttributeGroup with %v", attributes)
	result := make(Attributes)
	for _, a := range attributes {
		// log.Debugf("processing attribute element of type %T", a)
		if a, ok := a.(Attributes); ok {
			for k, v := range a {
				switch k {
				case AttrRole:
					result[k] = ElementRole{v}
				default:
					result[k] = v
				}
			}
		} else {
			return result, errors.Errorf("unable to process element of type '%[1]T': '%[1]s'", a)
		}
	}
	// log.Debugf("initialized a new AttributeGroup: %v", result)
	return result, nil
}

// NewGenericAttribute initializes a new ElementAttribute from the given key and optional value
func NewGenericAttribute(key string, value interface{}) (Attributes, error) {
	result := make(map[string]interface{})
	k := Apply(key,
		// remove surrounding quotes
		func(s string) string {
			return strings.Trim(s, "\"")
		},
		strings.TrimSpace)
	result[k] = nil
	if value, ok := value.(string); ok {
		v := Apply(value,
			// remove surrounding quotes
			func(s string) string {
				return strings.Trim(s, "\"")
			},
			strings.TrimSpace)
		if len(v) > 0 {
			result[k] = v
		}
	}
	// log.Debugf("initialized a new Attributes: %v", result)
	return result, nil
}

// NewQuoteAttributes initializes a new map of attributes for a verse paragraph
func NewQuoteAttributes(kind string, author, title interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{}, 3)
	switch kind {
	case "verse":
		result[AttrBlockKind] = Verse
	default:
		result[AttrBlockKind] = Quote
	}
	if author, ok := author.(string); ok {
		author = strings.TrimSpace(author)
		if len(author) > 0 {
			result[AttrQuoteAuthor] = author
		}
	}
	if title, ok := title.(string); ok {
		title = strings.TrimSpace(title)
		if len(title) > 0 {
			result[AttrQuoteTitle] = title
		}
	}
	return result, nil
}

// NewLiteralBlockAttribute initializes a new attribute map with a single entry for the `literal` kind of block
func NewLiteralBlockAttribute() (Attributes, error) {
	return Attributes{AttrBlockKind: Literal}, nil
}

// NewPassthroughBlockAttribute initializes a new attribute map with a single entry for the `passthrough` kind of block
func NewPassthroughBlockAttribute() (Attributes, error) {
	return Attributes{AttrBlockKind: Passthrough}, nil
}

// NewExampleBlockAttribute initializes a new attribute map with a single entry for the `example`` kind of block
func NewExampleBlockAttribute() (Attributes, error) {
	return Attributes{AttrBlockKind: Example}, nil
}

// NewListingBlockAttribute initializes a new attribute map with a single entry for the `listing`` kind of block
func NewListingBlockAttribute() (Attributes, error) {
	return Attributes{AttrBlockKind: Listing}, nil
}

// NewSourceAttributes initializes a new attribute map with two entries, one for the kind of element ("source") and another optional one for the language of the source code
func NewSourceAttributes(language interface{}, others ...interface{}) (Attributes, error) {
	result := Attributes{
		AttrBlockKind: Source,
	}
	if language, ok := language.(string); ok {
		result[AttrLanguage] = strings.TrimSpace(language)
	}
	for _, other := range others {
		result.Add(other)
	}
	return result, nil
}

// Attributes is a map[string]interface{} with some utility methods
type Attributes map[string]interface{}

// Set sets the key/value entry in the Attributes.
func (a Attributes) Set(key string, value interface{}) Attributes {
	if a == nil {
		a = Attributes{}
	}
	a[key] = value
	return a
}

// Add adds the given attributes to the current ones
func (a Attributes) Add(attrs interface{}) Attributes {
	if attrs, ok := attrs.(Attributes); ok {
		if len(attrs) == 0 {
			return a
		}
		if a == nil {
			a = Attributes{}
		}
		for k, v := range attrs {
			a[k] = v
			if k == AttrID {
				a[AttrCustomID] = true
			}
		}
	}
	return a
}

// Has returns the true if an entry with the given key exists
func (a Attributes) Has(key string) bool {
	_, ok := a[key]
	return ok
}

// HasOption returns true if the option is set.
func (a Attributes) HasOption(key string) bool {
	if opts, ok := a[AttrOptions].(map[string]bool); ok {
		key = strings.TrimPrefix(key, "%")
		return opts[key]
	}
	return false
}

// sets the value as a singular string value if it did not exist yet,
// or move the existing value in a slice and append the new one
func (a Attributes) append(key string, value interface{}) {
	v, found := a[key]
	if !found {
		a[key] = value
		return
	}
	switch v := v.(type) {
	case []interface{}:
		switch value := value.(type) {
		case []interface{}:
			a[key] = append(v, value...)
		default:
			a[key] = append(v, value) // just append the new value into the slice
		}
	default:
		switch value := value.(type) {
		case []interface{}:
			a[key] = append([]interface{}{v}, value...)
		default:
			a[key] = []interface{}{v, value} // move existing value in a slice, along with the new one
		}
	}
}

// GetAsString gets the string value for the given key (+ `true`),
// or empty string (+ `false`) if none was found
func (a Attributes) GetAsString(key string) (string, bool) {
	if value, found := a[key]; found {
		if value == nil {
			return "", false // nil here means attribute was reset
		}
		if value, ok := value.(string); ok {
			return value, true
		} else if v, ok := a[key]; ok {
			return fmt.Sprintf("%v", v), true
		}
	}
	return "", false
}

// GetAsStringWithDefault gets the string value for the given key,
// or returns the given default value
func (a Attributes) GetAsStringWithDefault(key, defaultValue string) string {
	if value, found := a[key]; found {
		if value == nil {
			return "" // nil present means attribute was reset
		}
		if value, ok := value.(string); ok {
			return value
		} else if v, ok := a[key]; ok {
			return fmt.Sprintf("%v", v)
		}
	}
	return defaultValue
}

// GetAsBool returns the value of the key as a bool, or `false` if the key did not exist
// or if its value was not a bool
func (a Attributes) GetAsBool(key string) bool {
	if v, ok := a[key]; ok {
		if v, ok := v.(bool); ok {
			return v
		}
	}
	return false
}

// AddNonEmpty adds the given attribute if its value is non-nil and non-empty
// TODO: raise a warning if there was already a name/value
func (a Attributes) AddNonEmpty(key string, value interface{}) {
	// do not add nil or empty values
	if value == "" {
		return
	}
	a[key] = value
}

// Positionals returns all positional attributes, ie, the values for the keys in the form of `positional-<int>`
func (a Attributes) Positionals() [][]interface{} {
	result := make([][]interface{}, 0, len(a))
	i := 0
	for {
		i++
		if arg, ok := a["positional-"+strconv.Itoa(i)].([]interface{}); ok {
			result = append(result, arg)
			continue
		}
		break
	}
	return result
}

// NewAttributes retrieves the ElementID, ElementTitle and ElementInlineLink from the given slice of attributes
func NewAttributes(attributes interface{}) (Attributes, error) {
	if attributes == nil {
		return nil, nil
	}
	switch attrs := attributes.(type) {
	case []interface{}:
		// nested case, because of the grammar syntax,
		// eg: `attributes:(ElementAttribute* LiteralAttribute ElementAttribute*)`
		// which is used to ensure that a `LiteralAttribute` element is set amongst the attributes
		if len(attrs) == 0 {
			return nil, nil
		}
		result := Attributes{}
		for _, a := range attrs {
			r, err := NewAttributes(a)
			if err != nil {
				return nil, err
			}
			for k, v := range r {
				result[k] = v
			}
		}
		return result, nil
	case Attributes:
		return attrs, nil
	case map[string]interface{}:
		if len(attrs) == 0 {
			return nil, nil
		}
		result := Attributes{}
		for k, v := range attrs {
			result[k] = v
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unexpected type of attributes: '%T'", attrs)
	}
}

// NewElementAttributes retrieves the attributes for elements.
// We always pass in an array of Attributes.  Special handling for certain attributes:
//
// AttrID - these are strings, the first occurrence of this wins, and we set AttrCustomID
// AttrRole - these are strings, we append them into an array
// AttrOptions - comma separated list, we split and put into a map
func NewElementAttributes(attributes ...interface{}) (Attributes, error) {
	if len(attributes) == 0 {
		return nil, nil
	}
	// if log.IsLevelEnabled(log.DebugLevel) {
	// 	log.Debug("new element attributes:")
	// 	spew.Fdump(log.StandardLogger().Out, attributes)
	// }
	var err error
	result := Attributes{}
	for _, item := range attributes {
		if item == nil {
			continue
		}
		var attrs Attributes
		switch item := item.(type) {
		case Attributes:
			attrs = item
		case []interface{}:
			attrs, err = NewElementAttributes(item...)
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unexpected type of attributes: %T", item)
		}
		for k, v := range attrs {
			switch k {
			case AttrID:
				// The first ID set wins.
				if !result.Has(AttrID) {
					result[AttrID] = v
					result[AttrCustomID] = true
				}
			case AttrRole:
				result.append(AttrRole, v)
			case AttrOptions, AttrOpts:
				if !result.Has(AttrOptions) {
					result[AttrOptions] = map[string]bool{}
				}
				m := result[AttrOptions].(map[string]bool)
				switch v := v.(type) {
				case string:
					for _, o := range strings.Split(v, ",") {
						m[o] = true
					}
				case map[string]bool:
					for o := range v {
						m[o] = v[o]
					}
				}
			default:
				result[k] = v
			}

		}
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

func resolveAlt(path Location) string {
	_, filename := filepath.Split(path.Stringify())
	ext := filepath.Ext(filename)
	if ext != "" {
		return strings.TrimSuffix(filename, ext)
	}
	return filename
}
