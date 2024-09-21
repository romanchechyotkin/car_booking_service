import { useSelector } from "react-redux";
import { Link } from "react-router-dom";
import { useEffect, useState } from "react";
import { axiosInstance } from "../../axios/axios";
import { CarPost } from "../CarPost/CarPost";
import "./billing.css";

export const BillingInfo = () => {
    return (
        <div>
            <h1>
                Billing System
            </h1>
        </div>
    )
}