// This wraps the H5G set of functions in the library
package h5g

/*
#cgo LDFLAGS: -lhdf5
#include <hdf5.h>
*/
import "C"

import (
	"github.com/valoox/h5go/core"
	"github.com/valoox/h5go/h5l"
	"github.com/valoox/h5go/h5p"
)

// Default group creation property list
const (
	DefaultCreate = Crt(h5p.Default)
	DefaultAccess = Acc(h5p.Default)
)

// Creates a new group creation property list
func Creation() (Crt, error) {
	id, err := h5p.Create(h5p.GROUP_CREATE)
	return Crt(id), err
}

// Creates a new group access property list
func Access() (Acc, error) {
	id, err := h5p.Create(h5p.GROUP_ACCESS)
	return Acc(id), err
}

// Id targeting a group creation property list
type Crt h5p.Property

// The Property list ID
func (self Crt) Id() h5p.Property { return h5p.Property(self) }

// Copies the property list
func (self Crt) Copy() (Crt, error) {
	id, err := h5p.Copy(self.Id())
	return Crt(id), err
}

// The class of the property list
func (self Crt) Class() h5p.Class { return h5p.GROUP_CREATE }

// Disposes of the resource
func (self Crt) Close() error { return h5p.Close(self.Id()) }

// Class for group access
type Acc h5p.Property

// The property list id
func (self Acc) Id() h5p.Property { return h5p.Property(self) }

// Copies the property list
func (self Acc) Copy() (Acc, error) {
	id, err := h5p.Copy(self.Id())
	return Acc(id), err
}

// The class of the property list
func (self Acc) Class() h5p.Class { return h5p.GROUP_ACCESS }

// Disposes of the resource
func (self Acc) Close() error { return h5p.Close(self.Id()) }

// Represents a Group id. This is represented similarly as all the
// other ids in the library, but uses this specific type to
// distinguish it from other types of Id.
type Group core.Id

// Implements the location interface, as this is a group
func (g Group) At() core.Id {
	return core.Id(g)
}

// The global id of the group
func (g Group) Id() core.Id { return core.Id(g) }

// Closes the group
func (g Group) Close() error {
	return Close(g)
}

// Returns a Group, raising an error if something is wrong
func try(id Group, ctxt string, args ...interface{}) (Group, error) {
	return id, core.Status(int(id), ctxt, args...)
}

// Wraps the H5Gcreate function
func Create(root core.Location, path core.Path, links h5l.Crt,
	create Crt, access Acc) (Group, error) {
	return try(Group(C.H5Gcreate2(C.hid_t(root.At()),
		C.CString(path.String()), C.hid_t(links),
		C.hid_t(create), C.hid_t(access))),
		"creating group at %s", path)
}

// Wraps the H5Gclose function
func Close(group Group) error {
	return core.Status(int(C.H5Gclose(C.hid_t(group))),
		"closing group")
}

// Wraps the H5Gopen function
func Open(at core.Location, pth core.Path, acc Acc) (Group, error) {
	return try(Group(C.H5Gopen2(C.hid_t(at.At()),
		C.CString(pth.String()), C.hid_t(acc))),
		"opening group at %s", pth)
}
