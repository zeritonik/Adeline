import { MainHeader } from "../header/Header";
import { Outlet } from "react-router-dom";

export default function MainLayout() {
    return (
        <>
            <MainHeader />
            <Outlet />
        </>
    )
}