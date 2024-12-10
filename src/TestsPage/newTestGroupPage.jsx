import GroupTestCreationForm from "../forms/GroupTestCreationForm";

export default function NewTestGroupPage() {
    return (
        <section className="section">
            <div style={{"width": "40%"}}>
                <GroupTestCreationForm />
            </div>
        </section>
    )
}