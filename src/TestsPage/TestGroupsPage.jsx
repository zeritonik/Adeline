import { useEffect, useState } from "react"
import { Link } from "react-router-dom"

import TestGroupCreationForm from "../forms/TestGroupCreationForm";
import { WidgetWithState, NoneState, LoadingState, SuccessState, ErrorState } from "../WidgetWithState"
import { getTestGroups } from "../api/tests"


export default function TestGroupsPage() {
    const [test_groups, setTestGroups] = useState([])
    const [new_displayed, setNewDisplayed] = useState(false)

    const [state, setState] = useState(NoneState)
    
    useEffect(() => {
        const func = async () => {
            setState(LoadingState)
            try {
                const json_data = await getTestGroups()
                setTestGroups(json_data)
                setState(SuccessState)
            } catch (error) {
                setState(ErrorState)
            }
        }
        func()
    }, [])
    
    return (
        <section className="section">
            <div className="section-popup" style={{display: new_displayed ? "block" : "none"}}>
                <button className="section-popup__close btn btn-warning" onClick={() => setNewDisplayed(false)} />
                <TestGroupCreationForm 
                    onSuccess={(json_data) => { 
                        setTestGroups([json_data,...test_groups])
                        setNewDisplayed(false)
                        setState(SuccessState)
                    }} 
                    onError={() => setState(ErrorState)}
                />
            </div>

            <WidgetWithState state={state}>
                <h2 className="section__title">Your test-groups:</h2>
                <div className="card-group" style={{gridTemplateColumns: "repeat(5, 1fr)", gap: "1.5rem"}}>
                    <a onClick={ () => { setNewDisplayed(true); setState(NoneState) } } className="card clickable" style={{justifyContent: "center"}}>
                        <h3 className="card__title">Create new test group</h3>
                    </a>
                    {test_groups.map((test_group) => {
                        return (
                            <Link to={`/profile/tests/${test_group.id}`} key={test_group.id} className="card clickable">
                                <h3 className="card__title">{test_group.name}</h3>
                                <div className="card__content">
                                    <p>time limit: {test_group.time_limit}ms</p>
                                    <p>memory limit: {test_group.memory_limit}mb</p>
                                    <p>tests count: {test_group.quantity_tests}</p>
                                </div>
                            </Link>
                        )
                    })}
                </div>
            </WidgetWithState>
        </section>
    )
}