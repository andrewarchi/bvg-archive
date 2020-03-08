package dom

import (
	"bytes"
	"io"
	"strings"

	escape "html"

	"golang.org/x/net/html"
)

type NodeMatcher func(*html.Node) bool

func GetElementByClass(node *html.Node, class string) *html.Node {
	return findNode(node, matchAttrContains("class", class))
}

func GetElementByID(node *html.Node, id string) *html.Node {
	return findNode(node, matchAttr("id", id))
}

func GetElementByType(node *html.Node, elemType string) *html.Node {
	return findNode(node, matchType(elemType))
}

func GetElementsByType(node *html.Node, elemType string) []*html.Node {
	return findAllNodes(node, matchType(elemType), nil)
}

func matchAttr(key, value string) NodeMatcher {
	return func(node *html.Node) bool {
		if node.Type == html.ElementNode {
			for _, attr := range node.Attr {
				if attr.Key == key && attr.Val == value {
					return true
				}
			}
		}
		return false
	}
}

func matchAttrContains(key, value string) NodeMatcher {
	return func(node *html.Node) bool {
		if node.Type == html.ElementNode {
			for _, attr := range node.Attr {
				if attr.Key == key {
					for _, v := range strings.Split(attr.Val, " ") {
						if v == value {
							return true
						}
					}
				}
			}
		}
		return false
	}
}

func matchType(elemType string) NodeMatcher {
	return func(node *html.Node) bool {
		return node.Type == html.ElementNode && node.Data == elemType
	}
}

func findNode(node *html.Node, matcher NodeMatcher) *html.Node {
	if node == nil {
		return nil
	}
	if matcher(node) {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if match := findNode(child, matcher); match != nil {
			return match
		}
	}
	return nil
}

func findAllNodes(node *html.Node, matcher NodeMatcher, nodes []*html.Node) []*html.Node {
	if node == nil {
		return nil
	}
	if matcher(node) {
		nodes = append(nodes, node)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		nodes = findAllNodes(node, matcher, nodes)
	}
	return nodes
}

func RenderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

func Unescape(s string) string {
	return escape.UnescapeString(s)
}
