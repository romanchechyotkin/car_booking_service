import { useSelector } from "react-redux";
import { Link } from "react-router-dom";
import { useEffect, useState } from "react";
import { axiosInstance } from "../../axios/axios";
import { CarPost } from "../CarPost/CarPost";
import "./feed.css";
import firstImage from '../../img/first.png';
import secondImage from '../../img/second.png';
import carImage from '../../img/car.png';
import firstCarImage from '../../img/car-first.png'

export const Feed = () => {
    const isVerified = useSelector((state) => state.user.isVerified);
    const isAuth = useSelector((state) => state.user.isAuth);
    const [cars, setCars] = useState([]);
    const [sortCriteria, setSortCriteria] = useState(null);
    

    const filters = [
        { id: null, label: "Новые" },
        { id: "prc.a", label: "Сначала дешевле" },
        { id: "prc.d", label: "Сначала дороже" },
        { id: "y.a", label: "Сначала старые" },
        { id: "y.d", label: "Сначала новые" }
    ];

    const fetchCars = async () => {
        try {
            const sortQuery = sortCriteria ? `?sort=${sortCriteria}` : "";
            const res = await axiosInstance.get(`/cars${sortQuery}`);
            if (res.data) {
                setCars(res.data);
            }
        } catch (e) {
            console.log(e);
        }
    };

    useEffect(() => {
        fetchCars();
    }, [sortCriteria]);

    return (
        <>  
            <div className="card-container">
                <div className="card-first">
                <div className="card-text">
                    <h1>Легко cдать!</h1>
                    <p>Сделайте фото авто, укажите прайс в час и получайте пассивный доход!</p>
                        <Link  to={`/createCar`} className="rent-button-first">
                            Сдать в аренду
                        </Link>
                    </div>
                    
                    <img src={firstCarImage} alt="White Sports Car" class="car-image"/>
                    
                </div>
                <div class="card">
                    <div className="card-text">
                    <h1>Легко арендовать!</h1>
                    <p>Выберите подходящую вам машину, дату и место для аренды и катайтесь!</p>
                    <a href="#" class="rent-button">Арендовать</a>
                    </div>
                
                  
                    <img src={carImage} alt="Silver Sports Car" class="car-image"/>
                </div>
        </div>

            {/* {isAuth && !isVerified && (
                <div>
                    <h2>You should verify yourself</h2>
                    <Link to={"/verify"} className="verify-button">Verify</Link>
                </div>
            )} */}
            <div className="main-content">
                    <div className="filter-column">
                    <h2>Фильтр</h2>
                    <ul className="filter-list">
                        {filters.map(filter => (
                            <li key={filter.id} className="filter-item">
                                <input
                                    type="radio"
                                    id={filter.id}
                                    name="filter"
                                    checked={sortCriteria === filter.id}
                                    onChange={() => setSortCriteria(filter.id)}
                                />
                                <label htmlFor={filter.id}>{filter.label}</label>
                            </li>
                        ))}
                    </ul>
                </div>

                <div className="posts">
                    {cars.length == 0 && <h1>Кажется, здесь пока что пусто..</h1>}
                    {cars && cars.map(c => (
                        <Link to={`/post/${c.car.id}`} key={c.car.id} className="link">
                            <div >
                                <CarPost car={c.car} user_id={c.user_id} />
                            </div>
                        </Link>
                    ))}
                </div>
                
            </div>
        </>
    );
};

