'use client';
import { useState } from 'react';
import '../styles/SidebarSection.css'; 

export default function SidebarSection({ title, children, defaultOpen = true }) {
    const [open, setOpen] = useState(defaultOpen);
    return (
        <section className="sidebar-section">
            <h2
            className="sidebar-section-title"
            onClick={() => setOpen(o => !o)}
            style={{ cursor: 'pointer' }}
            tabIndex={0}
            aria-expanded={open}
            >
            <span className="sidebar-arrow">{open ? '▼' : '►'} {title}</span>
        </h2>
        {open && children}
        </section>
    );
}
