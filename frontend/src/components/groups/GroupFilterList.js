// src/components/groups/GroupFilterList.js

import React from "react";
import "../../styles/SidebarSection.css"; // Optional: for consistent sidebar styling

export default function GroupFilterList({ filters, selectedFilter, onSelect }) {
    return (
        <ul className="group-filter-list">
        {filters.map((filter) => (
            <li
            key={filter.key}
            className={filter.key === selectedFilter ? "selected" : ""}
            >
            <button
                type="button"
                onClick={() => onSelect(filter.key)}
                className="group-filter-btn"
                aria-current={filter.key === selectedFilter ? "page" : undefined}
            >
                {filter.label}
            </button>
            </li>
        ))}
        </ul>
    );
}
