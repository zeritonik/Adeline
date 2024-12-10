import { Route, Routes } from "react-router-dom";
import { useState, useEffect, createContext } from "react";

import { loginUser } from "./api/base";

/* pages for routing */

import MainLayout from './MainPage/MainLayout';
import MainPageRouter from './MainPage/MainPageRouter';




import PageNotFound from './PageNotFound';

/* context */

export const UserContext = createContext(null);

export default function App() {
    const [user, setUser] = useState(null);

    useEffect(() => { (async () => {
        try {
            const json_data = await loginUser(null, null)
            setUser(json_data.login)
        } catch ( error ) {
            return
        }
    })() } , []);

    return (
        <UserContext.Provider value={[ user, setUser ]}>
            <Routes>
                <Route path="/" element={<MainLayout />} >
                    { MainPageRouter() }
                </Route>
                <Route path="*" element={<PageNotFound />} />
            </Routes>
        </UserContext.Provider>
    );
}