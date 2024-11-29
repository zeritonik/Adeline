import { useState, useEffect } from "react";

const VALIDATION_ERRORS = {
    login: 'Login is not valid',
    password: {
        length: 'Password is too short',
    },
    repeatPassword: 'Repeated password is different from password',
}

function validateLogin(login) {
    if (login.match(/^[a-zA-Z][a-zA-Z0-9_]+$/))
        return null;
    return VALIDATION_ERRORS.login;
}
function validatePassword(password) {
}

function RegistrationForm({ style }) {
    const [login, setLogin] = useState('');
    const [password, setPassword] = useState('');
    const [repeatPassword, setRepeatPassword] = useState('');

    const [error, setError] = useState(null);


    function handleSubmit(e) {
        e.preventDefault();

        setError(password !== repeatPassword ? VALIDATION_ERRORS.repeatPassword : null);
        console.log(login, password, repeatPassword);

        setLogin('');
        setPassword('');
        setRepeatPassword('');
    }

    return (
        <form className="form" style={style} onSubmit={handleSubmit}>
            <label>
                Login:
                <input type="text" name="login" placeholder="login" value={login} onChange={e => setLogin(e.target.value)} required />
            </label>
            <label>
                Password:
                <input type="password" name="password" placeholder="password" value={password} onChange={e => setPassword(e.target.value)} required />
            </label>
            <label>
                Repeat password:
                <input type="password" name="repeat-password" placeholder="repeat password" value={repeatPassword} onChange={e => setRepeatPassword(e.target.value)} required />
            </label>

            <button className={error ? 'btn-warning' : 'btn-ok'} type="submit">Зарегистрироваться</button>
        </form>
    );
}

export default RegistrationForm;