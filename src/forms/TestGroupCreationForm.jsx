import { useState, useEffect } from "react";
import { Form, FormGroup, ErrorsGroup } from "./Form"

import { createTestGroup } from "../api/tests";

function validateName(name) {
    if (name === "") {
        return null
    }
    let errors = []
    if (name.length < 3) {
        errors.push("Name is too short")
    }
    return errors
}

function validateTimeLimit(time_limit) {
    if (time_limit === "") {
        return null
    }
    let errors = []
    time_limit = parseInt(time_limit)
    if (isNaN(time_limit)) {
        errors.push("Time limit is not a number")
        return errors
    }
    if (time_limit <= 0) {
        errors.push("Time limit must be greater than 0")
    }
    return errors
}

function validateMemoryLimit(memory_limit) {
    if (memory_limit === "") {
        return null
    }
    let errors = []
    memory_limit = parseInt(memory_limit)
    if (isNaN(memory_limit)) {
        errors.push("Memory limit is not a number")
        return errors
    }
    if (memory_limit <= 0) {
        errors.push("Memory limit must be greater than 0")
    }
    return errors
}

function validateTests(tests) {
    let errors = []
    if (tests.length === 0) {
        errors.push("You must add at least one test")
    }
    return errors
}


export default function TestGroupCreationForm( {onSuccess=() => {}, onError=() => {}} ) {
    const [name, setName] = useState('');
    const [name_errors, setNameErrors] = useState(null);
    useEffect(() => {setNameErrors(validateName(name))}, [name])

    const [time_limit, setTimeLimit] = useState('');
    const [time_limit_errors, setTimeLimitErrors] = useState(null);
    useEffect(() => {setTimeLimitErrors(validateTimeLimit(time_limit))}, [time_limit])

    const [memory_limit, setMemoryLimit] = useState('');
    const [memory_limit_errors, setMemoryLimitErrors] = useState(null);
    useEffect(() => {setMemoryLimitErrors(validateMemoryLimit(memory_limit))}, [memory_limit])

    const [ test_id, setTestId ] = useState(0);
    const [ tests, setTests ] = useState([]);
    const [ tests_errors, setTestsErrors ] = useState(null);

    const [form_errors, setFormErrors] = useState([]);
    useEffect(() => {setFormErrors([])}, [name, time_limit, memory_limit])
    

    async function handleSubmit(e) {
        e.preventDefault();
        if ((name && name_errors.length) || (time_limit && time_limit_errors.length) || (memory_limit && memory_limit_errors.length) || (tests_errors && tests_errors.length)) {
            return
        }
        try {
            const json_data = await createTestGroup(name, parseInt(time_limit), parseInt(memory_limit), tests.map(t => {return {input: t.input, correct_output: t.output}} ))
            onSuccess(json_data)
        } catch (error) {
            setFormErrors(["Creation error " + error])
            onError()
        }
    }

    function onChangeInput(id, val) {
        setTests(tests.map(t => t.id === id ? {...t, input: val} : t))
    }

    function onChangeOutput(id, val) {
        setTests(tests.map(t => t.id === id ? {...t, output: val} : t))
    }

    function onDelete(id) {
        setTests(tests.filter(t => t.id !== id))
    }

    return (
        <Form errors={form_errors} onSubmit={handleSubmit}>
            <h2 className="form__title">Group test creation</h2>
            <ErrorsGroup errors={form_errors} />
            <FormGroup errors={name_errors}>
                <label className="label" htmlFor="name">Name</label>
                <ErrorsGroup errors={name_errors} />
                <input type="text" placeholder="name" id="name" className="input" value={name} onChange={e => setName(e.target.value)} required />
            </FormGroup>
            <FormGroup errors={time_limit_errors}>
                <label className="label" htmlFor="time_limit">Time limit</label>
                <ErrorsGroup errors={time_limit_errors} />
                <input type="text" placeholder="time limit in ms" id="time_limit" className="input" value={time_limit} onChange={e => setTimeLimit(e.target.value)} />
            </FormGroup>
            <FormGroup errors={memory_limit_errors}>
                <label className="label" htmlFor="memory_limit">Memory limit</label>
                <ErrorsGroup errors={memory_limit_errors} />
                <input type="text" placeholder="memory limit in kb" id="memory_limit" className="input" value={memory_limit} onChange={e => setMemoryLimit(e.target.value)} />
            </FormGroup>
            <FormGroup errors={null}>
                <label className="label" htmlFor="tests">Tests</label>
                <ErrorsGroup errors={null} />
                {tests.map(test => <TestForm key={test.id} id={test.id} onChangeInput={onChangeInput} onChangeOutput={onChangeOutput} onDelete={onDelete} />)}
                <button type="button" className="btn btn-ok" onClick={ () => { setTests([...tests, {id: test_id, input: "", output: ""}]); setTestId(test_id + 1) } }>Add test</button>
            </FormGroup>
            <button type="submit" className="btn">Create group test</button>
        </Form>
    )
}


function TestForm({ id, onDelete, onChangeInput, onChangeOutput }) {
    const [ open, setOpen ] = useState(false);

    return (
        <div className="card">
            <h3 className="card__title clickable" onClick={ () => setOpen(!open) }>Test {id}</h3>
            <div className="card__content" style={{display: open ? "flex" : "none"}}>
                <label className="label" htmlFor={`input-${id}`}>Input</label>
                <textarea id={`input-${id}`} className="input" onChange={e => onChangeInput(id, e.target.value)} style={{height: "15rem"}} />
                <label className="label" htmlFor={`output-${id}`}>Output</label>
                <textarea id={`output-${id}`} className="input" onChange={e => onChangeOutput(id, e.target.value)} style={{height: "15rem"}} />
            </div>
            <button className="btn btn-warning" onClick={ e => onDelete(id, e.target.value) }>Delete</button>
        </div>
    )
}