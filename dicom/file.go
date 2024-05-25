package dicom

import (
	"errors"
	"image"
	"io"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

// File represents a DICOM file
type File struct {
	ID string

	file io.ReadSeeker
	size int64

	// read-thru cache to avoid unnecessary parsing and image processing
	cache struct {
		dataset *dicom.Dataset
		image   *image.Gray
	}
}

// NewFile constructs a new DICOM file
func NewFile(
	id string,
	size int64,
	contents io.ReadSeeker,
) File {
	return File{
		ID:   id,
		size: size,
		file: contents,
	}
}

// Size returns the size of the DICOM file contents
func (d File) Size() int64 {
	return d.size
}

// Raw returns the raw DICOM file
func (d File) Raw() io.ReadSeeker {
	return d.file
}

// DataSet returns a parsed datastructure representing the data within
// the DICOM file
func (d File) DataSet() (*dicom.Dataset, error) {
	if d.cache.dataset != nil {
		return d.cache.dataset, nil
	}

	data, err := dicom.Parse(d.file, d.size, nil)
	if err != nil {
		return nil, err
	}

	d.cache.dataset = &data
	return d.cache.dataset, nil
}

// PNGGenerateOptions options for generating a PNG from a DICOM
type PNGGenerateOptions struct {
	shouldRemap bool
}

// PNGRemapValues option to remap pixel value ranges in PNG for better visibility
func PNGRemapPixels(shouldRemap bool) func(opts *PNGGenerateOptions) {
	return func(opts *PNGGenerateOptions) {
		opts.shouldRemap = shouldRemap
	}
}

// PNG returns a greyscale PNG of the DICOM file
func (d File) PNG(options ...func(opts *PNGGenerateOptions)) (*image.Gray, error) {
	if d.cache.image != nil {
		return d.cache.image, nil
	}

	var opts = PNGGenerateOptions{
		shouldRemap: true,
	}
	for _, opt := range options {
		opt(&opts)
	}

	dataSet, err := d.DataSet()
	if err != nil {
		return nil, err
	}
	images, err := generateImage(
		*dataSet,
		opts.shouldRemap,
	)
	if err != nil {
		return nil, err
	}

	if len(images) == 0 {
		return nil, errors.New("file does not contain any images")
	}

	d.cache.image = images[0]

	return d.cache.image, nil
}

// DICOMElementsLookup is a map of DICOM tags to their corresponding elements
type DICOMElementsLookup map[string]*dicom.Element

// FindElements given a list of tags, retrieves a lookup of elements by tag
func (d File) FindElements(tags ...tag.Tag) (DICOMElementsLookup, error) {
	dataSet, err := d.DataSet()
	if err != nil {
		return nil, err
	}

	elementsByTag := make(map[string]*dicom.Element)
	if len(tags) > 0 {
		for _, tag := range tags {
			// ignore error and assume element cannot be found
			element, _ := dataSet.FindElementByTag(tag)
			elementsByTag[tag.String()] = element
		}
	}

	return elementsByTag, nil
}

// AllElements retrieves a lookup of all elements by tag
func (d File) AllElements() (DICOMElementsLookup, error) {
	dataSet, err := d.DataSet()
	if err != nil {
		return nil, err
	}

	elementsByTag := make(map[string]*dicom.Element)
	for _, element := range dataSet.Elements {
		elementsByTag[element.Tag.String()] = element
	}

	return elementsByTag, nil
}
