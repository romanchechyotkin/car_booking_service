import React, {useState} from 'react';
import "./registration.css"
import {axiosInstance} from "../../axios/axios";
import {useNavigate} from "react-router-dom";

// "city": "string",
// "email": "string",
// "full_name": "string",
// "password": "string",
// "telephone_number": "string"

export const Registration = () => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [fullName, setFullName] = useState("");
    const [telephoneNumber, setTelephoneNumber] = useState("");
    const [city, setCity] = useState("");
    const navigate = useNavigate();


    const registration = async () => {
        try {
            const res = await axiosInstance.post("/auth/registration", JSON.stringify(
                {
                    "email": email,
                    "password": password,
                    "full_name": fullName,
                    "telephone_number": telephoneNumber,
                    "city": city
                }
            ))

            console.log(res)

            navigate("/login")
        } catch (e) {
            console.log(e)
        }
    }

    return (
        <div className={"registration"}>
            <h1>registration form</h1>
            <div className={"registration_form"}>
                <input type="text" placeholder={"email"} value={email} onChange={event => setEmail(event.target.value)} />
                <input type="text" placeholder={"password"} value={password} onChange={event => setPassword(event.target.value)} />
                <input type="text" placeholder={"full name"} value={fullName} onChange={event => setFullName(event.target.value)} />
                <input type="tel" placeholder={"telephone number"} value={telephoneNumber} onChange={event => setTelephoneNumber(event.target.value)} />
                <input type="text" placeholder={"city"} value={city} onChange={event => setCity(event.target.value)} />
                <button onClick={registration}>ok</button>
            </div>
        </div>
    )
};
