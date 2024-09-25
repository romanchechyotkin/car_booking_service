import React, { useState } from "react";
import "./ImageUploader.css"; 

export const ImageLoader = ({ files, setFiles }) => {
    const handleFilesChange = (e) => {
        const selectedFiles = Array.from(e.target.files);
        if (validateFiles(selectedFiles)) {
            setFiles([...files, ...selectedFiles]);
        }
    };

    const handleDrop = (e) => {
        e.preventDefault();
        if (files.length > 5) {
            alert("Можно загрузить не более 5-ти изображений.");
            return
        }
        
        const droppedFiles = Array.from(e.dataTransfer.files);
        if (validateFiles(droppedFiles)) {
            setFiles([...files, ...droppedFiles]);
        }
    };

    const validateFiles = (fileArray) => {
        const validTypes = ["image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"];
        const maxSize = 10 * 1024 * 1024; 
        let valid = true;

        fileArray.forEach((file) => {
            if (!validTypes.includes(file.type)) {
                alert("Неверный формат. Только JPEG, JPG, PNG, GIF, WEBP разрешены.");
                valid = false;
            }
            if (file.size > maxSize) {
                alert(".");
                valid = false;
            }
        });

        return valid;
    };

    const handleDragOver = (e) => {
        e.preventDefault();
    };

    return (
        <div className="image-uploader-container">
            <div
                className="upload-box"
                onDrop={handleDrop}
                onDragOver={handleDragOver}
            >
                <div className="upload-text">
                    <p>Выберите или перетащите фотографии в область</p>
                    <small>Форматы JPEG, JPG, PNG, WEBP или GIF, до 10 МБ каждый</small>
                </div>
                <input
                    type="file"
                    accept=".jpg,.jpeg,.png,.gif,.webp"
                    multiple
                    onChange={handleFilesChange}
                    className="file-input"
                />
                <button className="select-btn">Выбрать фотографии</button>
            </div>

            {files.length > 0 && (
                <div className="image-preview">
                    <h4>Добавленные фотографии:</h4>
                    <div className="preview-grid">
                        {files.map((file, index) => (
                            <div key={index} className="image-thumbnail">
                                <img src={URL.createObjectURL(file)} alt="Preview" />
                                <p>{file.name}</p>
                            </div>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
};
