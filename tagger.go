package main

import (
	"fmt"
	"log"
)

func tagImg(url string, vision Tagger) []Tag {
	result, err := vision.Tag(url)
	if err != nil {
		log.Println(fmt.Sprintf("While trying to tag %s got the following error: %s", url, err))
		return make([]Tag, 0)
	}
	response := make([]Tag, len(result.Tags))
	for i, visionTag := range result.Tags {
		response[i] = Tag{
			Tag: visionTag,
		}
	}
	return response
}
