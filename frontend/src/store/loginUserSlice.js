import { createSlice } from '@reduxjs/toolkit'

const initialState = {
    isAuth: false,
    isVerified: false,
    user: null,
    role: "",
}

export const userSlice = createSlice({
    name: 'user',
    initialState,
    reducers: {
        setUser: (state, action) => {
            state.user = action.payload
        },
        setIsAuth: (state) => {
            state.isAuth = true
        },
        setIsVerified: (state) => {
            state.isVerified = true
        },
        setRole: (state, action) => {
            state.role = action.payload
        },
        logout: (state) => {
            state.isAuth = false
            state.user = null
            state.role = ""
            state.isVerified = false
        },
    },
})

export const {actions: userActions} = userSlice;
export const {reducer: userReducer} = userSlice;