package main

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
	"hawx.me/code/assert"
)

func TestApply(t *testing.T) {
	testcases := map[string]struct {
		base, apply, expected string
	}{
		"apply empty": {
			base:     `<html><head></head><body><h1>Hey</h1></body></html>`,
			apply:    ``,
			expected: `<html><head></head><body><h1>Hey</h1></body></html>`,
		},
		"apply empty with doctype": {
			base:     `<!DOCTYPE html><html><head></head><body><h1>Hey</h1></body></html>`,
			apply:    ``,
			expected: `<!DOCTYPE html><html><head></head><body><h1>Hey</h1></body></html>`,
		},
		"applied to empty": {
			base:     ``,
			apply:    `<html><head></head><body><h1>Bye</h1></body></html>`,
			expected: `<html><head></head><body><h1>Bye</h1></body></html>`,
		},
		"applied to something": {
			base:     `<html><head></head><body><h1>Hello</h1></body></html>`,
			apply:    `<html><head></head><body><h1>Bye</h1></body></html>`,
			expected: `<html><head></head><body><h1>HelloBye</h1></body></html>`,
		},
		"new element": {
			base:     `<html><head></head><body><h1>Hello</h1></body></html>`,
			apply:    `<html><head></head><body><p>Bye</p></body></html>`,
			expected: `<html><head></head><body><h1>Hello</h1><p>Bye</p></body></html>`,
		},
		"inserted element": {
			base:     `<html><head></head><body><header>HEADER</header><!-- htmangl:insert --><footer>FOOTER</footer></body></html>`,
			apply:    `<html><head></head><body><p>CONTENT</p></body></html>`,
			expected: `<html><head></head><body><header>HEADER</header><p>CONTENT</p><footer>FOOTER</footer></body></html>`,
		},
		"copy elements": {
			base:     `<html><head><link rel="a" href="b"/><!-- htmangl:copy --></head><body><h1>Hello </h1></body></html>`,
			apply:    `<html><head><link rel="c" href="d"/></head><body><h1>Bye</h1></body></html>`,
			expected: `<html><head><link rel="a" href="b"/><link rel="c" href="d"/></head><body><h1>Hello Bye</h1></body></html>`,
		},
		"example": {
			base:     `<html lang="en"><head><meta charset="utf-8" /><title>My website</title><link rel="stylesheet" href="css/screen.css" type="text/css" /></head><body><header><h1>My website</h1></header><!-- htmangl:insert --><footer>Copyright me (this year)</footer></body></html>`,
			apply:    `<html><head><title> - Home</title></head><body><img src="hero"/><p>This is my website, welcome.</p></body></html>`,
			expected: `<html lang="en"><head><meta charset="utf-8"/><title>My website - Home</title><link rel="stylesheet" href="css/screen.css" type="text/css"/></head><body><header><h1>My website</h1></header><img src="hero"/><p>This is my website, welcome.</p><footer>Copyright me (this year)</footer></body></html>`,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			baseDoc, _ := html.Parse(strings.NewReader(tc.base))
			applyDoc, _ := html.Parse(strings.NewReader(tc.apply))

			assert.Equal(t, tc.expected, renderNode(apply(baseDoc, applyDoc)))
		})
	}
}

func TestOrderedMap(t *testing.T) {
	m := newOrderedMap[string, string]()

	m.Set("a", "1")
	m.Set("b", "2")

	v, ok := m.Get("a")
	assert.True(t, ok)
	assert.Equal(t, "1", v)
	v, ok = m.Get("b")
	assert.True(t, ok)
	assert.Equal(t, "2", v)
	assert.Equal(t, 2, m.Len())

	m.Set("c", "3")
	m.Set("a", "4")
	v, ok = m.Get("a")
	assert.True(t, ok)
	assert.Equal(t, "4", v)

	var got [][2]string
	for k, v := range m.Iter() {
		got = append(got, [2]string{k, v})
	}

	assert.Equal(t, [][2]string{
		{"b", "2"},
		{"c", "3"},
		{"a", "4"},
	}, got)

	assert.Equal(t, 3, m.Len())
}
