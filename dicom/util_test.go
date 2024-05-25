package dicom

import (
	"image"
	"reflect"
	"testing"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/frame"
	"github.com/suyashkumar/dicom/pkg/tag"
)

func Test_findBounds(t *testing.T) {
	type args struct {
		frame *frame.NativeFrame
	}
	tests := []struct {
		name string
		args args
		want domain1D
	}{
		{
			name: "returns zero domain for empty frame",
			args: args{
				frame: &frame.NativeFrame{},
			},
			want: domain1D{},
		},
		{
			name: "returns zero domain for empty frame variant 2",
			args: args{
				frame: &frame.NativeFrame{
					Data: [][]int{},
				},
			},
			want: domain1D{},
		},
		{
			name: "returns zero domain for empty frame variant 3",
			args: args{
				frame: &frame.NativeFrame{
					Data: [][]int{
						{},
					},
				},
			},
			want: domain1D{},
		},
		{
			name: "returns zero domain for empty frame variant 4",
			args: args{
				frame: &frame.NativeFrame{
					Data: [][]int{
						{},
						{},
					},
				},
			},
			want: domain1D{},
		},
		{
			name: "returns zero domain for empty frame variant 4",
			args: args{
				frame: &frame.NativeFrame{
					Data: [][]int{
						{},
						{},
					},
				},
			},
			want: domain1D{},
		},
		{
			name: "returns correct domain for 1 pixel frame",
			args: args{
				frame: &frame.NativeFrame{
					Data: [][]int{
						{1},
					},
				},
			},
			want: domain1D{
				min: 1,
				max: 1,
			},
		},
		{
			name: "returns correct domain for 2 pixel frame",
			args: args{
				frame: &frame.NativeFrame{
					Data: [][]int{
						{0},
						{5},
					},
				},
			},
			want: domain1D{
				min: 0,
				max: 5,
			},
		},
		{
			name: "returns correct domain for n pixel frame",
			args: args{
				frame: &frame.NativeFrame{
					Data: [][]int{
						{0},
						{5},
						{2},
						{9},
						{6},
						{12},
					},
				},
			},
			want: domain1D{
				min: 0,
				max: 12,
			},
		},
		{
			name: "returns correct domain for n pixel frame w/ negatives",
			args: args{
				frame: &frame.NativeFrame{
					Data: [][]int{
						{-2},
						{5},
						{2},
						{9},
						{6},
						{12},
					},
				},
			},
			want: domain1D{
				min: -2,
				max: 12,
			},
		},
		{
			name: "returns domain that ignores irrelevant indexes ",
			args: args{
				frame: &frame.NativeFrame{
					Data: [][]int{
						{-2, 99},
						{5, 22},
						{2, 4},
						{9, 10},
						{6, 1},
						{12, -4},
					},
				},
			},
			want: domain1D{
				min: -2,
				max: 12,
			},
		},
		{
			name: "returns correct domain for corrupt inner arrays",
			args: args{
				frame: &frame.NativeFrame{
					Data: [][]int{
						{-2, 99},
						{5, 22},
						{2, 4},
						{9, 10},
						{6, 1},
						{},
					},
				},
			},
			want: domain1D{
				min: -2,
				max: 9,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findBounds(tt.args.frame); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findBounds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mapValue(t *testing.T) {
	type args struct {
		value        int
		sourceDomain domain1D
		targetDomain domain1D
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "outputs zero for identical input and output domain",
			args: args{
				value: 0,
				sourceDomain: domain1D{
					min: 0,
					max: 1,
				},
				targetDomain: domain1D{
					min: 0,
					max: 1,
				},
			},
			want: 0,
		},
		{
			name: "outputs 1 for identical input and output domain",
			args: args{
				value: 1,
				sourceDomain: domain1D{
					min: 0,
					max: 1,
				},
				targetDomain: domain1D{
					min: 0,
					max: 1,
				},
			},
			want: 1,
		},
		{
			name: "outputs n for identical input and output domain",
			args: args{
				value: 1,
				sourceDomain: domain1D{
					min: 0,
					max: 5,
				},
				targetDomain: domain1D{
					min: 0,
					max: 5,
				},
			},
			want: 1,
		},
		{
			name: "outputs n for identical input and output domain variant 2",
			args: args{
				value: 4,
				sourceDomain: domain1D{
					min: 0,
					max: 5,
				},
				targetDomain: domain1D{
					min: 0,
					max: 5,
				},
			},
			want: 4,
		},
		{
			name: "maps 0 for similar input and output domains",
			args: args{
				value: 0,
				sourceDomain: domain1D{
					min: 0,
					max: 2,
				},
				targetDomain: domain1D{
					min: 1,
					max: 2,
				},
			},
			want: 1,
		},
		{
			name: "maps 0 for very different input and output domains",
			args: args{
				value: 0,
				sourceDomain: domain1D{
					min: 0,
					max: 2,
				},
				targetDomain: domain1D{
					min: 1,
					max: 999999,
				},
			},
			want: 1,
		},
		{
			name: "maps values at source start for similar input and output domains",
			args: args{
				value: 1,
				sourceDomain: domain1D{
					min: 0,
					max: 2,
				},
				targetDomain: domain1D{
					min: 1,
					max: 5,
				},
			},
			want: 3,
		},
		{
			name: "maps values at source start very different input and output domains",
			args: args{
				value: 1,
				sourceDomain: domain1D{
					min: 0,
					max: 2,
				},
				targetDomain: domain1D{
					min: 1,
					max: 999999,
				},
			},
			want: 500000,
		},

		{
			name: "maps values at from non-zero rooted source domains",
			args: args{
				value: 40,
				sourceDomain: domain1D{
					min: 10,
					max: 99,
				},
				targetDomain: domain1D{
					min: 1,
					max: 5,
				},
			},
			want: 2,
		},

		{
			name: "maps to inverted target domains",
			args: args{
				value: 0,
				sourceDomain: domain1D{
					min: 0,
					max: 5,
				},
				targetDomain: domain1D{
					min: 5,
					max: 0,
				},
			},
			want: 5,
		},

		{
			name: "clamps to target min if value is less than source min",
			args: args{
				value: -1,
				sourceDomain: domain1D{
					min: 0,
					max: 5,
				},
				targetDomain: domain1D{
					min: 0,
					max: 5,
				},
			},
			want: 0,
		},

		{
			name: "clamps to target max if value is greater than source max",
			args: args{
				value: 6,
				sourceDomain: domain1D{
					min: 0,
					max: 5,
				},
				targetDomain: domain1D{
					min: 0,
					max: 5,
				},
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapValue(tt.args.value, tt.args.sourceDomain, tt.args.targetDomain); got != tt.want {
				t.Errorf("mapValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateImage(t *testing.T) {

	emptyDicomPixelElement, _ := dicom.NewValue(dicom.PixelDataInfo{

		// Frames hold the processed PixelData frames (either Native or Encapsulated
		// PixelData).
		Frames: make([]*frame.Frame, 0),
	})

	validDicomPixelElement, _ := dicom.NewValue(dicom.PixelDataInfo{

		// Frames hold the processed PixelData frames (either Native or Encapsulated
		// PixelData).
		Frames: []*frame.Frame{
			{
				NativeData: frame.NativeFrame{
					Data: [][]int{{0}, {1}, {3}, {4}},
					Rows: 2,
					Cols: 2,
				},
			},
		},
	})

	type args struct {
		dataset dicom.Dataset
	}
	tests := []struct {
		name    string
		args    args
		want    []*image.Gray
		wantErr bool
	}{
		{
			name: "errors on empty dataset",
			args: args{
				dataset: dicom.Dataset{},
			},
			wantErr: true,
		},
		{
			name: "errors on missing element",
			args: args{
				dataset: dicom.Dataset{
					Elements: []*dicom.Element{
						{
							Tag: tag.Tag{
								Group:   0,
								Element: 0,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "errors on corrupt pixel data element",
			args: args{
				dataset: dicom.Dataset{
					Elements: []*dicom.Element{
						{
							Tag: tag.PixelData,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "returns no images for empty pixel element",
			args: args{
				dataset: dicom.Dataset{
					Elements: []*dicom.Element{
						{
							Tag:   tag.PixelData,
							Value: emptyDicomPixelElement,
						},
					},
				},
			},
			wantErr: false,
			want:    nil,
		},
		{
			name: "returns pixel-scaled image for valid pixel element",
			args: args{
				dataset: dicom.Dataset{
					Elements: []*dicom.Element{
						{
							Tag:   tag.PixelData,
							Value: validDicomPixelElement,
						},
					},
				},
			},
			wantErr: false,
			want: []*image.Gray{
				{
					Pix:    []uint8{0, 63, 191, 255},
					Stride: 2,
					Rect:   image.Rect(0, 0, 2, 2),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateImage(tt.args.dataset, true)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for idx := range tt.want {
				wantImages := tt.want[idx]
				gotImages := got[idx]

				for idx2 := range wantImages.Pix {
					if wantImages.Pix[idx2] != gotImages.Pix[idx2] {
						t.Errorf("generateImage() = %+v, want %+v", got, tt.want)
						return
					}
				}
			}
		})
	}
}
