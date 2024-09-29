import "./carpost.css"
import {STATIC} from "../../axios/axios";
import { useState } from "react";
import Lottie from 'react-lottie';
import heartAnimation from "../../animations/heart1.json"

export const CarPost = (props) => {
    const { car, user_id } = props;
    const [isLiked, setIsLiked] = useState(false);
    const [isStopped, setIsStopped] = useState(true); 

    const handleLikeClick = () => {
        setIsLiked(!isLiked);
        setIsStopped(false); 
    };

    const defaultOptions = {
        loop: false, 
        autoplay: false, 
        animationData: heartAnimation,
        rendererSettings: {
            preserveAspectRatio: 'xMidYMid slice',
        },
    };

    return (
        <div className="car-card">
            <div className="car-card-header">
                <h2>{car.brand} {car.model}</h2>
                <button className="favorite-icon" onClick={(event) => {
                    event.stopPropagation();
                    event.preventDefault();
                   
                }}>
                    <Lottie 
                        options={defaultOptions}
                        isStopped={isStopped}
                        width={50}
                        height={50}
                        eventListeners={[
                            {
                                eventName: 'complete',
                                callback: () => setIsStopped(true), 
                            },
                        ]}
                    />
                </button>
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