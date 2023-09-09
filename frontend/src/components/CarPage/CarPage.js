import "./carpgage.css"
import {useParams} from "react-router-dom";
import React, {useEffect, useState} from "react";
import {axiosInstance, STATIC} from "../../axios/axios";
import {Rate} from "../Rate/Rate";
import {useSelector} from "react-redux";

export const CarPage = () => {
    const user = useSelector((state) => state.user.user)
    const isAuth = useSelector((state) => state.user.isAuth)
    const params = useParams();
    const [car, setCar] = useState(null)
    const [rates, setRates] = useState([])
    const [comment, setComment] = useState("")
    const [rate, setRate] = useState(1)

    const fetchCarInfo = async () => {
        try {
            const res = await axiosInstance.get(`/cars/${params.id}`)
            setCar(res.data)
            console.log(res.data)
        } catch (e) {
            console.log(e)
        }
    }

    const fetchCarRates = async () => {
        try {
            const res = await axiosInstance.get(`/cars/${params.id}/rate`)
            console.log(res.data)
            setRates(res.data)
        } catch (e) {
            console.log(e)
        }
    }

    useEffect(() => {
        fetchCarInfo()
        fetchCarRates()
    }, []);

    const sendRate = async () => {
        try {
            const res = await axiosInstance.post(`/cars/${params.id}/rate`, JSON.stringify({
                "comment": comment,
                "rating": rate
            }))
            console.log(res.data)
            setRates(prevState => [...prevState, {
                "comment": comment,
                "rating": rate,
                "user": user.full_name
            }])
            setComment("")
        } catch (e) {
            console.log(e)
        }
    }

    return (
        <div className={"car_page"}>
            <div className={"car_info"}>
                {car !== null &&
                    <>
                        <h1>{car.brand} {car.model}, {car.year}</h1>
                        <div className={"car_page__images"}>
                            {car.images.map(i =>
                                <img key={i} src={STATIC+i} alt="img"/>
                            )}
                        </div>
                        <div>{car.pricePerDay}$</div>
                        <div>{car.rating}</div>
                    </>
                }
                <div className={"car_rates"}>
                    <h3>Car's rate</h3>
                    {isAuth ? <div className={"rate_form"}>
                        <input type="text" value={comment} onChange={e => setComment(e.target.value)}/>
                        <input type="number" min={1} max={5} value={rate} onChange={e => setRate(Number.parseInt(e.target.value))}/>
                        <button onClick={sendRate}>comment</button>
                    </div> : <div>to leave rate u need to be auth</div>}
                    {rates !== null && rates.map(r =>
                        <Rate rate={r} />
                    )}
                </div>
            </div>
            <div className={"reservation"}>
                <h2>reservation</h2>
            </div>
        </div>
    )
}