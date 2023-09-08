import "./verify.css";
import {useState} from "react";
import {axiosInstance} from "../../axios/axios";
import {useNavigate} from "react-router-dom";

export const Verify = () => {
    const [file, setFile] = useState()
    const [enable, setEnable] = useState(false)
    const navigate = useNavigate();

    const handleImage = (e) => {
        setFile(e.target.files[0])
        setEnable(true)
    }

    const send = async ()=> {
        const formData = new FormData()
        formData.append("image", file)
        try {
            const res = await axiosInstance.post("/users/verify", formData)
            console.log(res)
            navigate("/feed")
        } catch (e) {
            console.log(e)
        }
    }

    return (
        <div className={"verify"}>
            <h1>verify</h1>
            <input type="file" onChange={event => handleImage(event)} />
            <button disabled={!enable} onClick={send}>ok</button>
        </div>
    );
}