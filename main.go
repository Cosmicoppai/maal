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

// anime-name ep_no  -> Input Format

func HandleReadingError(err error, issueIn string) {
	if err != nil {
		log.Fatalln(fmt.Sprintf("Error while Reading the %s from the command-line: ", issueIn), err)
	}
}

var rootDomain = "https://ww1.gogoanime2.org"

func main() {
	for {
		var animeName, epNo string

		fmt.Println("Enter Anime-Name Episode-Number in \"name ep-no\" format ")
		fmt.Println("Example: (shingeki-no-kyojin 1) (shingeki-no-kyojin-dub 1), for movies enter 1 for episode (kimi-no-na-wa 1)")
		_, err := fmt.Scan(&animeName)
		HandleReadingError(err, "Anime-Name")
		_, err = fmt.Scan(&epNo)
		HandleReadingError(err, "Episode-Number")

		fmt.Printf("Playing %s episode %s \n", animeName, epNo)

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
			startIndex := iframeStartIndex + 30
			content = content[startIndex:]
			endIndex := strings.Index(content, "\" style=\"width: 100%;\"")
			src := rootDomain + content[:endIndex]
			fmt.Println("Link Collected: ", src)
			cmd := exec.Command("T://bootstrapper/mpv.exe", src)
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Println("Error while playing the video: ", err)
				return
			}
		} else {
			log.Println("Anime Doesn't exist, Recheck all the information!")
		}
	}
}
