import React from 'react';
import "./navbar.css";
import {userActions} from "../../store/loginUserSlice";
import {useDispatch, useSelector} from "react-redux";
import {Link} from "react-router-dom";

export const Navbar = () => {
    const isAuth = useSelector((state) => state.user.isAuth)
    const user = useSelector((state) => state.user.user)
    const dispatch = useDispatch();

    const logout = () => {
        localStorage.removeItem("access_token")
        localStorage.removeItem("user")
        dispatch(userActions.logout())
    }

    return (
        <nav className={"navbar"}>
            {isAuth ? <button onClick={logout}>log out</button> : <Link to={"/login"}>login</Link>}
            {user && <div>{user.email}</div>}
        </nav>
    )
};