import { useEffect, useState } from "react"

import { WidgetWithState, NoneState, LoadingState, SuccessState, ErrorState } from "../WidgetWithState"
import { getTestResults } from "../api/tests"
import { Link } from "react-router-dom"


export default function TestResultsPage() {
    const [test_results, setTestResults] = useState([])
    const [state, setState] = useState(NoneState)
    
    useEffect(() => {
        const func = async () => {
            setState(LoadingState)
            try {
                const json_data = await getTestResults()
                setTestResults(json_data)
                setState(SuccessState)
            } catch (error) {
                setState(ErrorState)
            }
        }
        func()
    }, [])
    
    return (
        <section className="section">
            <WidgetWithState state={state}>
                <h2 className="section__title">Your results:</h2>
                <div className="card-group" style={{ gridTemplateColumns: "repeat(5, 17rem)", gap: "1.5rem", justifyContent: "center"}}>
                    {test_results.map((test_result) => (
                        <div className="card" key={test_result.id}>
                            <h3 className="card__title">Test result {test_result.id}</h3>
                            <div className="card__content">
                                <p>Language: {test_result.language}.</p>
                                <p>Max time: {test_result.max_execution_time}.</p>
                                <p>Max memory: {test_result.max_memory}.</p>
                                <p className={test_result.verdict === "OK" ? "box-ok" : "box-warning"}>{test_result.verdict}</p>
                                <div className="carousel">
                                    <ul>
                                        {test_result.verdicts.map(({ id, verdict}) => (
                                            <li className={verdict === "OK" ? "box-ok" : "box-warning"} style={{width: "2.5rem", height: "2.5rem"}} key={id}>
                                                {verdict}
                                            </li>
                                        ))}
                                    </ul>
                                </div>
                                <Link to={`/profile/tests/${test_result.group_id}`}>View test group</Link>
                            </div>
                        </div>
                    ))}
                </div>
            </WidgetWithState>
        </section>
    )
}