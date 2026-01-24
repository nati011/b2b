'use client';
import { useEffect } from 'react';

const TawkChat = () => {
    useEffect(() => {
        const propertyId = process.env.NEXT_PUBLIC_TAWK_PROPERTY_ID;
        const widgetId = process.env.NEXT_PUBLIC_TAWK_WIDGET_ID;

        // Only initialize Tawk if both environment variables are set
        if (!propertyId || !widgetId) {
            console.warn('Tawk chat not initialized: Missing NEXT_PUBLIC_TAWK_PROPERTY_ID or NEXT_PUBLIC_TAWK_WIDGET_ID');
            return;
        }

        // Load Tawk script dynamically
        const script = document.createElement('script');
        script.async = true;
        script.src = `https://embed.tawk.to/${propertyId}/${widgetId}`;
        script.charset = 'UTF-8';
        script.setAttribute('crossorigin', '*');
        
        // Check if script already exists
        const existingScript = document.querySelector(`script[src*="tawk.to/${propertyId}"]`);
        if (existingScript) {
            return;
        }

        document.head.appendChild(script);

        return () => {
            // Cleanup: remove script on unmount if needed
            const scriptToRemove = document.querySelector(`script[src*="tawk.to/${propertyId}"]`);
            if (scriptToRemove) {
                scriptToRemove.remove();
            }
        };
    }, []);

    return null;
};

export default TawkChat;