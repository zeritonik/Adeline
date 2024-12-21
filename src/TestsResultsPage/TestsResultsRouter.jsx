import { Route } from "react-router-dom";

import TestResultsPage from "./TestResultsPage";

export default function TestPageRouter() {
    return <>
        <Route index element={<TestResultsPage />} />
        {/* <Route path=":id" element={<TestGroupPage />} /> */}
    </>
}