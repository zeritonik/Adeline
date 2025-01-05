import { register_url, login_url, logout_url } from "./settings"

export async function registerUser(login, nickname, password) {
    const response = await fetch(register_url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ login, nickname, password }),
    })
    if (response.status !== 200) {
        throw new Error("Registration error")
    }
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
    if (response.status === 401) {
        throw new Error("Invalid login or password")
    }
    const json_data = await response.json()
    return json_data
}


export async function logoutUser(all=false) {
    const response = await fetch(logout_url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ all }),
    })
    if (response.status !== 200) {
        throw new Error("Logout error")
    }
    const json_data = await response.json()
    return json_data
}