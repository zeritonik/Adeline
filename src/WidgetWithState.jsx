export const NoneState = 0
export const LoadingState = 1
export const SuccessState = 2
export const ErrorState = -1

export function WidgetWithState({ state, children }) {
    switch (state) {
        case NoneState:
            return <div className="widget-error"><h1>None state</h1></div>
        case LoadingState:
            return (
                <div className="widget-loading-container">
                    <span>Loading...</span>
                    <div className="widget-loading">
                        <div></div>
                        <div></div>
                        <div></div>
                        <div></div>
                    </div>
                </div>
            )
        case SuccessState:
            return children
        case ErrorState:
            return (
                <>
                    <div className="widget-error">
                        <h1>Something went wrong</h1>
                    </div>
                </>
            )
        default:
            return <div className="widget-error">
                <h1>Incorrect state!</h1>
            </div>
    }
}