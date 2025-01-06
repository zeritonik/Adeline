import { settings_url } from "./settings"

export async function getProfileSettings() {
    const response = await fetch(settings_url, {
        method: "GET",
    })
    const json_data = await response.json()
    return json_data
}

export async function postProfileSettings(settings) {
    let formData = new FormData()
    formData.append("login", settings.login)
    formData.append("nickname", settings.nickname)
    formData.append("avatar", settings.avatar)
    const response = await fetch(settings_url, {
        method: "POST",
        body: formData,
    })
    const json_data = await response.json()
    return json_data
}

