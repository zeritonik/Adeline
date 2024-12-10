import { Route } from "react-router-dom";

import MainPage from "./MainPage";
import RegistrationPage from "./RegistrationPage";
import LoginPage from "./LoginPage";

export default function MainPageRouter() {
    return <>
        <Route index element={<MainPage />} />
        <Route path='register' element={<RegistrationPage />} />
        <Route path='login' element={<LoginPage />} />
    </>
}