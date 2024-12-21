import { useEffect, useState } from "react"

import { WidgetWithState, NoneState, LoadingState, SuccessState, ErrorState } from "../WidgetWithState"
import { getTestResults } from "../api/tests"


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
                <div className="card-group" style={{ gridTemplateColumns: "repeat(4, 1fr)", gap: "1rem"}}>
                    {test_results.map((test_result) => (
                        <div className="card" key={test_result.id}>
                            <h3 className="card__title">{test_result.name}</h3>
                            <div className="card__content">
                                <p>Max time: {test_result.max_time}</p>
                                <p>Max memory: {test_result.max_memory}</p>
                                <ul className="carousel">
                                    {test_result.tests_results.map((test) => (
                                        <li className={test.verdict == "OK" ? "bg-ok" : "bg-warning"} style={{width: "1rem", height: "1rem"}} key={test.id}>
                                            {test.verdict}
                                        </li>
                                    ))}
                                </ul>
                            </div>
                        </div>
                    ))}
                </div>
            </WidgetWithState>
        </section>
    )
}