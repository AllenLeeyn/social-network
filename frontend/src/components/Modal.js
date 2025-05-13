import React from 'react';
import '../styles/Modal.css';

export default function Modal({ onClose, title, children }) {
    return (
        <div className="modal-overlay" onClick={onClose}>
            <div className="modal-content" onClick={e => e.stopPropagation()}>
                <header className="modal-header">
                <h2>{title}</h2>
                <button className="modal-close-btn" onClick={onClose} aria-label="Close modal">&times;</button>
                </header>
                <div className="modal-body">
                {children}
                </div>
            </div>
        </div>
    );
}
