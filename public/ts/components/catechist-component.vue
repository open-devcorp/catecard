<template>
  <div>
    <!-- Header actions -->
    <div class="flex justify-end">
      <button
        @click="openAddModal"
        class="flex flex-row items-center bg-[#1A388B] text-white px-3 py-2 rounded-lg text-sm cursor-pointer"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24"
             viewBox="0 0 24 24" fill="none" stroke="currentColor"
             stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
             class="lucide lucide-plus w-4 h-4 mr-2">
          <path d="M5 12h14"></path>
          <path d="M12 5v14"></path>
        </svg>
        Agregar Catecúmeno
      </button>
    </div>

    <!-- Modal Agregar/Editar -->
    <modal-component :show="showModal" @close="showModal = false">
      <template #default>
        <p class="text-lg font-semibold mb-1">{{ modalTitle }}</p>
        <p class="text-gray-400 text-sm">{{ modalSubtitle }}</p>
        <form class="space-y-4 mt-5" @submit.prevent>
          <!-- Campos ocultos -->
          <input type="hidden" v-model="form.id" />
          <input type="hidden" v-model="form.group_id" />

          <div>
            <label class="form-label">Nombre Completo</label>
            <input type="text" v-model.trim="form.full_name" required class="form-input" placeholder="Ej: Italo Luna" />
          </div>

          <div>
            <label class="form-label">Edad</label>
            <input type="text" v-model.trim="form.age" required class="form-input" placeholder="Ej: 15" />
            <p class="text-gray-400 text-[13px] mt-1">Edad entre 5 y 18 años</p>
          </div>

          <div>
            <label class="form-label">Grupo</label>
            <input type="text" :value="groupInputText" class="form-input focus:outline-none" readonly />
            <p class="text-gray-500 text-sm mt-5">El QR se podrá escanear solo 2 veces según coordinación</p>
          </div>

          <div class="flex justify-end space-x-2 mt-5">
            <button type="button" @click="showModal = false"
              class="px-3 py-2 bg-white text-sm rounded-lg hover:bg-gray-100 transition-colors border border-gray-300 cursor-pointer">
              Cancelar
            </button>
            <button type="button" @click="submit"
              class="flex flex-row items-center bg-[#1A388B] text-white px-3 py-2 rounded-lg text-sm cursor-pointer hover:bg-[#1A388B]/90 transition-colors">
              {{ submitLabel }}
            </button>
          </div>
        </form>
      </template>
    </modal-component>

    <!-- Modal QR creado -->
    <modal-component :show="showQrModal" @close="closeQrModal">
      <template #default>
        <div class="">
            <h3 class="text-[18px] font-semibold text-gray-900 mb-1">Código QR - {{ form.full_name }}</h3>
            <p class="text-sm text-gray-500">Este código QR permite el ingreso de <span class="font-semibold">{{ qrForumLeft }}</span> personas</p>

            <div v-if="qrClaimUrl" class="mt-5 flex flex-col items-center gap-3">
                <!-- Contenedor para QR con estilo -->
                <div class="p-4 rounded-2xl ring-1 ring-gray-200 bg-gray-50">
                    <div ref="qrContainer" class="w-56 h-56 flex items-center justify-center"></div>
                </div>
            </div>
            <div class="flex flex-col justify-center text-center mt-5 ">
                <p class="text-[13px] text-gray-500">Catecúmeno</p>
                <p class="text-base">{{ form.full_name }}</p>
            </div>
                
            <div v-if="qrForumLeft != null" class="mt-5 flex justify-center">
                <div class="flex flex-row items-center gap-6">
                    <div class="flex flex-col items-center text-center">
                    <div class="text-xs text-gray-500">Invitados restantes:</div>
                    <div class="inline-flex items-center justify-center rounded-md border border-gray-300 px-2 py-0.5 text-xs font-medium mt-1">
                        {{ qrForumLeft }}
                    </div>
                    </div>

                    <div class="flex flex-col items-center text-center">
          <div class="text-xs text-gray-500">Veces escaneadas:</div>
          <div class="inline-flex items-center justify-center rounded-md border border-gray-300 px-2 py-0.5 text-xs font-medium mt-1">
            {{ qrUsedCount ?? '…' }}
          </div>
                    </div>
                </div>
            </div>



          <div class="flex justify-end space-x-2 mt-10">
            <button @click="closeQrModal"
                    class="px-3 py-2 bg-white rounded-lg text-sm hover:bg-gray-100 transition-colors border border-gray-300 cursor-pointer">
              Cerrar
            </button>
            <button @click="downloadQr"
                      class="px-3 py-1.5 rounded-md bg-[#1A388B] text-white text-sm hover:bg-[#1A388B]/90 cursor-pointer">Descargar</button>
          </div>
        </div>
      </template>
    </modal-component>

        <!-- Modal de confirmación (no bloquea el hilo) -->
        <modal-component :show="showDeleteConfirm" @close="cancelDelete">
          <template #default>
            <div>
              <h3 class="text-[18px] font-semibold text-gray-900 mb-1">Eliminar catecúmeno</h3>
              <p class="text-sm text-gray-500 mb-4">Esta acción no se puede deshacer.</p>
              <div class="flex justify-end gap-2">
                <button @click="cancelDelete"
                        class="px-3 py-2 bg-white rounded-lg text-sm hover:bg-gray-100 transition-colors border border-gray-300 cursor-pointer">Cancelar</button>
                <button @click="performDelete"
                        class="px-3 sm:px-4 py-2 text-red-500 bg-red-500/20 rounded-lg text-sm font-medium hover:opacity-80 transition-all cursor-pointer">Eliminar</button>
              </div>
            </div>
          </template>
        </modal-component>

    <!-- Listado -->
    <div class="space-y-3 mt-6">
      <div v-if="loading" class="text-gray-500 text-sm">Cargando catecúmenos...</div>
      <template v-else>
        <div v-if="!catechumens.length" class="text-gray-500 text-sm">No hay catecúmenos registrados.</div>
        <div v-else id="catechumensList" class="space-y-3">
          <div v-for="c in catechumens" :key="c.id"
               class="bg-white p-3 sm:p-5 rounded-xl border border-gray-200 hover:shadow-md transition-all max-w-full">
            <div class="flex flex-row items-center justify-between gap-2 sm:gap-0">
              <div class="flex flex-row items-center gap-2 sm:gap-4 min-w-0 flex-1">
                <div class="w-12 h-12 rounded-xl bg-[#1A388B]/10 flex items-center justify-center shrink-0">
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" stroke="currentColor"
                       stroke-width="2" stroke-linecap="round" stroke-linejoin="round"
                       class="w-6 h-6 text-[#1A388B]">
                    <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path>
                    <circle cx="12" cy="7" r="4"></circle>
                  </svg>
                </div>
                <div class="flex flex-col gap-1">
                  <div class="flex flex-row gap-2 items-center">
                    <p class="font-normal text-sm sm:text-base truncate">{{ c.full_name }}</p>
                    <p class="text-xs sm:text-[13px] text-gray-500 block sm:hidden">- {{ c.age }} años</p>
                  </div>
                  <div class="flex flex-row items-center gap-2 sm:gap-3 flex-wrap">
                    <p class="text-xs bg-gray-200 px-2 rounded-lg py-1 truncate">{{ c.group?.name }}</p>
                    <p class="text-xs sm:text-[13px] text-gray-500 hidden sm:block">{{ c.age }} años</p>
                    <p class="text-xs sm:text-[13px] text-gray-500">• {{ (c as any).scanUsed ?? '…' }}/{{ (c as any).scanTotal ?? '…' }} escaneos</p>
                  </div>
                </div>
              </div>
              <div class="flex flex-row gap-1 sm:gap-3 shrink-0">
                <button
                  @click="openQr(c)"
                  class="inline-flex items-center justify-center h-8 w-8 sm:h-auto sm:w-auto sm:px-4 py-2 text-xs sm:text-sm bg-[#1A388B] text-white rounded-lg cursor-pointer"
                  title="Ver QR" aria-label="Ver QR">
                  <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-white lucide lucide-qr-code w-4 h-4 sm:mr-2" aria-hidden="true"><rect width="5" height="5" x="3" y="3" rx="1"></rect><rect width="5" height="5" x="16" y="3" rx="1"></rect><rect width="5" height="5" x="3" y="16" rx="1"></rect><path d="M21 16h-3a2 2 0 0 0-2 2v3"></path><path d="M21 21v.01"></path><path d="M12 7v3a2 2 0 0 1-2 2H7"></path><path d="M3 12h.01"></path><path d="M12 3h.01"></path><path d="M12 16v.01"></path><path d="M16 12h1"></path><path d="M21 12v.01"></path><path d="M12 21v-1"></path></svg>
                  <span class="hidden sm:inline">Ver QR</span>
                </button>
                <button
                  @click="openEdit(c)"
                  class="inline-flex items-center justify-center h-8 w-8 sm:h-auto sm:w-auto sm:px-4 py-2 text-xs sm:text-sm text-[#1A388B] bg-[#1A388B]/20 rounded-lg font-medium hover:opacity-80 transition-all cursor-pointer"
                  title="Editar" aria-label="Editar">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 sm:mr-2"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4 12.5-12.5z"/></svg>
                  <span class="hidden sm:inline">Editar</span>
                </button>
                <button
                  @click="askDelete(c)"
                  class="inline-flex items-center justify-center h-8 w-8 sm:h-auto sm:w-auto sm:px-4 py-2 text-xs sm:text-sm text-red-600 bg-red-500/20 rounded-lg font-medium hover:opacity-80 transition-all cursor-pointer"
                  title="Eliminar" aria-label="Eliminar">
                  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="w-4 h-4 sm:mr-2"><polyline points="3 6 5 6 21 6" /><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4 a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" /><line x1="10" y1="11" x2="10" y2="17" /><line x1="14" y1="11" x2="14" y2="17" /></svg>
                  <span class="hidden sm:inline">Eliminar</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed, onMounted, nextTick, watch } from 'vue'

