"use client";
import "../styles/globals.css";

export default function CategoriesList( {categories} ) {

    return (
        <ul className="categories">
        {categories.map((cat, index) => (
            <li key={index} className="category-item">
            <strong>{cat}</strong>
            </li>
        ))}
        </ul>
    );
}





