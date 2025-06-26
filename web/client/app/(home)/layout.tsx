import React from 'react';

export default function OverViewLayout({
    hero,
    product_grid,
    why_choose_us,
    testimonials
}: {
    hero: React.ReactNode;
    product_grid: React.ReactNode;
    why_choose_us: React.ReactNode;
    testimonials: React.ReactNode;
}) {
    return (
        <>
            {hero}
            {product_grid}
            {why_choose_us}
            {testimonials}
        </>

    );
}
