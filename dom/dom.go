package dom

import (
	"bytes"
	"io"

	escape "html"

	"golang.org/x/net/html"
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

func (node *Node) FindAttr(attr, value string) *Node {
	return node.Find(matchAttrEquals(attr, value))
}

func (node *Node) FindAttrAll(attr, value string) []*Node {
	return node.FindAll(matchAttrEquals(attr, value))
}

func (node *Node) Render() string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, (*html.Node)(node))
	return buf.String()
}

func Unescape(s string) string {
	return escape.UnescapeString(s)
}
