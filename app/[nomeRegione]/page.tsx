import React from "react"
import "../styles/globals.css"
import Navbar from "../components/Navbar"
import { redirect } from "next/navigation"
import { Region } from "../../addons"
export default async function Home({ params }: { params: { nomeRegione: string } }) {
    let res: Region | Response = await fetch(`https://schooldown.vercel.app/api/${params.nomeRegione.replaceAll("-", " ")}/`)
    let isValidRegion = res.status != 400 ? true : false
    let region = isValidRegion ? params.nomeRegione : await res.text()
    if (!isValidRegion) {
        redirect(`/${region.replaceAll(" ", "-")}`)
    }
    else {
        res = await res.json() as Region
    }
    const countdownInizio = Math.floor(res.inizioLezioni - Date.now() / 1000)
    const countdownFine = Math.floor(res.fineLezioni - Date.now() / 1000)
    let mesi = 0
    let settimane = 0
    let giorni = 0
    let ore = 0
    let minuti = 0
    let secondi = 0
    let restoMesi = 0
    let restoSettimane = 0
    let restoGiorni = 0
    let restoOre = 0
    let restoMinuti = 0
    if (countdownInizio < 0) {
        mesi = Math.floor(countdownFine / (2.628 * 10**6))
        restoMesi = Math.floor(countdownFine % (2.628 * 10**6))
        settimane = Math.floor(restoMesi / (6.048 * 10**5))
        restoSettimane = Math.floor(restoMesi % (6.048 * 10**5))
        giorni = Math.floor(restoSettimane / (8.64 * 10**4))
        restoGiorni = Math.floor(restoSettimane % (8.64 * 10**4))
        ore = Math.floor(restoGiorni / (3.6 * 10**3))
        restoOre = Math.floor(restoGiorni % (3.6 * 10**3))
        minuti = Math.floor(restoOre / (6.0 * 10))
        restoMinuti = Math.floor(restoOre % (6.0 * 10))
        secondi = Math.floor(restoMinuti)
    }
    else {
        mesi = Math.floor(countdownInizio / (2.628 * 10**6))
        restoMesi = Math.floor(countdownInizio % (2.628 * 10**6))
        settimane = Math.floor(restoMesi / (6.048 * 10**5))
        restoSettimane = Math.floor(restoMesi % (6.048 * 10**5))
        giorni = Math.floor(restoSettimane / (8.64 * 10**4))
        restoGiorni = Math.floor(restoSettimane % (8.64 * 10**4))
        ore = Math.floor(restoGiorni / (3.6 * 10**3))
        restoOre = Math.floor(restoGiorni % (3.6 * 10**3))
        minuti = Math.floor(restoOre / (6.0 * 10))
        restoMinuti = Math.floor(restoOre % (6.0 * 10))
        secondi = Math.floor(restoMinuti)
    }
    return (
        <>
            <Navbar />
            <h1>{region.replaceAll("-", " ")}</h1>
            <h1>{countdownInizio < 0 ? "La scuola finisce tra:" : "La scuola inizia tra:"}</h1>
            <h1>{`${mesi} mesi, ${settimane} settimane, ${giorni} giorni, ${ore} ore ${minuti} minuti ${secondi} secondi`}</h1>
        </>
    )
}