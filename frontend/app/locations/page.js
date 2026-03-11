'use client'

import dynamic from 'next/dynamic'

const LocationsMap = dynamic(() => import('@/components/LocationsMap'), {
  ssr: false,
  loading: () => (
    <div style={{ height: 'calc(100vh - 56px)', display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#888', fontSize: '14px' }}>
      Loading map...
    </div>
  ),
})

export default function LocationsPage() {
  return <LocationsMap />
}
