import "./Header.css"

import Logo from "../../logo/Logo"
import HeaderNavigation from "./HeaderNavigation"
import HeaderProfile from "./HeaderProfile"

function Header() {
    const HEADER_ITEMS = {
        navigation: [
          {href: "#About", text: "About"},
          {href: "#Examples", text: "Usage examples"},
          {href: "#Registration", text: "Try it"},
        ]
      }
    return (
        <header className="header">
            <HeaderNavigation items={HEADER_ITEMS.navigation} />
            <Logo />
            <div style={{display: "flex", justifyContent: "flex-end"}}>
                <HeaderProfile />
            </div>
        </header>
    )
}

export default Header