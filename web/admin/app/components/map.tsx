// components/Map.tsx
'use client'

import { useEffect } from 'react'
import { MapContainer, TileLayer, Marker, useMap } from 'react-leaflet'
import L from 'leaflet'
import icon from "leaflet/dist/images/marker-icon.png";
import iconShadow from "leaflet/dist/images/marker-shadow.png";

interface MapProps {
    markerPosition: [number, number]
    useCurrentLocation: boolean
    setMarkerPosition: (position: [number, number]) => void
    setFormData: (callback: (prev: any) => any) => void
}


const MapEvents = ({ useCurrentLocation, setMarkerPosition, setFormData }: Omit<MapProps, 'markerPosition'>) => {
    const map = useMap()

    useEffect(() => {
        if (!useCurrentLocation) {
            map.on('click', (e: L.LeafletMouseEvent) => {
                const { lat, lng } = e.latlng
                setMarkerPosition([lat, lng])
                setFormData(prev => ({ ...prev, location: [lat, lng] }))
            })
        }

        return () => {
            map.off('click')
        }
    }, [map, useCurrentLocation])

    return null
}

const Map = ({ markerPosition, useCurrentLocation, setMarkerPosition, setFormData }: MapProps) => {
    let DefaultIcon = L.icon({
        // @ts-ignore
        iconUrl: icon,
        // @ts-ignore
        shadowUrl: iconShadow,
    });

    L.Marker.prototype.options.icon = DefaultIcon;
    return (
        <MapContainer
            center={markerPosition}
            zoom={13}
            className="h-full w-full"
        >
            <TileLayer
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />
            <Marker position={markerPosition}
                icon={L.divIcon({
                    iconSize: [50, 50],
                    iconAnchor: [50 / 2, 50 + 9],
                    className: "text-6xl",
                    html: "📍",
                })}
            />
            <MapEvents
                useCurrentLocation={useCurrentLocation}
                setMarkerPosition={setMarkerPosition}
                setFormData={setFormData}
            />
        </MapContainer>
    )
}

export default Map