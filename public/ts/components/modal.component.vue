<template>
  <!-- Backdrop fade -->
  <Transition
    appear
    enter-active-class="transition-opacity ease-out duration-100"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
  leave-active-class="transition-opacity ease-in duration-200"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div
      v-if="show"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
      @click.self="close"
    >
      <!-- Panel scale/opacity -->
      <Transition
        enter-active-class="transition-[transform,opacity] duration-100 ease-[cubic-bezier(0.16,1,0.3,1)]"
        enter-from-class="opacity-0 scale-0"
        enter-to-class="opacity-100 scale-100"
        leave-active-class="transition-[transform,opacity] duration-200 ease-in"
        leave-from-class="opacity-100 scale-100"
        leave-to-class="opacity-0 scale-0"
      >
        <div
          v-if="show"
          role="dialog"
          aria-modal="true"
          style="will-change: transform, opacity;"
          class="relative bg-white w-full max-w-lg  p-6 rounded-lg"
        >
          <button
            @click="close"
            class="absolute top-4 right-4 text-gray-500 hover:text-gray-700"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-x"><path d="M18 6 6 18"></path><path d="m6 6 12 12"></path></svg>
          </button>
          <slot></slot>
        </div>
      </Transition>
    </div>
  </Transition>
</template>

<script setup lang="ts">
const props = defineProps<{ show: boolean }>()
const emit = defineEmits<{ (e: 'close'): void }>()
function close() {
  emit('close')
}
</script>
