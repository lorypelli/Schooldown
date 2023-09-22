import React from "react"
import "../styles/globals.css"
import Navbar from "../components/Navbar"
import { redirect } from "next/navigation"
import { Region } from "../../addons"
export default async function Home({ params }: { params: { nomeRegione: string } }) {
    let res: Region | Response = await fetch(`https://schooldown.vercel.app/api/${params.nomeRegione}/`)
    let isValidRegion = res.status != 400 ? true : false
    let region = isValidRegion ? params.nomeRegione : await res.text()
    if (!isValidRegion) {
        redirect(`/${region}`)
    }
    else {
        res = await res.json() as Region
    }
    const countdownInizio = Math.floor(res.inizioLezioni - Date.now() / 1000)
    const countdownFine = Math.floor(res.fineLezioni - Date.now() / 1000)
    return (
        <>
            <Navbar />
            <h1>{region.replaceAll("%20", " ")}</h1>
            <h1>{countdownInizio < 0 ? "La scuola finisce tra:" : "La scuola inizia tra:"}</h1>
            <h1>{countdownInizio > 0 ? countdownInizio : countdownFine}</h1>
        </>
    )
}