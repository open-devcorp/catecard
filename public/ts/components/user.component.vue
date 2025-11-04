<template>
  <section class="p-6 border border-gray-200 rounded-2xl">
    <!-- Header -->
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-4 gap-3">
      <div>
        <h2 class="font-semibold ">Usuarios</h2>
        <p class="text-[13px] text-gray-600">Gestiona catequistas y cuentas de escáner</p>
      </div>
      <div class="flex sm:items-start">

        <button
          @click="openCreate"
          class="flex flex-row items-center bg-[#1A388B] text-white px-3 py-2 rounded-lg text-sm cursor-pointer gap-3">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"
               fill="none" stroke="currentColor" stroke-width="2"
               stroke-linecap="round" stroke-linejoin="round"
               class="w-4 h-4">
            <path d="M5 12h14" />
            <path d="M12 5v14" />
          </svg>
          Nuevo Usuario
        </button>
      </div>
    </div>

 
    <!-- ========== SCANNERS ========== -->
    <div class="space-y-3 mb-10 mt-10">
      <h2 class="font-semibold ">
              Escáneres
            <span class="text-base text-gray-600">({{ scanners.length }})</span>
        </h2>

      <template v-if="loading.scanners">
        <div v-for="n in 2" :key="'s-skel-'+n" class="rounded-2xl border border-gray-200 bg-white p-4">
          <div class="flex items-center gap-4">
            <span class="w-9 h-9 rounded-xl bg-indigo-100 animate-pulse"></span>
            <div>
              <div class="h-4 w-40 bg-gray-200 rounded mb-2 animate-pulse"></div>
              <div class="h-3 w-56 bg-gray-200 rounded animate-pulse"></div>
            </div>
          </div>
        </div>
      </template>

      <div v-else-if="!scanners.length"
           class="rounded-2xl border border-gray-200 bg-white p-6 text-center text-gray-500">
        No hay cuentas de escáner.
      </div>

      <template v-else>
        <div v-for="u in scanners" :key="'s-'+u.id"
             class="rounded-2xl border border-gray-200 bg-white p-4 flex items-center justify-between hover:shadow-sm transition">
          <div class="flex items-center gap-4 min-w-0">
            <span class="w-9 h-9 rounded-xl bg-indigo-100 flex items-center justify-center text-indigo-700 font-medium">
              {{ u.full_name.charAt(0).toUpperCase() }}
            </span>
            <div class="min-w-0">
              <p class="text-[15px] font-medium truncate max-w-[160px] sm:max-w-none">{{ u.full_name }}</p>

            </div>
          </div>

          <div class="flex items-center gap-1 sm:gap-2">
            <button
              @click="openCredentials(u, 'scanner')"
              class="inline-flex items-center justify-center h-8 w-8 sm:w-auto sm:h-auto sm:px-3 sm:py-1.5 border border-gray-300 rounded-md bg-gray-100 transition text-xs sm:text-sm font-medium text-gray-800"
              title="Ver credenciales"
              aria-label="Ver credenciales">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"
                   fill="none" stroke="currentColor" stroke-width="2"
                   stroke-linecap="round" stroke-linejoin="round"
                   class="w-4 h-4 text-gray-700">
                <path d="M2.062 12.348a1 1 0 0 1 0-.696
                        10.75 10.75 0 0 1 19.876 0
                        1 1 0 0 1 0 .696
                        10.75 10.75 0 0 1-19.876 0" />
                <circle cx="12" cy="12" r="3" />
              </svg>
              <span class="hidden sm:inline">Ver Credenciales</span>
            </button>

            <button
              @click="confirmDelete(u, 'scanner')"
              class="inline-flex items-center justify-center h-8 w-8 rounded-md hover:bg-gray-100 text-red-600 transition"
              title="Eliminar">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"
                   fill="none" stroke="currentColor" stroke-width="2"
                   stroke-linecap="round" stroke-linejoin="round"
                   class="w-4 h-4">
                <polyline points="3 6 5 6 21 6" />
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4
                         a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
                <line x1="10" y1="11" x2="10" y2="17" />
                <line x1="14" y1="11" x2="14" y2="17" />
              </svg>
            </button>
          </div>
        </div>
      </template>
    </div>
       <!-- ========== CATEQUISTAS ========== -->
    <div class="space-y-3 mb-10">
        <h2 class="font-semibold ">
              Catequistas
            <span class="text-base text-gray-600">({{ catechists.length }})</span>
        </h2>

      <!-- Skeleton -->
      <template v-if="loading.catechists">
        <div v-for="n in 3" :key="'c-skel-'+n" class="rounded-2xl border border-gray-200 bg-white p-4">
          <div class="flex items-center gap-4">
            <span class="w-9 h-9 rounded-xl bg-emerald-100 animate-pulse"></span>
            <div>
              <div class="h-4 w-40 bg-gray-200 rounded mb-2 animate-pulse"></div>
              <div class="h-3 w-56 bg-gray-200 rounded animate-pulse"></div>
            </div>
          </div>
        </div>
      </template>

      <!-- Empty -->
      <div v-else-if="!catechists.length"
           class="rounded-2xl border border-gray-200 bg-white p-6 text-center text-gray-500">
        No hay catequistas.
      </div>

      <!-- Items -->
      <template v-else>
        <div v-for="u in catechists" :key="'c-'+u.id"
            class="rounded-2xl border border-gray-200 bg-white p-4 flex items-center justify-between hover:shadow-sm transition">
          <div class="flex items-center gap-4 min-w-0">
            <span class="w-9 h-9 rounded-xl bg-emerald-100 flex items-center justify-center text-emerald-700 font-medium">
              {{ u.full_name.charAt(0).toUpperCase() }}
            </span>
            <div class="min-w-0">
              <p class="text-[15px] font-medium truncate max-w-[160px] sm:max-w-none">{{ u.full_name }}</p>
              <div class="flex">
              <div class="flex gap-1 mt-1">
                <span v-for="g in u.groups" :key="g.id"
                      class="text-[11px] px-2 py-0.5 rounded-full bg-gray-100 text-gray-700">
                  {{ g.name }}
                </span>
              </div>
              </div>
             
            </div>
          </div>

          <div class="flex items-center gap-1 sm:gap-2">
            <button
                @click="openCredentials(u, 'catechist')"
                class="inline-flex items-center justify-center h-8 w-8 sm:w-auto sm:h-auto sm:px-3 sm:py-1.5 border border-gray-300 rounded-md bg-gray-100 transition text-xs sm:text-sm font-medium text-gray-800"
                title="Ver credenciales"
                aria-label="Ver credenciales"
              >
                <svg xmlns="http://www.w3.org/2000/svg"
                    fill="none" viewBox="0 0 24 24"
                    stroke="currentColor" stroke-width="2"
                    stroke-linecap="round" stroke-linejoin="round"
                    class="w-4 h-4">
                  <path d="M2.062 12.348a1 1 0 0 1 0-.696
                          10.75 10.75 0 0 1 19.876 0
                          1 1 0 0 1 0 .696
                          10.75 10.75 0 0 1-19.876 0" />
                  <circle cx="12" cy="12" r="3" />
                </svg>
                <span class="hidden sm:inline">Ver Credenciales</span>
            </button>


            <button
              @click="confirmDelete(u, 'catechist')"
              class="inline-flex items-center justify-center h-8 w-8 rounded-md hover:bg-gray-100 text-red-600 transition"
              title="Eliminar">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"
                  fill="none" stroke="currentColor" stroke-width="2"
                  stroke-linecap="round" stroke-linejoin="round"
                  class="w-4 h-4">
                <polyline points="3 6 5 6 21 6" />
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4
                        a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
                <line x1="10" y1="11" x2="10" y2="17" />
                <line x1="14" y1="11" x2="14" y2="17" />
              </svg>
            </button>
          </div>
        </div>
      </template>

    </div>

  </section>
  <!-- MODAL: Crear Nuevo Usuario -->
  <modal-component :show="showCreate" @close="closeCreate">
    <div class="pr-2">
      <h3 class="text-lg font-semibold ">Crear Nuevo Usuario</h3>
      <p class="text-sm text-gray-600 mt-1">Crea una cuenta de catequista o escáner</p>

      <form class="space-y-4 mt-5" @submit.prevent="submitCreate">
        <div>
          <label class="form-label">Nombre completo</label>
          <input v-model.trim="createForm.fullname" type="text" required
                 placeholder="Ej: María García"
                 class="form-input" />
        </div>

        <div>
          <label class="form-label">Nombre de usuario</label>
          <input v-model.trim="createForm.username" type="text" required
                 placeholder="Ej: mgarcia"
                 class="form-input" />
          <p class="text-[12px] text-gray-500 mt-1">Se generará una contraseña automáticamente.</p>
        </div>

        <div>
          <label class="block text-sm text-gray-700 mb-2">Tipo de cuenta</label>
          <div class="grid grid-cols-2 gap-2">
            <!-- Catequista -->
            <button
              type="button"
              @click="createForm.role = 'catechist'"
              :class="[
                'h-10 rounded-lg text-sm flex items-center justify-center gap-2 ring-1 transition',
                createForm.role === 'catechist'
                  ? 'bg-[#1A388B] text-white ring-[#1A388B]'
                  : 'bg-white  ring-gray-300 hover:bg-gray-50'
              ]"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="w-4 h-4"
              >
                <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path>
                <circle cx="12" cy="7" r="4"></circle>
              </svg>
              Catequista
            </button>

            <!-- Escáner -->
            <button
              type="button"
              @click="createForm.role = 'scanner'"
              :class="[
                'rounded-lg text-sm flex items-center justify-center gap-2 ring-1 transition',
                createForm.role === 'scanner'
                  ? 'bg-[#1A388B] text-white ring-[#1A388B]'
                  : 'bg-white  ring-gray-300 hover:bg-gray-50'
              ]"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="w-4 h-4"
              >
                <path d="M3 7V5a2 2 0 0 1 2-2h2"></path>
                <path d="M17 3h2a2 2 0 0 1 2 2v2"></path>
                <path d="M21 17v2a2 2 0 0 1-2 2h-2"></path>
                <path d="M7 21H5a2 2 0 0 1-2-2v-2"></path>
              </svg>
              Escáner
            </button>
          </div>

        </div>

        <p v-if="createError" class="text-sm text-red-600">{{ createError }}</p>

        <div class="flex justify-end gap-2 pt-2">
          <button type="button" @click="closeCreate"
                  class="px-3 py-2 bg-white text-sm rounded-lg hover:bg-gray-100 transition-colors border border-gray-300 cursor-pointer">
            Cancelar
          </button>
          <button type="submit" :disabled="submittingCreate"
                  class="flex flex-row items-center bg-[#1A388B] text-white px-3 py-2 rounded-lg text-sm cursor-pointer hover:bg-[#1A388B]/90 transition-colors">
            {{ submittingCreate ? 'Creando…' : 'Crear Usuario' }}
          </button>
        </div>
      </form>
    </div>
  </modal-component>

  <!-- MODAL: Credenciales -->
