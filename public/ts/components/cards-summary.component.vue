<template>
  <section
    class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 mt-8"
  >
    <!-- Groups -->
    <article
      class="bg-white rounded-xl border border-gray-200 p-6 flex flex-col justify-between transition hover:shadow-md"
    >
      <div class="flex justify-between items-center">
        <p class="text-sm font-medium text-gray-600">Total Grupos</p>
        <div
          class="rounded-full w-10 h-10 bg-blue-100 text-blue-800 flex items-center justify-center"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="w-5 h-5"
          >
            <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path>
            <path d="M16 3.128a4 4 0 0 1 0 7.744"></path>
            <path d="M22 21v-2a4 4 0 0 0-3-3.87"></path>
            <circle cx="9" cy="7" r="4"></circle>
          </svg>
        </div>
      </div>

      <div class="mt-8 text-center sm:text-left">
        <p class="text-4xl font-semibold">
          {{ loading.groups ? '...' : groups.length }}
        </p>
        <p class="text-sm text-slate-600">Grupos</p>
      </div>
    </article>

    <!-- Catechists -->
    <article
      class="bg-white rounded-xl border border-gray-200 p-6 flex flex-col justify-between transition hover:shadow-md"
    >
      <div class="flex justify-between items-center">
        <p class="text-sm font-medium text-gray-600">Total Catequistas</p>
        <div
          class="rounded-full w-10 h-10 bg-blue-100 text-blue-800 flex items-center justify-center"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="w-5 h-5"
          >
            <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path>
            <path d="M16 3.128a4 4 0 0 1 0 7.744"></path>
            <path d="M22 21v-2a4 4 0 0 0-3-3.87"></path>
            <circle cx="9" cy="7" r="4"></circle>
          </svg>
        </div>
      </div>

      <div class="mt-8 text-center sm:text-left">
        <p class="text-4xl font-semibold">
          {{ loading.catechists ? '...' : catechists.length }}
        </p>
        <p class="text-sm text-slate-600">Catequistas registrados</p>
      </div>
    </article>

    <!-- Catechumens -->
    <article
      class="bg-white rounded-xl border border-gray-200 p-6 flex flex-col justify-between transition hover:shadow-md"
    >
      <div class="flex justify-between items-center">
        <p class="text-sm font-medium text-gray-600">Total Catecúmenos</p>
        <div
          class="rounded-full w-10 h-10 bg-blue-100 text-blue-800 flex items-center justify-center"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="w-5 h-5"
          >
            <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"></path>
            <path d="M16 3.128a4 4 0 0 1 0 7.744"></path>
            <path d="M22 21v-2a4 4 0 0 0-3-3.87"></path>
            <circle cx="9" cy="7" r="4"></circle>
          </svg>
        </div>
      </div>

      <div class="mt-8 text-center sm:text-left">
        <p class="text-4xl font-semibold">
          {{ loading.catechumens ? '...' : catechumens.length }}
        </p>
        <p class="text-sm text-slate-600">Catecúmenos inscritos</p>
      </div>
    </article>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'

const groups = ref<any[]>([])
const catechists = ref<any[]>([])
const catechumens = ref<any[]>([])

const loading = reactive({
  groups: false,
  catechists: false,
  catechumens: false,
})

async function fetchData(url: string): Promise<any[]> {
  try {
    const resp = await fetch(url, {
      credentials: 'same-origin',
      headers: { Accept: 'application/json' },
    })
    if (!resp.ok) throw new Error(`HTTP ${resp.status}`)
    const data = await resp.json()
    if (Array.isArray(data)) return data
    return (
      data.items ||
      data.users ||
      data.groups ||
      data.catechists ||
      data.catechumens ||
      data.data ||
      []
    )
  } catch (err) {
    console.error('Error cargando', url, err)
    return []
  }
}

onMounted(async () => {
  loading.groups = true
  loading.catechists = true
  loading.catechumens = true

  try {
    const [g, c, m] = await Promise.all([
      fetchData('/all-groups'),
      fetchData('/all-catechists'),
      fetchData('/catechumens'),
    ])
    groups.value = g
    catechists.value = c
    catechumens.value = m
  } finally {
    loading.groups = false
    loading.catechists = false
    loading.catechumens = false
  }
})
</script>
