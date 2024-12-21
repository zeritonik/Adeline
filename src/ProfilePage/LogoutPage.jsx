import { useContext } from "react";

import { logoutUser } from "../api/base";
import { UserContext } from "../App";

export default function LogoutPage() {
    const [user, setUser] = useContext(UserContext);

    async function handleLogout() {
        await logoutUser()
        setUser(null)
    }

    return (
        <section className="section">
            <h2 className="section__title">Logout</h2>
            <div style={{display: "flex", gap: "1rem"}}>
                <button className="btn btn-primary" onClick={handleLogout}>Logout</button>
                <button className="btn btn-warning" onClick={() => setUser(null)}>Cancel</button>
            </div>
        </section>
    );
}