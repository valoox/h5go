H5go
====

Wrapping of the HDF5 library in Go, strongly inspired by h5py.

The package is composed of the core subpackage (which contains common
abstractions for the entire package) and subpackages corresponding to
the different HDF5 prefixes.

The goal is to build an entire layer of idiomatic Go abstraction on
top of it, so that HDF5 can be used from Go as h5py is used from
Python: without any knowledge of the actual C implementation. For the
moment, most of the operations are done through the h5X subpackages,
which usually simply wrap the C function using cgo.

This is a work in progress and contributions are welcome ! If you
have a nice idea to integrate HDF5 nicely with Go, feel free to
fork or propose your own contributions.
