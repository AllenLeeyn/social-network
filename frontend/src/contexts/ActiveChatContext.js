'use client';

import { createContext, useContext, useState, useEffect } from "react";

const ActiveChatContext = createContext();

export function ActiveChatProvider( {children} ) {
    const [ activeChat, setActiveChat ] = useState(null);

    // Add this useEffect for debugging:
    // useEffect(() => {
    //     console.log("Active chat changed:", activeChat);
    // }, [activeChat]);

    return (
        <ActiveChatContext.Provider value={{ activeChat, setActiveChat }}>
            {children}
        </ActiveChatContext.Provider>
    )
}

export function useActiveChat() {
    return useContext(ActiveChatContext);
}