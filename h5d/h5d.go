// This wraps the H5D* family of functions, for creating and
// manipulating datasets
package h5d

/*
#cgo LDFLAGS: -lhdf5
#include <hdf5.h>
*/
import "C"

import "fmt"

import (
	"github.com/valoox/h5go/core"
	"github.com/valoox/h5go/h5l"
	"github.com/valoox/h5go/h5p"
	"github.com/valoox/h5go/h5s"
	"github.com/valoox/h5go/h5t"
)

// Represents the layout of the data
type Layout int

// C conversion of the layout
func (self Layout) C() C.H5D_layout_t {
	return C.H5D_layout_t(self)
}

// Definition of the layouts
const (
	Compact    Layout = C.H5D_COMPACT
	Contiguous Layout = C.H5D_CONTIGUOUS
	Chunked    Layout = C.H5D_CHUNKED
)

// Default property lists
const (
	// Default create
	DefaultCreate = Crt(h5p.Default)
	// Default access
	DefaultAccess = Acc(h5p.Default)
	// Default transfer
	DefaultXfer = Xfer(h5p.Default)
)

// Creates a new dataset creation property list
func Creation() (Crt, error) {
	id, err := h5p.Create(h5p.DATASET_CREATE)
	return Crt(id), err
}

// Represents a dataset creation property list
type Crt h5p.Property

// The id of the property list
func (self Crt) Id() h5p.Property { return h5p.Property(self) }

// The class of the object
func (self Crt) Class() h5p.Class { return h5p.DATASET_CREATE }

// Closes the property list
func (self Crt) Close() error { return h5p.Close(self.Id()) }

// Copies the property list
func (self Crt) Copy() (Crt, error) {
	id, err := h5p.Copy(self.Id())
	return Crt(id), err
}

// Sets the size of the chunks to store a dataset
// This should have the same dimension as the dataset,
// and each dimension is the length of the chunk in that
// dimension.
func (self Crt) SetChunk(dims []int) error {
	ndims := len(dims)
	cdims := make([]C.hsize_t, ndims)
	for i, d := range dims {
		cdims[i] = C.hsize_t(d)
	}
	return core.Status(int(C.H5Pset_chunk(C.hid_t(self),
		C.int(ndims), &cdims[0])),
		"setting dataset chunk size")
}

// Gets the size of the chunks
func (self Crt) GetChunk() ([]int, error) {
	n := 8
	actual := n
	out := make([]C.hsize_t, n)
	for actual >= n {
		actual = int(C.H5Pget_chunk(C.hid_t(self),
			C.int(n),
			&out[0]))
		if err := core.Status(actual,
			"getting dataset chunk size"); err != nil {
			return nil, err
		}
		// Expanding n
		out = append(out, make([]C.hsize_t, n)...)
		n <<= 1
	}
	res := make([]int, actual)
	for i, d := range out[:actual] {
		res[i] = int(d)
	}
	return res, nil
}

// Sets the layout of the dataset
func (self Crt) SetLayout(layout Layout) error {
	return core.Status(int(C.H5Pset_layout(C.hid_t(self),
		layout.C())), "setting dataset layout")
}

// Gets the layout of the dataset
func (self Crt) GetLayout() (Layout, error) {
	layout := C.H5Pget_layout(C.hid_t(self))
	return Layout(layout), core.Status(int(layout),
		"getting layout")
}

// Creates a new property list for accessing a dataset
func Access() (Acc, error) {
	id, err := h5p.Create(h5p.DATASET_ACCESS)
	return Acc(id), err
}

// The property list for accessing a dataset
type Acc h5p.Property

// Implements the h5p.PList interface
func (self Acc) Id() h5p.Property { return h5p.Property(self) }
func (self Acc) Class() h5p.Class { return h5p.DATASET_ACCESS }
func (self Acc) Close() error     { return h5p.Close(self.Id()) }
func (self Acc) Copy() (Acc, error) {
	id, err := h5p.Copy(self.Id())
	return Acc(id), err
}

// Creates a new property list for transferring a dataset
func Transfer() (Xfer, error) {
	id, err := h5p.Create(h5p.DATASET_XFER)
	return Xfer(id), err
}

// Property list for transferring data
type Xfer h5p.Property

