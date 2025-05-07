// Placeholder - 
// If you need a global layout (applies to all pages), add /src/app/layout.jsx.

import React from 'react';
import "../styles/globals.css";
import Navbar from '../components/Navbar';
import Footer from '../components/Footer';

export default function RootLayout({ children }) {
    
return (
    <html lang="en">
        <body>
        <Navbar />
        <main>
            {children}
        </main>
        <Footer />
        </body>
    </html>
    );
}
