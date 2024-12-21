import { profile_tests_url } from "./settings"


export async function createTest(input, correctOutput) {
    const response = await fetch(profile_tests_url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ 
            type: "TEST", 
            input, 
            correctOutput 
        }),
    })
    return await response.json()
}


export async function createTestGroup(name, time_limit, memory_limit) {
    const response = await fetch(profile_tests_url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ 
            type: "GROUP_TEST", 
            name, 
            time_limit, 
            memory_limit 
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