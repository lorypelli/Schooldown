package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	app := fiber.New()
	app.Get("/", func(res *fiber.Ctx) error {
		return res.Redirect("http://localhost:4321/")
	})
	app.Get("/getData", func(res *fiber.Ctx) error {
		response, err := http.Get(fmt.Sprintf("https://www.fanpage.it/attualita/quando-inizia-la-scuola-regione-per-regione-le-date-e-il-calendario-%d-%d/", time.Now().Year(), (time.Now().Year()%100)+1))
		if err != nil {
			res.Status(500)
			return res.SendString(err.Error())
		} else if response.StatusCode == http.StatusNotFound {
			res.Status(404)
			return res.SendString("Not Found")
		} else {
			doc, err := goquery.NewDocumentFromReader(response.Body)
			if err != nil {
				res.Status(500)
				return res.SendString(err.Error())
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
				res.Status(200)
				return res.JSON(obj)
			}
		}
	})
	app.Get("/:nomeRegione", func(res *fiber.Ctx) error {
		nomeRegione := res.Params("nomeRegione")
		response, err := http.Get("http://localhost:8080/getData")
		if err != nil {
			res.Status(500)
			return res.SendString(err.Error())
		} else if response.StatusCode == http.StatusNotFound {
			res.Status(404)
			return res.SendString("Not Found")
		}
		body, err := io.ReadAll(response.Body)
		data := make(map[string]struct {
			InizioLezioni int64
			FineLezioni   int64
		})
		if err != nil {
			res.Status(500)
			return res.SendString(err.Error())
		} else {
			json.Unmarshal(body, &data)
			if data[nomeRegione].InizioLezioni > 0 {
				res.Status(200)
				return res.JSON(data[nomeRegione])
			} else {
				nomeRegione := make([]string, 0, len(data))
				randomIndex := rand.Intn(21)
				for k := range data {
					nomeRegione = append(nomeRegione, k)
				}
				res.Status(400)
				return res.SendString(nomeRegione[randomIndex])

			}
		}
	})
	app.Listen(":8080")
}
