import { Route } from "react-router-dom";

import NewTestGroupPage from "./newTestGroupPage";
import TestGroupPage from "./TestGroupPage";

export default function TestPageRouter() {
    return <>
        <Route path="new" element={<NewTestGroupPage />} />
        <Route path=":id" element={<TestGroupPage />} />
    </>
}