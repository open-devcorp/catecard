<template>
  <section class="p-6 border border-gray-200 rounded-2xl">
    <!-- Header -->
  <div class="flex flex-col sm:flex-row sm:justify-between gap-3">
    <div class="sm:pr-4 mb-4">
      <h2 class="text-base font-semibold text-gray-900">
        Grupos de Catequesis
        <span class="text-base text-gray-600">({{ totalLabel }})</span>
      </h2>
      <p class="text-[13px] text-gray-500">Gestiona catequistas y cuentas de escáner</p>
    </div>

    <div class="flex sm:items-start">
      <!-- Botón Nuevo Grupo -->
      <button
        @click="openAddModal"
        class="flex flex-row items-center bg-[#1A388B] text-white px-3 py-2 rounded-lg text-sm cursor-pointer gap-5">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"
            fill="none" stroke="currentColor" stroke-width="2"
            stroke-linecap="round" stroke-linejoin="round"
            class="w-4 h-4">
          <path d="M5 12h14" />
          <path d="M12 5v14" />
        </svg>
        <p class="text-sm">Nuevo Grupo</p>
      </button>
    </div>
  </div>



    <!-- Tabla -->
    <div class="overflow-x-auto rounded-xl bg-white mt-8">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b border-gray-200">
            <th class="px-2 py-2 text-left font-medium">Grupo</th>
            <th class="px-2 py-2 text-left font-medium">Líder/Catequista</th>
            <th class="px-2 py-2 text-center font-medium">Catecúmenos</th>
            <th class="px-3 sm:px-4 py-2 text-center font-medium">Límite</th>
            <th class="px-3 sm:px-4 py-2 text-right font-medium">Acciones</th>
          </tr>
        </thead>

        <!-- Skeleton -->
        <tbody v-if="loading">
          <tr v-for="n in 2" :key="n" class="border-t">
            <td class="px-2 py-2"><div class="h-4 w-36 bg-gray-200 rounded animate-pulse"></div></td>
            <td class="px-2 py-2"><div class="h-4 w-28 bg-gray-200 rounded animate-pulse"></div></td>
            <td class="px-2 py-2"><div class="h-4 w-20 bg-gray-200 rounded animate-pulse"></div></td>
            <td class="px-2 py-2"><div class="h-4 w-20 bg-gray-200 rounded animate-pulse"></div></td>
            <td class="px-2 py-2 text-right"><div class="h-8 w-20 bg-gray-200 rounded animate-pulse inline-block"></div></td>
          </tr>
        </tbody>

        <!-- Sin datos -->
        <tbody v-else-if="rows.length === 0">
          <tr class="border-t">
            <td colspan="5" class="px-4 py-6 text-center text-gray-500">
              No hay grupos que coincidan.
            </td>
          </tr>
        </tbody>

        <!-- Filas -->
        <tbody v-else>
          <tr v-for="r in rows" :key="r.id" class="border-t border-gray-200 hover:bg-gray-50 transition">
            <td class="px-2 py-2">
              <span class="text-gray-900 text-nowrap">{{ r.name }}</span>
            </td>
            <td class="px-2 py-2">
              <span class="text-gray-900 text-nowrap">{{ leaderName(r.id) }}</span>
            </td>
            <td class="px-2 py-2 text-center">
              <span class="inline-block border text-[12px] rounded-md border-gray-200 px-2">{{ countOf(r.id) }}</span>
            </td>
            <td class="px-2 py-2 text-center">
              <span class="inline-block text-gray-700">{{ r.limit_catechumens ?? '—' }}</span>
            </td>
            <td class="px-3 sm:px-4 py-2">
              <div class="flex justify-end gap-1 sm:gap-2">
                <!-- Ver -->
                <button
                  @click="openInfoModal(r.id)"
                  class="inline-flex items-center justify-center rounded-md h-8 w-8 hover:bg-gray-100 transition">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"
                      fill="none" stroke="currentColor" stroke-width="2"
                      stroke-linecap="round" stroke-linejoin="round"
                      class="w-4 h-4">
                    <path d="M2.062 12.348a1 1 0 0 1 0-.696
                            10.75 10.75 0 0 1 19.876 0
                            1 1 0 0 1 0 .696
                            10.75 10.75 0 0 1-19.876 0" />
                    <circle cx="12" cy="12" r="3" />
                  </svg>
                </button>

                <!-- Editar -->
                <button
                  @click="openEditModal(r.id)"
                  class="inline-flex items-center justify-center rounded-md h-8 w-8 hover:bg-gray-100 transition">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"
                      fill="none" stroke="currentColor" stroke-width="2"
                      stroke-linecap="round" stroke-linejoin="round"
                      class="w-4 h-4">
                    <path d="M12 3H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                    <path d="M18.375 2.625a1 1 0 0 1 3 3l-9.013 9.014a2 2 0 0 1-.853.505l-2.873.84a.5.5 0 0 1-.62-.62l.84-2.873a2 2 0 0 1 .506-.852z" />
                  </svg>
                </button>
                <!-- Eliminar -->
                <button
                  @click="confirmDelete(r.id)"
                  :disabled="deletingId === r.id"
                  class="inline-flex items-center justify-center rounded-md h-8 w-8 hover:bg-gray-100 transition text-red-600">
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-trash2 lucide-trash-2 w-4 h-4 text-destructive" aria-hidden="true"><path d="M10 11v6"></path><path d="M14 11v6"></path><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6"></path><path d="M3 6h18"></path><path d="M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>

  <!-- Modales -->
  <modal-component :show="showInfo" @close="closeInfo">
    <div class="pr-2">
      <h3 class="text-xl font-semibold text-gray-900">{{ infoTitle }}</h3>
      <p class="text-sm text-gray-600 mt-1 mb-6">Información detallada del grupo</p>

      <div v-if="loadingDetail" class="space-y-3">
        <div class="h-4 w-48 bg-gray-200 rounded animate-pulse"></div>
        <div class="h-3 w-56 bg-gray-200 rounded animate-pulse"></div>
      </div>

      <p v-else-if="detailError" class="text-sm text-red-600">{{ detailError }}</p>

      <div v-else class="space-y-5 text-[15px]">
        <div>
          <p class="font-semibold text-gray-900">Líder/Catequista</p>
          <p class="mt-1 text-gray-900">{{ catechistDisplay }}</p>
        </div>

        <div>
          <p class="font-semibold text-gray-900">Cantidad de Catecúmenos</p>
          <p class="mt-1 text-gray-900">{{ catechumensCount }}</p>
        </div>
         <div>
          <p class="font-semibold text-gray-900">Límite de catecúmenos</p>
          <p class="mt-1 text-gray-900">{{ limitCatechumens }}</p>
        </div>

        <div class="flex justify-end pt-4">
          <button @click="closeInfo"
                  class="px-3 py-2 bg-white text-sm rounded-lg hover:bg-gray-100 transition-colors border border-gray-300 cursor-pointer">
            Cerrar
          </button>
        </div>
      </div>
    </div>
  </modal-component>

  <modal-component :show="showEdit" @close="closeEdit">
    <div class="pr-2">
      <h3 class="text-xl font-semibold text-gray-900">Editar grupo</h3>
      <p class="text-sm mt-1 text-gray-600 mb-6">Modifica la información del grupo</p>

      <form class="space-y-4" @submit.prevent="submitEdit">
        <div>
          <label class="form-label">Nombre</label>
          <input v-model.trim="editForm.name" type="text" required
                 class="form-input" />
        </div>

      <div>
        <label class="block text-sm text-gray-600 mb-1">Límite de catecúmenos</label>
        <input v-model.number="editForm.limit_catechumens" type="number" min="0" required
        class="form-input" />
      </div>
        <div>
          <label class="form-label">Catequista</label>
          <label class="sr-only">Seleccione catequista</label>
          <select v-model.number="editForm.catechist_id" required
                  class="form-select">
            <option value="" disabled>Seleccione un catequista</option>
            <option v-for="c in catechistsOptions" :key="c.id" :value="c.id">{{ c.full_name }}</option>
          </select>
          <p v-if="loadingCatechistsOptions" class="text-xs text-gray-500 mt-1">Cargando catequistas…</p>
          <p v-else-if="!catechistsOptions.length" class="text-xs text-gray-500 mt-1">No hay catequistas disponibles.</p>
        </div>

        <p v-if="editError" class="text-sm text-red-600">{{ editError }}</p>
        <p v-if="editSuccess" class="text-sm text-emerald-700">{{ editSuccess }}</p>

        <div class="flex justify-end gap-2 pt-2">
          <button type="button" @click="closeEdit"
                  class="px-3 py-2 bg-white text-sm rounded-lg hover:bg-gray-100 transition-colors border border-gray-300 cursor-pointer">Cancelar</button>
          <button type="submit" :disabled="submittingEdit"
                  class="flex flex-row items-center bg-[#1A388B] text-white px-3 py-2 rounded-lg text-sm cursor-pointer hover:bg-[#1A388B]/90 transition-colors">
            {{ submittingEdit ? 'Guardando…' : 'Guardar cambios' }}
          </button>
        </div>
      </form>
    </div>
  </modal-component>

  <modal-component :show="showAdd" @close="closeAdd">
    <div class="pr-2">
      <h3 class="text-xl font-semibold text-gray-900">Nuevo grupo</h3>
      <p class="text-sm mt-1 text-gray-600 mb-6">Registra un nuevo grupo</p>

      <form class="space-y-4" @submit.prevent="submitAdd">
        <div>
          <label class="form-label">Nombre</label>
          <input v-model.trim="addForm.name" type="text" required
                 placeholder="Ej. Grupo San Marcos"
                 class="form-input" />
        </div>

        <div>
          <label class="form-label">Catequista </label>
          <select v-model.number="addForm.catechist_id" required
                  class="form-select">
            <option value="" disabled selected>Seleccione un catequista</option>
            <option v-for="c in catechistsOptions" :key="c.id" :value="c.id">{{ c.full_name }}</option>
          </select>
          <p v-if="loadingCatechistsOptions" class="text-xs text-gray-500 mt-1">Cargando catequistas…</p>
          <p v-else-if="!catechistsOptions.length" class="text-xs text-gray-500 mt-1">No hay catequistas disponibles.</p>
        </div>

        <div>
          <label class="block text-sm text-gray-600 mb-1">Límite de catecúmenos</label>
          <input v-model.number="addForm.limit_catechumens" type="number" min="0" required
                 placeholder="Ej. 30"
                 class="form-input" />
        </div>

        <p v-if="addError" class="text-sm text-red-600">{{ addError }}</p>
        <p v-if="addSuccess" class="text-sm text-emerald-700">{{ addSuccess }}</p>

        <div class="flex justify-end gap-2 pt-2">
          <button type="button" @click="closeAdd"
                  class="px-3 py-2 bg-white text-sm rounded-lg hover:bg-gray-100 transition-colors border border-gray-300 cursor-pointer">Cancelar</button>
          <button type="submit" :disabled="submittingAdd"
                  class="flex flex-row items-center bg-[#1A388B] text-white px-3 py-2 rounded-lg text-sm cursor-pointer hover:bg-[#1A388B]/90 transition-colors">
            {{ submittingAdd ? 'Creando…' : 'Crear grupo' }}
          </button>
        </div>
      </form>
    </div>
  </modal-component>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

