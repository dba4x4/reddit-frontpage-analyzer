package analyzer

import vision "github.com/ahmdrz/microsoft-vision-golang"

type tagger interface {
	Tag(url string) (vision.VisionResult, error)
}
