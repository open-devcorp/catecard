<template>
  <section class="max-w-5xl mx-auto">
    <!-- Header -->
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h2 class="text-[22px] font-semibold text-gray-900">Usuarios</h2>
        <p class="text-sm text-gray-500">Gestiona catequistas y cuentas de escáner</p>
      </div>
      <button
        @click="openCreate"
        class="h-9 px-3 rounded-lg bg-[#1A388B] text-white text-sm hover:bg-[#1A388B]/90">
        Nuevo Usuario
      </button>
    </div>

    <!-- CATEQUISTAS -->
    <div class="mb-8">
      <div class="flex items-center justify-between mb-3">
        <h3 class="text-lg font-semibold text-gray-900">
          Catequistas <span class="text-gray-500">({{ catechists.length }})</span>
        </h3>
        <button @click="loadCatechists" class="h-9 px-3 rounded-lg ring-1 ring-gray-300 text-sm hover:bg-gray-50">
          Actualizar
        </button>
      </div>

      <div class="overflow-x-auto rounded-xl border border-gray-200 bg-white">
        <table class="min-w-full text-sm">
          <thead class="bg-gray-50 text-gray-600">
            <tr>
              <th class="px-4 py-3 text-left font-medium">Nombre de usuario</th>
              <th class="px-4 py-3 text-left font-medium">Correo</th>
              <th class="px-4 py-3 text-left font-medium">Rol</th>
              <th class="px-4 py-3 text-right font-medium">Acciones</th>
            </tr>
          </thead>

          <!-- Skeleton -->
          <tbody v-if="loading.catechists">
            <tr v-for="n in 3" :key="'c-skel-'+n" class="border-t">
              <td class="px-4 py-3"><div class="h-4 w-44 bg-gray-200 rounded animate-pulse"></div></td>
              <td class="px-4 py-3"><div class="h-4 w-64 bg-gray-200 rounded animate-pulse"></div></td>
              <td class="px-4 py-3"><div class="h-4 w-20 bg-gray-200 rounded animate-pulse"></div></td>
              <td class="px-4 py-3 text-right"><div class="h-8 w-36 bg-gray-200 rounded animate-pulse inline-block"></div></td>
            </tr>
          </tbody>

          <!-- Datos -->
          <tbody v-else>
            <tr v-if="!catechists.length" class="border-t">
              <td colspan="4" class="px-4 py-6 text-center text-gray-500">No hay catequistas.</td>
            </tr>
            <tr v-for="u in catechists" :key="'c-'+u.id" class="border-t">
              <td class="px-4 py-3 text-gray-900">{{ u.username }}</td>
              <td class="px-4 py-3 text-gray-900">{{ u.email }}</td>
              <td class="px-4 py-3">
                <span class="inline-flex items-center rounded-full text-[11px] px-2 py-0.5 ring-1 ring-emerald-200 bg-emerald-50 text-emerald-700">
                  Catequista
                </span>
              </td>
              <td class="px-4 py-3">
                <div class="flex justify-end gap-2">
                  <button
                    @click="openCredentials(u, 'catechist')"
                    class="inline-flex items-center rounded-md text-xs px-3 py-1.5 ring-1 ring-blue-200 bg-blue-50 text-blue-700 hover:bg-blue-100">
                    Ver credenciales
                  </button>
                  <button
                    @click="confirmDelete(u, 'catechist')"
                    class="inline-flex items-center rounded-md text-xs px-3 py-1.5 ring-1 ring-red-200 bg-red-50 text-red-700 hover:bg-red-100">
                    Eliminar
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- SCANNERS -->
    <div>
      <div class="flex items-center justify-between mb-3">
        <h3 class="text-lg font-semibold text-gray-900">
          Escáneres <span class="text-gray-500">({{ scanners.length }})</span>
        </h3>
        <button @click="loadScanners" class="h-9 px-3 rounded-lg ring-1 ring-gray-300 text-sm hover:bg-gray-50">
          Actualizar
        </button>
      </div>

      <div class="overflow-x-auto rounded-xl border border-gray-200 bg-white">
        <table class="min-w-full text-sm">
          <thead class="bg-gray-50 text-gray-600">
            <tr>
              <th class="px-4 py-3 text-left font-medium">Nombre de usuario</th>
              <th class="px-4 py-3 text-left font-medium">Correo</th>
              <th class="px-4 py-3 text-left font-medium">Rol</th>
              <th class="px-4 py-3 text-right font-medium">Acciones</th>
            </tr>
          </thead>

          <!-- Skeleton -->
          <tbody v-if="loading.scanners">
            <tr v-for="n in 2" :key="'s-skel-'+n" class="border-t">
              <td class="px-4 py-3"><div class="h-4 w-44 bg-gray-200 rounded animate-pulse"></div></td>
              <td class="px-4 py-3"><div class="h-4 w-64 bg-gray-200 rounded animate-pulse"></div></td>
              <td class="px-4 py-3"><div class="h-4 w-20 bg-gray-200 rounded animate-pulse"></div></td>
              <td class="px-4 py-3 text-right"><div class="h-8 w-36 bg-gray-200 rounded animate-pulse inline-block"></div></td>
            </tr>
          </tbody>

          <!-- Datos -->
          <tbody v-else>
            <tr v-if="!scanners.length" class="border-t">
              <td colspan="4" class="px-4 py-6 text-center text-gray-500">No hay cuentas de escáner.</td>
            </tr>
            <tr v-for="u in scanners" :key="'s-'+u.id" class="border-t">
              <td class="px-4 py-3 text-gray-900">{{ u.username }}</td>
              <td class="px-4 py-3 text-gray-900">{{ u.email }}</td>
              <td class="px-4 py-3">
                <span class="inline-flex items-center rounded-full text-[11px] px-2 py-0.5 ring-1 ring-indigo-200 bg-indigo-50 text-indigo-700">
                  Escáner
                </span>
              </td>
              <td class="px-4 py-3">
                <div class="flex justify-end gap-2">
                  <button
                    @click="openCredentials(u, 'scanner')"
                    class="inline-flex items-center rounded-md text-xs px-3 py-1.5 ring-1 ring-blue-200 bg-blue-50 text-blue-700 hover:bg-blue-100">
                    Ver credenciales
                  </button>
                  <button
                    @click="confirmDelete(u, 'scanner')"
                    class="inline-flex items-center rounded-md text-xs px-3 py-1.5 ring-1 ring-red-200 bg-red-50 text-red-700 hover:bg-red-100">
                    Eliminar
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </section>

  <!-- MODAL: Crear Nuevo Usuario -->
  <modal-component :show="showCreate" @close="closeCreate">
    <div class="pr-2">
      <h3 class="text-xl font-semibold text-gray-900">Crear Nuevo Usuario</h3>
      <p class="text-sm text-gray-500">Crea una cuenta de catequista o escáner</p>

      <form class="space-y-4 mt-5" @submit.prevent="submitCreate">
        <div>
          <label class="block text-sm text-gray-600 mb-1">Nombre completo</label>
          <input v-model.trim="createForm.username" type="text" required
                 placeholder="Ej: María García"
                 class="w-full h-10 rounded-lg border border-gray-300 px-3 text-sm focus:ring-2 focus:ring-blue-200 focus:outline-none bg-gray-50" />
        </div>

        <div>
          <label class="block text-sm text-gray-600 mb-1">Correo electrónico</label>
          <input v-model.trim="createForm.email" type="email" required
                 placeholder="Ej: maria@catequesis.com"
                 class="w-full h-10 rounded-lg border border-gray-300 px-3 text-sm focus:ring-2 focus:ring-blue-200 focus:outline-none bg-gray-50" />
          <p class="text-[12px] text-gray-500 mt-1">Se generará una contraseña automáticamente.</p>
        </div>

        <div>
          <label class="block text-sm text-gray-700 mb-2">Tipo de cuenta</label>
          <div class="grid grid-cols-2 gap-2">
            <button type="button"
                    @click="createForm.role = 'catechist'"
                    :class="['h-10 rounded-lg text-sm flex items-center justify-center gap-2 ring-1',
                             createForm.role==='catechist' ? 'bg-[#1A388B] text-white ring-[#1A388B]' : 'bg-white text-gray-900 ring-gray-300 hover:bg-gray-50']">
              <span>👤</span> Catequista
            </button>
            <button type="button"
                    @click="createForm.role = 'scanner'"
                    :class="['h-10 rounded-lg text-sm flex items-center justify-center gap-2 ring-1',
                             createForm.role==='scanner' ? 'bg-[#1A388B] text-white ring-[#1A388B]' : 'bg-white text-gray-900 ring-gray-300 hover:bg-gray-50']">
              <span>🛰️</span> Escáner
            </button>
          </div>
        </div>

        <p v-if="createError" class="text-sm text-red-600">{{ createError }}</p>

        <div class="flex justify-end gap-2 pt-2">
          <button type="button" @click="closeCreate"
                  class="px-4 py-2 rounded-lg ring-1 ring-gray-300 bg-white hover:bg-gray-50">
            Cancelar
          </button>
          <button type="submit" :disabled="submittingCreate"
                  class="px-4 py-2 rounded-lg bg-[#1A388B] text-white hover:bg-[#1A388B]/90 disabled:opacity-60">
            {{ submittingCreate ? 'Creando…' : 'Crear Usuario' }}
          </button>
        </div>
      </form>
    </div>
  </modal-component>

  <!-- MODAL: Credenciales -->
  <modal-component :show="showModal" @close="closeModal">
    <div class="pr-2">
      <h3 class="text-xl font-semibold text-gray-900">Credenciales de Acceso</h3>
      <p class="text-sm text-gray-500">Comparte estas credenciales con {{ modalName }}</p>

      <div class="bg-gray-50 border border-gray-200 rounded-lg p-3 mt-4 text-[13px] text-gray-600">
        Guarda estas credenciales en un lugar seguro. El usuario necesitará estos datos para iniciar sesión.
      </div>

      <div class="space-y-5 mt-5 text-[15px]">
        <div>
          <p class="text-gray-600 text-sm">Nombre</p>
          <p class="mt-1 text-gray-900">{{ modalName }}</p>
        </div>

        <div>
          <p class="text-gray-600 text-sm">Tipo de Cuenta</p>
          <span v-if="modalType === 'catechist'"
                class="inline-flex items-center mt-1 rounded-full text-[11px] px-2 py-0.5 ring-1 ring-emerald-200 bg-emerald-50 text-emerald-700">
            Catequista
          </span>
          <span v-else
                class="inline-flex items-center mt-1 rounded-full text-[11px] px-2 py-0.5 ring-1 ring-indigo-200 bg-indigo-50 text-indigo-700">
            Escáner
          </span>
        </div>

        <div>
          <p class="text-gray-600 text-sm">Correo Electrónico</p>
          <div class="mt-1 flex items-center gap-2">
            <input :value="modalEmail" readonly
                   class="flex-1 h-10 rounded-lg border border-gray-300 px-3 text-sm focus:outline-none" />
            <button @click="copy(modalEmail)"
                    class="inline-flex items-center justify-center h-10 w-10 rounded-md ring-1 ring-gray-300 hover:bg-gray-50"
                    title="Copiar correo">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
              </svg>
            </button>
          </div>
        </div>

        <div>
          <p class="text-gray-600 text-sm">Contraseña</p>
          <div class="mt-1 flex items-center gap-2">
            <input :value="modalPassword" readonly
                   class="flex-1 h-10 rounded-lg border border-gray-300 px-3 text-sm focus:outline-none" />
            <button @click="copy(modalPassword)"
                    class="inline-flex items-center justify-center h-10 w-10 rounded-md ring-1 ring-gray-300 hover:bg-gray-50"
                    title="Copiar contraseña">
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
              </svg>
            </button>
          </div>
        </div>

        <div class="flex justify-end pt-2">
          <button @click="closeModal"
                  class="px-4 py-2 rounded-lg ring-1 ring-gray-300 bg-white hover:bg-gray-50">
            Cerrar
          </button>
        </div>
      </div>
    </div>
  </modal-component>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'

