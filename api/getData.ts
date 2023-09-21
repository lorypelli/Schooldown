import axios from 'axios'
import { load } from 'cheerio'
export default async (request: import("next").NextApiRequest, response: import("next").NextApiResponse) => {
    response.setHeader("Access-Control-Allow-Origin", "http://127.0.0.1:5173")
    if (request.method !== "GET") {
        response.setHeader("Content-Type", "text/plain")
        return response.status(405).send("Method not allowed")
    }
    else if (request.method === "GET") {
        const res = await axios.get(`https://www.fanpage.it/attualita/quando-inizia-la-scuola-regione-per-regione-le-date-e-il-calendario-${new Date().getFullYear()}-${parseInt(new Date().getFullYear().toString().slice(2)) + 1}/`)
        const $ = load(res.data)
        const date = $('div div div div ul li')
        const arr: string[] = []
        for (let i = 0; i < date.length; i++) {
            if (arr[i] && arr[i].split(":")[0] == "") {
                arr[i] = arr[i - 1] + arr[i]
                arr.splice(i - 1, 1)
            }
            date[i].childNodes.forEach(c => {
                if (c.type == 'text') {
                    arr.push(c.data)
                }
                else if (c.type == 'tag') {
                    c.childNodes.forEach(c => {
                        if (c.type == 'text') {
                            arr.push(c.data)
                        }
                    })
                }
            })
        }
        let obj = {}
        for (let i = 0; i < arr.length; i++) {
            let nomeRegione = arr[i].split(":")[0]
            let inizioLezioni = ""
            for (let c = 0; c < arr[i].split(";")[0].length; c++) {
                if (!isNaN(parseInt(arr[i].split(";")[0].charAt(c)))) {
                    inizioLezioni += arr[i].split(";")[0].charAt(c)
                }
            }
            let fineLezioni = ""
            for (let c = 0; c < arr[i].split(";")[1].length; c++) {
                if (!isNaN(parseInt(arr[i].split(";")[1].charAt(c)))) {
                    fineLezioni += arr[i].split(";")[1].charAt(c)
                }
            }
            obj[nomeRegione] = {
                "inizioLezioni": Math.floor(new Date(`Sep ${inizioLezioni}, ${new Date().getFullYear()}`).getTime() / 1000),
                "fineLezioni": Math.floor(new Date(`Jun ${fineLezioni}, ${new Date().getFullYear() + 1}`).getTime() / 1000)
            }
        }
        return response.send(obj)
    }
}