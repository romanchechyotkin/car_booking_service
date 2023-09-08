import {useEffect, useState} from "react";
import {axiosInstance, STATIC} from "../../axios/axios";

export const Admin = () => {
    const [users, setUsers] = useState([])

    const fetchUsers= async () => {
        try {
            const res= await axiosInstance.get("/users/verify")
            console.log(res.data)
            if (res.data !== null) {
                setUsers(prevState => [...prevState, ...res.data])
            }
        } catch (e) {
            console.log(e)
        }

    }

    useEffect(() => {
        fetchUsers()
    }, [])

    const verify = async (id) => {
        try {
            const res= await axiosInstance.post(`/users/verify/${id}`)
            console.log(res)
            setUsers(users.filter(u => u.user_id !== id))
        } catch (e) {
            console.log(e)
        }
    }

    return (
        <>
            <h1>admin panel</h1>
            {users !== null && users.map(u =>
                <div key={u.user_id}>
                    <img src={STATIC+u.filename} alt="img"/>
                    <button onClick={() => verify(u.user_id)}>verify</button>
                </div>
            )}
        </>
    )
}