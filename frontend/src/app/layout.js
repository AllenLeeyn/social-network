'use client';

import "../styles/globals.css";
import Navbar from '../components/Navbar';
import Footer from '../components/Footer';
import { WebSocketProvider } from '../contexts/WebSocketContext'
import { ActiveChatProvider } from "../contexts/ActiveChatContext";
import { NotificationsProvider } from '../contexts/NotificationsContext';



export default function RootLayout({ children }) {
    
return (
    <html lang="en">
        <body>
            <NotificationsProvider>
                <WebSocketProvider>
                    <ActiveChatProvider>
                        <Navbar />
                        <main>
                            {children}
                        </main>
                        <Footer />
                    </ActiveChatProvider>
                </WebSocketProvider>
            </NotificationsProvider>
        </body>
    </html>
    );
}
