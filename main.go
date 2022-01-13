package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	var rootDomain = "https://www1.gogoanime2.org"
	client := &http.Client{}
	request, err := http.NewRequest("GET", rootDomain+"/watch/shingeki-no-kyojin/1", nil)
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
		fmt.Println(content)
	}
}
