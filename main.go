// Command htmangl combines two HTML files.
//
// A way to use this is to make a "template" containing the basic outline of a
// page and then apply HTML files containing the content to it.
//
//	$ htmangl template.html home_partial.html > home.html
//
// Use the special comment "<!-- htmangl:insert -->" to replace with applied
// child nodes instead of appending to the parent.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	flag.Usage = func() { fmt.Fprint(os.Stderr, "usage: htmangl BASE APPLY\n") }
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	baseFile, err := os.Open(flag.Arg(0))
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	defer baseFile.Close()

	baseDoc, err := html.Parse(baseFile)
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}

	applyFile, err := os.Open(flag.Arg(1))
	if err != nil {
		return fmt.Errorf("read applied: %w", err)
	}
	defer applyFile.Close()

	applyDoc, err := html.Parse(applyFile)
	if err != nil {
		return fmt.Errorf("parse applied: %w", err)
	}

	newDoc := apply(baseDoc, applyDoc)
	return html.Render(os.Stdout, newDoc)
}

// apply may change baseDoc and/or applyDoc while producing the result.
func apply(baseDoc, applyDoc *html.Node) *html.Node {
	toApply := map[string]*html.Node{}

	for node := range applyDoc.ChildNodes() {
		toApply[node.Data] = node
	}

	if len(toApply) == 0 {
		return cloneTree(baseDoc)
	}

	if baseDoc.FirstChild == nil {
		return cloneTree(applyDoc)
	}

	newDoc := cloneNode(baseDoc)

	var preDirective, postDirective []*html.Node
	var seenDirective, seenCopyDirective bool

	for node := range baseDoc.ChildNodes() {
		if node.Type == html.CommentNode && strings.TrimSpace(node.Data) == "htmangl:insert" {
			seenDirective = true
			continue
		}

		if node.Type == html.CommentNode && strings.TrimSpace(node.Data) == "htmangl:copy" {
			seenCopyDirective = true
			break
		}

		if applyNode, ok := toApply[node.Data]; node.Type == html.ElementNode && ok {
			delete(toApply, node.Data)
			applied := apply(node, applyNode)
			if seenDirective {
				postDirective = append(postDirective, applied)
			} else {
				preDirective = append(preDirective, applied)
			}
		} else {
			applied := cloneTree(node)
			if seenDirective {
				postDirective = append(postDirective, applied)
			} else {
				preDirective = append(preDirective, applied)
			}
		}
	}

	if seenCopyDirective {
		for node := range baseDoc.ChildNodes() {
			if node.Type == html.CommentNode && strings.TrimSpace(node.Data) == "htmangl:copy" {
				continue
			}
			newDoc.AppendChild(cloneNode(node))
		}

		for node := range applyDoc.ChildNodes() {
			newDoc.AppendChild(cloneTree(node))
		}

		return newDoc
	} else {
		for _, node := range preDirective {
			newDoc.AppendChild(node)
		}

		for _, node := range toApply {
			newDoc.AppendChild(cloneTree(node))
		}

		for _, node := range postDirective {
			newDoc.AppendChild(node)
		}
	}

	return newDoc
}

func cloneNode(node *html.Node) *html.Node {
	m := &html.Node{
		Type:     node.Type,
		DataAtom: node.DataAtom,
		Data:     node.Data,
		Attr:     make([]html.Attribute, len(node.Attr)),
	}
	copy(m.Attr, node.Attr)
	return m
}

func cloneTree(node *html.Node) *html.Node {
	m := cloneNode(node)

	for {
		child := node.FirstChild
		if child == nil {
			break
		}
		node.RemoveChild(child)
		m.AppendChild(child)
	}

	return m
}

func renderNode(node *html.Node) string {
	var buf bytes.Buffer
	html.Render(&buf, node)
	return buf.String()
}
