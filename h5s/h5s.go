// This wraps the H5S* functions, which implement creation and
// manipulation of the dataspaces
package h5s

/*
#cgo LDFLAGS: -lhdf5
#include <hdf5.h>
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/valoox/h5go/core"
)

// An Id specifically for the dataspace objects
type Dataspace core.Id

// Closes the dataspace
func (ds Dataspace) Close() error {
	return Close(ds)
}

// Copies the dataspace
func (ds Dataspace) Copy() (Dataspace, error) {
	return Copy(ds)
}

// Selects the entire dataspace
func (ds Dataspace) Selection() Dataspace { return ds }

// Miscellaneous constants being used throughout the package
const (
	// Represents all data to be used, typically as an argument
	// in large read/write operations
	ALL Dataspace = C.H5S_ALL
)

// The selection operators
type OP int

const (
	SET OP = C.H5S_SELECT_SET // Sets the selection

	/** Hyperslab operators */
	AND  OP = C.H5S_SELECT_AND  // Intersection (AND)
	OR      = C.H5S_SELECT_OR   // Union (OR)
	XOR     = C.H5S_SELECT_XOR  // Exclusive union (XOR)
	NOTA    = C.H5S_SELECT_NOTA // A NotA B <=> {x ∈ B | x ∉ A}
	NOTB    = C.H5S_SELECT_NOTB // A NotB B <=> {x ∈ A | x ∉ B}

	/** Points operators */
	APPEND  OP = C.H5S_SELECT_APPEND  // Appends the points
	PREPEND    = C.H5S_SELECT_PREPEND // Prepends the points
)

// Represents the different classes of dataspaces
type Class int

// C representation of the class index
func (cls Class) C() C.H5S_class_t {
	return C.H5S_class_t(cls)
}

// The default classes of dataspaces
const (
	SCALAR Class = C.H5S_SCALAR // Scalar (single element)
	SIMPLE Class = C.H5S_SIMPLE // Simple (n-dims hyperrectangle)
	NULL   Class = C.H5S_NULL   // Null (no element)
)

// Returns the Id with a possible error if the id was negative
func try(id Dataspace, context string, args ...interface{}) (Dataspace, error) {
	return id, core.Status(int(id), fmt.Sprintf(context, args...))
}

// Creates a new dataspace from a given class
// Wraps H5Screate function
func Create(cls Class) (Dataspace, error) {
	return try(Dataspace(C.H5Screate(cls.C())), "creating datatype")
}

// Copy the dataspace
func Copy(id Dataspace) (Dataspace, error) {
	return try(Dataspace(C.H5Scopy(C.hid_t(id))), "copying dataspace")
}

// Disposes of the datatype. Wraps H5Sclose
func Close(id Dataspace) error {
	return core.Status(int(C.H5Sclose(C.hid_t(id))),
		"closing dataspace")
}

// Encodes the dataspace into a binary array
func Encode(id Dataspace) ([]byte, error) {
	var sze C.size_t
	if err := core.Status(int(C.H5Sencode(C.hid_t(id), nil,
		&sze)), "computing encoding size"); err != nil {
		return nil, err
	}
	bin := make([]byte, int(sze))
	return bin, core.Status(int(C.H5Sencode(C.hid_t(id),
		unsafe.Pointer(&bin[0]),
		&sze)),
		"encoding dataspace")
}

// Decodes the array of bytes into a dataspace
func Decode(bin []byte) (Dataspace, error) {
	if len(bin) > 0 {
		return try(Dataspace(C.H5Sdecode(unsafe.Pointer(&bin[0]))),
			"decoding dataspace")
	}
	return -1, fmt.Errorf("Empty array")
}

// Creates a new Null dataspace
func CreateNull() (Dataspace, error) {
	return try(Dataspace(C.H5Screate(NULL.C())),
		"creating null dataspace")
}

// Creates a new scalar dataspace
func CreateScalar() (Dataspace, error) {
	return try(Dataspace(C.H5Screate(SCALAR.C())),
		"creating scalar dataspace")
}

