import React, { useState } from 'react';
import "./registration.css"
import { axiosInstance } from "../../axios/axios";
import { useNavigate } from "react-router-dom";
import { useDispatch } from "react-redux";
import { loginUser } from '../SharedFunctions/login';
import ToastError from '../ErrorComponent/error-state';

export const Registration = () => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [fullName, setFullName] = useState("");
    const [telephoneNumber, setTelephoneNumber] = useState("");
    const [city, setCity] = useState("");
    const [error, setError] = useState("");
    const navigate = useNavigate();
    const dispatch = useDispatch();

    const registration = async () => {
        try {
            const res = await axiosInstance.post("/auth/registration", JSON.stringify(
                {
                    email,
                    password,
                    full_name: fullName,
                    telephone_number: telephoneNumber,
                    city,
                }
            ));

            console.log(res);
            
            if (res.status) {
                await loginUser(email, password, navigate, dispatch)
            }

        } catch (e) {
            console.log("cqtch", e);
            setError(e.response?.data.error);
        }
    };

    return (
        <div className={"registration"}>
            <ToastError error={error} />
            <h1>registration form</h1>
            <div className={"registration_form"}>
                <input type="text" placeholder={"email"} value={email} onChange={event => setEmail(event.target.value)} />
                <input type="password" placeholder={"password"} value={password} onChange={event => setPassword(event.target.value)} />
                <input type="text" placeholder={"full name"} value={fullName} onChange={event => setFullName(event.target.value)} />
                <input type="tel" placeholder={"telephone number"} value={telephoneNumber} onChange={event => setTelephoneNumber(event.target.value)} />
                <input type="text" placeholder={"city"} value={city} onChange={event => setCity(event.target.value)} />
                <button onClick={registration}>ok</button>
            </div>
        </div>
    );
};
