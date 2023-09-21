import axios from 'axios'
export default async (request: import("next").NextApiRequest, response: import("next").NextApiResponse) => {
    response.setHeader("Access-Control-Allow-Origin", "http://127.0.0.1:5173")
    if (request.method !== "GET") {
        response.setHeader("Content-Type", "text/plain")
        return response.status(405).send("Method not allowed")
    }
    else if (request.method === "GET") {
        const res = await axios.get(`https://schooldown.vercel.app/api/getData/`)
        if (res.data[request.query.nomeRegione as string]) {
            return response.send(res.data[request.query.nomeRegione as string])
        }
        else {
            const nomeRegione = Object.keys(res.data)
            let indexRegione = 1
            let randomIndex = Math.floor(Math.random() * (21 - 1 + 1) + 1)
            nomeRegione.forEach(region => {
                indexRegione++
                if (indexRegione == randomIndex) {
                    response.setHeader("Content-Type", "text/plain")
                    return response.status(400).send(region)
                }
            })
        }
    }
}