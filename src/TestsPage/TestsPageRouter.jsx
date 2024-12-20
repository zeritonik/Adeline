import { Route } from "react-router-dom";

import TestGroupsPage from "./TestGroupsPage";
import TestGroupPage from "./TestGroupPage";

export default function TestPageRouter() {
    return <>
        <Route index element={<TestGroupsPage />} />
        <Route path=":id" element={<TestGroupPage />} />
    </>
}