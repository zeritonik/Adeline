import "./Header.css"

import { useState } from "react"
import { Link } from "react-router-dom"

import Logo from "../logo/Logo"
import HeaderProfile from "./HeaderProfile"

export function MainHeader() {
    return (
        <header className="header">
            <ul className="header__item header__navigation">
                <li><Link to="/#About">About</Link></li>
                <li><Link to="/#Examples">Examples</Link></li>
                <li><Link to="/#Registration">Registration</Link></li>
            </ul>
            <div className="header__item">
                <Logo />
            </div>
            <div className="header__item" style={{display: "flex", justifyContent: "flex-end"}}>
                <HeaderProfile />
            </div>
        </header>
    )
}

export function ProfileHeader() {
    const [ sarchValue, setSarchValue ] = useState("")

    async function handleSearch(e) {
        e.preventDefault();
        console.log(sarchValue)
    }

    return (
        <header className="header">
        <ul className="header__item header__navigation">
            <li><Link to="/profile/tests">Your test groups</Link></li>
            <li><Link to="/profile/results">Your results</Link></li>
        </ul>
        <form action="" className="header__item" style={{display: "flex", justifyContent: "center", gap: "0.5rem"}} onSubmit={handleSearch}>
            <input className="input" type="text" placeholder="Search" value={sarchValue} onChange={e => setSarchValue(e.target.value)} />
            <button className="btn" type="submit">Search</button>
        </form>
        <div className="header__item" style={{display: "flex", justifyContent: "flex-end"}}>
            <HeaderProfile />
        </div>
    </header>
    )
}