<modal-component :show="showModal" @close="closeModal">
  <div class="pr-2">
    <h3 class="text-lg font-semibold ">Credenciales de Acceso</h3>
    <p class="text-sm text-gray-600 mt-1">Comparte estas credenciales con {{ modalName }}</p>

    <div class="border border-gray-200 rounded-lg py-3 px-4 mt-4 text-[13px] text-gray-600">
      Guarda estas credenciales en un lugar seguro. El usuario necesitará estos datos para iniciar sesión.
    </div>

    <div class="space-y-4 mt-5">
      <!-- Nombre completo -->
      <div>
        <p class="text-sm font-semibold">Nombre completo</p>
        <p class="mt-1 ">{{ modalFullName }}</p>
      </div>

      <!-- Tipo de Cuenta -->
      <div>
        <p class="text-sm font-semibold">Tipo de Cuenta</p>
        <span v-if="modalType === 'catechist'"
          class="inline-flex items-center mt-1 rounded-full text-[11px] px-2 py-0.5 ring-1 ring-emerald-200 bg-emerald-50 text-emerald-700">
          Catequista
        </span>
        <span v-else
          class="inline-flex items-center mt-1 rounded-full text-[11px] px-2 py-0.5 ring-1 ring-indigo-200 bg-indigo-50 text-indigo-700">
          Escáner
        </span>
      </div>

      <!-- Nombre de usuario (para login) -->
      <div>
        <p class="text-sm font-semibold">Nombre de usuario</p>
        <div class="mt-1 flex items-center gap-2">
          <input :value="modalUser?.username || ''" readonly class="form-input" />
          <button @click="copy(modalUser?.username || '')"
                  class="inline-flex items-center justify-center h-8 px-2 rounded-md ring-1 ring-gray-300 hover:bg-gray-50"
                  title="Copiar usuario">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-copy w-4 h-4" aria-hidden="true"><rect width="14" height="14" x="8" y="8" rx="2" ry="2"></rect><path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"></path></svg>
          </button>
        </div>
      </div>

      <!-- (Opcional) Puedes duplicar datos si quieres otro bloque de copia rápida del nombre completo -->

      <!-- Contraseña -->
      <div>
        <p class="text-sm font-semibold">Contraseña</p>
        <div class="mt-1 flex items-center gap-2">
          <input :value="modalPassword" readonly
                class="form-input" />
          <button @click="copy(modalPassword)"
                  class="inline-flex items-center justify-center h-8 px-2 rounded-md ring-1 ring-gray-300 hover:bg-gray-50"
                  title="Copiar contraseña">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-copy w-4 h-4" aria-hidden="true"><rect width="14" height="14" x="8" y="8" rx="2" ry="2"></rect><path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"></path></svg>
          </button>
        </div>
      </div>

      <!-- Botón cerrar -->
      <div class="flex justify-end pt-2">
        <button @click="closeModal"
                class="px-3 py-2 bg-white text-sm rounded-lg hover:bg-gray-100 transition-colors border border-gray-300 cursor-pointer">
          Cerrar
        </button>
      </div>
    </div>
  </div>
