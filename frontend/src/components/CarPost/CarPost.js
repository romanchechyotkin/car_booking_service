import "./carpost.css"
import {STATIC} from "../../axios/axios";
import { useState } from "react";
import Lottie from 'react-lottie';
import heartAnimation from "../../animations/heart1.json"

export const CarPost = (props) => {
    const { car, user_id } = props;
    const [isLiked, setIsLiked] = useState(false);
    const [isStopped, setIsStopped] = useState(true); // анимация изначально остановлена

    const handleLikeClick = () => {
        setIsLiked(!isLiked);
        setIsStopped(false); // запускаем анимацию при нажатии
    };

    const defaultOptions = {
        loop: false, // анимация не повторяется
        autoplay: false, // не запускается автоматически
        animationData: heartAnimation,
        rendererSettings: {
            preserveAspectRatio: 'xMidYMid slice',
        },
    };

    return (
        <div className="car-card">
            <div className="car-card-header">
                <h2>{car.brand} {car.model}</h2>
                <div className="favorite-icon" onClick={handleLikeClick}>
                    <Lottie 
                        options={defaultOptions}
                        height={40}
                        width={50}
                        isStopped={isStopped} // контролируем остановку анимации
                        eventListeners={[
                            {
                                eventName: 'complete',
                                callback: () => setIsStopped(true), // останавливаем после завершения
                            },
                        ]}
                    />
                </div>
            </div>
            <p className="car-type">{car.type}</p>
            <img
                src={STATIC + car.images[0]}
                alt={`${car.brand} ${car.model}`}
                className="car-image-box"
            />
            <div className="car-card-footer">
                <div className="price-feed">${car.pricePerDay}/day</div> 
                <button className="car-rent-button">Арендовать</button>
            </div>
        </div>
    );
};