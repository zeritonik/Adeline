import { useState, useEffect } from "react";
import { useParams } from "react-router-dom";

import { WidgetWithState, NoneState, LoadingState, SuccessState, ErrorState } from "../WidgetWithState"
import { getTestGroup } from "../api/tests";
import { useAnouth } from "../api/useAnouth";

export default function TestGroupPage() {
    useAnouth()

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


    async function onDeleteTest(id) {
        setTestGroup({
            ...testGroup, 
            tests: testGroup.tests.filter((test) => test.id !== id)
        })
    }

    
    return (
        <section className="section" style={{ width: "66%" }}>
            <WidgetWithState state={state}>
                <h2 className="section__title">{testGroup.name}</h2>
                <div className="card-group" style={{ display: "flex", flexDirection: "column", gap: "1rem"}}>
                    {testGroup.tests.map((test) => <Test key={test.id} id={test.id} onDelete={() => onDeleteTest(test.id)} />)}
                </div>
            </WidgetWithState>
        </section>
    )
}


function Test({ id, input, output, onDelete }) {
    const [ open, setOpen ] = useState(false);

    return (
        <div className="card">
            <div style={{display: "flex", justifyContent: "space-between"}}>
                <button className="btn" onClick={ () => setOpen(!open) }>{open ? "Close" : "Open"}</button>
                <h3 className="card__title">Test {id}</h3>
                <button className="btn btn-warning" onClick={ onDelete }>Delete</button>
            </div>
            <div className="card__content" style={{display: open ? "flex" : "none"}}>
                <textarea className="input" value={input} style={{height: "10rem"}} />
                <textarea className="input" value={output} style={{height: "10rem"}} />
            </div>
        </div>
    )
}