const props = defineProps({
  allGroupsUrl: { type: String, default: '/all-groups' },      // { groups: [...] }
  groupUrlBase: { type: String, default: '/admin/groups' },    // { groupInfo: { Group, Catechist, CatechumenSize } }
  editUrlBase:  { type: String, default: '/edit-group' },      // PUT /edit-group/{id} (FormValue)
  addUrl:       { type: String, default: '/add-group' },       // POST /add-group (FormValue)
  catechistsWithoutGroupUrl: { type: String, default: '/catechists-without-group' },
})

type Group = { id: number; name: string; catechist_id: number; limit_catechumens?: number | null }
type Catechist = { id?: number; username?: string; full_name?: string; role?: any }

/* ===== Estado base ===== */
const raw = ref<Group[]>([])
const loading = ref(true)

/* Detalles por fila (líder y conteo) */
const details = ref<Record<number, { catechist?: Catechist|null; count: number }>>({})

/* ===== Carga de lista ===== */
async function fetchGroups(url: string): Promise<Group[]> {
  const res = await fetch(url, { credentials: 'include', headers: { Accept: 'application/json' } })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  const payload = await res.json().catch(() => ({}))
  const list: any[] = Array.isArray(payload?.groups) ? payload.groups : []
  return list.map((g: any) => ({ id: g.id, name: g.name, catechist_id: g.catechist_id, limit_catechumens: g.limit_catechumens ?? null }))
}

