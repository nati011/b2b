'use client';
// @ts-ignore
import TawkMessengerReact from '@tawk.to/tawk-messenger-react';

const TawkChat = () => {
    return (
        <div className='z-50 bg-red'>
            <p>HERE</p>
            <TawkMessengerReact
                propertyId={process.env.NEXT_PUBLIC_TAWK_PROPERTY_ID}
                widgetId={process.env.NEXT_PUBLIC_TAWK_WIDGET_ID}
            />
        </div>
    );
};

export default TawkChat;