type User = {
  id: number
  username: string
  email: string
  password?: string
  role?: number | string
}

const props = defineProps({
  catechistsUrl: { type: String, default: '/all-catechists' },
  scannersUrl:   { type: String, default: '/all-scanners' },
  deleteCatechistUrlBase: { type: String, default: '/catechists' },
  deleteScannerUrlBase:   { type: String, default: '/scanners' },
  createUrl: { type: String, default: '/users' }, // POST /users
})

/* ===== Listas ===== */
const catechists = ref<User[]>([])
const scanners = ref<User[]>([])
const loading = reactive({ catechists: true, scanners: true })

/* ===== Modal Credenciales ===== */
const showModal = ref(false)
const modalType = ref<'catechist'|'scanner'>('catechist')
const modalUser = ref<User | null>(null)
const modalName = computed(() => modalUser.value?.username || 'Usuario')
const modalEmail = computed(() => modalUser.value?.email || '')
const modalPassword = computed(() => modalUser.value?.password || '')

/* ===== Modal Crear ===== */
const showCreate = ref(false)
const submittingCreate = ref(false)
const createError = ref<string|null>(null)
const createForm = ref<{ username: string; email: string; role: 'catechist'|'scanner' }>({
  username: '',
  email: '',
  role: 'catechist',
})