</modal-component>

</template>
<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'

/* ===== Tipos ===== */
type Group = { id: number; name: string }
type User = { id: number; username: string; full_name: string; password?: string; role?: number|string; groups?: Group[] }
type UserData = { User: User; Groups: Group[] }

/* ===== Props ===== */
const props = defineProps({
  catechistsUrl: { type: String, default: '/all-catechists' },
  scannersUrl: { type: String, default: '/all-scanners' },
  deleteCatechistUrlBase: { type: String, default: '/catechists' },
  deleteScannerUrlBase: { type: String, default: '/scanners' },
  deleteUserUrlBase: { type: String, default: '/delete-user' },
  createUrl: { type: String, default: '/users' },
})

/* ===== Estado ===== */
const catechists = ref<User[]>([])
const scanners = ref<User[]>([])
const loading = reactive({ catechists: true, scanners: true })

/* ===== Modal Credenciales ===== */
const showModal = ref(false)
const modalType = ref<'catechist'|'scanner'>('catechist')
const modalUser = ref<User | null>(null)
const modalName = computed(() => modalUser.value?.username || 'Usuario')
const modalFullName = computed(() => modalUser.value?.full_name || '')
const modalPassword = computed(() => modalUser.value?.password || '')
const modalGroups = computed(() => modalUser.value?.groups || [])

