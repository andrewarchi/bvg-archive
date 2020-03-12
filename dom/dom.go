package dom

import (
	"bytes"
	"io"
	"strings"

	escape "html"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Node html.Node

type NodeMatcher func(*Node) bool

func (node *Node) Find(matcher NodeMatcher) *Node {
	if node == nil {
		return nil
	}
	if matcher(node) {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if match := (*Node)(child).Find(matcher); match != nil {
			return match
		}
	}
	return nil
}

func (node *Node) FindAll(matcher NodeMatcher) []*Node {
	return node.findAll(matcher, nil)
}

func (node *Node) findAll(matcher NodeMatcher, nodes []*Node) []*Node {
	if node == nil {
		return nil
	}
	if matcher(node) {
		nodes = append(nodes, node)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		nodes = (*Node)(child).findAll(matcher, nodes)
	}
	return nodes
}

func (node *Node) FindID(id string) *Node {
	return node.Find(matchAttrEquals("id", id))
}

func (node *Node) FindIDAll(id string) []*Node {
	return node.FindAll(matchAttrEquals("id", id))
}

func (node *Node) FindClass(class string) *Node {
	return node.Find(matchAttrWord("class", class))
}

func (node *Node) FindClassAll(class string) []*Node {
	return node.FindAll(matchAttrWord("class", class))
}

func (node *Node) FindTag(tag string) *Node {
	return node.Find(matchTag(tag))
}

func (node *Node) FindTagAll(tag string) []*Node {
	return node.FindAll(matchTag(tag))
}

func (node *Node) FindTagAtom(tag atom.Atom) *Node {
	return node.Find(matchTagAtom(tag))
}

func (node *Node) FindTagAtomAll(tag atom.Atom) []*Node {
	return node.FindAll(matchTagAtom(tag))
}

func (node *Node) FindAttr(attr, value string) *Node {
	return node.Find(matchAttrEquals(attr, value))
}

func (node *Node) FindAttrAll(attr, value string) []*Node {
	return node.FindAll(matchAttrEquals(attr, value))
}

func (node *Node) LookupAttr(attr string) (string, bool) {
	if node != nil {
		for _, a := range node.Attr {
			if a.Key == attr {
				return a.Val, true
			}
		}
	}
	return "", false
}

func (node *Node) Render() string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, (*html.Node)(node))
	return buf.String()
}

// TextContent returns the text content of the node and its descendants.
func (node *Node) TextContent() string {
	if node == nil {
		return ""
	}
	switch node.Type {
	case html.TextNode, html.CommentNode:
		return node.Data
	case html.DoctypeNode:
	case html.ElementNode, html.DocumentNode:
		var b strings.Builder
		writeTextContent(&b, (*html.Node)(node))
		return b.String()
	}
	return ""
}

func writeTextContent(b *strings.Builder, node *html.Node) {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		switch child.Type {
		case html.TextNode:
			b.WriteString(child.Data)
		case html.ElementNode:
			writeTextContent(b, child)
		case html.DocumentNode, html.DoctypeNode, html.CommentNode:
		}
	}
}

func Unescape(s string) string {
	return escape.UnescapeString(s)
}
