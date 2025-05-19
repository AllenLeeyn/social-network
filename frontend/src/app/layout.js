'use client';

import "../styles/globals.css";
import Navbar from '../components/Navbar';
import Footer from '../components/Footer';
import { WebSocketProvider } from '../contexts/WebSocketContext'



export default function RootLayout({ children }) {
    
return (
    <html lang="en">
        <body>
            <WebSocketProvider>
                <Navbar />
                <main>
                    {children}
                </main>
                <Footer />
            </WebSocketProvider>
        </body>
    </html>
    );
}
