package h5go

import (
	"github.com/valoox/h5go/core"
	"github.com/valoox/h5go/h5d"
	"github.com/valoox/h5go/h5f"
	"github.com/valoox/h5go/h5g"
	"github.com/valoox/h5go/h5l"
)

var (
	FileAccess = h5f.DefaultAccess
	FileCreate = h5f.DefaultCreate
)

// Implements a generic location, i.e. features which are common
// between groups and files (such as the ability to add subgroups
// or files in them)
type loc struct {
	where core.Location // Embeds a location
	at    core.Path     // The path to this location in the file
	in    *File         // The file this belongs to
	// Default options
	lcreate h5l.Crt // Options for link creation
	laccess h5l.Acc // Options for link access
	gcreate h5g.Crt // Options for group creation
	gaccess h5g.Acc // Options for group access
	dcreate h5d.Crt // Options for dataset creation
	daccess h5d.Acc // Options for dataset access
}

// Initialises all options to the defaults
func (l *loc) defaults() (err error) {
	if l.lcreate, err = h5l.DefaultCreate.Copy(); err != nil {
		return err
	}
	if l.laccess, err = h5l.DefaultAccess.Copy(); err != nil {
		return err
	}
	if l.gcreate, err = h5g.DefaultCreate.Copy(); err != nil {
		return err
	}
	if l.gaccess, err = h5g.DefaultAccess.Copy(); err != nil {
		return err
	}
	if l.dcreate, err = h5d.DefaultCreate.Copy(); err != nil {
		return err
	}
	l.daccess, err = h5d.DefaultAccess.Copy()
	return err
}

// Copies the location parameters into a new location options
func (l *loc) copyTo(where core.Location, newpath core.Path) *loc {
	return &loc{
		where:   where,
		at:      newpath,
		in:      l.in,
		lcreate: l.lcreate,
		laccess: l.laccess,
		gcreate: l.gcreate,
		gaccess: l.gaccess,
		dcreate: l.dcreate,
		daccess: l.daccess,
	}
}

// Creates a new group at this location
func (l *loc) NewGroup(path core.Path) (Group, error) {
	gid, err := h5g.Create(l.where, path, l.lcreate,
		l.gcreate, l.gaccess)
	return Group{
		// Copies the location to the newly created group
		loc:   l.copyTo(gid, core.Join(l.at, path)),
		Group: gid,
	}, err
}

// Gets the group at the given path from this location
func (l *loc) Get(path core.Path) (Group, error) {
	gid, err := h5g.Open(l.where, path, l.gaccess)
	return Group{
		// Copies the location to the group
		loc:   l.copyTo(gid, core.Join(l.at, path)),
		Group: gid,
	}, err
}

// Wraps a h5g.Group handle and adds methods and features
type Group struct {
	*loc      // Embeds the location
	h5g.Group // The embedded Group handle
}

// Wraps an h5f.File object and adds convenience accesses
type File struct {
	*loc            // Embeds the location
	h5f.File        // The embedded File handle
	path     string // The path to the file
}

// Opens a file, stating whether it is read-only (rw = false) or
// if it can be edited (rw = true)
func Open(path string, rw bool) (*File, error) {
	var flag h5f.Flag
	if rw {
		flag = h5f.RW
	} else {
		flag = h5f.RO
	}
	fid, err := h5f.Open(path, flag, FileAccess)
	out := &File{
		loc:  new(loc),
		File: fid,
		path: path,
	}
	out.loc.at = ""
	out.loc.where = out
	out.loc.in = out
	if err != nil {
		return out, err
	}
	return out, out.defaults()
}

// Creates a new file. If the overwrite (`ow`) boolean is set,
// overwrite any existing file.
// Otherwise, this raises an error if the file already exists
func Create(path string, ow bool) (*File, error) {
	var flag h5f.Flag
	if ow {
		flag = h5f.TRUNC
	} else {
		flag = h5f.CREATE
	}
	fid, err := h5f.Create(path, flag, FileCreate, FileAccess)
	out := &File{
		loc:  new(loc),
		File: fid,
		path: path,
	}
	out.loc.at = ""
	out.loc.where = out
	out.loc.in = out
	if err != nil {
		return out, err
	}
	return out, out.defaults()
}
