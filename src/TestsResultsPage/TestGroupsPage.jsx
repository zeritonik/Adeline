import { useEffect, useState } from "react"
import { Link } from "react-router-dom"

import { WidgetWithState, NoneState, LoadingState, SuccessState, ErrorState } from "../WidgetWithState"
import { getGroupTests } from "../api/tests"


export default function TestGroupsPage() {
    const [test_groups, setTestGroups] = useState([])

    const [state, setState] = useState(NoneState)
    
    useEffect(() => {
        const func = async () => {
            setState(LoadingState)
            try {
                // const response = await getGroupTests(0, 10)
                // if (response.status === 200) {
                    
                // }
                setTimeout(
                    () => {
                        setTestGroups(
                            Array(10).fill().map(
                                (_, i) => { return {
                                    id: i,
                                    name: `Test group ${i}`,
                                    time_limit: 100 * i,
                                    memory_limit: 100 * i,
                                    tests_count: i + 1
                                }}
                            )
                        )
                        setState(SuccessState)
                    },
                    1000
                )
            } catch (error) {
                setState(ErrorState)
            }
        }
        func()
    }, [])
    
    return (
        <WidgetWithState state={state}>
            <section className="section">
                <h2 className="section__title">Your test-groups:</h2>
                <div className="card-group" style={{gridTemplateColumns: "repeat(5, 1fr)", gap: "1.5rem"}}>
                    {test_groups.map((test_group) => {
                        return (
                            <Link to={`/profile/tests/${test_group.id}`} key={test_group.id} className="card clickable">
                                <h3 className="card__title">{test_group.name}</h3>
                                <div className="card__content">
                                    <p>time limit: {test_group.time_limit}</p>
                                    <p>memory limit: {test_group.memory_limit}</p>
                                    <p>tests count: {test_group.tests_count}</p>
                                </div>
                            </Link>
                        )
                    })}
                </div>
            </section>
        </WidgetWithState>
    )
}