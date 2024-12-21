const protocol = "http"
const host = "localhost"
const port = 8080

export const baseUrl = protocol + "://" + host + ":" + port + "/api"

export const register_url = baseUrl + "/register"
export const login_url = baseUrl + "/login"
export const logout_url = baseUrl + "/logout"

export const profile_url = baseUrl + "/profile"
export const settings_url = profile_url + "/settings"

export const tests_url = baseUrl + "/tests"
export const profile_tests_url = baseUrl + "/profile/tests"