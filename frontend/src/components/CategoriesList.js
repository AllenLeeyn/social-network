"use client";
import "../styles/globals.css";

export default function CategoriesList( {categories, loading, error, onCategoryClick} ) {

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error: {error}</div>;
    if (!categories || categories.length === 0) return <div>No categories found.</div>;

    return (
        <ul className="categories">
        {categories.map((cat, index) => (
            <li 
                key={cat.id || index} 
                className="category-item"
                onClick={() => onCategoryClick && onCategoryClick(cat.name)}
                style={{cursor: onCategoryClick ? 'pointer' : 'default' }}
            >
                <strong>{cat.name}</strong>
            </li>
        ))}
        </ul>
    );
}





