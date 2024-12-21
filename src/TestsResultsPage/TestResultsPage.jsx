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

        // setState(LoadingState)
        // setTimeout(() => { setTestResults([
        //         { id: 1, max_time: 100, max_memory: 100, tests_results: ['OK', 'OK', 'OK', 'WA', 'TL', 'OK', 'CE'] },
        //         { id: 2, max_time: 100, max_memory: 100, tests_results: ['OK', 'OK', 'OK', 'WA', 'TL', 'OK', 'CE'] },
        //         { id: 3, max_time: 100, max_memory: 100, tests_results: ['OK', 'OK', 'OK', 'WA', 'TL', 'OK', 'CE'] },
        //         { id: 4, max_time: 100, max_memory: 100, tests_results: ['OK', 'OK', 'OK', 'WA', 'TL', 'OK', 'CE'] },
        //         { id: 5, max_time: 100, max_memory: 100, tests_results: ['OK', 'OK', 'OK', 'WA', 'TL', 'OK', 'CE'] },
        //     ])
        //     setState(SuccessState)
        // }, 1000)
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
                                <p>Max time: {test_result.max_time}</p>
                                <p>Max memory: {test_result.max_memory}</p>
                                <div className="carousel">
                                    <ul>
                                        {test_result.tests_results.map((verdict, i) => (
                                            <li className={verdict === "OK" ? "box-ok" : "box-warning"} style={{width: "2.5rem", height: "2.5rem"}} key={i}>
                                                {verdict}
                                            </li>
                                        ))}
                                    </ul>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            </WidgetWithState>
        </section>
    )
}