import { useState } from "react";

function HeaderProfile() {
    const [open, setOpen] = useState(false);
    return (
        <div class="header__profile" onMouseEnter={() => setOpen(true)} onMouseLeave={() => setOpen(false)}>

            <img src="img/user.svg" alt="" onClick={() => setOpen(!open)} />

            <div class="dropdown" style={{display: (open ? 'block' : 'none')}}>
                <ul>
                    <li><a href="profile">Profile</a></li>
                    <li><a href="login">Log in</a></li>
                    <li><a href="registration">Register</a></li>
                </ul>
            </div>

        </div>
    );
}

export default HeaderProfile