import {configureStore} from '@reduxjs/toolkit'
import {userReducer} from "./loginUserSlice";

export const store = configureStore({
    reducer: {
        user: userReducer,
    },
})