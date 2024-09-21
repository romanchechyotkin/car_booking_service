import React, { useState } from 'react';
import "./navbar.css";
import { userActions } from "../../store/loginUserSlice";
import { useDispatch, useSelector } from "react-redux";
import { Link } from "react-router-dom";

export const Navbar = () => {
    const dispatch = useDispatch();
    const { isAuth, role } = useSelector((state) => state.user);

    const logout = () => {
        dispatch(userActions.logout());
        localStorage.clear();
    };

    const [searchText, setSearchText] = useState("");

    const handleSearchChange = (e) => {
        setSearchText(e.target.value);
    };

    const clearSearch = () => {
        setSearchText("");
    };

    return (
        <nav className="navbar">
            <div className="navbar_logo">
                <Link to="/">CarBook</Link>
            </div>
            
            {/* Добавляем поисковую строку */}
            <div className="search-container">
                <input 
                    type="text" 
                    className="search-input"
                    placeholder="Search..."
                    value={searchText}
                    onChange={handleSearchChange}
                />
                {searchText && (
                    <button className="clear-button" onClick={clearSearch}>
                        ✖
                    </button>
                )}
                <span className="search-icon">🔍</span>
            </div>

            <ul className="navbar_links">
                {isAuth && role !== "ADMIN" && (
                    <>
                        <li>
                            <Link to="/feed">Feed</Link>
                        </li>
                        <li>
                            <Link to="/verify">Verify</Link>
                        </li>
                    </>
                )}
                {isAuth && role === "ADMIN" && (
                    <>
                        <li>
                            <Link to="/admin">Admin Panel</Link>
                        </li>
                    </>
                )}
                {!isAuth && (
                    <>
                        <li>
                            <Link to="/login">Login</Link>
                        </li>
                        <li>
                            <Link to="/registration">Register</Link>
                        </li>
                    </>
                )}
                {isAuth && (
                    <li>
                        <button onClick={logout}>Logout</button>
                    </li>
                )}
            </ul>
        </nav>
    );
};
