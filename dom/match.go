package dom

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var spacePattern = regexp.MustCompile("[ \t\n\r\f]+")

// matchAttr matches elements with an attribute name of attr.
// CSS: [attr]
func matchAttr(attr string) NodeMatcher {
	return func(node *Node) bool {
		if node.Type == html.ElementNode {
			for _, a := range node.Attr {
				if a.Key == attr {
					return true
				}
			}
		}
		return false
	}
}

// matchAttrExact matches elements with an attribute name of attr whose
// value is exactly value.
// CSS: [attr=value]
func matchAttrEquals(attr, value string) NodeMatcher {
	return func(node *Node) bool {
		if node.Type == html.ElementNode {
			for _, a := range node.Attr {
				if a.Key == attr && a.Val == value {
					return true
				}
			}
		}
		return false
	}
}

// matchAttrWord matches elements with an attribute name of attr
// whose value is a whitespace-separated list of words, one of which is
// exactly value.
// CSS: [attr~=value]
func matchAttrWord(attr, value string) NodeMatcher {
	return func(node *Node) bool {
		if node.Type == html.ElementNode {
			for _, a := range node.Attr {
				if a.Key == attr {
					for _, v := range spacePattern.Split(a.Val, -1) {
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

// matchAttrHyphen matches elements with an attribute name of attr whose
// value can be exactly value or can begin with value immediately
// followed by a hyphen.
// CSS: [attr|=value]
func matchAttrHyphen(attr, value string) NodeMatcher {
	return func(node *Node) bool {
		if node.Type == html.ElementNode {
			for _, a := range node.Attr {
				if a.Key == attr && strings.HasPrefix(a.Val, value) &&
					(len(a.Val) == len(value) || a.Val[len(value)] == '-') {
					return true
				}
			}
		}
		return false
	}
}

// matchAttrPrefix matches elements with an attribute name of attr whose
// value is prefixed by prefix.
// CSS: [attr^=prefix]
func matchAttrPrefix(attr, prefix string) NodeMatcher {
	return func(node *Node) bool {
		if node.Type == html.ElementNode {
			for _, a := range node.Attr {
				if a.Key == attr && strings.HasPrefix(a.Val, prefix) {
					return true
				}
			}
		}
		return false
	}
}

// matchAttrSuffix matches elements with an attribute name of attr whose
// value is suffixed by suffix.
// CSS: [attr$=value]
func matchAttrSuffix(attr, suffix string) NodeMatcher {
	return func(node *Node) bool {
		if node.Type == html.ElementNode {
			for _, a := range node.Attr {
				if a.Key == attr && strings.HasSuffix(a.Val, suffix) {
					return true
				}
			}
		}
		return false
	}
}

// matchAttrContains matches elements with an attribute name of attr
// whose value contains at least one occurrence of value within the
// string.
// [attr*=value]
func matchAttrContains(attr, value string) NodeMatcher {
	return func(node *Node) bool {
		if node.Type == html.ElementNode {
			for _, a := range node.Attr {
				if a.Key == attr && strings.Contains(a.Val, value) {
					return true
				}
			}
		}
		return false
	}
}

func matchTag(tag string) NodeMatcher {
	return func(node *Node) bool {
		return node.Type == html.ElementNode && node.Data == tag
	}
}

func matchTagAtom(tag atom.Atom) NodeMatcher {
	return func(node *Node) bool {
		return node.Type == html.ElementNode && node.DataAtom == tag
	}
}
