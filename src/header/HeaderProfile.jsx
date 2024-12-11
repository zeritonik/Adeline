import { useState, useContext } from "react";
import { Link } from "react-router-dom";

import { UserContext } from "../App";


function HeaderProfileMenu() {
    const [user] = useContext(UserContext);

    if (user) {
        return (
            <ul className="header__profile__menu">
                <li><Link to="/profile">Profile</Link></li>
                <li><Link to="/profile/settings">Settings</Link></li>
                <li><Link to="/logout">Logout</Link></li>
            </ul>
        );
    }
    return (
        <ul className="header__profile__menu">
            <li><Link to="/login">Login</Link></li>
            <li><Link to="/register">Register</Link></li>
        </ul>
    );
}

export default function HeaderProfile() {
    const [open, setOpen] = useState(false);
    return (
        <div className={open ? "header__profile open" : "header__profile"}
             onMouseEnter={() => setOpen(true)} onMouseLeave={() => setOpen(false)}>
            <div className="header__profile__avatar">
                <img src="img/user.svg" alt="" />   
            </div>
           <HeaderProfileMenu />
        </div>
    );
}