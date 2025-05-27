'use client';

import React from 'react';

export default function NotificationFilterList({ filters, selectedFilter, onSelect }) {
    return (
        <ul className="notification-filter-list">
            {filters.map(filter => (
                <li
                    key={filter.key}
                    className={`notification-filter-item${selectedFilter === filter.key ? ' active' : ''}`}
                    onClick={() => onSelect(filter.key)}
                    style={{ cursor: 'pointer'}}
                >
                    {filter.label}
                </li>
            ))}
        </ul>
    );
}
