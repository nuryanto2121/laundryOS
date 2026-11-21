<template>
  <div class="space-y-2">
    <div class="flex items-center justify-between">
      <label class="block text-xs font-bold text-slate-700 flex items-center gap-1.5">
        <MapPin class="w-3.5 h-3.5 text-brand-600" />
        <span>Pilih Titik Lokasi Cabang & Geofence (Peta)</span>
      </label>
      <button
        type="button"
        class="inline-flex items-center gap-1 text-[11px] font-bold text-brand-600 hover:text-brand-700 bg-brand-50 hover:bg-brand-100 px-2 py-1 rounded-lg transition-colors cursor-pointer"
        @click="getCurrentLocation"
      >
        <Navigation class="w-3 h-3" />
        <span>Gunakan Lokasi GPS Saya</span>
      </button>
    </div>

    <!-- Map Container -->
    <div
      ref="mapContainer"
      class="w-full h-64 rounded-2xl border border-slate-200 overflow-hidden shadow-inner relative z-0"
    ></div>

    <div class="flex items-center justify-between text-[11px] text-slate-500 bg-slate-50 p-2 rounded-xl border border-slate-100">
      <div class="flex items-center gap-2">
        <span class="font-mono text-slate-700 font-bold">Lat: {{ currentLat.toFixed(6) }}</span>
        <span class="text-slate-300">|</span>
        <span class="font-mono text-slate-700 font-bold">Lng: {{ currentLng.toFixed(6) }}</span>
      </div>
      <span class="text-brand-600 font-semibold">Radius Geofence: {{ radiusMeters }}m</span>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { MapPin, Navigation } from 'lucide-vue-next'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

// Fix Leaflet marker icon asset resolution in Vite
import markerIcon2x from 'leaflet/dist/images/marker-icon-2x.png'
import markerIcon from 'leaflet/dist/images/marker-icon.png'
import markerShadow from 'leaflet/dist/images/marker-shadow.png'

delete L.Icon.Default.prototype._getIconUrl
L.Icon.Default.mergeOptions({
  iconUrl: markerIcon,
  iconRetinaUrl: markerIcon2x,
  shadowUrl: markerShadow,
})

const props = defineProps({
  latitude: {
    type: Number,
    default: -6.113425,
  },
  longitude: {
    type: Number,
    default: 106.91423,
  },
  radiusMeters: {
    type: Number,
    default: 25,
  },
})

const emit = defineEmits(['update:latitude', 'update:longitude'])

const mapContainer = ref(null)
const currentLat = ref(props.latitude || -6.113425)
const currentLng = ref(props.longitude || 106.91423)

let map = null
let marker = null
let circle = null

const updatePosition = (lat, lng, pan = false) => {
  currentLat.value = lat
  currentLng.value = lng
  emit('update:latitude', Number(lat.toFixed(6)))
  emit('update:longitude', Number(lng.toFixed(6)))

  if (marker) {
    marker.setLatLng([lat, lng])
  }
  if (circle) {
    circle.setLatLng([lat, lng])
    circle.setRadius(props.radiusMeters || 25)
  }
  if (pan && map) {
    map.panTo([lat, lng])
  }
}

const getCurrentLocation = () => {
  if (!navigator.geolocation) {
    alert('Browser tidak mendukung geolokasi GPS')
    return
  }
  navigator.geolocation.getCurrentPosition(
    (pos) => {
      const lat = pos.coords.latitude
      const lng = pos.coords.longitude
      updatePosition(lat, lng, true)
      if (map) {
        map.setView([lat, lng], 17)
      }
    },
    (err) => {
      alert('Gagal mengambil lokasi GPS: ' + err.message)
    },
    { enableHighAccuracy: true, timeout: 8000 }
  )
}

onMounted(() => {
  if (!mapContainer.value) return

  const initialLat = props.latitude || -6.113425
  const initialLng = props.longitude || 106.91423

  map = L.map(mapContainer.value).setView([initialLat, initialLng], 16)

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; OpenStreetMap contributors',
    maxZoom: 19,
  }).addTo(map)

  // Circle Geofence indicator
  circle = L.circle([initialLat, initialLng], {
    color: '#6366f1',
    fillColor: '#818cf8',
    fillOpacity: 0.25,
    radius: props.radiusMeters || 25,
  }).addTo(map)

  // Draggable pin marker
  marker = L.marker([initialLat, initialLng], {
    draggable: true,
  }).addTo(map)

  marker.on('dragend', () => {
    const latlng = marker.getLatLng()
    updatePosition(latlng.lat, latlng.lng)
  })

  map.on('click', (e) => {
    updatePosition(e.latlng.lat, e.latlng.lng)
  })

  // Invalidate map size after animation/render
  setTimeout(() => {
    if (map) {
      map.invalidateSize()
    }
  }, 250)
})

watch(
  () => [props.latitude, props.longitude],
  ([newLat, newLng]) => {
    if (newLat && newLng && (newLat !== currentLat.value || newLng !== currentLng.value)) {
      currentLat.value = newLat
      currentLng.value = newLng
      if (marker) marker.setLatLng([newLat, newLng])
      if (circle) circle.setLatLng([newLat, newLng])
      if (map) map.panTo([newLat, newLng])
    }
  }
)

watch(
  () => props.radiusMeters,
  (newRadius) => {
    if (circle && newRadius > 0) {
      circle.setRadius(newRadius)
    }
  }
)

onUnmounted(() => {
  if (map) {
    map.remove()
    map = null
  }
})
</script>
