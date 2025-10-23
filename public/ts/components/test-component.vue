<template>
  <div class="flex flex-col items-center justify-center bg-gray-50">
    <!-- Botón que abre el modal -->
    <button
      @click="openModal = true"
      class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition"
    >
      Abrir Modal
    </button>

    <!-- Modal -->
    <div
      v-if="openModal"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
    >
      <div
        class="bg-white rounded-2xl shadow-lg p-6 w-80 text-center relative"
      >
        <h2 class="text-xl font-semibold mb-3 text-gray-800">Hola 👋</h2>
        <p class="text-gray-600 mb-5">
          Este es un ejemplo de modal con Vue y Tailwind.
        </p>

        <button
          @click="openModal = false"
          class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition"
        >
          Cerrar
        </button>

        <!-- Botón de cierre (ícono X) -->
        <button
          @click="openModal = false"
          class="absolute top-3 right-3 text-gray-500 hover:text-gray-700"
        >
          ✕
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const openModal = ref(false)
</script>
