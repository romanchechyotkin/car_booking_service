import React, {useState} from 'react';
import "./login.css"
import {axiosInstance} from "../../axios/axios";
import {useNavigate} from "react-router-dom";
import {useDispatch} from "react-redux";
import { loginUser } from '../SharedFunctions/login';


export const Login = () => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const dispatch = useDispatch();
    const navigate = useNavigate();

    const login = () => {
        loginUser(email, password, navigate, dispatch)
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
