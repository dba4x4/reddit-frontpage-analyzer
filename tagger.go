package main

import (
	"fmt"
	"log"

	vision "github.com/ahmdrz/microsoft-vision-golang"
	"github.com/spf13/viper"
)

func tagImg(url string) []Tag {
	vision, err := vision.New(viper.GetString("microsoft.key"))
	if err != nil {
		log.Fatalln(err)
	}
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
