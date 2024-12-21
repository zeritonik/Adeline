import { Route } from "react-router-dom";

import TestGroupsPage from "./TestGroupsPage";
import NewTestGroupPage from "./NewTestGroupPage";
import TestGroupPage from "./TestGroupPage";

export default function TestPageRouter() {
    return <>
        <Route index element={<TestGroupsPage />} />
        <Route path="new" element={<NewTestGroupPage />} />
        <Route path=":id" element={<TestGroupPage />} />
    </>
}