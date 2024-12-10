function HeaderNavigation({items}) {
    return (
        <nav class="header__navigation">
            <ul>
                {items.map(item => <HeaderNavigationItem href={item.href}>{item.text}</HeaderNavigationItem>)}
            </ul>
        </nav>
    )
}

function HeaderNavigationItem({href, children}) {
    return (
        <li>
            <a href={href}>{children}</a>
        </li>
    )
}

export default HeaderNavigation