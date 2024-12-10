import { baseUrl } from "./settings"

const profile_url = baseUrl + "/profile"

export async function getProfileSettings() {
    const response = await fetch(profile_url + "/settings", {
        method: "GET",
    })
    const json_data = await response.json()
    return json_data
}

export async function postProfileSettings(settings) {
    const response = await fetch(profile_url + "/settings", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(settings),
    })
    const json_data = await response.json()
    return json_data
}