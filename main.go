package main

import (
	"encoding/json"
	"fmt"
	"github.com/gboddin/goccm"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var resultPerPage = 500
var downloadCcm = goccm.New(20)

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("sonarcon <url> <action> [options...]")
	}
	baseUrl := strings.TrimSuffix(os.Args[1], "/")
	switch os.Args[2] {
	case "lp":
		lp(baseUrl)
	case "ls":
		if len(os.Args) < 4 {
			log.Fatalln("sonarcon <url> ls <project>")
		}
		ls(baseUrl, os.Args[3], 1)
	case "dump":
		if len(os.Args) < 5 {
			log.Fatalln("sonarcon <url> ls <project> <directory>")
		}
		dump(baseUrl, os.Args[3], 1, os.Args[4])
		downloadCcm.WaitAllDone()
	default:
		log.Fatalln("Valid command are lp/ls/dump")
	}
}
func lp(baseUrl string) {
	resp, err := http.Get(baseUrl + "/api/components/search_projects")
	if err != nil {
		log.Fatalln(err)
	}
	var sonarResp CompResponse
	jsonDecoder := json.NewDecoder(resp.Body)
	err = jsonDecoder.Decode(&sonarResp)
	if err != nil {
		log.Fatalln(err)
	}
	for _, baseComp := range sonarResp.Components {
		fmt.Println(baseComp.Key)
	}

}

func ls(baseUrl, projectKey string, page int) {
	resp, err := http.Get(baseUrl + "/api/components/tree?p=" + strconv.Itoa(page) + "&component=" + projectKey + "&qualifiers=FIL,UTS&ps=" + strconv.Itoa(resultPerPage))
	if err != nil {
		log.Fatalln(err)
	}
	var sonarResp CompResponse
	jsonDecoder := json.NewDecoder(resp.Body)
	err = jsonDecoder.Decode(&sonarResp)
	if err != nil {
		log.Fatalln(err)
	}
	for _, baseComp := range sonarResp.Components {
		fmt.Println(baseComp.Key)
	}
	if sonarResp.Paging != nil && page*resultPerPage < sonarResp.Paging.Total {
		page++
		ls(baseUrl, projectKey, page)
	}
}

func dump(baseUrl, projectKey string, page int, directory string) {
	resp, err := http.Get(baseUrl + "/api/components/tree?p=" + strconv.Itoa(page) + "&component=" + projectKey + "&qualifiers=FIL,UTS&ps=" + strconv.Itoa(resultPerPage))
	if err != nil {
		log.Fatalln(err)
	}
	var sonarResp CompResponse
	jsonDecoder := json.NewDecoder(resp.Body)
	err = jsonDecoder.Decode(&sonarResp)
	if err != nil {
		log.Fatalln(err)
	}
	for _, baseComp := range sonarResp.Components {
		downloadCcm.Wait()
		go download(baseUrl, projectKey, baseComp.Key, directory)
	}
	if sonarResp.Paging != nil && page*resultPerPage < sonarResp.Paging.Total {
		page++
		dump(baseUrl, projectKey, page, directory)
	}
}

func download(baseUrl, projectKey, fileKey, directory string) {
	defer downloadCcm.Done()
	filePath := cleanAbsolutePath(
		strings.Replace(fileKey, projectKey+":", "", 1))
	finalPath := filepath.Join(directory, filePath)
	basePath := filepath.Dir(finalPath)
	log.Printf("Downloading to %s", finalPath)
	info, err := os.Stat(basePath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(basePath, 0700)
		if err != nil && !os.IsExist(err) {
			log.Fatalln(err)
		}
	} else if !info.IsDir() {
		log.Fatalf("%s found but not a directory", basePath)
	}
	file, err := os.Create(finalPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	resp, err := http.Get(baseUrl + "/api/sources/raw?key=" + fileKey)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
}

func cleanAbsolutePath(path string) string {
	path = "/" + path
	path = filepath.Clean(path)
	path = strings.TrimPrefix(path, "/")
	return path
}

type Paging struct {
	PageIndex int `json:"pageIndex"`
	Total int `json:"total"`
}
type CompResponse struct {
	Paging     *Paging `json:"paging,omitempty"`
	Components []struct {
		Id           string `json:"id"`
		Key          string `json:"key"`
		Organization string `json:"organization"`
		Name         string `json:"name"`
		Qualifier    string `json:"qualifier"`
	}
}