async function load() {
  try {
    loading.value = true
    raw.value = await fetchGroups(props.allGroupsUrl)
    await preloadDetails() // ← precarga líder y conteo para la tabla
  } catch (e) {
    console.error('all-groups error:', e)
    raw.value = []
    details.value = {}
  } finally {
    loading.value = false
  }
}
onMounted(load)

async function reloadAll() { await load() }

/* ===== Precarga detalles por grupo ===== */
async function fetchGroupInfo(id: number) {
  const res = await fetch(`${props.groupUrlBase}/${id}`, { credentials: 'include', headers: { Accept: 'application/json' } })
  if (!res.ok) throw new Error(`HTTP ${res.status}`)
  const payload = await res.json().catch(() => ({}))
  const gi = payload?.groupInfo ?? {}
  const g = gi?.Group ?? {}
  const c = gi?.Catechist ?? null
  return {
    group: { id: g.id, name: g.name, catechist_id: g.catechist_id,limit_catechumens: g.limit_catechumens } as Group,
    catechist: c ? { id: c.id, username: c.username, full_name: (c.full_name || ''), role: c.role } as any : null,
    count: Number(gi?.CatechumenSize ?? 0),
  }
}

async function preloadDetails() {
  const map: Record<number, { catechist?: Catechist|null; count: number }> = {}
  await Promise.all(
    raw.value.map(async (g) => {
      try {
        const { catechist, count } = await fetchGroupInfo(g.id)
        map[g.id] = { catechist, count }
      } catch {
        map[g.id] = { catechist: null, count: 0 }
      }
    })
  )
  details.value = map
}

