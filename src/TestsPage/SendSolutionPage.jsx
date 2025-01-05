import { useParams } from "react-router-dom";
import { SendForm } from "../forms/SendForm"


export default function SendSolutionPage() {
    const { id } = useParams();

    return <section className="section">
        <div style={{"width": "40%"}}>
            <SendForm test_group_id={id} />
        </div>
    </section>
}