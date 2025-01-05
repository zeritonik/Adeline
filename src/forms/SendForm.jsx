import { useState } from "react";

import { ErrorsGroup, Form, FormGroup } from "./Form";
import { sendSolution } from "../api/tests";


export function SendForm({ test_group_id }) {
    const [ language , setLanguage ] = useState('python');
    const [ source, setSource ] = useState('');

    const [ form_errors, setFormErrors ] = useState(null);

    async function handleSubmit(e) {
        e.preventDefault();
        console.log("hui")
        try {
            await sendSolution(language, source)
        } catch (error) {
            setFormErrors(["Send error: " + error])
        }
    }
    
    return (
        <Form errors={form_errors} onSubmit={handleSubmit}>
            <h2 className="form__title">Send solution</h2>
            <ErrorsGroup errors={form_errors} />
            <FormGroup>
                <label className="label" htmlFor="testGroupId">Test group id</label>
                <input type="text" id="testGroupId" className="input" value={test_group_id} disabled required />
            </FormGroup>
            <FormGroup>
                <label className="label" htmlFor="language">Language</label>
                <select id="language" className="input" value={language} onChange={e => setLanguage(e.target.value)} required>
                    <option value="python">Python</option>
                    <option value="c++">C++</option>
                </select>
            </FormGroup>
            <FormGroup>
                <label className="label" htmlFor="source">Source</label>
                <textarea type="text" id="source" className="input textarea" value={source} onChange={e => setSource(e.target.value)} required />
            </FormGroup>
            <button type="submit" className="btn">Send</button>
        </Form>
    )
}