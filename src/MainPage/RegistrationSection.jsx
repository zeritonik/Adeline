import RegistrationForm from "../forms/RegistrationForm";
import TestCreationForm from "../forms/TestCreationForm";

function RegistrationSection() {
    return (
        <section className="section" id="Registration">
            <div style={{"width": "40%"}}>
                <RegistrationForm />
                <TestCreationForm />
            </div>
        </section>
    )
}

export default RegistrationSection;