import React,  { useEffect } from 'react';
import '../styles/Modal.css';

export default function Modal({ onClose, title, children }) {

    // Lock background scroll when modal is open
    useEffect(() => {
        document.body.style.overflow = 'hidden';
        return () => {
            document.body.style.overflow = '';
        };
    }, []);

    return (
        <div className="modal-overlay" onClick={onClose}>
            <div className='modal-inner-box'>
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
        </div>
    );
}
