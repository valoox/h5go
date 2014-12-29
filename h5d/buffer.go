package h5d

import "unsafe"

import (
	"github.com/valoox/h5go/h5s"
	"github.com/valoox/h5go/h5t"
)

// This general interface is required for both reading and writing
// It allows the buffer to define its own datatype
type Typed interface {
	// Produces (or gets) the type of the data
	Type() (h5t.Datatype, error)
}

// This general interface states that the object has a shape
type Shaped interface {
	// Produces the shape of the object
	// The simplest is to return h5s.ALL
	Shape() (h5s.Dataspace, error)
}

// Simple structure stating a flat shape. As long as the read/write
// pointers are correct, this should be enough for most cases
type Flat struct{}

// Returns a flat shape for the interface
// As long as the read/write pointers are all right, this should
// be enough in most cases
func (f Flat) Shape() (h5s.Dataspace, error) { return h5s.ALL, nil }

// Corresponds to a view on an HDF5 object, both typed and shaped
type view interface {
	Typed
	Shaped
}

// The interface for Go objects used as input to HDF5.
// If the structure implements this interface,
// it can be used to write data to a HDF5 file
type Input interface {
	// Gets the buffer reader for the input
	ReadPtr() unsafe.Pointer
}

// Represents an input buffer, from which data can be read, but
// not necessarily written
type IBuffer interface {
	view  // The buffer is a view
	Input // The buffer can act as an input
}

// The interface for Go objects used as output from HDF5
// If the structure implements this interface,
// it can be used to read data from a HDF5 file
type Output interface {
	// Gets the buffer writer for the output
	WritePtr() unsafe.Pointer
}

// Represents an output buffer, where data can be written, but
// not necessarily read
type OBuffer interface {
	view   // The buffer is a view
	Output // The buffer can act as an output
}

// Represents the general interface for the hyperblobs of data which
// can be used as buffer by the library.
// If a user defined structure implements this interface, it can be
// transparently used to read or write data from the HDF5 files.
type Buffer interface {
	view   // The buffer is a view
	Input  // Data can be read from the interface
	Output // Data can be written to the interface
}

// Wraps a pointer and type into a Buffer object
type wptr struct {
	dtype h5t.Datatype   // The datatype
	p     unsafe.Pointer // The pointer
}

// Implements the Buffer interface
func (w wptr) Type() (h5t.Datatype, error)   { return w.dtype, nil }
func (w wptr) Shape() (h5s.Dataspace, error) { return h5s.CreateScalar() }
func (w wptr) ReadPtr() unsafe.Pointer       { return w.p }
func (w wptr) WritePtr() unsafe.Pointer      { return w.p }

// Wraps the datatype and pointer into a buffer
func Wrap(dtype h5t.Datatype, ptr unsafe.Pointer) Buffer {
	return &wptr{
		dtype: dtype,
		p:     ptr,
	}
}
