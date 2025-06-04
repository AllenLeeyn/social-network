'use client';

import "../styles/globals.css";
import Navbar from '../components/Navbar';
import Footer from '../components/Footer';
import { WebSocketProvider } from '../contexts/WebSocketContext'
import { ActiveChatProvider } from "../contexts/ActiveChatContext";
import { NotificationsProvider } from '../contexts/NotificationsContext';
import { ToastContainer } from 'react-toastify';
import VantaBackground from '../components/VantaBackground';
import 'react-toastify/dist/ReactToastify.css';

export default function RootLayout({ children }) {
    
return (
    <html lang="en">
        <body>
            <VantaBackground/>
            <ToastContainer position="top-right" autoClose={1500} />
            <NotificationsProvider>
                <ActiveChatProvider>
                    <WebSocketProvider>
                        <Navbar />
                        <main>
                            {children}
                        </main>
                        <Footer />
                    </WebSocketProvider>
                </ActiveChatProvider>
            </NotificationsProvider>
        </body>
    </html>
    );
}
