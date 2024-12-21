import TestGroupCreationForm from "../forms/TestGroupCreationForm";

export default function NewTestGroupPage() {
    return (
        <section className="section">
            <div style={{"width": "40%"}}>
                <TestGroupCreationForm />
            </div>
        </section>
    )
}