import { useState, useContext } from "react";

import { ErrorsGroup, Form, FormGroup } from "./Form";
import { loginUser } from "../api/base";

import { UserContext } from "../App";


export default function LoginForm() {
    const [login, setLogin] = useState('');
    const [password, setPassword] = useState('');
    const [form_errors, setFormErrors] = useState(null);

    const [ setUser ] = useContext(UserContext)

    async function handleSubmit(e) {
        e.preventDefault();
        
        try {
            const json_data = await loginUser(login, password)
            console.log(json_data)
            setUser(json_data.login) // set user in context
        } catch (error) {
            setFormErrors(["Login error: " + error])
        }
    }

    return (
        <Form errors={form_errors} onSubmit={handleSubmit}>
            <h2 className="form__title">Login</h2>
            <ErrorsGroup errors={form_errors} />
            <FormGroup>
                <label className="label" htmlFor="login">Login</label>
                <input type="text" id="login" className="input" value={login} onChange={e => setLogin(e.target.value)} />
            </FormGroup>
            <FormGroup>
                <label className="label" htmlFor="password">Password</label>
                <input type="password" id="password" className="input" value={password} onChange={e => setPassword(e.target.value)} />
            </FormGroup>
            <button type="submit" className="btn">Login</button>
        </Form>
    )
}