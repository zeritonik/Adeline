import RegistrationForm from "../forms/RegistrationForm";

function RegistrationSection() {
    return (
        <section className="section" id="Registration">
            <div style={{"width": "40%"}}>
                <RegistrationForm />
            </div>
        </section>
    )
}

export default RegistrationSection;