// Creates a simple N-dimensional array dataspace for
// the data (what HDF5 calls a 'simple' dataspace) using
// the provided dimensions
// The Shape provided will contain the dimensions of the
// ND array. Maxs is optional, and nil can be passed to
// set the limit to the current dimension (i.e. redimension
// is forbidden). Otherwise, it provides limits to the
// expansion of the dataspace in the given dimensions to
// its number.
// If maxs is < 0, the value will be replaced by 'unlimited'
// and the array can be expanded in this dimension without limit
func CreateSimple(shape []int, maxs []int) (Dataspace, error) {
	n := len(shape)
	current := make([]C.hsize_t, n)
	for i, dim := range shape {
		current[i] = C.hsize_t(dim)
	}
	var mx *C.hsize_t
	if maxs != nil {
		if len(maxs) != n {
			return -1, fmt.Errorf("Invalid max length: expecting %v, got %v instead.", n, len(maxs))
		}
		a := make([]C.hsize_t, len(maxs))
		for i, m := range maxs {
			if m < 0 {
				a[i] = 0xffffffff & C.H5S_UNLIMITED
			} else {
				a[i] = C.hsize_t(m)
			}
		}
		mx = &a[0]
	}
	return try(Dataspace(C.H5Screate_simple(C.int(n), &current[0], mx)),
		"creating simple dataspace")
}

// The C coordinates for the pointer
func ccoords(args []uint) *C.hsize_t {
	if args == nil || len(args) == 0 {
		return nil
	}
	cargs := make([]C.hsize_t, len(args))
	for i, arg := range args {
		cargs[i] = C.hsize_t(arg)
	}
	return &cargs[0]
}

// Represents a selection (either a point or a hyperslab)
type Selection interface {
	// The h5s.Dataspace of the selection
	Selection() Dataspace
}

// Represents a hyperslab selection
type Hyperslab Dataspace

// The id of the hyperslab
func (h Hyperslab) Dataspace() Dataspace { return Dataspace(h) }

// Gets the selection provided
// Wraps the H5Sselect_hyperslab function
func (h Hyperslab) Ref(selector OP, start, stride, count, block []uint) error {
	return core.Status(int(C.H5Sselect_hyperslab(C.hid_t(h),
		C.H5S_seloper_t(selector),
		ccoords(start),
		ccoords(stride),
		ccoords(count),
		ccoords(block))), "selecting hyperslab")
}

// Wraps the H5Sselect_hyperslab function with operator SET
func (h Hyperslab) Set(start, stride, count, block []uint) error {
	return h.Ref(SET, start, stride, count, block)
}

// Wraps the H5Sselect_hyperslab function with operator OR
func (h Hyperslab) Union(start, stride, count, block []uint) error {
	return h.Ref(OR, start, stride, count, block)
}

// Wraps the H5Sselect_hyperslab function with operator XOR
func (h Hyperslab) XUnion(start, stride, count, block []uint) error {
	return h.Ref(XOR, start, stride, count, block)
}

// Wraps the H5Sselect_hyperslab function with operator AND
func (h Hyperslab) Inter(start, stride, count, block []uint) error {
	return h.Ref(AND, start, stride, count, block)
}

// Wraps the H5Sselect_hyperslab function with operator NOTA
func (h Hyperslab) Neg(start, stride, count, block []uint) error {
	return h.Ref(OR, start, stride, count, block)
}

// Wraps the H5Sselect_hyperslab function with operator NOTB
func (h Hyperslab) Excl(start, stride, count, block []uint) error {
	return h.Ref(OR, start, stride, count, block)
}

// Selects a set of points
type Points Dataspace

// The Dataspace of the points
func (pt Points) Dataspace() Dataspace { return Dataspace(pt) }

// Wraps the H5Sselect_elements function
func (pt Points) Ref(op OP, coords [][]uint) error {
	if len(coords) == 0 {
		return nil
	}
	c := make([]C.hsize_t, len(coords)*len(coords[0]))
	k := 0
	for _, pt := range coords {
		for _, x := range pt {
			c[k] = C.hsize_t(x)
			k++
		}
	}
	err := core.Status(int(
		C.H5Sselect_elements(C.hid_t(pt),
			C.H5S_seloper_t(op),
			C.size_t(len(coords)),
			&c[0])), "selecting point")
	return err
}

// Sets the point to a new value
func (pt Points) Set(coords ...[]uint) error {
	return pt.Ref(SET, coords)
}

// Appends the points
func (pt Points) Append(coords ...[]uint) error {
	return pt.Ref(APPEND, coords)
}

// Prepresents the points
func (pt Points) Prepend(coords ...[]uint) error {
	return pt.Ref(PREPEND, coords)
}