/* ===== Modal Crear ===== */
const showCreate = ref(false)
const submittingCreate = ref(false)
const createError = ref<string|null>(null)
const createForm = ref<{ fullname: string; username: string; role: 'catechist'|'scanner' }>({
  fullname: '',
  username: '',
  role: 'catechist',
})

/* ===== Fetchers ===== */
async function loadCatechists() {
  loading.catechists = true
  try {
    const res = await fetch(props.catechistsUrl, { headers: { Accept: 'application/json' }, credentials: 'include' })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data: UserData[] = await res.json().catch(() => [])
    catechists.value = Array.isArray(data) ? data.map(d => ({ ...d.User, groups: d.Groups || [] })) : []
  } catch (e) {
    console.error('catechists error', e)
    catechists.value = []
  } finally {
    loading.catechists = false
  }
}

async function loadScanners() {
  loading.scanners = true
  try {
    const res = await fetch(props.scannersUrl, { headers: { Accept: 'application/json' }, credentials: 'include' })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data: UserData[] = await res.json().catch(() => [])
    scanners.value = Array.isArray(data) ? data.map(d => ({ ...d.User, groups: d.Groups || [] })) : []
  } catch (e) {
    console.error('scanners error', e)
    scanners.value = []
  } finally {
    loading.scanners = false
  }
}

