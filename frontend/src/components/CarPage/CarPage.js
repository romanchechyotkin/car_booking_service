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
    const [reservations, setReservations] = useState([])

    const [startHour, setStartHour] = useState("")
    const [endHour, setEnbHour] = useState("")

    const [startDate, setStartDate] = useState()
    const [endDate, setEndDate] = useState()

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

    const fetchCarReservations = async () => {
        try {
            const res = await axiosInstance.get(`/cars/${params.id}/reservations`)
            console.log(res.data)
            setReservations(res.data)
        } catch (e) {
            console.log(e)
        }
    }

    useEffect(() => {
        fetchCarInfo()
        fetchCarRates()
        fetchCarReservations()
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

    const reserve = async () => {
        const start = new Date(startDate)
        const end = new Date(endDate)

        let startDay = start.getDate()
        let startMonth = start.getMonth()
        let startYear = start.getFullYear()
        console.log(startHour, startDay, startMonth+1, startYear)

        let endDay = end.getDate()
        let endMonth = end.getMonth()
        let endYear = end.getFullYear()
        console.log(endHour, endDay, endMonth+1, endYear)

        try {
            const res = await axiosInstance.post(`/cars/${params.id}/rent`, JSON.stringify({
                "start_date": `${addZero(startHour)} ${addZero(startDay)}.${addZero(startMonth+1)}.${startYear}`,
                "end_date": `${addZero(endHour)} ${addZero(endDay)}.${addZero(endMonth+1)}.${endYear}`,
            }))
            console.log(res)
            window.location.reload()
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
                <div className={"reservation_dates"}>       {/*    15 02.01.2006  */}
                    <div className={"reservation_start"}>
                        <label htmlFor={"start"}>start date</label>
                        <input type="number" min={0} max={23} onChange={e => setStartHour(e.target.value)} />
                        <input onChange={e => setStartDate(e.target.value)} type="date" id={"start"}/>
                    </div>
                    <div className={"reservation_end"}>
                        <label htmlFor={"end"}>end date</label>
                        <input type="number" min={0} max={23} onChange={e => setEnbHour(e.target.value)} />
                        <input onChange={e => setEndDate(e.target.value)} type="date" id={"end"}/>
                    </div>
                    <button onClick={reserve}>ok</button>
                </div>
                {reservations !== "no reservations" ? reservations.map(r =>
                        <div key={r.start_date}>Start:{r.start_date} End:{r.end_date}</div>)
                    :
                    <div>no reservations</div>
                }
            </div>
        </div>
    )
}

function addZero(num) {
    if (String(num).length === 1) {
        return `0${num}`
    }

    return `${num}`
}