package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var AnimeNotExist = errors.New("anime Source Doesn't Exist")

func getVideoUrl(client *http.Client, animeName string, epNo string) (error, string) {
	url := fmt.Sprintf("%s/watch/%s/%s", rootDomain, animeName, epNo)
	err, content := makeReq(client, url)
	if err != nil {
		return err, ""
	}

	iframeStartIndex := strings.Index(content, "<iframe id=\"playerframe\" src=")
	if iframeStartIndex == -1 {
		return errors.New("unable to find IFrame[Level 1]"), ""
	}
	startIndex := iframeStartIndex + 30
	content = content[startIndex:]
	endIndex := strings.Index(content, "\" style=\"width: 100%;\"")
	videoSrc := content[:endIndex]
	if strings.HasPrefix(videoSrc, "//") {
		videoSrc = "https:" + videoSrc
		err, content = makeReq(client, videoSrc)
		if err != nil {
			return err, ""
		}
		log.Println(content)
		videoFrameStartIndex := strings.Index(content, "<video class=\"jw-video jw-reset\" ")
		if videoFrameStartIndex == -1 {
			return errors.New("unable to find IFrame[Level 2]"), ""
		}
		startIndex = strings.Index(content[videoFrameStartIndex:], "src=")
		endIndex = strings.Index(content[startIndex:], "\" style></video>")
		videoSrc = content[startIndex:endIndex]
	} else {
		videoSrc = rootDomain + videoSrc
	}
	return nil, videoSrc
}

func makeReq(client *http.Client, url string) (error, string) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		HandleError(err, "Error while creating http.NewRequest: ")
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36")
	resp, err := client.Do(request)
	if err != nil {
		return err, ""
	}
	if resp.StatusCode == 200 {
		dataInBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err, ""
		}

		content := string(dataInBytes)
		return nil, content

	} else {
		return AnimeNotExist, ""
	}
}
