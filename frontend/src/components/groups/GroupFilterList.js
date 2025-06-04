// src/components/groups/GroupFilterList.js

import React from "react";
import "../../styles/SidebarSection.css"; 

export default function GroupFilterList({ filters, selectedFilter, onSelect }) {
    return (
        <ul className="group-filter-list">
        {filters.map((filter) => (
            <li
            key={filter.key}
            className={
                filter.key === selectedFilter 
                    ? "group-filter-item active"
                    : "group-filter-item"
            }
            
            onClick={() => onSelect(filter.key)}
            tabIndex={0}
            role="button"
        >
            {filter.label}
            </li>
        ))}
        </ul>
    );
}
