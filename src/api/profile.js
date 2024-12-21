import { settings_url } from "./settings"

export async function getProfileSettings() {
    const response = await fetch(settings_url, {
        method: "GET",
    })
    const json_data = await response.json()
    return json_data
}

export async function postProfileSettings(settings) {
    const response = await fetch(settings_url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(settings),
    })
    const json_data = await response.json()
    return json_data
}

