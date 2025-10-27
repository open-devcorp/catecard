<template>
  <section class="bg-white rounded-xl border border-gray-200 p-6 space-y-6">
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-[22px] font-semibold text-gray-900">
        Escaneos
           <span class="text-base text-gray-600">({{ groups.length }})</span>
        </h2>
 
    </div>

    <!-- Lista (cards) -->
    <div class="space-y-3">
      <!-- Skeleton -->
      <template v-if="loading">
        <div v-for="n in 3" :key="n" class="rounded-2xl border border-gray-200 bg-white p-4">
          
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-4">
              <span class="w-9 h-9 rounded-xl bg-emerald-100 animate-pulse"></span>
              <div>
                <div class="h-4 w-40 bg-gray-200 rounded mb-2 animate-pulse"></div>
                <div class="h-3 w-56 bg-gray-200 rounded animate-pulse"></div>
              </div>
            </div>
          </div>
        </div>
      </template>

      <!-- Empty -->
      <div
        v-else-if="!groups.length"
        class="rounded-2xl border border-gray-200 bg-white p-6 text-center text-gray-500"
      >
        No hay registros aún.
      </div>

      <!-- Items -->
      <template v-else>
        <div
          v-for="g in groups"
          :key="g.key"
          class="rounded-2xl border border-gray-200 bg-white p-4"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-4">
<span class="w-12 h-12 rounded-xl bg-emerald-100 flex items-center justify-center">
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"
       fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"
       stroke-linejoin="round"
       class="w-6 h-6 text-green-800"
       aria-hidden="true">
    <path d="M3 7V5a2 2 0 0 1 2-2h2"></path>
    <path d="M17 3h2a2 2 0 0 1 2 2v2"></path>
    <path d="M21 17v2a2 2 0 0 1-2 2h-2"></path>
    <path d="M7 21H5a2 2 0 0 1-2-2v-2"></path>
  </svg>
</span>


              <div>
                <div class="flex items-center gap-3">
                  <p class="text-[17px] font-medium text-gray-900">{{ g.name }}</p>
                </div>
                <p class="text-[13px] text-gray-600 mt-1">
                  <span class="font-medium">Escaneado por:</span> {{ g.usersText }}
                </p>
                <p v-if="g.lastAt" class="text-[12px] text-gray-500 mt-0.5">
                  Último: {{ timeAgo(g.lastAt) }}
                </p>
              </div>
            </div>
            <div class="flex items-center gap-2">
              <span
                class="inline-flex items-center rounded-full text-xs px-2 py-0.5 ring-1 ring-blue-200 bg-blue-50 text-blue-700"
              >
                {{ g.count }} scan{{ g.count > 1 ? 's' : '' }}
              </span>
            </div>
          </div>
        </div>
      </template>
    </div>
  </section>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

// props mínimas: endpoint (por defecto /admin/scans)
const props = defineProps({
  endpoint: { type: String, default: '/admin/scans' },
  showSkeleton: { type: Boolean, default: true },
  // autoRefresh (opcional): vuelve a pedir data cada N ms
  autoRefreshMs: { type: Number, default: 0 },
})

// 👉 seed para precarga inmediata (se reemplaza con la data real apenas llega)
const seed = [
  {
    catechumen: { id: 101, full_name: 'Juan Pérez' },
    user: { username: 'scanner-1' },
    scanned_at: new Date(Date.now() - 1000 * 60 * 12).toISOString(), // hace 12 min
  },
  {
    catechumen: { id: 102, full_name: 'María López' },
    user: { username: 'scanner-2' },
    scanned_at: new Date(Date.now() - 1000 * 60 * 45).toISOString(), // hace 45 min
  },
  {
    catechumen: { id: 101, full_name: 'Juan Pérez' },
    user: { username: 'scanner-3' },
    scanned_at: new Date(Date.now() - 1000 * 60 * 2).toISOString(), // hace 2 min
  },
]

const CACHE_KEY = 'admin_scans_cache_v1'

// estado
const loading = ref(true)
const raw = ref(
  // usa cache si existe, si no seed
  (() => {
    try {
      const cached = JSON.parse(localStorage.getItem(CACHE_KEY) || 'null')
      return Array.isArray(cached) ? cached : seed
    } catch {
      return seed
    }
  })()
)

// helpers
function pick(o, keys, def = undefined) {
  for (const k of keys) if (o && o[k] != null) return o[k]
  return def
}
function parseDate(v) {
  const d = new Date(v)
  return isNaN(d) ? null : d
}
function timeAgoDate(d) {
  const s = (Date.now() - d.getTime()) / 1000
  if (s < 60) return 'hace unos segundos'
  if (s < 3600) return `hace ${Math.floor(s / 60)} min`
  if (s < 86400) return `hace ${Math.floor(s / 3600)} h`
  return d.toLocaleString()
}
function timeAgo(d) {
  return d ? timeAgoDate(d) : '—'
}

// agrupación computada
const groups = computed(() => {
  const map = new Map()
  for (const it of raw.value) {
    const cate = pick(it, ['catechumen', 'Catechumen'], {}) || {}
    const user = pick(it, ['user', 'User'], {}) || {}
    const cid = pick(cate, ['id', 'ID', 'Id'])
    const name = pick(cate, ['full_name', 'fullName', 'name', 'Name'], 'Sin nombre')
    const uname = pick(user, ['username', 'Username', 'name', 'Name'], 'N/A')
    const at = parseDate(pick(it, ['scanned_at', 'scannedAt', 'created_at', 'createdAt', 'timestamp']))

    const key = cid != null && cid !== '' ? `id:${cid}` : `name:${String(name).toLowerCase()}`
    if (!map.has(key)) {
      map.set(key, { key, cid: cid ?? null, name, count: 0, users: new Set(), lastAt: null })
    }
    const e = map.get(key)
    e.count += 1
    e.users.add(uname)
    if (!e.lastAt || (at && at > e.lastAt)) e.lastAt = at
  }

  const arr = Array.from(map.values()).map((g) => ({
    ...g,
    usersText: Array.from(g.users).join(', ') || 'N/A',
  }))

  // ordenar: más escaneos y más reciente
  arr.sort((a, b) => {
    if (b.count !== a.count) return b.count - a.count
    return (b.lastAt?.getTime() || 0) - (a.lastAt?.getTime() || 0)
  })
  return arr
})

// fetch real (reemplaza seed/cache)
async function load() {
  try {
    loading.value = props.showSkeleton
    const res = await fetch(props.endpoint, {
      method: 'GET',
      credentials: 'include',
      headers: { Accept: 'application/json' },
      cache: 'no-store',
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const payload = await res.json()
    const data = Array.isArray(payload)
      ? payload
      : Array.isArray(payload?.data)
      ? payload.data
      : Array.isArray(payload?.records)
      ? payload.records
      : []
    raw.value = data
    // guarda en cache para próximas visitas
    try { localStorage.setItem(CACHE_KEY, JSON.stringify(data)) } catch {}
  } catch (e) {
    console.error('Error cargando scans:', e)
    // si falla, mantén cache/seed (raw ya tiene algo)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  // muestra seed/cache inmediato y luego reemplaza
  load()
  // auto refresh opcional
  if (props.autoRefreshMs > 0) {
    setInterval(load, props.autoRefreshMs)
  }
})

// export para template (en <script setup> las funciones ya están en alcance)
</script>
