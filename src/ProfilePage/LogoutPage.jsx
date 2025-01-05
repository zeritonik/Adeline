import { useContext } from "react";

import { logoutUser } from "../api/base";
import { useNext } from "../api/useNext";
import { UserContext } from "../App";

export default function LogoutPage() {
    const [ user, setUser ] = useContext(UserContext);
    const [ next, navigate ] = useNext();

    async function handleLogout() {
        await logoutUser()
        setUser(null)
        navigate()
    }

    return (
        <section className="section">
            <h2 className="section__title">Logout</h2>
            <div style={{display: "flex", gap: "1rem"}}>
                <button className="btn btn-primary" onClick={handleLogout}>Logout</button>
                <button className="btn btn-warning" onClick={navigate}>Cancel</button>
            </div>
        </section>
    );
}