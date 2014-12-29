// Wraps the  H5L* functions of the library, for link manipulation
// and creation
package h5l

/*
#cgo LDFLAGS: -lhdf5
#include <hdf5.h>
*/
import "C"

import (
	"fmt"

	"github.com/valoox/h5go/core"
	"github.com/valoox/h5go/h5p"
)

const (
	// Default link creation property list
	DefaultCreate = Crt(h5p.Default)
	// Default link access property list
	DefaultAccess = Acc(h5p.Default)
)

// Creates a new list for link creation
func Creation() (Crt, error) {
	id, err := h5p.Create(h5p.LINK_CREATE)
	return Crt(id), err
}

// Link creation property list
type Crt h5p.Property

func (self Crt) Id() h5p.Property { return h5p.Property(self) }
func (self Crt) Close() error     { return h5p.Close(self.Id()) }
func (self Crt) Class() h5p.Class { return h5p.LINK_CREATE }
func (self Crt) Copy() (Crt, error) {
	id, err := h5p.Copy(self.Id())
	return Crt(id), err
}

// Creates a new property list for link access
func Access() (Acc, error) {
	id, err := h5p.Create(h5p.LINK_ACCESS)
	return Acc(id), err
}

type Acc h5p.Property

func (self Acc) Id() h5p.Property { return h5p.Property(self) }
func (self Acc) Close() error     { return h5p.Close(self.Id()) }
func (self Acc) Class() h5p.Class { return h5p.LINK_ACCESS }
func (self Acc) Copy() (Acc, error) {
	id, err := h5p.Copy(self.Id())
	return Acc(id), err
}

// Represents the Id of a link
type Link core.Id

// Tries to return the Link, returning an error  if it was <0
func try(id Link, msg string, args ...interface{}) (Link, error) {
	return id, core.Status(int(id), fmt.Sprintf(msg, args...))
}

// Wraps the H5Lcreate_hard function
func Hard(src core.Location, name core.Path,
	link core.Location, lname core.Path,
	create Crt, access Acc) (Link, error) {
	return try(Link(C.H5Lcreate_hard(C.hid_t(src.At()),
		C.CString(name.String()), C.hid_t(link.At()),
		C.CString(lname.String()),
		C.hid_t(create), C.hid_t(access))),
		"creating hard link from %s to %s", name, lname)
}

// Wraps the H5Lcreate_soft function
func Soft(target core.Path, loc core.Location, link core.Path,
	create Crt, access Acc) (Link, error) {
	return try(Link(C.H5Lcreate_soft(C.CString(target.String()),
		C.hid_t(loc.At()), C.CString(link.String()),
		C.hid_t(create), C.hid_t(access))),
		"creating hard link from %s to %s", target, link)
}

// Wraps the H5Lcopy function
func Copy(src core.Location, rel core.Path,
	dest core.Location, drel core.Path,
	create Crt, access Acc) (Link, error) {
	return try(Link(C.H5Lcopy(C.hid_t(src.At()),
		C.CString(rel.String()), C.hid_t(dest.At()),
		C.CString(drel.String()),
		C.hid_t(create), C.hid_t(access))),
		"copying link from %s to %s", rel, drel)
}

// Wraps the H5Lmove function
func Move(src core.Location, rel core.Path,
	dest core.Location, drel core.Path,
	create Crt, access Acc) (Link, error) {
	return try(Link(C.H5Lmove(C.hid_t(src.At()),
		C.CString(rel.String()), C.hid_t(dest.At()),
		C.CString(drel.String()),
		C.hid_t(create), C.hid_t(access))),
		"moving link from %s to %s", rel, drel)
}

// Wraps the H5Ldelete function
func Delete(src core.Location, name core.Path, prop Acc) error {
	return core.Status(int(C.H5Ldelete(C.hid_t(src.At()),
		C.CString(name.String()), C.hid_t(prop))),
		"deleting link %s", name)
}
