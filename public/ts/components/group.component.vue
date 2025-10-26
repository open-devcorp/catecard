<template>
  <section class="max-w-5xl mx-auto">
    <!-- Header -->
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-[22px] font-semibold text-gray-900">
        Grupos <span class="text-gray-500">({{ totalLabel }})</span>
      </h2>
      <div class="flex items-center gap-2">
        <!-- New -->
        <button
          @click="openAddModal"
          class="h-9 px-3 rounded-lg bg-[#1A388B] text-white text-sm hover:bg-[#1A388B]/90"
        >
          Nuevo
        </button>
        <button
          @click="reloadAll"
          class="h-9 px-3 rounded-lg ring-1 ring-gray-300 text-sm hover:bg-gray-50"
        >
          Actualizar
        </button>
      </div>
    </div>

    <!-- Tabla -->
    <div class="overflow-x-auto rounded-xl border border-gray-200 bg-white">
      <table class="min-w-full text-sm">
        <thead class="bg-gray-50 text-gray-600">
          <tr>
            <th class="px-4 py-3 text-left font-medium">Nombre</th>
            <th class="px-4 py-3 text-left font-medium">Líder</th>
            <th class="px-4 py-3 text-left font-medium">Cantidad de catecúmenos</th>
            <th class="px-4 py-3 text-right font-medium">Acciones</th>
          </tr>
        </thead>

        <!-- Skeleton -->
        <tbody v-if="loading">
          <tr v-for="n in 3" :key="n" class="border-t">
            <td class="px-4 py-3"><div class="h-4 w-44 bg-gray-200 rounded animate-pulse"></div></td>
            <td class="px-4 py-3"><div class="h-4 w-36 bg-gray-200 rounded animate-pulse"></div></td>
            <td class="px-4 py-3"><div class="h-4 w-24 bg-gray-200 rounded animate-pulse"></div></td>
            <td class="px-4 py-3 text-right"><div class="h-8 w-28 bg-gray-200 rounded animate-pulse inline-block"></div></td>
          </tr>
        </tbody>

        <!-- Sin datos -->
        <tbody v-else-if="rows.length === 0">
          <tr class="border-t">
            <td colspan="4" class="px-4 py-6 text-center text-gray-500">
              No hay grupos que coincidan.
            </td>
          </tr>
        </tbody>

        <!-- Filas -->
        <tbody v-else>
          <tr v-for="r in rows" :key="r.id" class="border-t">
            <td class="px-4 py-3">
              <span class="text-gray-900">{{ r.name }}</span>
            </td>
            <td class="px-4 py-3">
              <span class="text-gray-900">{{ leaderName(r.id) }}</span>
            </td>
            <td class="px-4 py-3">
              <span class="text-gray-900">{{ countOf(r.id) }}</span>
            </td>
            <td class="px-4 py-3">
              <div class="flex justify-end gap-2">
                <button
                  @click="openInfoModal(r.id)"
                  class="inline-flex items-center rounded-md text-xs px-3 py-1.5 ring-1 ring-blue-200 bg-blue-50 text-blue-700 hover:bg-blue-100">
                  Ver
                </button>
                <button
                  @click="openEditModal(r.id)"
                  class="inline-flex items-center rounded-md text-xs px-3 py-1.5 ring-1 ring-gray-300 hover:bg-gray-50">
                  Editar
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>

  <!-- Modal INFO -->
  <modal-component :show="showInfo" @close="closeInfo">
    <div class="pr-2">
      <h3 class="text-xl font-semibold text-gray-900">{{ infoTitle }}</h3>
      <p class="text-sm text-gray-400 mb-6">Información detallada del grupo</p>

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

        <div class="flex justify-end pt-4">
          <button @click="closeInfo"
                  class="px-4 py-2 rounded-lg ring-1 ring-gray-300 bg-white hover:bg-gray-50">
            Cerrar
          </button>
        </div>
      </div>
    </div>
  </modal-component>

  <!-- Modal EDIT -->
  <modal-component :show="showEdit" @close="closeEdit">
    <div class="pr-2">
      <h3 class="text-xl font-semibold text-gray-900">Editar grupo</h3>
      <p class="text-sm text-gray-400 mb-6">Modifica la información del grupo</p>

      <form class="space-y-4" @submit.prevent="submitEdit">
        <div>
          <label class="block text-sm text-gray-600 mb-1">Nombre</label>
          <input v-model.trim="editForm.name" type="text" required
                 class="w-full h-9 rounded-lg border border-gray-300 px-3 text-sm focus:ring-2 focus:ring-blue-200 focus:outline-none" />
        </div>

        <div>
          <label class="block text-sm text-gray-600 mb-1">Catequista ID</label>
          <input v-model.number="editForm.catechist_id" type="number" min="1" required
                 class="w-full h-9 rounded-lg border border-gray-300 px-3 text-sm focus:ring-2 focus:ring-blue-200 focus:outline-none" />
        </div>

        <p v-if="editError" class="text-sm text-red-600">{{ editError }}</p>
        <p v-if="editSuccess" class="text-sm text-emerald-700">{{ editSuccess }}</p>

        <div class="flex justify-end gap-2 pt-2">
          <button type="button" @click="closeEdit"
                  class="px-4 py-2 rounded-lg ring-1 ring-gray-300 bg-white hover:bg-gray-50">Cancelar</button>
          <button type="submit" :disabled="submittingEdit"
                  class="px-4 py-2 rounded-lg bg-[#1A388B] text-white hover:bg-[#1A388B]/90 disabled:opacity-60">
            {{ submittingEdit ? 'Guardando…' : 'Guardar cambios' }}
          </button>
        </div>
      </form>
    </div>
  </modal-component>

  <!-- Modal ADD (NUEVO) -->
  <modal-component :show="showAdd" @close="closeAdd">
    <div class="pr-2">
      <h3 class="text-xl font-semibold text-gray-900">Nuevo grupo</h3>
      <p class="text-sm text-gray-400 mb-6">Registra un nuevo grupo</p>

      <form class="space-y-4" @submit.prevent="submitAdd">
        <div>
          <label class="block text-sm text-gray-600 mb-1">Nombre</label>
          <input v-model.trim="addForm.name" type="text" required
                 placeholder="Ej. Grupo San Marcos"
                 class="w-full h-9 rounded-lg border border-gray-300 px-3 text-sm focus:ring-2 focus:ring-blue-200 focus:outline-none" />
        </div>

        <div>
          <label class="block text-sm text-gray-600 mb-1">Catequista ID</label>
          <input v-model.number="addForm.catechist_id" type="number" min="1" required
                 placeholder="Ej. 12"
                 class="w-full h-9 rounded-lg border border-gray-300 px-3 text-sm focus:ring-2 focus:ring-blue-200 focus:outline-none" />
        </div>

        <p v-if="addError" class="text-sm text-red-600">{{ addError }}</p>
        <p v-if="addSuccess" class="text-sm text-emerald-700">{{ addSuccess }}</p>

        <div class="flex justify-end gap-2 pt-2">
          <button type="button" @click="closeAdd"
                  class="px-4 py-2 rounded-lg ring-1 ring-gray-300 bg-white hover:bg-gray-50">Cancelar</button>
          <button type="submit" :disabled="submittingAdd"
                  class="px-4 py-2 rounded-lg bg-[#1A388B] text-white hover:bg-[#1A388B]/90 disabled:opacity-60">
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
})