/* ===== Helpers para la tabla ===== */
const rows = computed(() => raw.value)
const totalLabel = computed(() => (loading.value ? '...' : String(raw.value.length)))

function leaderName(id: number) {
  const c = details.value[id]?.catechist as any
  return c?.full_name || c?.username || (c?.id ? `ID ${c.id}` : '—') || '—'
}
function countOf(id: number) {
  const n = details.value[id]?.count
  return typeof n === 'number' ? n : '—'
}

/* ===== Modales INFO ===== */
const showInfo = ref(false)
const loadingDetail = ref(false)
const detailError = ref<string | null>(null)
const groupName = ref<string>('')
const catechist = ref<Catechist | null>(null)
const catechumensCount = ref<number>(0)
const limitCatechumens= ref<number | null>(null)

const infoTitle = computed(() => groupName.value || 'Grupo')
const catechistDisplay = computed(() => leaderName(currentInfoId.value || 0))

const currentInfoId = ref<number | null>(null)

async function openInfoModal(id: number) {
  showInfo.value = true
  loadingDetail.value = true
  detailError.value = null
  currentInfoId.value = id
  try {
    const { group, catechist: c, count } = await fetchGroupInfo(id)
    groupName.value = group.name
    catechist.value = c
    catechumensCount.value = count
    limitCatechumens.value = group.limit_catechumens ?? null
    // sincroniza detalles de la tabla
    details.value[id] = { catechist: c, count }
  } catch (e) {
    console.error('group detail error:', e)
    detailError.value = 'No se pudo cargar la información del grupo.'
  } finally {
    loadingDetail.value = false
  }
}
function closeInfo() { showInfo.value = false }

/* ===== Editar ===== */
const showEdit = ref(false)
const editForm = ref<{ id: number|null; name: string; catechist_id: number|null; limit_catechumens: number|null }>({ id: null, name: '', catechist_id: null, limit_catechumens: null })
const editError = ref<string|null>(null)
const editSuccess = ref<string|null>(null)
const submittingEdit = ref(false)

/* ===== Eliminar ===== */
const deletingId = ref<number | null>(null)
const deleteError = ref<string | null>(null)

async function confirmDelete(id: number) {
  // simple confirmation dialog
  const ok = window.confirm('¿Eliminar este grupo? Esta acción no se puede deshacer.')
  if (!ok) return

  deleteError.value = null
  deletingId.value = id
  try {
    const res = await fetch(`/delete-group/${id}`, {
      method: 'DELETE',
      credentials: 'include',
      headers: { Accept: 'application/json' },
    })
    if (!res.ok) {
      // try to read structured error
      const body = await res.json().catch(() => ({}))
      const msg = body?.message || body?.error || `HTTP ${res.status}`
      throw new Error(msg)
    }

    // Quita del listado y detalles
    raw.value = raw.value.filter(g => g.id !== id)
    const copy = { ...details.value }
    delete copy[id]
    details.value = copy

    // Cierra modal de info si estaba abierto sobre este id
    if (currentInfoId.value === id) closeInfo()
  } catch (e) {
    console.error('delete group error:', e)
    deleteError.value = 'No se pudo eliminar el grupo.'
  } finally {
    deletingId.value = null
  }
}

async function openEditModal(id: number) {
  showEdit.value = true
  editError.value = null
  editSuccess.value = null
  submittingEdit.value = false
  // Cargar opciones en paralelo y la info del grupo
  const loadOptionsP = loadCatechistsWithoutGroup()
  const { group, catechist: currentCatechist } = await fetchGroupInfo(id)
  await loadOptionsP

    // Asegura que el catequista actual esté en las opciones (puede no estar si tiene grupo)
    if (currentCatechist && !catechistsOptions.value.find(c => c.id === currentCatechist.id)) {
      const cc: any = currentCatechist as any
      catechistsOptions.value.unshift({ id: cc.id, username: cc.username || `ID ${cc.id}`, full_name: cc.full_name || '' })
  }

  editForm.value = { id: group.id, name: group.name, catechist_id: group.catechist_id, limit_catechumens: (group.limit_catechumens ?? null) }
}