/* ===== Tipos ===== */
interface Group { id: number; name: string }
interface Catechumen {
  id: number
  full_name: string
  age: string | number
  group: Group
  // enriquecidos en FE
  scanUsed?: number
  scanTotal?: number
}

/* ===== Props (URLs y contexto) ===== */
const props = defineProps({
  listUrl: { type: String, default: '/catechumens' },
  addUrl: { type: String, default: '/add-catechumen' },
  editUrlBase: { type: String, default: '/catechumen' },
  deleteUrlBase: { type: String, default: '/delete-catechumen' },
  groupName: { type: String, default: '' },
  groupId: { type: Number, default: 0 },
})

/* ===== Estado ===== */
const loading = ref(true)
const catechumens = ref<Catechumen[]>([])
const showModal = ref(false)
const isEditMode = ref(false)
const showQrModal = ref(false)
const form = ref<{ id: number | null; group_id: number | null; full_name: string; age: string | number }>({
  id: null,
  group_id: props.groupId || null,
  full_name: '',
  age: ''
})

// Confirmación de eliminación no bloqueante
const showDeleteConfirm = ref(false)
const pendingDeleteId = ref<number | null>(null)

/* ===== UI texts ===== */
const modalTitle = computed(() => (isEditMode.value ? 'Editar Catecúmeno' : 'Agregar Catecúmeno'))
const modalSubtitle = computed(() => (isEditMode.value ? 'Modifica la información del catecúmeno' : 'Completa la información del catecúmeno'))
const submitLabel = computed(() => (isEditMode.value ? 'Actualizar catecúmeno' : 'Agregar catecúmeno'))
const groupInputText = computed(() => {
  // en edición, preferir nombre del grupo de la fila si existe
  if (isEditMode.value) {
    const c = catechumens.value.find(x => x.id === form.value.id)
    if (c?.group?.name) return c.group.name
  }
  return props.groupName || ''
})

