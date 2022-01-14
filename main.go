package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// anime-name ep_no sub/dub

var rootDomain = "https://ww1.gogoanime2.org"

func main() {
	arguments := os.Args
	var animeName, epNo string
	if len(arguments) == 1 {
		_, err := fmt.Scan(&animeName)
		if err != nil {
			log.Fatalln("Error while Reading the animeName from the command-line: ", err)
		}
		_, err = fmt.Scan(&epNo)
		if err != nil {
			log.Fatalln("Error while Reading the epNo from the command-line: ", err)
		}
	} else {
		animeName = arguments[1]
		epNo = arguments[2]
	}
	client := &http.Client{}
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/watch/%s/%s", rootDomain, animeName, epNo), nil)
	if err != nil {
		log.Println(err)
		return
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36")
	resp, err := client.Do(request)
	if err != nil {
		log.Println("Error while making the request", err)
		return
	}
	if resp.StatusCode == 200 {
		// _, err = io.Copy(os.Stdout, resp.Body)
		dataInBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error while Reading the Body: ", err)
			return
		}
		content := string(dataInBytes)
		iframeStartIndex := strings.Index(content, "<iframe id=\"playerframe\" src=")
		if iframeStartIndex == -1 {
			fmt.Println(content, iframeStartIndex)
			log.Println("Video src doesn't exist")
			return
		}
		index := iframeStartIndex + 30
		src := rootDomain + content[index:index+15]
		fmt.Println(src)
		cmd := exec.Command("T://bootstrapper/mpv.exe", src)
		// cmd.Path = "T://bootstrapper"
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Println("Error while playing the video: ", err)
			return
		}
	}
}
