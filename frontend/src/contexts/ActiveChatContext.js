'use client';

import { createContext, useContext, useState } from "react";

const ActiveChatContext = createContext();

export function ActiveChatProvider( {children} ) {
    const [ activeChat, setActiveChat ] = useState(null);

    return (
        <ActiveChatContext.Provider value={{ activeChat, setActiveChat }}>
            {children}
        </ActiveChatContext.Provider>
    )
}

export function useActiveChat() {
    return useContext(ActiveChatContext);
}