/* ===== Data Loaders ===== */
async function loadCatechumens(opts?: { silent?: boolean }) {
  const silent = !!opts?.silent
  if (!silent) loading.value = true
  try {
    const res = await fetch(props.listUrl, { headers: { Accept: 'application/json' }, credentials: 'include' })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json().catch(() => ({} as any))
    const list: any[] = Array.isArray(data?.catechumens) ? data.catechumens : []
    catechumens.value = list.map((c: any) => ({
      id: c.id,
      full_name: c.full_name,
      age: c.age,
      group: c.group ? { id: c.group.id, name: c.group.name } : { id: c.group_id || 0, name: '' }
    }))
    // Enriquecer con métricas de escaneos en background
    void enrichScanStats()
  } catch (e) {
    console.error('Error cargando catecúmenos:', e)
    catechumens.value = []
  } finally {
    if (!silent) loading.value = false
  }
}

onMounted(loadCatechumens)

/* ===== Acciones ===== */
function openAddModal() {
  isEditMode.value = false
  form.value = { id: null, group_id: props.groupId || null, full_name: '', age: '' }
  showModal.value = true
}

function openEdit(c: Catechumen) {
  isEditMode.value = true
  form.value = {
    id: c.id,
    group_id: (c.group?.id ?? props.groupId) || null,
    full_name: c.full_name,
    age: c.age
  }
  showModal.value = true
}

