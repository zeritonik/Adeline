import { useState, useEffect } from "react";

import { useAnouth } from "../api/useAnouth";
import { getProfileSettings, postProfileSettings } from "../api/profile";
import { NoneState, LoadingState, SuccessState, ErrorState, WidgetWithState } from "../WidgetWithState";
import { ErrorsGroup, Form, FormGroup} from "./Form"


function validateForm(nickname, avatar) {
    let errors = []
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

function validateAvatar(avatar) {
    if (avatar === null) {
        return null
    }
    let errors = []
    return errors
}


export default function ProfileSettingsForm() {
    useAnouth()
    const [login, setLogin] = useState('')

    const [nickname, setNickName] = useState('');
    const [nickname_errors, setNicknameErrors] = useState(null);
    useEffect(() => {setNicknameErrors(validateNickName(nickname))}, [nickname])

    const [avatar, setAvatar] = useState(null);
    const [avatar_errors, setAvatarErrors] = useState(null);
    useEffect(() => {setAvatarErrors(validateAvatar(avatar))}, [avatar])
    console.log(avatar)

    const [form_errors, setFormErrors] = useState(null);
    useEffect(() => {setFormErrors(validateForm(nickname, avatar))}, [nickname, avatar])
    const [status, setStatus] = useState(NoneState);

    useEffect(() => {
        const func = async () => {
            setStatus(LoadingState)
            try {
                const json_data = await getProfileSettings()
                loadSettings(json_data)
                setStatus(SuccessState)
            } catch (error) {
                setStatus(ErrorState)
            }
        }
        func()
    }, [])

    function loadSettings(settings) {
        setLogin(settings.login)
        setNickName(settings.nickname)
        setAvatar(settings.avatar)
    }

    async function handleSubmit(e) {
        e.preventDefault();
        if (
            form_errors === null || nickname_errors === null || avatar_errors === null ||
            form_errors.length > 0 || nickname_errors.length > 0 || avatar_errors.length > 0
        ) {
            return;
        }

        setStatus(LoadingState)
        try {
            const json_data = await postProfileSettings({ 
                login: login,
                avatar: avatar,
                nickname: nickname
            })
            loadSettings(json_data)
            setStatus(SuccessState)
        } catch (error) {
            setFormErrors(["Registration error " + error])
            setStatus(ErrorState)
        }
    }

    async function handleReset(e) {
        e.preventDefault();

        setStatus(LoadingState)
        try {
            const json_data = await getProfileSettings()
            loadSettings(json_data)
            setStatus(SuccessState)
        } catch (error) {
            setStatus(ErrorState)
        }
    }


    return (
        <WidgetWithState state={status}>
            <Form errors={form_errors} onSubmit={handleSubmit}>
                <h2 className="form__title">Profile settings</h2>
                <ErrorsGroup errors={form_errors} />
                <FormGroup>
                    <label className="label" htmlFor="login">Login remainder:</label>
                    <input className="input" id="login" type="text" value={login} disabled />
                </FormGroup>
                <FormGroup errors={nickname_errors}>
                    <label className="label" htmlFor="nickname">Nick name:</label>
                    <ErrorsGroup errors={nickname_errors} />
                    <input className="input" id="nickname" type="text" name="nickname" value={nickname} onChange={e => setNickName(e.target.value)} />
                </FormGroup>
                <FormGroup errors={avatar_errors}>
                    <label className="label" htmlFor="avatar">Avatar:</label>
                    { avatar_errors !== null && avatar_errors.length === 0 && <img style={{width: '128px', height: '128px'}} src={avatar} alt="You have no avatar..." />}
                    <ErrorsGroup errors={avatar_errors} />
                    <input className="input" id="avatar" type="file" name="avatar" onChange={e => setAvatar(URL.createObjectURL(e.target.files[0]))} />
                </FormGroup>
                <button className="btn" type="submit">Сохранить</button>
                <button className="btn btn-neutral" type="reset" onClick={handleReset}>Сбросить</button>
            </Form>
        </WidgetWithState>
    )
}