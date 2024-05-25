package dicom

import (
	"errors"
	"image"
	"image/color"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/frame"
	"github.com/suyashkumar/dicom/pkg/tag"
)

// domain1D represents a one-dimensional domain
type domain1D struct {
	min int
	max int
}

// targetPixelDomain default target domain for a greyscale image
func targetPixelDomain() domain1D {
	return domain1D{
		min: 0,
		max: 255,
	}
}

// findBounds returns the domain that bounds all pixel values in the provided
// frame
func findBounds(frame *frame.NativeFrame) domain1D {

	if len(frame.Data) == 0 {
		return domain1D{}
	}

	if len(frame.Data[0]) == 0 {
		return domain1D{}
	}

	min := frame.Data[0][0]
	max := frame.Data[0][0]
	for _, dataRow := range frame.Data {

		if len(dataRow) == 0 {
			continue
		}

		val := dataRow[0]
		if val < min {
			min = val
		}
		if val > max {
			max = val
		}
	}
	return domain1D{min, max}
}

// mapValue utility to map a value within a source domain to a target domain
func mapValue(
	value int,
	sourceDomain domain1D,
	targetDomain domain1D,
) int {
	if value < sourceDomain.min {
		value = sourceDomain.min
	}
	if value > sourceDomain.max {
		value = sourceDomain.max
	}

	targetDomainRange := targetDomain.max - targetDomain.min
	sourceDomainRange := sourceDomain.max - sourceDomain.min

	return targetDomain.min + (targetDomainRange)*(value-sourceDomain.min)/(sourceDomainRange)
}

// generateImage utility to generate an 8-bit normalized image from a
// dicom dataset. generateImage optionally re-maps the pixel range to allow for better
// image visibility. See https://github.com/suyashkumar/dicom/issues/301 for
// more details
func generateImage(dataset dicom.Dataset, shouldRemap bool) ([]*image.Gray, error) {
	pixelDataElement, err := dataset.FindElementByTag(tag.PixelData)
	if err != nil {
		return nil, err
	}

	var pixelDataInfo dicom.PixelDataInfo
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = errors.New("failed to get pixel data info for DICOM file")
			}
		}()

		// recover from potential panic as an error
		pixelDataInfo = dicom.MustGetPixelDataInfo(pixelDataElement.Value)
	}()
	if err != nil {
		return nil, err
	}

	pixelTargetDomain := targetPixelDomain()

	var images []*image.Gray
	for _, fr := range pixelDataInfo.Frames {
		// get source domain of pixels in the frame
		pixelSourceDomain := findBounds(&fr.NativeData)

		// generate a blank greyscale image
		outImage := image.NewGray(
			image.Rect(
				0,
				0,
				fr.NativeData.Cols,
				fr.NativeData.Rows,
			),
		)

		for idx := 0; idx < len(fr.NativeData.Data); idx++ {
			x := idx % fr.NativeData.Cols
			y := idx / fr.NativeData.Cols

			pixelValue := 0
			if len(fr.NativeData.Data[idx]) > 0 {
				pixelValue = fr.NativeData.Data[idx][0]
			}

			// remap pixel value to target domain (0-255)
			mappedPixelValue := pixelValue
			if shouldRemap {
				mappedPixelValue = mapValue(
					pixelValue,
					pixelSourceDomain,
					pixelTargetDomain,
				)
			}

			// write to image
			outImage.SetGray(
				x,
				y,
				color.Gray{
					Y: uint8(mappedPixelValue),
				},
			)
		}

		images = append(images, outImage)
	}

	return images, nil
}
