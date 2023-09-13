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
    const [isLast, setIsLast] = useState(true)
    const [isAscPrice, setIsAscPrice] = useState(false)
    const [isDescPrice, setIsDescPrice] = useState(false)
    const [isAscYear, setIsAscYear] = useState(false)
    const [isDescYear, setIsDescYear] = useState(false)

    const setLast = () => {
        setIsLast(true)
        setIsAscPrice(false)
        setIsDescPrice(false)
        setIsAscYear(false)
        setIsDescYear(false)
    }

    const setPriceAsc = () => {
        setIsLast(false)
        setIsAscPrice(true)
        setIsDescPrice(false)
        setIsAscYear(false)
        setIsDescYear(false)
    }

    const setPriceDesc = () => {
        setIsLast(false)
        setIsAscPrice(false)
        setIsDescPrice(true)
        setIsAscYear(false)
        setIsDescYear(false)
    }

    const setYearAsc = () => {
        setIsLast(false)
        setIsAscPrice(false)
        setIsDescPrice(false)
        setIsAscYear(true)
        setIsDescYear(false)
    }

    const setYearDesc = () => {
        setIsLast(false)
        setIsAscPrice(false)
        setIsDescPrice(false)
        setIsAscYear(false)
        setIsDescYear(true)
    }

    // SORT_BY_ASC_PRICE  = "prc.a"
    // SORT_BY_DESC_PRICE = "prc.d"
    // SORT_BY_ASC_YEAR   = "y.a"
    // SORT_BY_DESC_YEAR  = "y.d"

    // const dispatch = useDispatch();

    const fetchCars= async () => {
        try {
            let res
            if (isLast) {
                res= await axiosInstance.get("/cars")
            } else if (isAscYear) {
                res = await axiosInstance.get("/cars?sort=y.a")
            } else if (isDescYear) {
                res = await axiosInstance.get("/cars?sort=y.d")
            } else if (isAscPrice) {
                res = await axiosInstance.get("/cars?sort=prc.a")
            } else if (isDescPrice) {
                res = await axiosInstance.get("/cars?sort=prc.d")
            }

            console.log(res.data)
            if (res.data !== null) {
                setCars(res.data)
            }
        } catch (e) {
            console.log(e)
        }

    }

    useEffect(() => {
        fetchCars()
    }, [isAscYear, isAscPrice, isDescPrice, isDescYear, isLast])

    return (
        <>
            <h1>feed</h1>
            {isAuth && !isVerified &&
                <div>
                    <h2>you should verify yourself</h2>
                    <Link to={"/verify"}>verify</Link>
                </div>
            }
            <div className={"feed_sort"}>
                <div onClick={setLast}>last</div>
                <div onClick={setPriceAsc}>price asc</div>
                <div onClick={setPriceDesc}>price desc</div>
                <div onClick={setYearAsc}>year asc</div>
                <div onClick={setYearDesc}>year desc</div>
            </div>
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