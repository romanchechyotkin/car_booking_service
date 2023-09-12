import React, {useEffect} from 'react';
import {BrowserRouter, Navigate, Route, Routes} from "react-router-dom";
import {useDispatch, useSelector} from "react-redux";
import {userActions} from "./store/loginUserSlice";

import {Registration} from "./components/Registration/Registration";
import {Login} from "./components/Login/Login";
import {Navbar} from "./components/Navbar/Navbar";
import {Admin} from "./components/Admin/Admin";
import {Feed} from "./components/Feed/Feed";
import {Verify} from "./components/Verify/Verify";
import {CarPage} from "./components/CarPage/CarPage";

export const App = () => {
    const isAuth = useSelector((state) => state.user.isAuth)
    // const isVerified = useSelector((state) => state.user.isVerified)
    const role = useSelector((state) => state.user.role)
    const dispatch = useDispatch()

    useEffect(() => {
        const token = JSON.parse(localStorage.getItem('access_token'))
        if (!token) {
            dispatch(userActions.logout)
        } else {
            const user = JSON.parse(localStorage.getItem('user'))
            dispatch(userActions.setUser(user))
            dispatch(userActions.setRole(user.role))
            dispatch(userActions.setIsAuth())
            if (user.is_verified) {
                dispatch(userActions.setIsVerified())
            }
        }
    }, [dispatch]);

    return (
        <BrowserRouter>
            <Navbar />
            <Routes>
                {role !== "ADMIN" && <>
                        <Route path="*" element={<Navigate to={isAuth ? "/" : "/login"} />} />
                        <Route path={"/"} element={isAuth ? <Navigate to={"/feed"} /> : <Navigate to={"/login"} />} />
                        <Route path={"/login"} element={isAuth ? <Navigate to={"/"} /> : <Login />} />
                        <Route path={"/registration"} element={isAuth ? <Navigate to={"/"} /> : <Registration />} />
                        <Route path={"/feed"} element={<Feed />} />
                        <Route path={"/verify"} element={<Verify />} />
                        <Route path={"/post/:id"} element={<CarPage />} />
                    </>
                }
                {role === "ADMIN" &&
                    <>
                        <Route path={"*"} element={isAuth ? <Navigate to={"/admin"} /> : <Navigate to={"/login"} />} />
                        <Route path={"/"} element={isAuth ? <Navigate to={"/admin"} /> : <Navigate to={"/login"} />} />
                        <Route path={"/login"} element={isAuth ? <Navigate to={"/admin"} /> : <Login />} />
                        <Route path={"/registration"} element={isAuth ? <Navigate to={"/admin"} /> : <Registration />} />
                        <Route path={"/admin"} element={isAuth ? <Admin /> : <Navigate to={"/login"} />} />
                    </>
                }
            </Routes>
        </BrowserRouter>
    )
}


