import axios from "axios"
import { redirect } from "next/navigation"
export default async function Home() {
    const res = await axios.get(`https://schooldown.vercel.app/api/getData/`)
    const nomeRegione = Object.keys(res.data)
    let indexRegione = 1
    let randomIndex = Math.floor(Math.random() * (21 - 1 + 1) + 1)
    nomeRegione.forEach(region => {
        indexRegione++
        if (indexRegione == randomIndex) {
            redirect(`/${region}`)
        }
    })
}