<template>
  <nav
    class="sticky top-0 z-50 backdrop-blur-md bg-white/80 border-b border-gray-100 px-4 sm:px-8 py-3 flex items-center justify-between"
  >
    <!-- Left: Logo + Title -->
    <div class="flex items-center gap-3">
      <div
        class="w-10 h-10 rounded-xl bg-[#1A388B]/10 flex items-center justify-center"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="w-5 h-5 text-[#1A388B]"
        >
          <path d="M12 7v14" />
          <path
            d="M3 18a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1h5a4 4 0 0 1 4 4 4 4 0 0 1 4-4h5a1 1 0 0 1 1 1v13a1 1 0 0 1-1 1h-6a3 3 0 0 0-3 3 3 3 0 0 0-3-3z"
          />
        </svg>
      </div>

      <div class="leading-tight">
        <h1
          class="text-[15px] sm:text-lg font-semibold text-gray-900 tracking-tight"
        >
          Sistema de Catequesis
        </h1>
        <p class="text-xs text-gray-500 hidden sm:block">
          Parroquia San José Obrero
        </p>
      </div>
    </div>

    <!-- Right: Logout Button -->
    <button
      class="flex items-center gap-2 px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-lg  hover:bg-gray-200 transition-all duration-150 active:scale-[0.98]"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="18"
        height="18"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
        class="lucide lucide-log-out"
      >
        <path d="m16 17 5-5-5-5" />
        <path d="M21 12H9" />
        <path
          d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"
        />
      </svg>
      <span>Salir</span>
    </button>
  </nav>
</template>
