import { Route, Routes } from "react-router-dom";
import { useState, useEffect, createContext } from "react";

import { loginUser } from "./api/base";

/* pages for routing */

import MainLayout from './MainPage/MainLayout';
import MainPageRouter from './MainPage/MainPageRouter';

import ProfileLayout from './ProfilePage/ProfileLayout';
import ProfilePageRouter from './ProfilePage/ProfilePageRouter';

import TestsPageRouter from './TestsPage/TestsPageRouter';




import PageNotFound from './PageNotFound';

/* context */

export const UserContext = createContext(null);

export default function App() {
    const [user, setUser] = useState(null);

    useEffect(() => { (async () => {
        try {
            const json_data = await loginUser(null, null)
            setUser(json_data)
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
                <Route path="/profile" element={<ProfileLayout />} >
                    { ProfilePageRouter() }
                </Route>
                <Route path="/tests" element={<ProfileLayout />} >
                    { TestsPageRouter() }
                </Route>
                <Route path="*" element={<PageNotFound />} />
            </Routes>
        </UserContext.Provider>
    );
}