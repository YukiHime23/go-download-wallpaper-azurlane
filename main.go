package main

import (
	crawal_sqlite "craw-al/db"
	"craw-al/function"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	ApiListWallpaperAzurLane    = "https://azurlane.yo-star.com/api/admin/special/public-list?page_index=1&page_num=1200&type=1"
	DomainLoadWallpaperAzurLane = "https://webusstatic.yo-star.com/"
)

type ResponseApi struct {
	StatusCode int     `json:"statusCode"`
	Data       ResData `json:"data"`
}

type ResData struct {
	Count int         `json:"count"`
	Rows  []Wallpaper `json:"rows"`
}

type Wallpaper struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Cover       string `json:"cover"`
	Works       string `json:"works"`
	Type        int    `json:"type"`
	Sort        int    `json:"sort_index"`
	PublishTime int    `json:"publish_time"`
	New         bool   `json:"new"`
}

type AzurLane struct {
	FileName    string `json:"file_name"`
	IdWallpaper int    `json:"id_wallpaper"`
	Url         string `json:"url"`
}

func main() {

	res, err := http.Get(ApiListWallpaperAzurLane)
	if err != nil {
		log.Fatal("call api error: ", err)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("read body error: ", err)
	}

	var resApi ResponseApi
	if err = json.Unmarshal(resBody, &resApi); err != nil {
		log.Fatal("json Unmarshal error: ", err)
	}

	pathFile := "AzurLane"
	if err = os.MkdirAll(pathFile, os.ModePerm); err != nil {
		log.Fatal("mkdir file error: ", err)
	}

	var idExist []int

	db := crawal_sqlite.GetDb()
	// get id exist
	err = db.QueryRow("SELECT id_wallpaper FROM azur_lane").Scan(&idExist)

	for _, row := range resApi.Data.Rows {
		if function.IntInArray(idExist, row.ID) {
			continue
		}

		var al AzurLane
		al.Url = DomainLoadWallpaperAzurLane + row.Works
		al.FileName = strings.ReplaceAll(row.Title+" ("+row.Artist+").jpeg", "/", "-")
		al.IdWallpaper = row.ID
		if err = function.DownloadFile(al.Url, al.FileName, pathFile); err != nil {
			log.Fatal("download file error: ", err)
		}
		insertData := "INSERT INTO azur_lane VALUES (?, ?, ?)"
		_, err = db.Exec(insertData, al.Url, al.FileName, pathFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	defer db.Close()
}
