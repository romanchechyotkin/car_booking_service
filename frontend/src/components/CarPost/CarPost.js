import "./carpost.css"
import {STATIC} from "../../axios/axios";

export const CarPost = (props) => {
    const {car, user_id} = props

    return (
        <div key={car.id} className={"post"}>
            <div>{car.brand} {car.model}</div>
            <div className={"post_images"}>
                <img src={STATIC+car.images[0]} alt="img"/>
            </div>
            <div>{car.pricePerDay}$</div>
        </div>
    )
}