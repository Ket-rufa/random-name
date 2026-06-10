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
import AppHeader from './components/layout/AppHeader.vue';
import CustomizePanel from './components/settings/CustomizePanel.vue';
import BaseToast from './components/common/BaseToast.vue';

const store = useWheelStore();
const { initTheme, attachSystemListener } = useTheme();
const isSettingsOpen = ref<boolean>(false);

onMounted(() => {
  // Initialize theme FIRST before anything else renders
  initTheme();
  attachSystemListener();
  // Load local configurations on startup
  store.loadFromLocalStorage();
});
</script>

<style>
html, body, #app {
  height: 100%;
  margin: 0;
  padding: 0;
}
</style>
