import { register_url, login_url } from "./settings"

export async function registerUser(login, password) {
    const response = await fetch(register_url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ login, password }),
    })
    const json_data = await response.json()
    return json_data
}

export async function loginUser(login, password) {
    const response = await fetch(login_url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ login, password }),
    })
    const json_data = await response.json()
    return json_data
}