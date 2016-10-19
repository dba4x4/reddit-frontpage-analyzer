package test

import (
	"reflect"
	"testing"

	vision "github.com/ahmdrz/microsoft-vision-golang"
)

const firstImage = "http://test.com/firstimage.jpg"
const secondImage = "http://test.com/secondImage.jpg"

type MockedVision struct {
}

func (mV MockedVision) Tag(url string) (vision.VisionResult, error) {
	switch url {
	case firstImage:
		return vision.VisionResult{
			Tags: []vision.Tag{
				vision.Tag{
					Name:       "Person",
					Confidence: 0.95,
				},
			},
		}, nil
	case secondImage:
		return vision.VisionResult{
			Tags: []vision.Tag{
				vision.Tag{
					Name:       "Dog",
					Confidence: 0.95,
				},
				vision.Tag{
					Name:       "Grass",
					Confidence: 0.75,
				},
			},
		}, nil
	}
	return vision.VisionResult{}, nil
}

func Test_tagImg(t *testing.T) {
	type args struct {
		url    string
		vision Tagger
	}
	tests := []struct {
		name string
		args args
		want []Tag
	}{
		{
			"Person image",
			args{
				firstImage,
				MockedVision{},
			},
			[]Tag{
				Tag{
					Tag: vision.Tag{
						Name:       "Person",
						Confidence: 0.95,
					},
				},
			},
		},
		{
			"Dog running on grass image",
			args{
				secondImage,
				MockedVision{},
			},
			[]Tag{
				Tag{
					Tag: vision.Tag{
						Name:       "Dog",
						Confidence: 0.95,
					},
				},
				Tag{
					Tag: vision.Tag{
						Name:       "Grass",
						Confidence: 0.75,
					},
				},
			},
		},
		{
			"Unkonwn image",
			args{
				"http://test.com/unknown.jpg",
				MockedVision{},
			},
			[]Tag{},
		},
	}
	for _, tt := range tests {
		if got := tagImg(tt.args.url, tt.args.vision); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. tagImg() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
