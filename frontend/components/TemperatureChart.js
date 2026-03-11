'use client'

import { useState, useEffect, useRef } from 'react'
import {
  AreaChart, Area, XAxis, YAxis, Tooltip,
  ResponsiveContainer, CartesianGrid,
} from 'recharts'
import styles from './TemperatureChart.module.css'

const LOOKBACKS = [
  { label: '1D', value: '24h' },
  { label: '1W', value: '1wk' },
  { label: '1M', value: '1mo' },
  { label: '3M', value: '3mo' },
]

const BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL

function formatXAxis(timeStr, lookback) {
  const d = new Date(timeStr)
  if (lookback === '24h') {
    return d.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false, timeZone: 'UTC' })
  }
  if (lookback === '1wk') {
    return d.toLocaleDateString('en-US', { weekday: 'short', timeZone: 'UTC' })
  }
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', timeZone: 'UTC' })
}

function formatTooltipTime(timeStr, lookback) {
  const d = new Date(timeStr)
  if (lookback === '24h') {
    return d.toLocaleString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false, timeZone: 'UTC' }) + ' UTC'
  }
  return d.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric', timeZone: 'UTC' })
}

function CustomTooltip({ active, payload, lookback }) {
  if (!active || !payload?.length) return null
  const { time, value } = payload[0].payload
  return (
    <div className={styles.tooltip}>
      <span className={styles.tooltipTemp}>{value.toFixed(2)}°C</span>
      <span className={styles.tooltipTime}>{formatTooltipTime(time, lookback)}</span>
    </div>
  )
}

export default function TemperatureChart({ locationId }) {
  const [chartData, setChartData] = useState([])
  const [lookback, setLookback] = useState('24h')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [sseConnected, setSseConnected] = useState(false)

  // Refs so the SSE handler always reads the latest values without reconnecting
  const lookbackRef = useRef(lookback)
  const liveAppendedRef = useRef(false)
  const eventSourceRef = useRef(null)

  useEffect(() => {
    lookbackRef.current = lookback
  }, [lookback])

  // Fetch historical data whenever locationId or lookback changes
  useEffect(() => {
    setLoading(true)
    setError(null)

    fetch(`${BASE_URL}/temperature/${locationId}?lookback=${lookback}`)
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`)
        return res.json()
      })
      .then((json) => {
        setChartData((json.data ?? []).map((d) => ({ time: d.time, value: d.val })))
        setLoading(false)
      })
      .catch((err) => {
        setError(err.message)
        setLoading(false)
      })
  }, [locationId, lookback])

  // SSE — reconnects only when locationId changes; lookback logic handled via ref
  useEffect(() => {
    eventSourceRef.current?.close()

    const es = new EventSource(`${BASE_URL}/streams/?location_id=${locationId}`)
    eventSourceRef.current = es

    es.onopen = () => setSseConnected(true)

    es.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data)
        const point = { time: msg.time, value: parseFloat(msg.val) }

        setChartData((prev) => {
          if (lookbackRef.current === '24h') {
            return [...prev, point]
          }
          // For longer lookbacks: keep updating only the last point
          return [...prev.slice(0, -1), point]
        })
      } catch (e) {
        console.error(`[${locationId}] SSE parse error:`, e)
      }
    }

    es.onerror = () => setSseConnected(false)

    return () => es.close()
  }, [locationId])

  const currentTemp = chartData.at(-1)?.value ?? null
  const startTemp = chartData[0]?.value ?? null
  const delta = currentTemp !== null && startTemp !== null ? currentTemp - startTemp : null
  const isPositive = delta === null || delta >= 0
  const tickInterval = Math.max(1, Math.floor(chartData.length / 7))

  return (
    <div className={styles.card}>
      {/* Header */}
      <div className={styles.header}>
        <div className={styles.locationRow}>
          <span className={styles.location}>{locationId}</span>
          <span className={`${styles.sseDot} ${sseConnected ? styles.sseLive : styles.sseDown}`} title={sseConnected ? 'Live' : 'Disconnected'} />
        </div>
        <div className={styles.tempRow}>
          {currentTemp !== null && (
            <span className={styles.currentTemp}>{currentTemp.toFixed(2)}°C</span>
          )}
          {delta !== null && (
            <span className={`${styles.delta} ${isPositive ? styles.positive : styles.negative}`}>
              {isPositive ? '+' : ''}{delta.toFixed(2)}°C
            </span>
          )}
        </div>
      </div>

      {/* Chart */}
      <div className={styles.chartWrap}>
        {loading && <div className={styles.placeholder}>Loading...</div>}
        {error && <div className={styles.placeholderError}>{error}</div>}
        {!loading && !error && (
          <ResponsiveContainer width="100%" height={180}>
            <AreaChart data={chartData} margin={{ top: 4, right: 4, left: -24, bottom: 0 }}>
              <defs>
                <linearGradient id={`grad-${locationId}`} x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor={isPositive ? '#00c805' : '#ff4444'} stopOpacity={0.25} />
                  <stop offset="95%" stopColor={isPositive ? '#00c805' : '#ff4444'} stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid vertical={false} stroke="#2a2a2a" />
              <XAxis
                dataKey="time"
                tickFormatter={(v) => formatXAxis(v, lookback)}
                interval={tickInterval}
                tick={{ fill: '#888', fontSize: 11 }}
                axisLine={false}
                tickLine={false}
              />
              <YAxis
                domain={['auto', 'auto']}
                tick={{ fill: '#888', fontSize: 11 }}
                axisLine={false}
                tickLine={false}
                tickFormatter={(v) => `${v}°`}
              />
              <Tooltip content={<CustomTooltip lookback={lookback} />} />
              <Area
                type="monotone"
                dataKey="value"
                stroke={isPositive ? '#00c805' : '#ff4444'}
                strokeWidth={1.5}
                fill={`url(#grad-${locationId})`}
                dot={false}
                isAnimationActive={false}
              />
            </AreaChart>
          </ResponsiveContainer>
        )}
      </div>

      {/* Lookback buttons */}
      <div className={styles.lookbacks}>
        {LOOKBACKS.map((lb) => (
          <button
            key={lb.value}
            className={`${styles.lbBtn} ${lookback === lb.value ? styles.lbActive : ''}`}
            onClick={() => setLookback(lb.value)}
          >
            {lb.label}
          </button>
        ))}
      </div>
    </div>
  )
}
