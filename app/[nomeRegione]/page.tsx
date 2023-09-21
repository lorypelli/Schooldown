"use client"
import React from "react"
import "../styles/globals.css"
import Navbar from "../components/Navbar"
export default async function Home({ params }: { params: { nomeRegione: string } }) {
    const res = await fetch(`https://schooldown.vercel.app/api/${params.nomeRegione}/`)
    let isValidRegion = res.status != 400 ? true : false
    return (
        <>
            <Navbar />
            <h1>{isValidRegion ? params.nomeRegione : await res.text()}</h1>
        </>
    )
}