type Group = { id: number; name: string; catechist_id: number }
type Catechist = { id?: number; username?: string; email?: string; role?: any }

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
  return list.map((g: any) => ({ id: g.id, name: g.name, catechist_id: g.catechist_id }))
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
    group: { id: g.id, name: g.name, catechist_id: g.catechist_id } as Group,
    catechist: c ? { id: c.id, username: c.username, email: c.email, role: c.role } : null,
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
  const c = details.value[id]?.catechist
  return c?.username || c?.email || (c?.id ? `ID ${c.id}` : '—') || '—'
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
const editForm = ref<{ id: number|null; name: string; catechist_id: number|null }>({ id: null, name: '', catechist_id: null })
const editError = ref<string|null>(null)
const editSuccess = ref<string|null>(null)
const submittingEdit = ref(false)

async function openEditModal(id: number) {
  showEdit.value = true
  editError.value = null
  editSuccess.value = null
  submittingEdit.value = false
  const { group } = await fetchGroupInfo(id)
  editForm.value = { id: group.id, name: group.name, catechist_id: group.catechist_id }
}

function closeEdit() { showEdit.value = false }

async function submitEdit() {
  editError.value = null
  editSuccess.value = null
  const f = editForm.value
  if (!f.id || !f.name || !f.catechist_id) {
    editError.value = 'Completa todos los campos.'
    return
  }

  submittingEdit.value = true
  try {
    const body = new URLSearchParams({
      id: String(f.id),
      name: f.name,
      catechist_id: String(f.catechist_id),
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
    if (idx !== -1) raw.value[idx] = { id: f.id, name: f.name, catechist_id: f.catechist_id }

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
const addForm = ref<{ name: string; catechist_id: number|null }>({ name: '', catechist_id: null })
const addError = ref<string|null>(null)
const addSuccess = ref<string|null>(null)
const submittingAdd = ref(false)

function openAddModal() {
  showAdd.value = true
  addForm.value = { name: '', catechist_id: null }
  addError.value = null
  addSuccess.value = null
  submittingAdd.value = false
}

function closeAdd() { showAdd.value = false }

async function submitAdd() {
  addError.value = null
  addSuccess.value = null
  const f = addForm.value
  if (!f.name || !f.catechist_id) {
    addError.value = 'Completa todos los campos.'
    return
  }

  submittingAdd.value = true
  try {
    const body = new URLSearchParams({
      name: f.name,
      catechist_id: String(f.catechist_id),
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
