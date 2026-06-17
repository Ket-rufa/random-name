<template>
  <div class="min-h-screen flex flex-col bg-slate-50 text-slate-900 dark:bg-slate-950 dark:text-slate-100 transition-colors duration-300">
    <!-- App Header -->
    <AppHeader @open-settings="isSettingsOpen = true" />

    <!-- Main Content Area -->
    <main class="flex-1 flex flex-col overflow-hidden">
      <RouterView />
    </main>

    <!-- Settings Dialog Customize Panel -->
    <CustomizePanel :is-open="isSettingsOpen" @close="isSettingsOpen = false" />

    <!-- App Footer with Visitor Statistics -->
    <footer class="py-4 px-6 border-t border-slate-100 dark:border-slate-800/60 bg-white dark:bg-slate-900/40 text-center text-xs text-slate-500 dark:text-slate-400 flex items-center justify-between flex-wrap gap-4 z-10 shrink-0">
      <div>
        © 2026 Vòng Quay May Mắn. Tất cả quyền được bảo lưu.
      </div>
      <div v-if="totalVisits !== null" class="flex items-center gap-1.5 font-bold text-slate-700 dark:text-slate-350 bg-slate-100/60 dark:bg-slate-800/60 px-3 py-1 rounded-full border border-slate-200/40 dark:border-slate-700/40">
        <svg class="h-3.5 w-3.5 text-primary-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
        </svg>
        <span>Lượt truy cập: {{ formatNumber(totalVisits) }}</span>
      </div>
    </footer>

    <!-- Global Toast Notifications -->
    <BaseToast
      :message="store.toastMessage"
      :type="store.toastType"
      @close="store.toastMessage = null"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { RouterView } from 'vue-router';
import { useWheelStore } from './stores/wheelStore';
import { useTheme } from './composables/useTheme';
import { wheelApi } from './api/wheelApi';
import AppHeader from './components/layout/AppHeader.vue';
import CustomizePanel from './components/settings/CustomizePanel.vue';
import BaseToast from './components/common/BaseToast.vue';

const store = useWheelStore();
const { initTheme, attachSystemListener } = useTheme();
const isSettingsOpen = ref<boolean>(false);
const totalVisits = ref<number | null>(null);

const formatNumber = (num: number) => {
  return new Intl.NumberFormat('vi-VN').format(num);
};

onMounted(async () => {
  // Initialize theme FIRST before anything else renders
  initTheme();
  attachSystemListener();
  // Load local configurations on startup
  store.loadFromLocalStorage();

  // Record page visit
  try {
    const res = await wheelApi.recordVisit();
    if (res.data && res.data.success) {
      totalVisits.value = res.data.data.totalVisits;
    }
  } catch (err) {
    // Fail silently
  }
});
</script>

<style>
html, body, #app {
  height: 100%;
  margin: 0;
  padding: 0;
}
</style>
