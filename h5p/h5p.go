// This wraps the property lists of the HDF5 library
package h5p

/*
#cgo LDFLAGS: -lhdf5
#include <hdf5.h>
#include "init.h"
*/
import "C"
import (
	"fmt"

	"github.com/valoox/h5go/core"
)

// Represents the default property list
// This is an untyped constant, and should therefore work with most
// implementations of the property lists in all other packages.
const Default = C.H5P_DEFAULT

// The class of the property list
type Class int

// C representation of the class id
func (self Class) C() C.hid_t { return C.hid_t(self) }

// The name of the class
func (self Class) Name() string { return names[int(self)] }

// The different classes
// These are initialised with the library, and are therefore
// not constant. They should be correctly initialised via the
// 'init' method of the package
var (
	OBJECT_CREATE,
	FILE_CREATE,
	FILE_ACCESS,
	DATASET_CREATE,
	DATASET_ACCESS,
	DATASET_XFER,
	FILE_MOUNT,
	GROUP_CREATE,
	GROUP_ACCESS,
	DATATYPE_CREATE,
	DATATYPE_ACCESS,
	STRING_CREATE,
	ATTRIBUTE_CREATE,
	OBJECT_COPY,
	LINK_CREATE,
	LINK_ACCESS Class
)

// When initialising the package, sets the classes with
// the values being used by the HDF5 library
func init() {
	C.init()
	OBJECT_CREATE = Class(C.OBJECT_CREATE)
	FILE_CREATE = Class(C.FILE_CREATE)
	FILE_ACCESS = Class(C.FILE_ACCESS)
	DATASET_CREATE = Class(C.DATASET_CREATE)
	DATASET_ACCESS = Class(C.DATASET_ACCESS)
	DATASET_XFER = Class(C.DATASET_XFER)
	FILE_MOUNT = Class(C.FILE_MOUNT)
	GROUP_CREATE = Class(C.GROUP_CREATE)
	GROUP_ACCESS = Class(C.GROUP_ACCESS)
	DATATYPE_CREATE = Class(C.DATATYPE_CREATE)
	DATATYPE_ACCESS = Class(C.DATATYPE_ACCESS)
	STRING_CREATE = Class(C.STRING_CREATE)
	ATTRIBUTE_CREATE = Class(C.ATTRIBUTE_CREATE)
	OBJECT_COPY = Class(C.OBJECT_COPY)
	LINK_CREATE = Class(C.LINK_CREATE)
	LINK_ACCESS = Class(C.LINK_ACCESS)
}

// The names of the classes
var names = [...]string{
	"user", "root", "object_create", "file_create",
	"file_access", "dataset_create", "dataset_access",
	"dataset_xfer", "file_mount", "group_create", "group_access",
	"datatype_create", "datatype_access", "string_create",
	"attribute_create", "object_copy", "link_create",
	"link_access",
}

// Returns the Id, raising an error if needed
func try(id Property, ctxt string, params ...interface{}) (Property, error) {
	return id, core.Status(int(id), fmt.Sprintf(ctxt, params...))
}

// Represents the Id of a property list
type Property core.Id

// The Id implements the PList interface, but more
// specific implementations are usually preferred
func (self Property) Id() Property            { return self }
func (self Property) Copy() (Property, error) { return Copy(self) }
func (self Property) Close() error            { return Close(self) }

// This function PANICS if the class cannot be retrieved
func (self Property) Class() Class {
	cls, err := GetClass(self)
	if err != nil {
		panic(err)
	}
	return cls
}

// Creates the object with the given class
func Create(cls Class) (Property, error) {
	return try(Property(C.H5Pcreate(cls.C())), "creating property list")
}

// Copies the property list in C
func Copy(id Property) (Property, error) {
	return try(Property(C.H5Pcopy(C.hid_t(id))),
		"copying property list")
}

// Closes the property list in C
func Close(id Property) error {
	return core.Status(int(C.H5Pclose(C.hid_t(id))), "closing property list")
}

// Gets the class of the C id
func GetClass(id Property) (Class, error) {
	i := C.H5Pget_class(C.hid_t(id))
	return Class(i), core.Status(int(i), "getting property list class")
}
