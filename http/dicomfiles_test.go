package http

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/suyashkumar/dicom/pkg/tag"
)

func Test_dicomFiles_parseTagQuery(t *testing.T) {
	type args struct {
		query url.Values
	}
	tests := []struct {
		name    string
		args    args
		want    []tag.Tag
		wantErr bool
	}{
		{
			name: "returns empty slice for empty query params",
			args: args{
				query: url.Values{
					"": {""},
				},
			},
		},
		{
			name: "returns empty slice for empty query params variant 2",
			args: args{
				query: url.Values{
					"": {""},
				},
			},
		},
		{
			name: "errors for malformed tag query param",
			args: args{
				query: url.Values{
					"tag": {"("},
				},
			},
			wantErr: true,
		},
		{
			name: "errors for malformed tag query param variant 2",
			args: args{
				query: url.Values{
					"tag": {")"},
				},
			},
			wantErr: true,
		},
		{
			name: "errors for malformed tag query param variant 3",
			args: args{
				query: url.Values{
					"tag": {"()"},
				},
			},
			wantErr: true,
		},
		{
			name: "errors for malformed tag query param variant 4",
			args: args{
				query: url.Values{
					"tag": {")("},
				},
			},
			wantErr: true,
		},
		{
			name: "errors for malformed tag query param variant 5",
			args: args{
				query: url.Values{
					"tag": {"(0,1)(1,1)"},
				},
			},
			wantErr: true,
		},
		{
			name: "returns correct tag for single non-padded tag",
			args: args{
				query: url.Values{
					"tag": {"(0,1)"},
				},
			},
			wantErr: false,
			want: []tag.Tag{
				{Group: 0000, Element: 0001},
			},
		},
		{
			name: "returns correct tag for single padded tag",
			args: args{
				query: url.Values{
					"tag": {"(0000,0001)"},
				},
			},
			wantErr: false,
			want: []tag.Tag{
				{Group: 0000, Element: 0001},
			},
		},
		{
			name: "returns correct tag for multiple non-padded tag",
			args: args{
				query: url.Values{
					"tag": {"(0,1)", "(5,2)"},
				},
			},
			wantErr: false,
			want: []tag.Tag{
				{Group: 0000, Element: 0001},
				{Group: 0005, Element: 0002},
			},
		},
		{
			name: "returns correct tag for single padded tag",
			args: args{
				query: url.Values{
					"tag": {"(0000,0001)", "(0005,0002)"},
				},
			},
			wantErr: false,
			want: []tag.Tag{
				{Group: 0000, Element: 0001},
				{Group: 0005, Element: 0002},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := dicomFiles{}
			got, err := d.parseTagQuery(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("dicomFiles.parseTagQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dicomFiles.parseTagQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
