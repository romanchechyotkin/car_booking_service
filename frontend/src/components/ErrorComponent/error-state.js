import React, { useEffect } from 'react';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

const ToastError = ({ error }) => {
    useEffect(() => {
        if (error != "") {
            toast.error(error);
        }
    }, [error]);

    return (
        <ToastContainer
            position="top-center" 
            autoClose={5000} 
            hideProgressBar
            newestOnTop
            closeOnClick
            rtl={false}
            pauseOnFocusLoss
            draggable
            pauseOnHover
            style={{ 
                zIndex: 9999,
                width: 500,
                display: 'flex',
                flexDirection: 'column'
             }} 
x           />
    );
};

export default ToastError;
