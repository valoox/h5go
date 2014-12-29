// This package wraps the H5I* family of functions, for creating and
// manipulating identifiers in the file
package h5i

/*
#cgo LDFLAGS: -lhdf5
#include <hdf5.h>
*/
import "C"
import (
	"github.com/valoox/h5go/core"
)

// Returns the path of an object in a file
// Wraps the H5Iget_name function
func GetName(id core.Id) (core.Path, error) {
	var out *C.char
	sze := C.H5Iget_name(C.hid_t(id), out, 1)
	if err := core.Status(int(sze), "getting name"); err != nil {
		return "", err
	} else if err := core.Status(int(C.H5Iget_name(C.hid_t(id),
		out, C.size_t(sze+1))), "getting name"); err != nil {
		return "", err
	}
	return core.Path(C.GoString(out)), nil
}
