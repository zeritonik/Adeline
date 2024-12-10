import { Route } from "react-router-dom";

import ProfilePage from "./ProfilePage";
import ProfileSettingsPage from "./ProfileSettingsPage";

export default function ProfilePageRouter() {
    return <>
        <Route index element={<ProfileSettingsPage />} />
    </>
}