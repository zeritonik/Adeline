import { profile_tests_url, profile_results_url } from "./settings"


export async function createTestGroup(name, time_limit, memory_limit, tests) {
    const response = await fetch(profile_tests_url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ 
            type: "GROUP_TEST", 
            name, 
            time_limit, 
            memory_limit,
            tests
        }),
    })
    return await response.json()
}


export async function getTestGroups() {
    const response = await fetch(profile_tests_url, {
        method: "GET",
    })
    return await response.json()
}


export async function getTestGroup(id) {
    const response = await fetch(profile_tests_url + `/${id}`, {
        method: "GET",
    })
    return await response.json()
}


export async function getTestResults() {
    const response = await fetch(profile_results_url, {
        method: "GET",
    })
    return await response.json()
}

export async function getTestResult(id) {
    const response = await fetch(profile_results_url + `/${id}`, {
        method: "GET",
    })
    return await response.json()
}