function askDelete(c: Catechumen) {
  pendingDeleteId.value = c.id
  showDeleteConfirm.value = true
}
function cancelDelete() {
  showDeleteConfirm.value = false
  pendingDeleteId.value = null
}
async function performDelete() {
  if (!pendingDeleteId.value) return
  const id = pendingDeleteId.value
  cancelDelete()
  await remove(id, true)
}

// Abrir modal QR para un catecúmeno existente
async function openQr(c: Catechumen) {
  try {
    // 1) Traer catecúmeno por ID para obtener qr_id
    const res = await fetch(`${props.editUrlBase}/${c.id}`, { headers: { Accept: 'application/json' }, credentials: 'include' })
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json().catch(() => ({} as any))
    const cate = data?.catechumen || {}
    const qrId = Number(cate.qr_id ?? NaN)
    if (Number.isNaN(qrId)) {
      alert('QR no encontrado para este catecúmeno')
      return
    }

    // 2) Opcional: traer datos del QR para mostrar cupos restantes
    try {
      const r2 = await fetch(`/qr/${qrId}`, { headers: { Accept: 'application/json' }, credentials: 'include' })
      if (r2.ok) {
        const qr = await r2.json().catch(() => ({} as any))
        const forum = Number(qr?.forum ?? qr?.Forum ?? NaN)
        const used = Number(qr?.count ?? qr?.Count ?? NaN)
        qrForumLeft.value = Number.isNaN(forum) ? null : forum
        qrUsedCount.value = Number.isNaN(used) ? null : used
      }
    } catch {}

    // 3) Preparar modal
    form.value.full_name = c.full_name
    qrClaimUrl.value = `${window.location.origin}/qr/${qrId}/claim`
    showQrModal.value = true
    await nextTick()
    await ensureQrInstance(qrClaimUrl.value)
  } catch (e) {
    console.error('Error abriendo QR del catecúmeno:', e)
    alert('No se pudo abrir el QR. Inténtalo de nuevo.')
  }
}


