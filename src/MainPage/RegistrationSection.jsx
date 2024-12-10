import LoginForm from "../LoginForm";
import RegistrationForm from "../RegistrationForm";

function RegistrationSection() {
    return (
        <section className="section" id="Registration">
            <h2>Try it!</h2>
            <RegistrationForm style={{ width: '50%' }} />
        </section>
    )
}

export default RegistrationSection;