import {axiosInstance} from "../../axios/axios";
import {userActions} from "../../store/loginUserSlice";

export const loginUser = async (email, password, navigate, dispatch) => {
    try {
        const res = await axiosInstance.post("/auth/login", JSON.stringify(
            {
                "email": email,
                "password": password,
            }
        ));

        console.log(res);

        localStorage.setItem('access_token', JSON.stringify(res.data.access_token));
        localStorage.setItem('refresh_token', JSON.stringify(res.data.refresh_token));
        localStorage.setItem('user', JSON.stringify(res.data.user));

        dispatch(userActions.setUser(res.data.user));
        dispatch(userActions.setRole(res.data.user.role));
        dispatch(userActions.setIsAuth());

        if (res.data.user.is_verified === true) {
            dispatch(userActions.setIsVerified());
        }

        navigate("/feed");
    } catch (e) {
        console.log(e);
    }
};
