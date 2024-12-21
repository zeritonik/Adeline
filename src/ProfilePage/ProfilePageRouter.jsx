import { Route } from "react-router-dom";

import LogoutPage from "./LogoutPage";
import ProfilePage from "./ProfilePage";
import ProfileSettingsPage from "./ProfileSettingsPage";

export default function ProfilePageRouter() {
    return <>
        <Route index element={<ProfilePage />} />
        <Route path="settings" element={<ProfileSettingsPage />} />
        <Route path="logout" element={<LogoutPage />} />
    </>
}