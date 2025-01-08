import { useState, useEffect, useContext } from "react";

import { UserContext } from "../App";

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

async function processAvatar(avatar, avatar_url, setAvatarUrl) {
    if (avatar_url) {
        URL.revokeObjectURL(avatar_url)
    }

    if (typeof(avatar) === "string") {
        setAvatarUrl(avatar)
        return 
    }

    if (typeof(avatar) === "object" && avatar instanceof File) {
        const img = new Image();
        img.src = URL.createObjectURL(avatar);
        img.onload = () => {
            const canvas = document.createElement('canvas');
            canvas.width = 128;
            canvas.height = 128;
    
            const ctx = canvas.getContext('2d');
            ctx.drawImage(img, 0, 0, 128, 128);
            setAvatarUrl(canvas.toDataURL("image/png"));
        }
        return
    }
}


export default function ProfileSettingsForm() {
    const [user, setUser] = useContext(UserContext);

    const [login, setLogin] = useState(user.login)

    const [nickname, setNickName] = useState(user.nickname);
    const [nickname_errors, setNicknameErrors] = useState(null);
    useEffect(() => {setNicknameErrors(validateNickName(nickname))}, [nickname])

    const [avatar_file, setAvatarFile] = useState(user.avatar);
    const [avatar_url, setAvatarUrl] = useState(null);
    useEffect(() => { processAvatar(avatar_file, avatar_url, setAvatarUrl) }, [avatar_file])

    const [form_errors, setFormErrors] = useState(null);
    useEffect(() => {setFormErrors(validateForm(nickname, avatar_file))}, [nickname, avatar_file])

    const [status, setStatus] = useState(SuccessState);

    useEffect(() => {
        setLogin(user.login);
        setNickName(user.nickname);
        setAvatarFile(user.avatar);
    }, [user]);

    async function handleSubmit(e) {
        e.preventDefault();
        if (
            form_errors === null || nickname_errors === null ||
            form_errors.length > 0 || nickname_errors.length > 0
        ) {
            return;
        }

        setStatus(LoadingState)
        try {
            const json_data = await postProfileSettings({ 
                login: login,
                avatar: await fetch(avatar_url).then(res => res.blob()),
                nickname: nickname
            })
            setUser(json_data)
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
            setUser(json_data)
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
                <FormGroup>
                    <label className="label" htmlFor="avatar">Avatar:</label>
                    { avatar_url && <img className="avatar" src={avatar_url} alt="You have no avatar..." /> }
                    <input className="input" id="avatar" type="file" name="avatar" onChange={e => setAvatarFile(e.target.files[0])} />
                </FormGroup>
                <button className="btn" type="submit">Сохранить</button>
                <button className="btn btn-neutral" type="reset" onClick={handleReset}>Сбросить</button>
            </Form>
        </WidgetWithState>
    )
}