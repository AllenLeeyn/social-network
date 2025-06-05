"use client";
import { useEffect, useRef } from "react";
import * as THREE from "three";
import NET from "vanta/dist/vanta.net.min";

export default function VantaBackground({ children }) {
    const vantaRef = useRef(null);
    const vantaEffect = useRef(null);

    useEffect(() => {
        if (!vantaEffect.current) {
        vantaEffect.current = NET({
        el: vantaRef.current,
        THREE,
        mouseControls: true,
        touchControls: false,
        gyroControls: false,

        points: 2.0,              // default is 12.0, lower = fewer dots
        maxDistance: 30.0,        // default is 20.0, lower = fewer connections
        spacing: 20.0,    

        minHeight: 50.0,
        minWidth: 50.0,
        scale: 0.25,
        scaleMobile: 1.0,
        color: 0x33FFFF,
        backgroundColor: 0x111111, // optional: set a background color
        });
    }
    return () => {
        if (vantaEffect.current) {
        vantaEffect.current.destroy();
        vantaEffect.current = null;
        }
    };
    }, []);

    return (
    <div ref={vantaRef} style={{ position: "fixed", inset: 0, zIndex: -1 }}>
        {children}
    </div>
    );
}