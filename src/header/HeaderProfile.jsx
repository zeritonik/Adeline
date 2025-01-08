import { useState, useContext, useEffect } from "react";
import { Link } from "react-router-dom";

import { UserContext } from "../App";


function HeaderProfileMenu() {
    const [user] = useContext(UserContext);

    if (user) {
        return (
            <ul className="header__profile__menu">
                <li><Link to="/profile">Profile</Link></li>
                <li><Link to="/profile/settings">Settings</Link></li>
                <li><Link to="/profile/logout">Logout</Link></li>
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
    const [user] = useContext(UserContext);
    const [open, setOpen] = useState(false);
    return (
        <div className={open ? "header__profile open" : "header__profile"}
             onMouseEnter={() => setOpen(true)} onMouseLeave={() => setOpen(false)}>
            {user && user.avatar && <img key={user} className="header__profile__avatar avatar-small" src={user.avatar} alt="" />}
           <HeaderProfileMenu />
        </div>
    );
}