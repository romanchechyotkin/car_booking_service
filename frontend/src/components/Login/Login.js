import React, {useState} from 'react';
import "./login.css"
import {axiosInstance} from "../../axios/axios";
import {useNavigate} from "react-router-dom";
import {userActions} from "../../store/loginUserSlice";
import {useDispatch} from "react-redux";


export const Login = () => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const dispatch = useDispatch();
    const navigate = useNavigate();

    const login = async () => {
        try {
            const res = await axiosInstance.post("/auth/login", JSON.stringify(
                {
                    "email": email,
                    "password": password,
                }
            ))

            console.log(res)

            localStorage.setItem('access_token', JSON.stringify(res.data.access_token))
            localStorage.setItem('refresh_token', JSON.stringify(res.data.refresh_token))
            localStorage.setItem('user', JSON.stringify(res.data.user))

            dispatch(userActions.setUser(res.data.user))
            dispatch(userActions.setRole(res.data.user.role))
            dispatch(userActions.setIsAuth())
            if (res.data.user.is_verified === true) {
                dispatch(userActions.setIsVerified())
            }

            navigate("/feed")
        } catch (e) {
            console.log(e)
        }
    }

    return (
        <div className={"login"}>
            <h1>login form</h1>
            <div className={"login_form"}>
                <input type="text" placeholder={"email"} value={email} onChange={event => setEmail(event.target.value)} />
                <input type="text" placeholder={"password"} value={password} onChange={event => setPassword(event.target.value)} />
                <button onClick={login}>ok</button>
            </div>
        </div>
    )
};
