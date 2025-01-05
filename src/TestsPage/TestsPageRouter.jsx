import { Route } from "react-router-dom";

import TestGroupsPage from "./TestGroupsPage";
import TestGroupPage from "./TestGroupPage";
import SendSolutionPage from "./SendSolutionPage";

export default function TestPageRouter() {
    return <>
        <Route index element={<TestGroupsPage />} />
        <Route path=":id/send" element={<SendSolutionPage />} />
        <Route path=":id" element={<TestGroupPage />} />
    </>
}