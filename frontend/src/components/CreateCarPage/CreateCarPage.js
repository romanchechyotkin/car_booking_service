import { useState } from "react";
import { axiosInstance } from "../../axios/axios";
import { ImageLoader } from "../ImageLoader/ImageLoader";
import "./CreateCarPage.css";
import {useNavigate} from "react-router-dom";

export const CreateCarPage =() => {
    const [carId, setCarId] = useState("");
    const [brand, setBrand] = useState("");
    const [model, setModel] = useState("");
    const [price, setPrice] = useState("");
    const [year, setYear] = useState("");
    const [photo, setPhoto] = useState([]);
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();

        const formData = new FormData();
        formData.append("id", carId);
        formData.append("brand", brand);
        formData.append("model", model);
        formData.append("price", price);
        formData.append("year", year);

        photo.forEach((file, index) => {
            formData.append(`image`, file);
        });

        try {
            setLoading(true);
            
            const response = await axiosInstance.post("/cars", formData, {
                headers: {
                    "Content-Type": "multipart/form-data"
                }
            });

            console.log("машина загружена:", response.data);
            setLoading(false);
            
            navigate("/feed");
        } catch (error) {
            console.error("Error:", error);
            setLoading(false);
        }
    };

    return (
        <div className="create-post-container">
            <h2>Подать объявление</h2>
            <form onSubmit={handleSubmit} className="post-form">
                <div className="form-field">
                    <label htmlFor="id">Номера автомобиля</label>
                    <input
                        type="text"
                        id="id"
                        value={carId}
                        onChange={(e) => setCarId(e.target.value)}
                        placeholder="Ввести номера авто"
                        required
                    />
                </div>

                <div className="form-field">
                    <label htmlFor="Бренд машины">Бренд</label>
                    <input
                        type="text"
                        id="brand"
                        value={brand}
                        onChange={(e) => setBrand(e.target.value)}
                        placeholder="Ввести бренд"
                        required
                    />
                </div>

                <div className="form-field">
                    <label htmlFor="model">Модель</label>
                    <input
                        type="text"
                        id="model"
                        value={model}
                        onChange={(e) => setModel(e.target.value)}
                        placeholder="Ввести модель"
                        required
                    />
                </div>

                <div className="form-field">
                    <label htmlFor="price">Цена в час (BYN)</label>
                    <input
                        type="number"
                        id="price"
                        value={price}
                        onChange={(e) => setPrice(e.target.value)}
                        placeholder="Ввести цену"
                        required
                    />
                </div>

                <div className="form-field">
                    <label htmlFor="year">Год выпуска</label>
                    <input
                        type="number"
                        id="year"
                        value={year}
                        onChange={(e) => setYear(e.target.value)}
                        placeholder="Ввести год выпуска"
                        required
                    />
                </div>

                <ImageLoader files={photo} setFiles={setPhoto} />

                <button type="submit" className="submit-btn" disabled={loading}>
                    {loading ? "Загружается" : "Создать"}
                </button>
            </form>
        </div>)
}