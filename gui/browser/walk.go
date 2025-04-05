package browser

import (
	"strings"

	"golang.org/x/net/html"
)

func walkDOM(n *html.Node, depth int, lines *[]Line) {
	switch n.Type {
	case html.ElementNode:
		// *lines = append(*lines, Line{Text: n.Data, Depth: depth})
		// for _, a := range n.Attr {
		// 	if a.Key == "class" {
		// 		fmt.Println(n.Data, a.Key, a.Val)
		// 	}
		// }
	case html.TextNode:
		trimmed := strings.TrimSpace(n.Data)
		if trimmed != "" {
			*lines = append(*lines, Line{Text: trimmed, Depth: depth})
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkDOM(c, depth+1, lines)
	}
}
