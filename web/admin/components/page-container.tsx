import React from 'react';

interface PageContainerProps {
    children: React.ReactNode;
    className?: string;
}

const PageContainer: React.FC<PageContainerProps> = ({ children, className = '' }) => {
    return (
        <div className={`w-full max-w-none ${className}`}>
            {children}
        </div>
    );
};

export default PageContainer;