function closeEdit() { showEdit.value = false }

async function submitEdit() {
  editError.value = null
  editSuccess.value = null
  const f = editForm.value
  if (!f.id || !f.name || !f.catechist_id || f.limit_catechumens == null) {
    editError.value = 'Completa todos los campos.'
    return
  }

  submittingEdit.value = true
  try {
    const body = new URLSearchParams({
      id: String(f.id),
      name: f.name,
      catechist_id: String(f.catechist_id),
      limit_catechumens: String(f.limit_catechumens),
    })

    const res = await fetch(`${props.editUrlBase}/${f.id}`, {
      method: 'PUT',
      credentials: 'include',
      headers: { 'Accept': 'application/json', 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
      body,
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    await res.json().catch(() => ({}))

    // Actualiza lista base
    const idx = raw.value.findIndex(g => g.id === f.id)
  if (idx !== -1) raw.value[idx] = { id: f.id, name: f.name, catechist_id: f.catechist_id, limit_catechumens: f.limit_catechumens }

    // Relee detalles y abre modal de info con la data actualizada
    showEdit.value = false
    await openInfoModal(f.id)
  } catch (e) {
    console.error('edit error:', e)
    editError.value = 'No se pudo guardar los cambios.'
  } finally {
    submittingEdit.value = false
  }
}

/* ===== Agregar (NUEVO) ===== */
const showAdd = ref(false)
const addForm = ref<{ name: string; catechist_id: number|null; limit_catechumens: number|null }>({ name: '', catechist_id: null, limit_catechumens: null })
const addError = ref<string|null>(null)
const addSuccess = ref<string|null>(null)
const submittingAdd = ref(false)

/* Opciones de catequistas (para selector) */
const catechistsOptions = ref<Array<{ id: number; username: string; full_name?: string }>>([])
const loadingCatechistsOptions = ref(false)

async function loadCatechistsWithoutGroup() {
  loadingCatechistsOptions.value = true
  try {
    const res = await fetch(props.catechistsWithoutGroupUrl, { credentials: 'include', headers: { Accept: 'application/json' } })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const payload = await res.json().catch(() => ([]))
    const list: any[] = Array.isArray(payload) ? payload : (Array.isArray(payload?.users) ? payload.users : [])
  catechistsOptions.value = list.map((u: any) => ({ id: u.id, username: u.username, full_name: u.full_name || '' }))
  } catch (e) {
    console.error('load catechists without group error', e)
    catechistsOptions.value = []
  } finally {
    loadingCatechistsOptions.value = false
  }
}

function openAddModal() {
  showAdd.value = true
  addForm.value = { name: '', catechist_id: null, limit_catechumens: null }
  addError.value = null
  addSuccess.value = null
  submittingAdd.value = false
  // cargar opciones de catequistas sin grupo antes de mostrar
  loadCatechistsWithoutGroup()
}

function closeAdd() { showAdd.value = false }

async function submitAdd() {
  addError.value = null
  addSuccess.value = null
  const f = addForm.value
  if (!f.name || !f.catechist_id || f.limit_catechumens == null) {
    addError.value = 'Completa todos los campos.'
    return
  }

  submittingAdd.value = true
  try {
    const body = new URLSearchParams({
      name: f.name,
      catechist_id: String(f.catechist_id),
      limit_catechumens: String(f.limit_catechumens),
    })

    const res = await fetch(props.addUrl, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Accept': 'application/json', 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
      body,
    })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const payload = await res.json().catch(() => ({}))

    // intenta obtener el id creado desde distintas formas posibles del payload
    const newId = Number(payload?.group?.id ?? payload?.id ?? payload?.Group?.id ?? NaN)

    addSuccess.value = 'Grupo creado correctamente.'
    // Sincroniza lista
    await reloadAll()

    // Si tenemos id, abre modal de info del nuevo grupo
    if (!Number.isNaN(newId)) {
      showAdd.value = false
      await openInfoModal(newId)
    } else {
      // si la API no devuelve id, solo cierra
      showAdd.value = false
    }
  } catch (e) {
    console.error('add error:', e)
    addError.value = 'No se pudo crear el grupo.'
  } finally {
    submittingAdd.value = false
  }
}
</script>
