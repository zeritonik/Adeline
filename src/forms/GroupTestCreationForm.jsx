import { useState, useEffect } from "react";
import { Form, FormGroup, ErrorsGroup } from "./Form"
import { createGroupTest } from "../api/tests";

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


export default function GroupTestCreationForm() {
    const [name, setName] = useState('');
    const [name_errors, setNameErrors] = useState(null);
    useEffect(() => {setNameErrors(validateName(name))}, [name])

    const [time_limit, setTimeLimit] = useState('');
    const [time_limit_errors, setTimeLimitErrors] = useState(null);
    useEffect(() => {setTimeLimitErrors(validateTimeLimit(time_limit))}, [time_limit])

    const [memory_limit, setMemoryLimit] = useState('');
    const [memory_limit_errors, setMemoryLimitErrors] = useState(null);
    useEffect(() => {setMemoryLimitErrors(validateMemoryLimit(memory_limit))}, [memory_limit])

    const [form_errors, setFormErrors] = useState([]);
    useEffect(() => {setFormErrors([])}, [name, time_limit, memory_limit])

    async function handleSubmit(e) {
        e.preventDefault();
        if ((name && name_errors.length) || (time_limit && time_limit_errors.length) || (memory_limit && memory_limit_errors.length)) {
            return
        }
        try {
            const json_data = await createGroupTest(name, parseInt(time_limit), parseInt(memory_limit))
            console.log(json_data) // TODO handle response
        } catch (error) {
            setFormErrors(["Creation error " + error])
        }
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
            <button type="submit" className="btn">Create group test</button>
        </Form>
    )
}