import { useState, useEffect } from "react";
import { useParams } from "react-router-dom";

import { WidgetWithState, NoneState, LoadingState, SuccessState, ErrorState } from "../WidgetWithState"

export default function TestGroupPage() {
    const { id } = useParams();
    const [ testGroup, setTestGroup ] = useState({ name: `Test group ${id}`, tests: [] })
    const [ state, setState ] = useState(NoneState)

    useEffect(() => {
        const func = async () => {
            setState(LoadingState)
            try {
                // const response = await getGroupTest(id)
                // if (response.status === 200) {
                    
                // }
                // setState(SuccessState)
                setTimeout(
                    () => {
                        setTestGroup({
                            id: id,
                            name: `Test group ${id}`,
                            tests: Array(10).fill().map(
                                (_, i) => { return {
                                    id: i,
                                }}
                            )
                        })
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