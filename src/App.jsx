import { Route, Routes, Navigate, useHref } from "react-router-dom";
import { useState, useEffect, createContext } from "react";

import { loginUser } from "./api/base";

/* pages for routing */

import MainLayout from './MainPage/MainLayout';
import MainPageRouter from './MainPage/MainPageRouter';

import ProfileLayout from './ProfilePage/ProfileLayout';
import ProfilePageRouter from './ProfilePage/ProfilePageRouter';

import TestsPageRouter from './TestsPage/TestsPageRouter';

import TestsResultsRouter from './TestsResultsPage/TestsResultsRouter';

import PageNotFound from './PageNotFound';

/* context */

export const UserContext = createContext(null);

export default function App() {
    const [user, setUser] = useState(null);
    const href = useHref();

    useEffect(() => { (async () => {
        try {
            const json_data = await loginUser(null, null)
            setUser(json_data)
        } catch ( error ) {
            return
        }
    })() } , []);

    useEffect(() => {
        console.log("БЛЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯЯ")
    }, [user])

    return (
        <UserContext.Provider value={[ user, setUser ]}>
            <Routes>
                <Route path="/" element={<MainLayout />} >
                    { MainPageRouter() }
                </Route>
                <Route path="/profile" element={!user ? <Navigate to={"/login?next=" + href} /> : <ProfileLayout />} >
                    { ProfilePageRouter() }
                    <Route path="tests">
                        { TestsPageRouter() }
                    </Route>
                    <Route path="results">
                        { TestsResultsRouter() }
                    </Route>
                </Route>
                <Route path="*" element={<PageNotFound />} />
            </Routes>
        </UserContext.Provider>
    );
}