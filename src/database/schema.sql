CREATE DATABASE IF NOT EXISTS CustomerDB;
USE CustomerDB;

CREATE TABLE IF NOT EXISTS Contact (
    id INT PRIMARY KEY,
    phoneNumber VARCHAR(255),
    email VARCHAR(255),
    linkedId INT,
    linkPrecedence ENUM('secondary', 'primary'),
    createdAt DATETIME,
    updatedAt DATETIME,
    deletedAt DATETIME
);
