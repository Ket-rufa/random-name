<template>
  <header class="h-16 px-6 border-b border-slate-100 dark:border-slate-800/60 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md flex items-center justify-between z-40 relative transition-all duration-300">
    <!-- Logo & Brand Title -->
    <div class="flex items-center gap-3 cursor-pointer" @click="goHome">
      <div class="w-10 h-10 rounded-xl overflow-hidden shadow-lg shadow-purple-500/30 flex items-center justify-center bg-indigo-950">
        <img src="/favicon.svg" alt="Logo" class="w-8 h-8" />
      </div>
      <h1 class="text-base font-black tracking-tight bg-linear-to-r from-slate-900 to-slate-700 dark:from-white dark:to-slate-300 bg-clip-text text-transparent hidden sm:block">
        VÒNG QUAY MAY MẮN
      </h1>
    </div>

    <!-- Right Control Actions -->
    <div class="flex items-center gap-2">
      <!-- Create New -->
      <BaseButton v-if="store.canEdit" variant="ghost" size="sm" @click="handleNew">
        <template #icon><PlusIcon class="h-4 w-4" /></template>
        <span class="hidden md:inline">Tạo mới</span>
      </BaseButton>

      <!-- Save -->
      <BaseButton v-if="store.canEdit" variant="ghost" size="sm" @click="handleSave" :loading="store.isLoading">
        <template #icon><SaveIcon class="h-4 w-4" /></template>
        <span class="hidden md:inline">Lưu server</span>
      </BaseButton>

      <!-- Share -->
      <BaseButton variant="ghost" size="sm" @click="handleShare">
        <template #icon><Share2Icon class="h-4 w-4" /></template>
        <span>Chia sẻ</span>
      </BaseButton>

      <!-- Customize -->
      <BaseButton v-if="store.canEdit" variant="ghost" size="sm" @click="emit('open-settings')">
        <template #icon><SlidersIcon class="h-4 w-4" /></template>
        <span>Cấu hình</span>
      </BaseButton>

      <!-- Quick Sound Toggle -->
      <button
        @click="toggleSound"
        class="p-2 rounded-xl text-slate-500 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800/60 transition-colors cursor-pointer"
        :title="store.settings.enableSound ? 'Tắt âm thanh' : 'Bật âm thanh'"
      >
        <Volume2Icon v-if="store.settings.enableSound" class="h-5 w-5" />
        <VolumeXIcon v-else class="h-5 w-5 text-rose-500" />
      </button>

      <!-- Light/Dark Mode Toggle -->
      <button
        @click="toggle"
        class="p-2 rounded-xl text-slate-500 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800/60 transition-colors cursor-pointer relative overflow-hidden"
        title="Chuyển giao diện sáng/tối"
        aria-label="Toggle dark mode"
      >
        <!-- Sun icon shown in dark mode (click to go light) -->
        <Transition name="icon-swap">
          <SunIcon v-if="isDark" key="sun" class="h-5 w-5 text-amber-400" />
          <MoonIcon v-else key="moon" class="h-5 w-5" />
        </Transition>
      </button>
    </div>
  </header>
</template>

<script setup lang="ts">
import { useWheelStore } from '../../stores/wheelStore';
import { useTheme } from '../../composables/useTheme';
import BaseButton from '../common/BaseButton.vue';
import { PlusIcon, SaveIcon, Share2Icon, SlidersIcon, Volume2Icon, VolumeXIcon, SunIcon, MoonIcon } from 'lucide-vue-next';
import { useRouter } from 'vue-router';

const store = useWheelStore();
const router = useRouter();
const { isDark, toggle } = useTheme();

const emit = defineEmits<{
  (e: 'open-settings'): void;
}>();

const goHome = () => {
  router.push('/');
};

const handleNew = () => {
  if (confirm('Bạn có chắc chắn muốn tạo mới vòng quay? Thao tác này sẽ xóa dữ liệu hiện tại.')) {
    store.resetWheel();
    store.showToast('Đã khởi tạo vòng quay mới', 'success');
    router.push('/');
  }
};

const handleSave = async () => {
  store.showToast('Đang kết nối server để lưu...', 'info');
  await store.saveWheelToServer();
};

const handleShare = () => {
  if (store.shareCode) {
    const shareUrl = `${window.location.origin}/wheels/${store.shareCode}`;
    navigator.clipboard.writeText(shareUrl)
      .then(() => {
        store.showToast('Đã sao chép liên kết chia sẻ vào clipboard!', 'success');
      })
      .catch(() => {
        store.showToast(`Liên kết: ${shareUrl}`, 'info');
      });
  } else {
    store.showToast('Vui lòng lưu vòng quay lên server trước khi chia sẻ!', 'info');
  }
};

const toggleSound = () => {
  store.settings.enableSound = !store.settings.enableSound;
  store.saveToLocalStorage();
  store.showToast(store.settings.enableSound ? 'Đã bật âm thanh' : 'Đã tắt âm thanh', 'info');
};
</script>

<style scoped>
/* Smooth icon swap animation */
.icon-swap-enter-active,
.icon-swap-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
  position: absolute;
}

.icon-swap-enter-from {
  opacity: 0;
  transform: rotate(-90deg) scale(0.5);
}

.icon-swap-leave-to {
  opacity: 0;
  transform: rotate(90deg) scale(0.5);
}
</style>