/* ===== Fetchers ===== */
async function loadCatechists() {
  try {
    loading.catechists = true
    const res = await fetch(props.catechistsUrl, { headers: { Accept: 'application/json' }, credentials: 'include', cache: 'no-store' })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json().catch(() => [])
    catechists.value = Array.isArray(data) ? data : (data?.users || data?.catechists || [])
  } catch (e) {
    console.error('catechists error', e); catechists.value = []
  } finally {
    loading.catechists = false
  }
}
async function loadScanners() {
  try {
    loading.scanners = true
    const res = await fetch(props.scannersUrl, { headers: { Accept: 'application/json' }, credentials: 'include', cache: 'no-store' })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json().catch(() => [])
    scanners.value = Array.isArray(data) ? data : (data?.users || data?.scanners || [])
  } catch (e) {
    console.error('scanners error', e); scanners.value = []
  } finally {
    loading.scanners = false
  }
}
onMounted(async () => { await Promise.all([loadCatechists(), loadScanners()]) })

/* ===== Modal Credenciales ===== */
function openCredentials(u: User, type: 'catechist'|'scanner') {
  modalUser.value = { ...u }
  modalType.value = type
  showModal.value = true
}
function closeModal() { showModal.value = false; modalUser.value = null }

/* ===== Crear Usuario ===== */
function openCreate() {
  createError.value = null
  createForm.value = { username: '', email: '', role: 'catechist' }
  showCreate.value = true
}
function closeCreate() { showCreate.value = false }

