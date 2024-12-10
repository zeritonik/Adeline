export function Form({ errors, onSubmit, children }) {
    let classes = "form"
    if ( errors !== null ) {
        classes += errors.length ? " form--error" : " form--ok"
    }
    return (
        <form className={classes} onSubmit={onSubmit}>
            {children}
        </form>
    )
}


export function FormGroup({ errors=null, children }) {
    let classes = "form__group"
    if ( errors !== null ) {
        classes += errors.length ? " form__group--error" : " form__group--ok"
    }

    return <div className={classes}>
                {children}
            </div>
}

export function ErrorsGroup({ errors }) {
    if (errors === null)
        return <></>
    return (
        <ul className="form__errors-group">
            {errors.map(error => <li key={error}>{error}</li>)}
        </ul>
    )
}