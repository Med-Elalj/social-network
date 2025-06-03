package upload

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

func validateSecureSVG(r io.Reader) error {
	decoder := xml.NewDecoder(r)

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("XML parse error: %w", err)
		}

		switch t := token.(type) {
		case xml.Directive:
			// Reject any directive, e.g., <!DOCTYPE ...>
			return errors.New("invalid SVG: XML directives (e.g. DOCTYPE) are not allowed")

		case xml.ProcInst:
			// Optionally reject processing instructions like <?xml-stylesheet ...?>
			return errors.New("invalid SVG: XML processing instructions are not allowed")

		case xml.StartElement:
			// Reject <script> elements
			if _, is := disallowedTags[strings.ToLower(t.Name.Local)]; is {
				return fmt.Errorf("invalid SVG: <%s> elements are not allowed", t.Name.Local)
			}
			// Reject any attributes that start with "on" (e.g. onclick, onload)
			for _, attr := range t.Attr {
				if _, is := dangerousAttr[strings.ToLower(attr.Name.Local)]; is || strings.HasPrefix(strings.ToLower(attr.Name.Local), "on") {
					return fmt.Errorf("invalid SVG: disallowed attribute: %s", attr.Name.Local)
				}
			}
		}
	}

	return nil // Safe SVG
}

var disallowedTags = map[string]struct{}{
	"script":        {},
	"foreignobject": {},
	"style":         {},
	"iframe":        {},
	"object":        {},
	"embed":         {},
}

var dangerousAttr = map[string]struct{}{
	"style":  {},
	"srcdoc": {},
	"xmlns":  {},
}

// Content-Security-Policy: default-src 'none'; img-src 'self'; style-src 'none'; script-src 'none'
