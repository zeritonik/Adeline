import { useNavigate, useSearchParams } from "react-router-dom"

export function useNext(next="/") {
    const [params] = useSearchParams()
    const navigate = useNavigate()

    next = params.get("next") ? params.get("next") : next
    return [ next, navigate ]
}