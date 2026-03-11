'use client'

import { useEffect, useState, useRef, useMemo } from 'react'
import { MapContainer, TileLayer, Polygon, Marker, Tooltip } from 'react-leaflet'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import styles from './LocationsMap.module.css'

const BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL

const CENTER   = [34.2, -111.9]
const ZOOM     = 8
const MIN_ZOOM = 7
const MAX_ZOOM = 14

// Keeps the user within a generous Arizona-area bounding box so the
// fog-of-war edge stays visible but they can't wander too far.
const MAX_BOUNDS = [
  [29.5, -116.8],
  [38.6, -107.0],
]

// Large outer ring covering the renderable Web-Mercator world, used as the
// base layer of the fog-of-war mask polygon.
const WORLD_RING = [
  [-85.051, -180],
  [ 85.051, -180],
  [ 85.051,  180],
  [-85.051,  180],
]

// Simplified Arizona state boundary [lat, lng].
// Western border follows the Colorado River; all other borders are approximate.
const ARIZONA_RING = [
  [37.004, -114.050], // NW
  [37.004, -109.045], // NE (Four Corners)
  [31.332, -109.045], // SE
  [31.335, -111.072], // S border jog
  [32.497, -114.815], // SW (near Yuma)
  [32.718, -114.720], // Colorado River – Yuma
  [33.426, -114.548], // Colorado River
  [34.479, -114.332], // Colorado River – Parker
  [35.295, -114.727], // Colorado River – Needles area
  [35.796, -114.640], // Colorado River
  [36.154, -114.300], // near Hoover Dam
  [36.841, -114.050], // NW approach
]

function createMarkerIcon() {
  return L.divIcon({
    className: 'mkr-wrapper',
    html: '<div class="mkr"><div class="mkr-dot"></div><div class="mkr-ring"></div></div>',
    iconSize: [28, 28],
    iconAnchor: [14, 14],
    tooltipAnchor: [16, -8],
  })
}

function LocationMarker({ location, temp, pulseCount }) {
  const leafletMarkerRef = useRef(null)
  const icon = useMemo(() => createMarkerIcon(), [])

  // Imperatively add/remove the pulse class on the Leaflet DOM element so we
  // don't need to recreate the icon (which would cause a flicker).
  useEffect(() => {
    if (pulseCount === 0 || !leafletMarkerRef.current) return
    const el = leafletMarkerRef.current.getElement()
    if (!el) return
    // Remove first to restart the animation if events arrive back-to-back.
    el.classList.remove('mkr-pulsing')
    void el.offsetWidth // force reflow
    el.classList.add('mkr-pulsing')
    const timer = setTimeout(() => el.classList.remove('mkr-pulsing'), 900)
    return () => clearTimeout(timer)
  }, [pulseCount])

  return (
    <Marker
      position={[location.latitude, location.longitude]}
      icon={icon}
      eventHandlers={{ add: (e) => { leafletMarkerRef.current = e.target } }}
    >
      <Tooltip direction="top" offset={[0, -6]}>
        <div className={styles.tooltip}>
          <span className={styles.tooltipName}>{location.location_id}</span>
          {temp !== null
            ? <span className={styles.tooltipTemp}>{temp.toFixed(2)}°C</span>
            : <span className={styles.tooltipWaiting}>waiting for data…</span>
          }
        </div>
      </Tooltip>
    </Marker>
  )
}

export default function LocationsMap() {
  const [locations, setLocations]     = useState([])
  const [temps, setTemps]             = useState({})
  const [pulseCounts, setPulseCounts] = useState({})

  // Fetch device metadata on mount
  useEffect(() => {
    fetch(`${BASE_URL}/locations`)
      .then((res) => res.json())
      .then((data) => {
        setLocations(data)
        const initial = {}
        data.forEach((loc) => { initial[loc.location_id] = 0 })
        setPulseCounts(initial)
      })
      .catch((err) => console.error('[LocationsMap] Failed to fetch locations:', err))
  }, [])

  // Open one SSE connection per location once the metadata is loaded
  useEffect(() => {
    if (!locations.length) return

    const sources = locations.map((loc) => {
      const es = new EventSource(`${BASE_URL}/streams/?location_id=${loc.location_id}`)

      es.onmessage = (event) => {
        try {
          const msg = JSON.parse(event.data)
          setTemps((prev)       => ({ ...prev, [loc.location_id]: parseFloat(msg.val) }))
          setPulseCounts((prev) => ({ ...prev, [loc.location_id]: (prev[loc.location_id] ?? 0) + 1 }))
        } catch (e) {
          console.error(`[LocationsMap] SSE parse error (${loc.location_id}):`, e)
        }
      }

      es.onerror = () => {} // EventSource will auto-retry; suppress console noise

      return es
    })

    return () => sources.forEach((es) => es.close())
  }, [locations])

  return (
    <div className={styles.container}>
      <MapContainer
        center={CENTER}
        zoom={ZOOM}
        minZoom={MIN_ZOOM}
        maxZoom={MAX_ZOOM}
        maxBounds={MAX_BOUNDS}
        maxBoundsViscosity={1.0}
        scrollWheelZoom
        className={styles.map}
      >
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>'
          url="https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png"
        />

        {/* Fog of war: world polygon with Arizona punched out as a hole.
            fillRule "evenodd" makes the inner ring render as a transparent hole. */}
        <Polygon
          positions={[WORLD_RING, ARIZONA_RING]}
          pathOptions={{
            fillColor: '#0f0f0f',
            fillOpacity: 0.65,
            stroke: false,
            fillRule: 'evenodd',
          }}
        />

        {locations.map((loc) => (
          <LocationMarker
            key={loc.location_id}
            location={loc}
            temp={temps[loc.location_id] ?? null}
            pulseCount={pulseCounts[loc.location_id] ?? 0}
          />
        ))}
      </MapContainer>
    </div>
  )
}
