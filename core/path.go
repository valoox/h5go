package core

import "path"

// The separator for the paths
const sep = "/"

// Joins the two paths
func Join(p1, p2 Path) Path {
	return Path(path.Join(string(p1), string(p2)))
}

// Represents a path inside a file.
// A path is typically a single string, but it is always of the form:
// <atom>/<atom>/.../<atom>
// so this type provides some methods for extracting/splitting the
// path in its atomic forms
type Path string

// String representation of the path
func (p Path) String() string { return string(p) }

// Joins the two paths
func (p1 Path) Join(p2 Path) Path {
	return Join(p1, p2)
}
