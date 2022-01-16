package main

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// anime-name ep_no  -> Input Format

var rootDomain = "https://ww1.gogoanime2.org"

func main() {
	myFigure := figure.NewColorFigure("MAAL", "doh", "blue", true)
	myFigure.Print()
	fmt.Println()
	videoPlayer, err := installPlayer()
	if err != nil {
		log.Fatalln("Please Install the mpv videoPlayer(https://mpv.io) and add it to the PATH")
	}
	for {
		var animeName, epNo string

		fmt.Println(colorGreen, "====> Enter Anime-Name Episode-Number in \"name ep-no\" format ")
		fmt.Println("----> Example: (shingeki-no-kyojin 1) (shingeki-no-kyojin-dub 1), for movies enter 1 for episode (kimi-no-na-wa 1)")
		_, err := fmt.Scan(&animeName)
		HandleError(err, "Error while Reading the Anime-Name from the command-line: ")
		_, err = fmt.Scan(&epNo)
		HandleError(err, "Error while Reading the Episode-Number from the command-line: ")

		fmt.Print(colorBlue)
		fmt.Printf("Playing %s episode %s \n", animeName, epNo)

		client := &http.Client{}
		request, err := http.NewRequest("GET", fmt.Sprintf("%s/watch/%s/%s", rootDomain, animeName, epNo), nil)
		HandleError(err, "Error while creating http.NewRequest: ")

		request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36")
		resp, err := client.Do(request)
		HandleError(err, "Error while making the request: ")

		if resp.StatusCode == 200 {
			// _, err = io.Copy(os.Stdout, resp.Body)
			dataInBytes, err := ioutil.ReadAll(resp.Body)
			HandleError(err, "Error while Reading the Body: ")

			content := string(dataInBytes)
			iframeStartIndex := strings.Index(content, "<iframe id=\"playerframe\" src=")
			if iframeStartIndex == -1 {
				log.Println(colorRed, "Video src doesn't exist")
				return
			}
			startIndex := iframeStartIndex + 30
			content = content[startIndex:]
			endIndex := strings.Index(content, "\" style=\"width: 100%;\"")
			src := rootDomain + content[:endIndex]
			fmt.Println(colorCyan, "Link Collected: ", src)
			cmd := exec.Command(videoPlayer, src)
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Println(colorRed, "Error while playing the video: ", err)
				return
			}
		} else {
			log.Println(colorRed, "Anime Doesn't exist, Recheck all the information!")
			log.Printf("Verify information on %s", rootDomain)
		}
	}
}