// Implements the h5p.PList interface
func (self Xfer) Id() h5p.Property { return h5p.Property(self) }
func (self Xfer) Class() h5p.Class { return h5p.DATASET_XFER }
func (self Xfer) Close() error     { return h5p.Close(self.Id()) }
func (self Xfer) Copy() (Xfer, error) {
	id, err := h5p.Copy(self.Id())
	return Xfer(id), err
}

// Represents an Id specifically for datasets
type Dataset core.Id

// The dataspace for this dataset
// Wraps the H5Dget_space function
func (d Dataset) Shape() (h5s.Dataspace, error) {
	out := h5s.Dataspace(C.H5Dget_space(C.hid_t(d)))
	return out, core.Status(int(out),
		"getting associated dataspace of dataset %v",
		d)
}

// The datatype for this dataset
// Wraps the H5Dget_type function
func (d Dataset) Type() (h5t.Datatype, error) {
	out := h5t.Datatype(C.H5Dget_type(C.hid_t(d)))
	return out, core.Status(int(out),
		"getting associated datatype of dataset %v",
		d)
}

// The HDF5 Id for this dataset
func (d Dataset) Id() core.Id { return core.Id(d) }

// Closes the dataset
// Wraps the H5Dclose function
func (d Dataset) Close() error {
	return core.Status(int(C.H5Dclose(C.hid_t(d))),
		"closing dataset")
}

// Sets the dimensions of the dataset to the given dimensions
// Returns an error if the dataset is not chunked, or if the
// dimensions are superior to the maximum dimension available
// Wraps the H5Dset_extent function
func (d Dataset) SetDims(dims []int) error {
	cdims := make([]C.hsize_t, len(dims))
	for i, x := range dims {
		cdims[i] = C.hsize_t(x)
	}
	return core.Status(int(C.H5Dset_extent(C.hid_t(d),
		&cdims[0])),
		"setting array extent")
}

// Writes the content of the buffer in the dataset
// Wraps the H5Dwrite function
func (d Dataset) Write(data IBuffer, selection h5s.Dataspace, xfr Xfer) error {
	T, err := data.Type()
	if err != nil {
		return err
	}
	defer T.Close()
	ishape, err := data.Shape()
	if err != nil {
		return err
	}
	defer ishape.Close()
	return core.Status(int(C.H5Dwrite(
		C.hid_t(d),
		C.hid_t(T),
		C.hid_t(ishape),
		C.hid_t(selection),
		C.hid_t(xfr),
		data.ReadPtr())),
		"writing data to dataset")
}

// Reads the data into the provided buffer
func (d Dataset) Read(data OBuffer, selection h5s.Dataspace, xfr Xfer) error {
	T, err := data.Type()
	if err != nil {
		return err
	}
	defer T.Close()
	oshape, err := data.Shape()
	if err != nil {
		return err
	}
	defer oshape.Close()
	return core.Status(int(C.H5Dread(
		C.hid_t(d),
		C.hid_t(T),
		C.hid_t(oshape),
		C.hid_t(selection),
		C.hid_t(xfr),
		data.WritePtr())),
		"reading data from dataset")
}

// Tries to return the status, raising and error if it is negative
func try(id Dataset, context string, args ...interface{}) (Dataset, error) {
	return id, core.Status(int(id), fmt.Sprintf(context, args...))
}

// Creates a new dataset with all the provided information
// Wraps the H5Dcreate function
func Create(at core.Location, name core.Path, dtype h5t.Datatype,
	dspace h5s.Dataspace,
	link h5l.Crt, c Crt, a Acc) (Dataset, error) {
	return try(Dataset(C.H5Dcreate2(C.hid_t(at.At()),
		C.CString(name.String()),
		C.hid_t(dtype), C.hid_t(dspace),
		C.hid_t(link),
		C.hid_t(c), C.hid_t(a))),
		"creating dataset at %s", name)

}

// Opens an existing location, from a root location and a path
// Wraps the H5Dopen function
func Open(at core.Location, name core.Path, access Acc) (Dataset, error) {
	return try(Dataset(C.H5Dopen2(C.hid_t(at.At()),
		C.CString(name.String()), C.hid_t(access))),
		"opening dataset at %s", name)
}
