import {useSelector} from "react-redux";
import {Link} from "react-router-dom";
import {useEffect, useState} from "react";
import {axiosInstance} from "../../axios/axios";
import {CarPost} from "../CarPost/CarPost";
import "./feed.css"

export const Feed = () => {
    const isVerified = useSelector((state) => state.user.isVerified)
    const isAuth = useSelector((state) => state.user.isAuth)
    const [cars, setCars] = useState([])
    // const dispatch = useDispatch();

    const fetchCars= async () => {
        try {
            const res= await axiosInstance.get("/cars")
            console.log(res.data)
            if (res.data !== null) {
                setCars(prevState => [...prevState, ...res.data])
            }
        } catch (e) {
            console.log(e)
        }

    }

    useEffect(() => {
        fetchCars()
    }, [])

    return (
        <>
            <h1>feed</h1>
            {isAuth && !isVerified &&
                <div>
                    <h2>you should verify yourself</h2>
                    <Link to={"/verify"}>verify</Link>
                </div>
            }
            <div className={"posts"}>
                {cars !== null && cars.map(c =>
                    <Link to={`/post/${c.car.id}`} >
                        <CarPost car={c.car} user_id={c.user_id} />
                    </Link>
                )}
            </div>

        </>
    )
}