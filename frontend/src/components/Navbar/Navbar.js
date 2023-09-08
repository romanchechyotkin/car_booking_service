import React from 'react';
import "./navbar.css";
import {userActions} from "../../store/loginUserSlice";
import {useDispatch, useSelector} from "react-redux";

export const Navbar = () => {
    const isAuth = useSelector((state) => state.user.isAuth)
    const dispatch = useDispatch();

    const logout = () => {
        localStorage.removeItem("access_token")
        localStorage.removeItem("user")
        dispatch(userActions.logout())
    }

    return (
        <nav className={"navbar"}>
            {isAuth && <button onClick={logout}>log out</button>}
        </nav>
    )
};