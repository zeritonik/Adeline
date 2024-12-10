import { tests_url } from "./settings"

export async function createTest(input, correctOutput) {
    const response = await fetch(tests_url, {
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
    const json_data = await response.json()
    return json_data
}


export async function createGroupTest(name, time_limit, memory_limit) {
    const response = await fetch(tests_url, {
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
    const json_data = await response.json()
    return json_data
}