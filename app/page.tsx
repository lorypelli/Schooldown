import React from "react"
import axios from "axios"
import "./styles/globals.css"
import Navbar from "./components/Navbar"
import { redirect } from "next/navigation"
export default async function Home() {
    const res = await axios.get(`https://schooldown.vercel.app/api/getData/`)
    const nomeRegione = Object.keys(res.data)
    let indexRegione = 1
    let randomIndex = Math.floor(Math.random() * (21 - 1 + 1) + 1)
    let randomRegion = ""
    nomeRegione.forEach(region => {
        indexRegione++
        if (indexRegione == randomIndex) {
            randomRegion = region
            redirect(`/${region}`)
        }
    })
    return (
        <>
            <Navbar />
            <h1>{randomRegion.replaceAll("%20", " ")}</h1>
        </>
    )
}