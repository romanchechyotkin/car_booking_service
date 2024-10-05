import { Link } from "react-router-dom";
import { useEffect, useState } from "react";
import { axiosInstance } from "../../axios/axios";
import { CarPost } from "../CarPost/CarPost";
import "./feed.css";
import carImage from '../../img/car.png';
import firstCarImage from '../../img/car-first.png'
import strelkiImage from '../../img/feed/strelki.png'


export const Feed = () => {
    const [cars, setCars] = useState([]);
    const [sortCriteria, setSortCriteria] = useState(null);
    
    const filters = [
        { id: "new", label: "Новые" },
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
                setCars(res.data.cars);
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
                        <Link to={`/createCar`} className="rent-button-first">
                            Сдать в аренду
                        </Link>
                    </div>
                    <img src={firstCarImage} alt="White Sports Car" className="car-image"/>
                </div>
                
                <div className="card">
                    <div className="card-text">
                        <h1>Легко арендовать!</h1>
                        <p>Выберите подходящую вам машину, дату и место для аренды и катайтесь!</p>
                        <a href="#main-content" className="rent-button">Арендовать</a> {}
                    </div>
                    <img src={carImage} alt="Silver Sports Car" className="car-image"/>
                </div>
            </div>
    
            <div id="main-content" className="main-content">
                {/* <div className="filter-column">
                    <h1>Фильтр</h1>
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
                </div> */}

                    <div className="main-header-rental-container">
                        <div className="main-header-rental-start">
                        <div className="main-header-radio-group">
                            <input type="radio" id="start" name="rental" checked/>
                            <label for="start">Старт аренды</label>
                        </div>

                            <div class="main-header-rental-fields">
                                <div class="main-header-field">
                                    <label>Место</label>
                                    <select>
                                        <option>Выбери город</option>
                                    </select>
                                </div>
                                <div className="main-header-divider"></div>
                                <div class="main-header-field">
                                    <label>Дата</label>
                                    <select>
                                        <option>Выбери дату</option>
                                    </select>
                                </div>
                                <div className="main-header-divider"></div>
                                <div class="main-header-field">
                                    <label>Время</label>
                                    <select>
                                        <option>Выбери время</option>
                                    </select>
                                </div>
                            </div>
                        </div>

                        <div className="main-header-swap-button">
                            <button>
                                <img src={strelkiImage}/>
                            </button>
                        </div>

                        <div class="main-header-rental-end">
                        <div className="main-header-radio-group">
                            <input type="radio" id="end" name="rental"/>
                            <label for="end">Конец аренды</label>
                            </div>
                            <div class="main-header-rental-fields">
                                <div class="main-header-field">
                                    <label>Место</label>
                                    <select>
                                        <option>Выбери город</option>
                                    </select>
                                </div>
                                <div className="main-header-divider"></div>
                                <div class="main-header-field">
                                    <label>Дата</label>
                                    <select>
                                        <option>Выбери дату</option>
                                    </select>
                                </div>
                                <div className="main-header-divider"></div>
                                <div class="main-header-field">
                                    <label>Время</label>
                                    <select>
                                        <option>Выбери время</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                    </div>

    
                <div className="posts">
                    {cars.length === 0 && <h1>Кажется, здесь пока что пусто..</h1>}
                    {cars && cars.map(c => (
                        <Link to={`/post/${c.car.id}`} key={c.car.id} className="link">
                            <div>
                                <CarPost car={c.car} user_id={c.user_id} />
                            </div>
                        </Link>
                    ))}
                </div>
            </div>
        </>
    );
    
};

