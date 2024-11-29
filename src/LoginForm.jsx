import { useState } from "react";

function LoginForm({ style }) {
    const [login, setLogin] = useState('');
    const [password, setPassword] = useState('');

    const [error, setError] = useState(null);

    function handleSubmit(e) {
        e.preventDefault();
        
        setLogin('');
        setPassword('');
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
            <button className={error ? 'btn-error' : 'btn-ok'} type="submit">Войти</button>
        </form>
    )
}

export default LoginForm