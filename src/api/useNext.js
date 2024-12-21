import { useNavigate, useParams } from "react-router-dom"

export function useNext(next="/") {
    const params = useParams()
    const navigate = useNavigate()

    console.log(params.next)
    next = params.next ? params.next : next
    return [ next, navigate ]
}