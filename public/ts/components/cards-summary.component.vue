<template>
  <section class="grid grid-cols-1 sm:grid-cols-3 gap-8 mt-8">
    <!-- Groups -->
    <article class="bg-white rounded-xl border border-gray-300 p-6">
      <div class="flex justify-between items-center">
        <p>Total Grupos</p>
        <div class="aspect-square rounded-3xl size-10 bg-blue-300 text-blue-800"></div>
      </div>
      <div class="mt-12">
        <p class="text-4xl">
          {{ loading.groups ? '...' : groups.length }}
        </p>
        <p class="text-sm text-slate-600">Grupos</p>
      </div>
    </article>

    <!-- Catechists -->
    <article class="bg-white rounded-xl border border-gray-300 p-6">
      <div class="flex justify-between items-center">
        <p>Total Catequistas</p>
        <div class="aspect-square rounded-3xl size-10 bg-blue-300 text-blue-800"></div>
      </div>
      <div class="mt-12">
        <p class="text-4xl">
          {{ loading.catechists ? '...' : catechists.length }}
        </p>
        <p class="text-sm text-slate-600">Catequistas registrados</p>
      </div>
    </article>

    <!-- Catechumens -->
    <article class="bg-white rounded-xl border border-gray-300 p-6">
      <div class="flex justify-between items-center">
        <p>Total Catecúmenos</p>
        <div class="aspect-square rounded-3xl size-10 bg-blue-300 text-blue-800"></div>
      </div>
      <div class="mt-12">
        <p class="text-4xl">
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