async function submit() {
  const f = form.value
  if (!f.full_name || !f.age) {
    alert('Por favor completa todos los campos')
    return
  }

  try {
    let res: Response
    if (isEditMode.value) {
      if (!f.id) throw new Error('Falta ID')
      const body = new URLSearchParams({
        id: String(f.id),
        full_name: String(f.full_name),
        age: String(f.age),
        group_id: String(f.group_id || '')
      })
      res = await fetch(`${props.editUrlBase}/${f.id}`, {
        method: 'PUT',
        headers: { 'Accept': 'application/json', 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
        credentials: 'include',
        body
      })
    } else {
      const body = new URLSearchParams({
        full_name: String(f.full_name),
        age: String(f.age)
      })
      res = await fetch(props.addUrl, {
        method: 'POST',
        headers: { 'Accept': 'application/json', 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8' },
        credentials: 'include',
        body
      })
    }

    if (!res.ok) {
      const msg = await readErrorMessage(res)
      throw new Error(msg)
    }

  const payload: any = await res.json().catch(() => ({}))
  // refrescar lista sin mostrar spinner para evitar parpadeo
  await loadCatechumens({ silent: true })
  void enrichScanStats()

    if (!isEditMode.value) {
      // Si es creación, mostrar modal con QR
      const qr = payload?.qr || null
      const qrId = Number(qr?.id ?? qr?.ID ?? NaN)
      const forum = Number(qr?.forum ?? qr?.Forum ?? NaN)
      const used = Number(qr?.count ?? qr?.Count ?? 0)
      if (!Number.isNaN(qrId)) {
        qrForumLeft.value = Number.isNaN(forum) ? null : forum
        qrUsedCount.value = Number.isNaN(used) ? 0 : used
        qrClaimUrl.value = `${window.location.origin}/qr/${qrId}/claim`
  showModal.value = false
  showQrModal.value = true
  // Render QR con estilo cuando el contenedor exista
  await nextTick()
  await ensureQrInstance(qrClaimUrl.value)
      } else {
        // fallback: solo cerrar
        showModal.value = false
      }
    } else {
      // Edición: solo cerrar
      showModal.value = false
    }
  } catch (e: any) {
    console.error('Error guardando catecúmeno:', e)
    const msg = String(e?.message || e || '')
    alert(msg)
  }
}

async function remove(id: number, skipConfirm?: boolean) {
  if (!id) return
  if (!skipConfirm) {
    if (!confirm('¿Seguro que deseas eliminar este catecúmeno?')) return
  }

  try {
    const res = await fetch(`${props.deleteUrlBase}/${id}`, { method: 'DELETE', headers: { Accept: 'application/json' }, credentials: 'include' })
    if (!res.ok) {
      const msg = await readErrorMessage(res)
      throw new Error(msg)
    }
    await loadCatechumens()
  } catch (e: any) {
    console.error('Error eliminando catecúmeno:', e)
    alert('No se pudo eliminar el catecúmeno: ' + (e?.message || e))
  }
}

// ===== QR helpers =====
const qrClaimUrl = ref<string | null>(null)
const qrForumLeft = ref<number | null>(null)
const qrUsedCount = ref<number | null>(null)
const qrContainer = ref<HTMLElement | null>(null)
let qrStyling: any = null

async function ensureQrInstance(data: string | null) {
  if (!data) return
  try {
    if (!qrStyling) {
      const mod = await import('qr-code-styling')
      const QRCodeStyling = (mod as any).default || (mod as any)
      qrStyling = new QRCodeStyling({
        width: 256,
        height: 256,
        type: 'svg',
        data,
        qrOptions: { errorCorrectionLevel: 'Q' },
        dotsOptions: { color: '#1A388B', type: 'rounded' },
        cornersSquareOptions: { color: '#1A388B', type: 'extra-rounded' },
        cornersDotOptions: { color: '#1A388B', type: 'dot' },
        backgroundOptions: { color: '#F3F4F6' }, // gray-100
      })
    } else {
      qrStyling.update({ data })
    }
    if (qrContainer.value) {
      qrContainer.value.innerHTML = ''
      qrStyling.append(qrContainer.value)
    }
  } catch (e) {
    console.warn('QR styling load error, falling back to link only', e)
  }
}

function closeQrModal() {
  showQrModal.value = false
  qrClaimUrl.value = null
  qrForumLeft.value = null
  qrUsedCount.value = null
}

async function copy(text: string | null) {
  try { await navigator.clipboard.writeText(text || '') } catch { /* ignore */ }
}

function downloadQr() {
  const raw = String(form.value.full_name || 'catecumeno')
  // Quitar acentos y caracteres no válidos, reemplazar por _ y poner minúsculas
  const cleaned = raw
    .normalize('NFD').replace(/[\u0300-\u036f]/g, '') // sin acentos
    .replace(/[^a-zA-Z0-9]+/g, '_')                    // no alfanum -> _
    .replace(/^_+|_+$/g, '')                           // trim _
    .toLowerCase()
  const name = `qr-${cleaned || 'catecumeno'}`
  const data = qrClaimUrl.value || ''
  ;(async () => {
    try {
      // Genera un QR en alta resolución SOLO para descarga, sin afectar el tamaño en pantalla
      const mod = await import('qr-code-styling')
      const QRCodeStyling = (mod as any).default || (mod as any)
      const exportQr = new QRCodeStyling({
        width: 1024,
        height: 1024,
        type: 'svg',
        data,
        qrOptions: { errorCorrectionLevel: 'Q' },
        dotsOptions: { color: '#1A388B', type: 'rounded' },
        cornersSquareOptions: { color: '#1A388B', type: 'extra-rounded' },
        cornersDotOptions: { color: '#1A388B', type: 'dot' },
        backgroundOptions: { color: '#FFFFFF' },
      })
      exportQr.download({ name, extension: 'png' })
    } catch {
      try { qrStyling?.download({ name, extension: 'png' }) } catch { /* ignore */ }
    }
  })()
}

// Si el modal se vuelve a abrir o cambia el enlace, renderiza el QR después del próximo tick
watch([showQrModal, qrClaimUrl], async ([isOpen, link]) => {
  if (isOpen && link) {
    await nextTick()
    await ensureQrInstance(String(link))
  }
})

// Enriquecimiento de métricas de escaneos para el listado
async function enrichScanStats() {
  const items = catechumens.value.slice()
  const limit = 4
  let i = 0
  async function worker() {
    while (i < items.length) {
      const idx = i++
      const c = items[idx]
      if (!c) { continue }
      try {
        const r1 = await fetch(`${props.editUrlBase}/${c.id}`, { headers: { Accept: 'application/json' }, credentials: 'include' })
        if (!r1.ok) continue
        const data = await r1.json().catch(() => ({} as any))
        const cate = data?.catechumen || {}
        const qrId = Number(cate.qr_id ?? NaN)
        if (!Number.isNaN(qrId)) {
          const r2 = await fetch(`/qr/${qrId}`, { headers: { Accept: 'application/json' }, credentials: 'include' })
          if (r2.ok) {
            const qr = await r2.json().catch(() => ({} as any))
            const forum = Number(qr?.forum ?? qr?.Forum ?? NaN)
            const used = Number(qr?.count ?? qr?.Count ?? NaN)
            const total = (Number.isNaN(forum) ? 0 : forum) + (Number.isNaN(used) ? 0 : used)
            // aplicar directamente sobre el objeto reactivo
            const row = catechumens.value.find(x => x.id === c.id) as any
            if (row) { row.scanUsed = Number.isNaN(used) ? undefined : used; row.scanTotal = total || undefined }
          }
        }
      } catch {}
    }
  }
  // lanzar N workers
  await Promise.all(Array.from({ length: Math.min(limit, items.length) }, () => worker()))
}

// Extrae un mensaje de error legible desde una Response (JSON o texto)
async function readErrorMessage(res: Response): Promise<string> {
  // 1) Intentar leer JSON directamente
  try {
    const j: any = await res.clone().json()
    const msg = j?.message || j?.error || j?.msg
    if (msg) return String(msg)
  } catch {}

  // 2) Leer como texto y volver a intentar parsear/extraer
  let text = ''
  try { text = await res.text() } catch {}
  if (text) {
    // si es JSON en texto
    try {
      const j2: any = JSON.parse(text)
      const msg2 = j2?.message || j2?.error || j2?.msg
      if (msg2) return String(msg2)
    } catch {}
    // regex simple para "message": "..."
    const m = text.match(/"message"\s*:\s*"([^"\\]*(?:\\.[^"\\]*)*)"/)
    if (m?.[1]) return m[1].replace(/\\"/g, '"')
    // mapear frase conocida en inglés a español
    if (/Group has reached its catechumen limit/i.test(text)) {
      return 'Ya se alcanzó el límite de catecúmenos para tu grupo'
    }
    return text
  }

  // 3) Fallback por status
  if (res.status === 409) return 'Ya se alcanzó el límite de catecúmenos para tu grupo'
  return `HTTP ${res.status}`
}
</script>
