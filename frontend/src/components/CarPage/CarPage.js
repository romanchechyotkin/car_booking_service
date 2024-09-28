import "./carpgage.css"
import {Link, useParams} from "react-router-dom";
import React, {useEffect, useState} from "react";
import {axiosInstance, STATIC} from "../../axios/axios";
import {Rate} from "../Rate/Rate";
import {useSelector} from "react-redux";
import Calendar from 'react-calendar';
import 'react-calendar/dist/Calendar.css';
import { BillingInfo } from "../BillingInfo/BillingInfo";
import {useNavigate} from "react-router-dom";

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
    const [endHour, setEndHour] = useState("")

    const [startDate, setStartDate] = useState()
    const [endDate, setEndDate] = useState()
   
    const navigate = useNavigate();


    const [marketingChecked, setMarketingChecked] = useState(false);
    const [termsChecked, setTermsChecked] = useState(false);

    
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
            if (e.response?.status === 404) {
                setRate("Пока что нет отзывов")
            }
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

    //перенести функцию в биллинг
    
    const handleReserve = () => {
        navigate(`/billing/${car.id}`);
    };


    return (
            <div className={"car_page"}>
        {car !== null &&
            <div className="car_page__content">
                <div className="car_page__image">
                    <img src={STATIC + car.images[0]} alt={`${car.brand} ${car.model}`} />
                </div>

                <div className="car_page__details">
                    

                    <div className="reservation">
                    <h1>{car.brand} {car.model}</h1>
                            <div className={"car_rates"}>
                            {isAuth ? <div className={"rate_form"}>
                            <h3>{rate}</h3>
                            </div> : <div>to leave rate u need to be auth</div>}
                            
                        </div>
                        <p className="car_description">{car.description} Lorem  inputdsjfskhdfajd  odshflkahdjf hsh akdfh askh dkfhskd hfh sakh fksajhdjkf  hjfhsakldh fkahkshd fkahkshd shdakjlfaslkhdfj askh
                            sdjafhlkjsahlkdfhsadkf hskjdlafha hfskhkhehu fhru. Hhfhejdke lopaskd ffjdl
                        </p>

                        <div className="year-container">
                            <div className=""></div>
                            <p> Год выпуска </p>
                            <p className="year">
                                
                                    {car.year}
                            </p>
                        </div>
                        <div className="footer">
                            <div className="car_price">
                                <span className="price">${car.pricePerDay}/day</span>
                            </div>
                            <div>
                        <button onClick={handleReserve}>
                            OK
                        </button>

            {/* Показываем BillingInfo только после нажатия */}
            
        </div>
                        </div>
                        {/* <div className="reservation_dates">
                            <div className="reservation_start">
                                <label htmlFor="start">Start date</label>
                                <input
                                    type="number"
                                    min={0}
                                    max={23}
                                    onChange={e => setStartHour(e.target.value)}
                                    placeholder="Hour"
                                />
                                <input
                                    onChange={e => setStartDate(e.target.value)}
                                    type="date"
                                    id="start"
                                />
                            </div>
                            <div className="reservation_end">
                                <label htmlFor="end">End date</label>
                                <input
                                    type="number"
                                    min={0}
                                    max={23}
                                    onChange={e => setEndHour(e.target.value)}
                                    placeholder="Hour"
                                />
                                <input
                                    onChange={e => setEndDate(e.target.value)}
                                    type="date"
                                    id="end"
                                />
                            </div>
                            <div className="footer">
                            <div className="car_price">
                                <span className="price">${car.pricePerDay}/day</span>
                                <span className="old_price">$100.00</span>
                            </div>
                            <button onClick={reserve}>OK</button>
                            </div>
                        </div> */}
                    </div>
                </div>
            </div>
        }
        
        
    </div>

    )
}

