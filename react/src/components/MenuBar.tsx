import { Link } from 'react-router-dom';
import { useAuth0 } from '@auth0/auth0-react';
import { LoginButton } from './LoginButton';
import { LogoutButton } from './LogoutButton';

export const MenuBar = () => {
    const { isAuthenticated, isLoading } = useAuth0();

    const renderAuthenticated = () => {
        return (
            <>
                <div className="menu-bar__auth">
                    <LogoutButton />
                </div>
                <ul>
                    <li><Link to="/">Landing</Link></li>
                    <li><Link to="/home">Home</Link></li>
                    <li><Link to="/profile">Profile</Link></li>
                    <li><Link to="/about">About</Link></li>
                </ul>
            </>
        );
    }

    if (isLoading) {
        return;
    }

    return (
        <div className="menu-bar">
            {isAuthenticated ? renderAuthenticated() : <LoginButton />}
        </div>
    );
};

export default MenuBar;
