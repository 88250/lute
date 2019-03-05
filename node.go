// Lute - A structural markdown engine.
// Copyright (C) 2019, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package lute

// Node represents a node in ast. https://github.com/syntax-tree/mdast
type Node struct {
	NodeType
	Position Pos
	tr       *Tree
}

// NodeType identifies the type of a parse tree node.
type NodeType int

const (
	NodeParent NodeType = iota
	NodeLiteral
	NodeHTML
	NodeCode
	NodeText
	NodeInlineCode
	NodeEmphasis
	NodeStrong
	NodeImage
)

// Nodes.

// ParentNode holds a sequence of nodes.
type ParentNode struct {
	Node
	Children []Node // element nodes in lexical order
}

type LiteralNode struct {
	Node
	Value string
}

type HTMLNode struct {
	LiteralNode
}

type CodeNode struct {
	LiteralNode
	Lang string
	Meta string
}

type TextNode struct {
	LiteralNode
}

type InlineCode struct {
	LiteralNode
}

type Emphasis struct {
	ParentNode
	Children []Node
}

type Strong struct {
	ParentNode
	Children []Node
}

type Delete struct {
	ParentNode
	Children []Node
}

type Break struct {
	Node
}

type Image struct {
	Node
	Resource
	Alternative
}

type ImageReference struct {
	
}

// Mixins.

type Resource struct {
	URL   string
	Title string
}

type Association struct {
	Identifier string
	Label      string
}

type Reference struct {
	ReferenceType referenceType
	Association
}

type Alternative struct {
	Alt string
}

// Enumerations.

type alignType string     // "left" | "right" | "center" | null
type referenceType string // "shortcut" | "collapsed" | "full"
