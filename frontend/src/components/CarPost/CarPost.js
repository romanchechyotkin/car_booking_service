import "./carpost.css"
import {STATIC} from "../../axios/axios";

export const CarPost = (props) => {
    const {car, user_id} = props

    return (


        // <div key={car.id} className={"post"}>
        //     {/* <div>{car.brand} {car.model}</div>
        //     <div className={"post_images"}>
        //         <img src={STATIC+car.images[0]} alt="img"/>
        //     </div>
        //     <div>{car.pricePerDay}$</div> */}
        

        //     <div className="card">

        //     </div>
        // </div>

        <div className="car-card">
            <div className="car-card-header">
                <h3>{car.brand} {car.model}</h3>
                <span className="favorite-icon">❤️</span>
            </div>
            <p className="car-type">{car.type}</p>
            <img
                src={STATIC+car.images[0]}
                alt={`${car.brand} ${car.model}`}
            
                className="car-image-box"
            />
            
            <div className="car-card-footer">
                <div className="price">${car.pricePerDay}/day</div> 
                <button className="rent-button">Арендовать</button>
            </div>
        </div>
    )
}