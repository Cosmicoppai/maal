package main

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
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
	videoPlayer, err := installPlayer()
	if err != nil {
		log.Fatalln("Please Install the mpv videoPlayer(https://mpv.io) and add it to the PATH")
	}
	client := &http.Client{}
	fmt.Print(colorGreen)
	fmt.Println("====> Enter Anime-Name Episode-Number in \"name ep-no\" format ")
	fmt.Println("----> Example: (shingeki-no-kyojin 1) (shingeki-no-kyojin-dub 1), for movies enter 1 for episode (kimi-no-na-wa 1)")
	for {
		fmt.Println(strings.Repeat("==", 80))
		var animeName, epNo string
		_, err := fmt.Scan(&animeName)
		if err != nil {
			HandleError(err, "Error while Reading the Anime-Name from the command-line: ")
			continue
		}
		_, err = fmt.Scan(&epNo)
		if err != nil {
			HandleError(err, "Error while Reading the Episode-Number from the command-line: ")
			continue
		}

		fmt.Print(colorBlue)
		fmt.Printf("Playing %s episode %s \n", animeName, epNo)

		err, src := getVideoUrl(client, animeName, epNo)
		if err != nil {
			HandleError(err, "")
			continue
		}
		fmt.Println(colorCyan, "Link Collected: ", src)
		cmd := exec.Command(videoPlayer, src)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Println(colorRed, "Error while playing the video: ", err)
		}
	}
}