const ROLE_MAP: Record<'catechist'|'scanner', number> = { catechist: 1, scanner: 3 }

async function submitCreate() {
  createError.value = null
  const f = createForm.value
  if (!f.username || !f.email) { createError.value = 'Completa nombre y correo.'; return }

  submittingCreate.value = true
  try {
    const body = new URLSearchParams({
      username: f.username,
      email: f.email,
      role: String(ROLE_MAP[f.role]),
    })
    const res = await fetch(props.createUrl, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Accept': 'application/json', 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
      body,
    })
    if (!res.ok) {
      const msg = `HTTP ${res.status}`
      createError.value = 'No se pudo crear el usuario. ' + msg
      throw new Error(msg)
    }
    const payload = await res.json().catch(() => ({}))
    const u = payload?.user || {}

    const created: User = {
      id: u.id, username: u.username, email: u.email, password: u.password, role: u.role
    }

    // Añadir a la tabla correspondiente
    if (f.role === 'catechist') catechists.value = [created, ...catechists.value]
    else scanners.value = [created, ...scanners.value]

    // Abrir credenciales con el usuario recién creado
    modalType.value = f.role
    openCredentials(created, f.role)

    // Cerrar el modal de crear
    showCreate.value = false
  } catch (e) {
    console.error('create error:', e)
    if (!createError.value) createError.value = 'No se pudo crear el usuario.'
  } finally {
    submittingCreate.value = false
  }
}

/* ===== Copiar ===== */
async function copy(text: string) {
  try { await navigator.clipboard.writeText(text || '') }
  catch (e) { console.error('clipboard error', e) }
}

/* ===== Delete ===== */
async function doDelete(type: 'catechist'|'scanner', id: number) {
  const base = type === 'catechist' ? '/catechists' : '/scanners'
  const res = await fetch(`${base}/${id}`, { method: 'DELETE', credentials: 'include', headers: { Accept: 'application/json' } })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  if (type === 'catechist') catechists.value = catechists.value.filter(u => u.id !== id)
  else scanners.value = scanners.value.filter(u => u.id !== id)
}
async function confirmDelete(u: User, type: 'catechist'|'scanner') {
  if (!confirm('¿Eliminar esta cuenta? Esta acción no se puede deshacer.')) return
  try { await doDelete(type, u.id) }
  catch (e) { console.error('delete error', e); alert('No se pudo eliminar.') }
}
async function handleDeleteFromModal() {
  if (!modalUser.value) return
  const type = modalType.value
  const id = modalUser.value.id
  if (!confirm('¿Eliminar esta cuenta? Esta acción no se puede deshacer.')) return
  try { await doDelete(type, id); closeModal() }
  catch (e) { console.error('delete error', e); alert('No se pudo eliminar.') }
}
</script>
