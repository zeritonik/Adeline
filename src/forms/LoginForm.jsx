import { useState, useContext } from "react";
import { useNavigate, useSearchParams } from "react-router-dom"

import { ErrorsGroup, Form, FormGroup } from "./Form";
import { loginUser } from "../api/base";

import { UserContext } from "../App";


export default function LoginForm() {
    const [login, setLogin] = useState('');
    const [password, setPassword] = useState('');
    const [form_errors, setFormErrors] = useState(null);

    const [ user, setUser ] = useContext(UserContext)

    const [searchParams] = useSearchParams();
    const navigate = useNavigate();

    if (user) {
        navigate(searchParams.get("next") || "/profile")
    }

    async function handleSubmit(e) {
        e.preventDefault();
        
        try {
            const json_data = await loginUser(login, password)
            setUser(json_data) // set user in context
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
                <input type="text" id="login" className="input" value={login} onChange={e => setLogin(e.target.value)} required />
            </FormGroup>
            <FormGroup>
                <label className="label" htmlFor="password">Password</label>
                <input type="password" id="password" className="input" value={password} onChange={e => setPassword(e.target.value)} required />
            </FormGroup>
            <button type="submit" className="btn">Login</button>
        </Form>
    )
}