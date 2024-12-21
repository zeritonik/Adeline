import { useNavigate } from "react-router-dom"
import { useContext } from "react"
import { UserContext } from "../App"

export function useAnouth() {
    const navigate = useNavigate()
    const [ user, _ ] = useContext(UserContext)

    if (!user) {
        navigate("/login")
    }
}