onMounted(async () => {
  await Promise.all([loadCatechists(), loadScanners()])
})

/* ===== Modal Credenciales ===== */
function openCredentials(u: User, type: 'catechist'|'scanner') {
  modalUser.value = { ...u }
  modalType.value = type
  showModal.value = true
}
function closeModal() {
  showModal.value = false
  modalUser.value = null
}

/* ===== Crear Usuario ===== */
function openCreate() {
  createError.value = null
  createForm.value = { fullname: '', username: '', role: 'catechist' }
  showCreate.value = true
}
function closeCreate() {
  showCreate.value = false
}

const ROLE_MAP: Record<'catechist'|'scanner', number> = {
  catechist: 1,
  scanner: 3,
}

async function submitCreate() {
  createError.value = null
  const f = createForm.value
  if (!f.fullname || !f.username) {
    createError.value = 'Completa nombre completo y nombre de usuario.'
    return
  }

  submittingCreate.value = true
  try {
    const body = new URLSearchParams({
      username: f.username,
      full_name: f.fullname,
      role: String(ROLE_MAP[f.role]),
    })

    const res = await fetch(props.createUrl, {
      method: 'POST',
      credentials: 'include',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
      },
      body,
    })

    // ⚙️ Intentamos parsear el JSON aunque el status no sea 200
    let payload: any = null
    try {
      payload = await res.json()
    } catch {
      payload = null
    }

    // Normalizar distintas formas de respuesta del backend
    // Algunos endpoints devuelven { user: {...} } y otros { User: {...}, Groups: [...] }
    const respUser = payload?.user ?? payload?.User ?? null
    const respGroups = payload?.Groups ?? payload?.groups ?? []

    if (!res.ok && !respUser) {
      throw new Error(`HTTP ${res.status}`)
    }

    if (!respUser) {
      // nada útil en la respuesta
      throw new Error('No user returned from server')
    }

    const created: User = {
      ...respUser,
      password: respUser.password || '',
      groups: respGroups || [],
    }

    // 🔄 Refrescar la lista según el tipo y luego mostrar modal con las credenciales
    if (f.role === 'catechist') {
      await loadCatechists()
    } else {
      await loadScanners()
    }

    // ✅ Mostrar modal de credenciales usando la información devuelta por el servidor
    openCredentials(created, f.role)
    showCreate.value = false

  } catch (e) {
    console.error('create error', e)
    if (!createError.value) createError.value = 'No se pudo crear el usuario (revisa si ya existe).'
  } finally {
    submittingCreate.value = false
  }
}

/* ===== Copiar ===== */
async function copy(text: string) {
  try {
    await navigator.clipboard.writeText(text || '')
  } catch (e) {
    console.error('clipboard error', e)
  }
}

/* ===== Delete ===== */
async function doDelete(type: 'catechist'|'scanner', id: number) {
  // Prefer explicit delete-user endpoint when available (server routes use /delete-user/{id})
  const base = props.deleteUserUrlBase || (type === 'catechist' ? props.deleteCatechistUrlBase : props.deleteScannerUrlBase)
  const res = await fetch(`${base}/${id}`, {
    method: 'DELETE',
    credentials: 'include',
    headers: { Accept: 'application/json' },
  })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)

  if (type === 'catechist') catechists.value = catechists.value.filter(u => u.id !== id)
  else scanners.value = scanners.value.filter(u => u.id !== id)
}

async function confirmDelete(u: User, type: 'catechist'|'scanner') {
  if (!confirm('¿Eliminar esta cuenta? Esta acción no se puede deshacer.')) return
  try {
    await doDelete(type, u.id)
  } catch (e) {
    console.error('delete error', e)
    alert('No se pudo eliminar.')
  }
}

</script>
