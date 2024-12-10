const protocol = "http"
const host = "localhost"
const port = 8080

export const baseUrl = protocol + "://" + host + ":" + port
export const tests_url = baseUrl + "/tests"
export const register_url = baseUrl + "/register"
export const login_url = baseUrl + "/login"