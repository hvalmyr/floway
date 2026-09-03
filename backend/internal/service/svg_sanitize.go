package service

import (
	"encoding/xml"
	"errors"
	"strings"
)

// svgAllowedElements/svgAllowedAttrs are a deliberately small allowlist —
// enough for the icon shapes this site actually needs (paths, basic
// shapes, groups, simple gradients), not "every SVG feature". Anything not
// on the list is dropped rather than rejected outright, so one unsupported
// decorative element in an otherwise-fine icon doesn't block the whole
// upload.
//
// This is the only line of defense against a malicious upload: icons are
// stored as raw markup and rendered with v-html (see AppIcon.vue), which
// executes anything left in — <script>, event-handler attributes
// (onload/onclick/...), and external references (href/xlink:href to
// anything other than a local "#fragment") are exactly what this strips.
var svgAllowedElements = map[string]bool{
	"svg": true, "g": true, "path": true, "circle": true, "rect": true,
	"line": true, "polyline": true, "polygon": true, "ellipse": true,
	"defs": true, "use": true, "title": true, "desc": true,
	"linearGradient": true, "radialGradient": true, "stop": true,
	"clipPath": true, "mask": true,
}

var svgAllowedAttrs = map[string]bool{
	"viewBox": true, "width": true, "height": true, "fill": true,
	"stroke": true, "stroke-width": true, "stroke-linecap": true,
	"stroke-linejoin": true, "stroke-dasharray": true, "fill-rule": true,
	"clip-rule": true, "opacity": true, "fill-opacity": true,
	"stroke-opacity": true, "transform": true, "d": true, "cx": true,
	"cy": true, "r": true, "rx": true, "ry": true, "x": true, "y": true,
	"x1": true, "y1": true, "x2": true, "y2": true, "points": true,
	"offset": true, "stop-color": true, "stop-opacity": true,
	"gradientUnits": true, "gradientTransform": true, "id": true,
	"clip-path": true, "mask": true, "xmlns": true,
}

// sanitizeSVG parses the input as XML and re-serializes only the elements/
// attributes on the allowlist above, dropping everything else (including
// the element itself for anything not in svgAllowedElements — a <script>
// tag's closing tag is also just dropped, not "closed early" some unsafe
// way, since the whole subtree is skipped). `href`/`xlink:href` are kept
// only when they point at a local "#fragment" (needed for <use> to
// reference a <defs> shape) — anything else (an http(s) URL, a
// "javascript:" URI) is stripped, since it'd otherwise let an uploaded icon
// fetch or execute something outside the markup itself.
//
// Returns an error if the input isn't well-formed XML or has no root <svg>
// element at all — those aren't things a repair pass should paper over.
func sanitizeSVG(input string) (string, error) {
	decoder := xml.NewDecoder(strings.NewReader(input))
	// SVGs commonly declare entities in a DOCTYPE the standard decoder
	// doesn't resolve by default; icons don't need custom entities, so
	// this only has to not error out on a harmless declaration.
	decoder.Strict = false
	decoder.Entity = xml.HTMLEntity

	var out strings.Builder
	var skipDepth int // >0 while inside a disallowed element's subtree
	sawRoot := false

	for {
		tok, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return "", err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			name := t.Name.Local
			if skipDepth > 0 || !svgAllowedElements[name] {
				skipDepth++
				continue
			}
			if name == "svg" {
				sawRoot = true
			}
			out.WriteByte('<')
			out.WriteString(name)
			for _, attr := range t.Attr {
				attrName := attr.Name.Local
				// t.Name.Local strips the namespace prefix, so this also
				// catches the legacy "xlink:href" form.
				if attrName == "href" {
					if !strings.HasPrefix(attr.Value, "#") {
						continue
					}
					out.WriteString(` href="`)
					out.WriteString(escapeAttr(attr.Value))
					out.WriteByte('"')
					continue
				}
				if !svgAllowedAttrs[attrName] {
					continue
				}
				out.WriteByte(' ')
				out.WriteString(attrName)
				out.WriteString(`="`)
				out.WriteString(escapeAttr(attr.Value))
				out.WriteByte('"')
			}
			out.WriteByte('>')

		case xml.EndElement:
			if skipDepth > 0 {
				skipDepth--
				continue
			}
			if !svgAllowedElements[t.Name.Local] {
				continue
			}
			out.WriteString("</")
			out.WriteString(t.Name.Local)
			out.WriteByte('>')

		case xml.CharData:
			if skipDepth == 0 {
				out.WriteString(escapeText(string(t)))
			}

		// Comments, processing instructions, and directives (DOCTYPE) are
		// silently dropped — none are needed to render an icon, and a
		// comment is a common place to smuggle disallowed content past a
		// naive sanitizer.
		default:
			continue
		}
	}

	if !sawRoot {
		return "", errors.New("no <svg> root element found")
	}
	return out.String(), nil
}

func escapeAttr(s string) string {
	var b strings.Builder
	// strings.Builder.Write never errors, so xml.EscapeText can't either
	// (it only ever fails if the underlying writer does).
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}

func escapeText(s string) string {
	var b strings.Builder
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}
