import {useLocation} from 'react-router-dom';

export default function PageNotFound() {
    console.log(useLocation())
    return (
        <section className="section">
            <h2 className="section__title">404</h2>
        </section>
    )
}