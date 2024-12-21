import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";

import { ErrorsGroup, Form, FormGroup} from './Form'

import { registerUser } from '../api/base'


function validateForm(login, password, repeatPassword) {
    if (login === "" || password === "" || repeatPassword === "") {
        return null
    }
    let errors = []
    if (password !== repeatPassword) {
        errors.push("Passwords don't match")
    }
    return errors
}

function validateLogin(login) {
    if (login === "") {
        return null
    }
    let errors = []
    if (login.length < 3) {
        errors.push("Login is too short")
    } else if (login.length > 16) {
        errors.push("Login is too long")
    }
    return errors
}

function validateNickName(nickname) {
    if (nickname === "") {
        return null
    }
    let errors = []
    if (nickname.length < 3) {
        errors.push("Nick name is too short")
    } else if (nickname.length > 16) {
        errors.push("Nick name is too long")
    }
    return errors
}

function validatePassword(password) {
    if (password === "") {
        return null
    }
    let errors = []
    if (password.length < 6) {
        errors.push("Password is too short. Minimum 6 symbols")
    }
    return errors
}


export default function RegistrationForm() {
    const [login, setLogin] = useState('');
    const [login_errors, setLoginErrors] = useState(null);
    useEffect(() => {setLoginErrors(validateLogin(login))}, [login])

    const [password, setPassword] = useState('');
    const [password_errors, setPasswordErrors] = useState(null);
    useEffect(() => {setPasswordErrors(validatePassword(password))}, [password])
    const [repeatPassword, setRepeatPassword] = useState('');

    const [nickname, setNickName] = useState('');
    const [nickname_errors, setNicknameErrors] = useState(null);
    useEffect(() => {setNicknameErrors(validateNickName(nickname))}, [nickname])

    const [form_errors, setFormErrors] = useState(null);
    useEffect(() => {setFormErrors(validateForm(login, password, repeatPassword))}, [login, password, repeatPassword])

    const navigate = useNavigate();


    async function handleSubmit(e) {
        e.preventDefault();
        if (!form_errors || nickname_errors.length || login_errors.length || password_errors.length || form_errors.length) {
            return
        }
        try {
            await registerUser(login, password)
            navigate("/login")
        } catch (error) {
            setFormErrors(["Registration error " + error])
        }
    }

    return <Form errors={form_errors} onSubmit={handleSubmit}>
        <h2 className="form__title">Registration</h2>
        <ErrorsGroup errors={form_errors} />
        <FormGroup errors={login_errors}>
            <label htmlFor="login" className="label">Login:</label>
            <ErrorsGroup errors={login_errors} />
            <input id="login" type="text" className="input" name="login" placeholder="login" value={login} onChange={e => setLogin(e.target.value)} required />
        </FormGroup>
        <FormGroup errors={nickname_errors}>
            <label htmlFor="nickname" className="label">Nick name:</label>
            <ErrorsGroup errors={nickname_errors} />
            <input id="nickname" type="text" className="input" name="nickname" placeholder="nickname" value={nickname} onChange={e => setNickName(e.target.value)} required />
        </FormGroup>
        <FormGroup errors={password_errors}>
            <label htmlFor="password" className="label">Password:</label>
            <ErrorsGroup errors={password_errors} />
            <input id="password" type="password" className="input" name="password" placeholder="password" value={password} onChange={e => setPassword(e.target.value)} required />
        </FormGroup>
        <FormGroup>
            <label htmlFor="repeatPassword" className="label">Repeat password:</label>
            <input id="repeatPassword" type="password" className="input" name="repeatPassword" placeholder="repeat password" value={repeatPassword} onChange={e => setRepeatPassword(e.target.value)} required />
        </FormGroup>
        <button className="btn" type="submit">Register</button>
    </Form>
}