package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	req := gin.Default()
	req.NoRoute(func(res *gin.Context) {
		res.Redirect(307, "http://localhost:4321/")
	})
	req.GET("/getData", func(res *gin.Context) {
		res.Header("Content-Type", "application/json; charset=UFT-8")
		response, err := http.Get(fmt.Sprintf("https://www.fanpage.it/attualita/quando-inizia-la-scuola-regione-per-regione-le-date-e-il-calendario-%d-%d/", time.Now().Year(), (time.Now().Year()%100)+1))
		if err != nil {
			res.Header("Content-Type", "text/plain")
			res.String(500, err.Error())
		} else if response.StatusCode == http.StatusNotFound {
			res.Header("Content-Type", "text/plain")
			res.String(404, "Not Found")
		} else {
			doc, err := goquery.NewDocumentFromReader(response.Body)
			if err != nil {
				res.Header("Content-Type", "text/plain")
				res.String(500, err.Error())
			} else {
				obj := make(map[string]struct {
					InizioLezioni int64
					FineLezioni   int64
				})
				doc.Find("div div div div ul li").Each(func(i int, s *goquery.Selection) {
					nomeRegione := strings.Split(s.Text(), ":")[0]
					inizioLezioni := ""
					for c := 0; c < len(strings.Split(s.Text(), ";")[0]); c++ {
						character, err := strconv.ParseFloat(string((strings.Split(s.Text(), ";")[0])[c]), 64)
						if err != nil {
							character = -1
						}
						if character >= 0 {
							inizioLezioni += string((strings.Split(s.Text(), ";")[0])[c])
						}
					}
					inizioLezioniInt, _ := strconv.Atoi(inizioLezioni)
					fineLezioni := ""
					for c := 0; c < len(strings.Split(s.Text(), ";")[1]); c++ {
						character, err := strconv.ParseFloat(string((strings.Split(s.Text(), ";")[1])[c]), 64)
						if err != nil {
							character = -1
						}
						if character >= 0 {
							fineLezioni += string((strings.Split(s.Text(), ";")[1])[c])
						}
					}
					fineLezioniInt, _ := strconv.Atoi(fineLezioni)
					obj[nomeRegione] = struct {
						InizioLezioni int64
						FineLezioni   int64
					}{
						InizioLezioni: time.Date(time.Now().Year(), time.September, inizioLezioniInt, 0, 0, 0, 0, time.UTC).Unix(),
						FineLezioni:   time.Date(time.Now().Year()+1, time.June, fineLezioniInt, 0, 0, 0, 0, time.UTC).Unix(),
					}
				})
				res.JSON(200, obj)
			}
		}
	})
	req.GET("/:nomeRegione", func(res *gin.Context) {
		res.Header("Content-Type", "application/json; charset=UFT-8")
		nomeRegione := res.Param("nomeRegione")
		response, err := http.Get("http://localhost:8080/getData")
		if err != nil {
			res.Header("Content-Type", "text/plain")
			res.String(500, err.Error())
		} else if response.StatusCode == http.StatusNotFound {
			res.Header("Content-Type", "text/plain")
			res.String(404, "Not Found")
		}
		body, err := io.ReadAll(response.Body)
		data := make(map[string]struct {
			InizioLezioni int64
			FineLezioni   int64
		})
		if err != nil {
			res.Header("Content-Type", "text/plain")
			res.String(500, err.Error())
		} else {
			json.Unmarshal(body, &data)
			if data[nomeRegione].InizioLezioni > 0 {
				res.JSON(200, data[nomeRegione])
			} else {
				nomeRegione := make([]string, 0, len(data))
				randomIndex := rand.Intn(21)
				for k := range data {
					nomeRegione = append(nomeRegione, k)
				}
				res.Header("Content-Type", "text/plain")
				res.String(400, nomeRegione[randomIndex])

			}
		}
	})
	req.Run()
}
