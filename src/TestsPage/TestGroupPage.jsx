import { useState, useEffect } from "react";
import { Link, useParams } from "react-router-dom";

import { WidgetWithState, NoneState, LoadingState, SuccessState, ErrorState } from "../WidgetWithState"
import { getTestGroup } from "../api/tests";

export default function TestGroupPage() {
    const { id } = useParams();
    const [ testGroup, setTestGroup ] = useState({ name: `Test group ${id}`, tests: [] })
    const [ state, setState ] = useState(NoneState)

    useEffect(() => {
        setState(LoadingState)
        const func = async () => {
            setState(LoadingState)
            try {
                const json_data = await getTestGroup(id)
                setTestGroup(json_data)
                setState(SuccessState)
            } catch (error) {
                setState(ErrorState)
            }
        }
        func()
    }, [ id ])

    
    return (
        <section className="section" style={{ width: "66%" }}>
            <WidgetWithState state={state}>
                <h2 className="section__title">{testGroup.name}</h2>
                <p>Time limit: {testGroup.time_limit}ms.</p>
                <p>Memory limit: {testGroup.memory_limit}mb.</p>
                <div className="card-group" style={{ display: "flex", flexDirection: "column", gap: "1rem"}}>
                    {testGroup.tests.map((test) => <Test key={test.id} id={test.id} input={test.input} output={test.correct_output} />)}
                </div>
                <Link className="btn" to="send">Send solution!</Link>
            </WidgetWithState>
        </section>
    )
}


function Test({ id, input, output }) {
    const [ open, setOpen ] = useState(false);

    return (
        <div className="card">
            <div style={{display: "flex", justifyContent: "space-between"}}>
                <button className="btn" onClick={ () => setOpen(!open) }>{open ? "Close" : "Open"}</button>
                <h3 className="card__title">Test {id}</h3>
            </div>
            <div className="card__content" style={{display: open ? "flex" : "none"}}>
                <textarea className="input" value={input} style={{height: "10rem"}} disabled />
                <textarea className="input" value={output} style={{height: "10rem"}} disabled />
            </div>
        </div>
    )
}