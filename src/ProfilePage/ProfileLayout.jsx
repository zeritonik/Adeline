import { Outlet } from "react-router-dom"

import { ProfileHeader } from "../header/Header"

export default function ProfileLayout() {
    return (
        <>
            <ProfileHeader />
            <Outlet />
        </>
    )
}