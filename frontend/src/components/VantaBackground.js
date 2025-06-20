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

        points: 2,
        maxDistance: 70,
        spacing: 60,    

        minHeight: 200.0,
        minWidth: 200.0,
        scale: 2,
        scaleMobile: 1.0,
        color: 0x33FFFF,
        backgroundColor: 0x111111,
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