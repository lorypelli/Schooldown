import React from "react"
import "../styles/globals.css"
import Navbar from "../components/Navbar"
import { redirect } from "next/navigation"
export default async function Home({ params }: { params: { nomeRegione: string } }) {
    const res = await fetch(`https://schooldown.vercel.app/api/${params.nomeRegione}/`)
    let isValidRegion = res.status != 400 ? true : false
    let region = isValidRegion ? params.nomeRegione : await res.text()
    if (!isValidRegion) {
        redirect(`/${region}`)
    }
    return (
        <>
            <Navbar />
            <h1>{region.replaceAll("%20", " ")}</h1>
        </>
    )
}