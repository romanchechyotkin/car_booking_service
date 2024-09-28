import { useSelector } from "react-redux";
import { useEffect, useState } from "react";
import { axiosInstance, STATIC } from "../../axios/axios";
import "./billing.css";
import visaImage from '../../img/visa-mastercard-logo.png';
import paypalImage from '../../img/paypal.png';
import bitcoinImage from '../../img/bitcoin.png';
import { useParams} from "react-router-dom";

export const BillingInfo = () => {

    const [startHour, setStartHour] = useState("")
    const [endHour, setEndHour] = useState("")

    const [startDate, setStartDate] = useState()
    const [endDate, setEndDate] = useState()
    const [reservations, setReservations] = useState([])
    const [rates, setRates] = useState([])
    const [car, setCar] = useState({})
    const [rate, setRate] = useState(1)
    const user = useSelector((state) => state.user.user)
    const isAuth = useSelector((state) => state.user.isAuth)

    const params = useParams();

    const [marketingChecked, setMarketingChecked] = useState(false);
    const [termsChecked, setTermsChecked] = useState(false);

    const handleMarketingChange = () => {
        setMarketingChecked(!marketingChecked);
    };

    const handleTermsChange = () => {
        setTermsChecked(!termsChecked);
    };

    const fetchCarInfo = async () => {
        try {
            const res = await axiosInstance.get(`/cars/${params.id}`)
            console.log(res.data)
            setCar(res.data)
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

    useEffect(() => {
        fetchCarInfo()
        fetchCarRates()
    }, []);
    
    return (
        <div>
            <div className="main-billing-section">
                <div className="left-section">
                    <div className="billing-card">
                        <div class="billing-header">
                        <h2>Billing Info</h2>
                        <span class="step">Step 1 of 4</span>
                    </div>
                        <p>Please enter your billing info</p>
                            <form class="billing-form">
                        <div>
                            <label for="name">Name</label>
                            <input type="text" id="name" placeholder="Your name"/>
                        </div>
                        <div>
                            <label for="phone">Phone Number</label>
                            <input type="text" id="phone" placeholder="Phone number"/>
                        </div>
                        <div>
                            <label for="address">Address</label>
                            <input type="text" id="address" placeholder="Address"/>
                        </div>
                        <div>
                            <label for="city">Town / City</label>
                            <input type="text" id="city" placeholder="Town or city"/>
                        </div>
                        </form>
                    </div>

                    <div className="rental-info">
                    <div className="header">
                        <h2>Rental Info</h2>
                        <p>Please select your rental date</p>
                        <span>Step 2 of 4</span>
                    </div>

                    <div className="form-section">
                        <div className="radio-group">
                            <input type="radio" id="pickup" name="rentalType" checked/>
                            <label for="pickup">Pick - Up</label>
                        </div>

                        <div className="form-grid">
                            <div className="form-field">
                                <label for="pickup-location">Locations</label>
                                <select id="pickup-location">
                                    <option>Select your city</option>
                                </select>
                            </div>

                            <div className="form-field">
                                <label for="pickup-date">Date</label>
                                <select id="pickup-date">
                                    <option>Select your date</option>
                                </select>
                            </div>

                            <div className="form-field">
                                <label for="pickup-time">Time</label>
                                <select id="pickup-time">
                                    <option>Select your time</option>
                                </select>
                            </div>
                        </div>

                        <div className="radio-group">
                            <input type="radio" id="dropoff" name="rentalType"/>
                            <label for="dropoff">Drop - Off</label>
                        </div>

                        <div className="form-grid">
                            <div className="form-field">
                                <label for="dropoff-location">Место</label>
                                <select id="dropoff-location">
                                    <option>Выберите свой город</option>
                                </select>
                            </div>

                            <div className="form-field">
                                <label for="dropoff-date">Дата</label>
                                <select id="dropoff-date">
                                    <option>Выберите дату</option>
                                </select>
                            </div>

                            <div className="form-field">
                                <label for="dropoff-time">Время</label>
                                <select id="dropoff-time">
                                    <option>Выберите время</option>
                                </select>
                            </div>
                        </div>
                        </div>
                    </div>
                
            
                    <div className="payment-method">
                <div className="header">
                    <h2>Способ оплаты</h2>
                    <p>Введите данные для оплаты</p>
                    <span>Step 3 of 4</span>
                </div>

                <div className="form-section">
                    <div className="radio-group">
                        <input type="radio" id="credit-card" name="payment" checked/>
                        <label for="credit-card">Credit Card</label>
                        <img src= {visaImage} alt="Visa & Mastercard" class="card-logo"/>
                    </div>

                    <div className="form-grid">
                        <div className="form-field">
                            <label for="card-number">Номер карты</label>
                            <input type="text" id="card-number" placeholder="Card number"/>
                        </div>

                        <div className="form-field">
                            <label for="expiration-date">Срок действия</label>
                            <input type="text" id="expiration-date" placeholder="DD / MM / YY"/>
                        </div>

                        <div className="form-field">
                            <label for="card-holder">Имя держателя</label>
                            <input type="text" id="card-holder" placeholder="Card holder"/>
                        </div>

                        <div className="form-field">
                            <label for="cvc">CVC</label>
                            <input type="text" id="cvc" placeholder="CVC"/>
                        </div>
                    </div>
                    
                    <div className="radio-group">
                        <input type="radio" id="paypal" name="payment"/>
                        <label for="paypal">PayPal</label>
                        <img src={paypalImage} alt="PayPal" class="payment-logo"/>
                    </div>

                    <div className="radio-group">
                        <input type="radio" id="bitcoin" name="payment"/>
                        <label for="bitcoin">Bitcoin</label>
                        <img src={bitcoinImage} alt="Bitcoin" class="payment-logo"/>
                        </div>
                    </div>
                    </div>

                    <div className="confirmation-container">
                <div className="confirmation-header">
                    <h2>Confirmation</h2>
                    <p className="step-info">Step 4 of 4</p>
                </div>
                <p>We are getting to the end. Just a few clicks and your rental is ready!</p>

                <div className="checkbox-section">
                    <label className="checkbox-label">
                        <input
                            type="checkbox"
                            checked={marketingChecked}
                            onChange={handleMarketingChange}
                        />
                        I agree with sending marketing and newsletter emails. No spam, promised!
                    </label>

                    <label className="checkbox-label">
                        <input
                            type="checkbox"
                            checked={termsChecked}
                            onChange={handleTermsChange}
                        />
                        I agree with our terms and conditions and privacy policy.
                    </label>
                </div>

                <button className="rent-button" disabled={!termsChecked}>Rent Now</button>

                <div className="data-safe">
                    <div className="security-icon">
                        
                    </div>
                    <div>
                        <p>All your data are safe</p>
                        <p>We are using the most advanced security to provide you the best experience ever.</p>
                    </div>
                </div>
                    </div>
                </div>

                <div className="rental-summary-container">
                    <h2>Rental Summary</h2>
                    <p className="rental-description">
                        Prices may change depending on the length of the rental and the price of your rental car.
                    </p>

                    <div className="car-info">
                    {car?.images?.length > 0 ? (
                            <img src={STATIC + car.images[0]} alt="Nissan GT-R" className="car-image-billing" />
                        ) : (
                            <p>No image available</p>
                        )}  
                        <div className="car-details">
                            <h2>{car.brand} {car.model}</h2>
                            <div className="review-section">
                                <span className="stars">★★★★☆</span>
                                <span className="reviewer-info">440+ Reviewer</span>
                            </div>
                        </div>
                    </div>

                    <div className="price-details">
                        <div className="subtotal">
                            <span>Subtotal</span>
                            <span> {car.pricePerDay} BYN/day</span>
                        </div>
                        <div className="tax">
                            <span>Tax</span>
                            <span>$0</span>
                        </div>
                    </div>

                    <div className="promo-code-section">
                        <input
                            type="text"
                            placeholder="Apply promo code"
                           
                            className="promo-code-input"
                         
                        />
                        <button  className="apply-promo-button" >
                           
                        </button>
                    </div>

                    <div className="total-price">
                        <span>Total Rental Price</span>
                        <span>{car.pricePerDay} BYN/day</span>
                    </div>

                    <p className="price-disclaimer">Overall price and includes rental discount</p>
                </div>
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