package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/", func(res *fiber.Ctx) error {
		response, err := http.Get("http://localhost:8080/api/getData")
		if err != nil {
			return res.Status(500).SendString(err.Error())
		}
		body, err := io.ReadAll(response.Body)
		data := make(map[string]struct {
			InizioLezioni int64
			FineLezioni   int64
		})
		if err != nil {
			return res.Status(500).SendString(err.Error())
		} else {
			json.Unmarshal(body, &data)
			nomeRegione := make([]string, 0, len(data))
			randomIndex := rand.Intn(21)
			for k := range data {
				nomeRegione = append(nomeRegione, k)
			}
			return res.Redirect(fmt.Sprintf("/%s", nomeRegione[randomIndex]))
		}
	})
	app.Get("/:nomeRegione", func(res *fiber.Ctx) error {
		nomeRegione, err := url.QueryUnescape(res.Params("nomeRegione"))
		if err != nil {
			return res.Status(500).SendString(err.Error())
		}
		response, err := http.Get(fmt.Sprintf("http://localhost:8080/api/%s", nomeRegione))
		if err != nil {
			return res.Status(500).SendString(err.Error())
		}
		var isValidRegion bool
		if response.StatusCode != 400 {
			isValidRegion = true
		} else {
			isValidRegion = false
		}
		var region string
		data := struct {
			InizioLezioni int64
			FineLezioni   int64
		}{}
		if response.StatusCode == 404 {
			region = "Invalid"
		} else {
			if isValidRegion {
				region = nomeRegione
				body, err := io.ReadAll(response.Body)
				if err != nil {
					return res.Status(500).SendString(err.Error())
				} else {
					json.Unmarshal(body, &data)
				}
			} else {
				body, err := io.ReadAll(response.Body)
				if err != nil {
					return res.Status(500).SendString(err.Error())
				} else {
					region = string(body)
					return res.Redirect(fmt.Sprintf("/%s", region))
				}
			}
		}
		countdownInizio := int(data.InizioLezioni - time.Now().Unix())
		countdownFine := int(data.FineLezioni - time.Now().Unix())
		var mesi int
		var settimane int
		var giorni int
		var ore int
		var minuti int
		var secondi int
		var restoMesi int
		var restoSettimane int
		var restoGiorni int
		var restoOre int
		var restoMinuti int
		if int(countdownInizio) < 0 {
			mesi = int(math.Floor(float64(countdownFine) / (2.628 * math.Pow10(6))))
			restoMesi = countdownFine % int(2.628*math.Pow10(6))
			settimane = int(math.Floor(float64(restoMesi) / (6.048 * math.Pow10(5))))
			restoSettimane = restoMesi % int(6.048*math.Pow10(5))
			giorni = int(math.Floor(float64(restoSettimane) / (8.64 * math.Pow10(4))))
			restoGiorni = restoSettimane % int(8.64*math.Pow10(4))
			ore = int(math.Floor(float64(restoGiorni) / (3.6 * math.Pow10(3))))
			restoOre = restoGiorni % int(3.6*math.Pow10(3))
			minuti = int(math.Floor(float64(restoOre) / (6.0 * 10)))
			restoMinuti = restoOre % (6.0 * 10)
			secondi = restoMinuti
		} else {
			mesi = int(math.Floor(float64(countdownInizio) / (2.628 * math.Pow10(6))))
			restoMesi = countdownInizio % int(2.628*math.Pow10(6))
			settimane = int(math.Floor(float64(restoMesi) / (6.048 * math.Pow10(5))))
			restoSettimane = restoMesi % int(6.048*math.Pow10(5))
			giorni = int(math.Floor(float64(restoSettimane) / (8.64 * math.Pow10(4))))
			restoGiorni = restoSettimane % int(8.64*math.Pow10(4))
			ore = int(math.Floor(float64(restoGiorni) / (3.6 * math.Pow10(3))))
			restoOre = restoGiorni % int(3.6*math.Pow10(3))
			minuti = int(math.Floor(float64(restoOre) / (6.0 * 10)))
			restoMinuti = restoOre % (6.0 * 10)
			secondi = restoMinuti
		}
		template.New("index.html").ParseFiles("index.html")
		return res.Render("index.html", struct {
			Region                                                         string
			CountdownInizio, Mesi, Settimane, Giorni, Ore, Minuti, Secondi int
		}{
			Region:          region,
			CountdownInizio: countdownInizio,
			Mesi:            mesi,
			Settimane:       settimane,
			Giorni:          giorni,
			Ore:             ore,
			Minuti:          minuti,
			Secondi:         secondi,
		})
	})
	app.Get("/api", func(res *fiber.Ctx) error {
		return res.Redirect("http://localhost:8080/")
	})
	app.Get("/api/getData", func(res *fiber.Ctx) error {
		response, err := http.Get(fmt.Sprintf("https://www.fanpage.it/attualita/quando-inizia-la-scuola-regione-per-regione-le-date-e-il-calendario-%d-%d/", time.Now().Year(), (time.Now().Year()%100)+1))
		if err != nil {
			return res.Status(500).SendString(err.Error())
		} else if response.StatusCode == http.StatusNotFound {
			return res.Status(404).SendString("Not Found")
		} else {
			doc, err := goquery.NewDocumentFromReader(response.Body)
			if err != nil {
				return res.Status(500).SendString(err.Error())
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
				return res.Status(200).JSON(obj)
			}
		}
	})
	app.Get("/api/:nomeRegione", func(res *fiber.Ctx) error {
		nomeRegione, err := url.QueryUnescape(res.Params("nomeRegione"))
		if err != nil {
			return res.Status(500).SendString(err.Error())
		}
		response, err := http.Get("http://localhost:8080/api/getData")
		if err != nil {
			return res.Status(500).SendString(err.Error())
		} else if response.StatusCode == http.StatusNotFound {
			return res.Status(404).SendString("Not Found")
		}
		body, err := io.ReadAll(response.Body)
		data := make(map[string]struct {
			InizioLezioni int64
			FineLezioni   int64
		})
		if err != nil {
			return res.Status(500).SendString(err.Error())
		} else {
			json.Unmarshal(body, &data)
			if data[nomeRegione].InizioLezioni > 0 {
				return res.Status(200).JSON(data[nomeRegione])
			} else {
				nomeRegione := make([]string, 0, len(data))
				randomIndex := rand.Intn(21)
				for k := range data {
					nomeRegione = append(nomeRegione, k)
				}
				return res.Status(400).SendString(nomeRegione[randomIndex])
			}
		}
	})
	app.Listen(":8080")
}
