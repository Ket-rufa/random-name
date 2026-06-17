<template>
  <div class="flex-1 flex flex-col md:flex-row gap-6 p-6 overflow-hidden content-height">
    <!-- Loading State -->
    <div v-if="store.isLoading" class="flex-1 flex flex-col items-center justify-center bg-white dark:bg-slate-900 rounded-3xl border border-slate-100 dark:border-slate-800/60 shadow-sm p-6 min-h-[400px]">
      <div class="flex flex-col items-center gap-3">
        <svg class="animate-spin h-10 w-10 text-primary-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <span class="text-sm font-semibold text-slate-500 dark:text-slate-400">Đang tải vòng quay...</span>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="store.isError" class="flex-1 flex flex-col items-center justify-center bg-white dark:bg-slate-900 rounded-3xl border border-slate-100 dark:border-slate-800/60 shadow-sm p-6 min-h-[400px] gap-4">
      <div class="w-16 h-16 rounded-full bg-rose-50 dark:bg-rose-950/20 flex items-center justify-center text-rose-500 text-3xl">
        ⚠️
      </div>
      <h3 class="text-lg font-bold text-slate-800 dark:text-slate-200">Không tìm thấy vòng quay</h3>
      <p class="text-sm text-slate-400 dark:text-slate-500 text-center max-w-md">
        {{ store.errorMessage || 'Đường dẫn chia sẻ này không tồn tại hoặc đã bị xóa.' }}
      </p>
      <BaseButton variant="primary" @click="goHome">
        Về trang chủ
      </BaseButton>
    </div>

    <!-- Standard Dashboard Viewport -->
    <template v-else>
      <!-- Left Column: Canvas -->
      <div class="flex-1 flex flex-col items-center justify-center bg-white dark:bg-slate-900 rounded-3xl border border-slate-100 dark:border-slate-800/60 shadow-sm p-6 relative min-h-[380px] transition-all duration-300">
        <!-- Permission warning banners -->
        <div v-if="store.permission === 'view'" class="absolute top-4 px-4 py-1.5 rounded-full bg-slate-100 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-xs text-slate-650 dark:text-slate-350 font-semibold z-10">
          🔒 Vòng quay ở chế độ chỉ xem. Bạn không thể quay.
        </div>
        <div v-else-if="store.permission === 'spin'" class="absolute top-4 px-4 py-1.5 rounded-full bg-blue-50 dark:bg-blue-950/20 border border-blue-200 dark:border-blue-900 text-xs text-blue-700 dark:text-blue-300 font-semibold z-10">
          🎮 Vòng quay ở chế độ chỉ quay. Bạn không thể sửa danh sách.
        </div>

        <WheelCanvas
          :entries="store.entries"
          :settings="store.settings"
          :rotation-angle="rotationAngle"
          :is-spinning="isSpinning"
          :get-entry-color="store.getEntryColor"
          @spin="handleSpin"
        />
      </div>

      <!-- Right Column -->
      <div class="w-full md:w-[400px] flex flex-col gap-6 h-full min-h-[400px]">
        <!-- Entry list (Disabled if not editor) -->
        <EntryEditor class="flex-1 min-h-[250px]" />
        
        <!-- Local spin history log -->
        <SpinHistory />
      </div>

      <!-- Result Modal -->
      <BaseModal :is-open="store.isResultModalOpen" @close="closeResultModal" size="sm">
        <template #title>
          <div class="text-center w-full">🎉 KẾT QUẢ</div>
        </template>

        <div class="text-center py-6 flex flex-col items-center justify-center gap-4">
          <span class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-widest">Lựa chọn trúng giải</span>
          <h2 class="text-3xl font-black bg-linear-to-r from-primary-600 to-indigo-500 bg-clip-text text-transparent wrap-break-word max-w-full px-4">
            {{ store.selectedResult?.label }}
          </h2>
          <div class="w-16 h-1 bg-linear-to-r from-primary-600 to-indigo-500 rounded-full my-2"></div>
        </div>

        <template #footer>
          <div class="flex items-center justify-center gap-3 w-full">
            <BaseButton variant="secondary" @click="closeResultModal" class="flex-1">
              Đóng
            </BaseButton>
            <BaseButton v-if="store.canEdit" variant="danger" @click="removeWinnerFromList" class="flex-1">
              Xóa khỏi danh sách
            </BaseButton>
          </div>
        </template>
      </BaseModal>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useWheelStore } from '../stores/wheelStore';
import WheelCanvas from '../components/wheel/WheelCanvas.vue';
import EntryEditor from '../components/entries/EntryEditor.vue';
import SpinHistory from '../components/history/SpinHistory.vue';
import BaseModal from '../components/common/BaseModal.vue';
import BaseButton from '../components/common/BaseButton.vue';
import { pickSecureRandom } from '../composables/useSecureRandom';
import { useWheelAnimation } from '../composables/useWheelAnimation';
import { useAudio } from '../composables/useAudio';
import { useConfetti } from '../composables/useConfetti';
import { wheelApi } from '../api/wheelApi';

const store = useWheelStore();
const route = useRoute();
const router = useRouter();

const { rotationAngle, isSpinning, spin, stop } = useWheelAnimation();
const { playVictory } = useAudio();
const { fireConfetti } = useConfetti();

const goHome = () => {
  router.push('/');
};

const fetchWheelData = async (code: string) => {
  store.isLoading = true;
  store.isError = false;
  store.errorMessage = '';
  
  try {
    const res = await wheelApi.getWheel(code);
    if (res.data && res.data.success) {
      const { wheel, editToken: localEditToken } = res.data.data;
      store.setWheelDataFromServer(wheel, localEditToken);
      store.showToast('Tải vòng quay thành công!', 'success');
    } else {
      store.isError = true;
      store.errorMessage = res.data.message || 'Lỗi tải vòng quay từ server';
    }
  } catch (err: any) {
    store.isError = true;
    store.errorMessage = err.response?.data?.message || 'Không thể kết nối đến máy chủ.';
    
    // Fallback if backend is unavailable: Load default local state
    store.showToast('Không thể kết nối máy chủ. Đang dùng dữ liệu cục bộ.', 'error');
    store.loadFromLocalStorage();
  } finally {
    store.isLoading = false;
  }
};

const handleSpin = () => {
  if (isSpinning.value || store.permission === 'view') return;
  if (store.entries.length === 0) {
    store.showToast('Vui lòng thêm lựa chọn trước khi quay!', 'error');
    return;
  }

  const winner = pickSecureRandom(store.entries, store.settings.enableWeights);
  
  store.isSpinning = true;
  spin(
    store.entries,
    store.settings.enableWeights,
    winner,
    store.settings.spinDuration,
    store.settings.volume,
    store.settings.enableTickSound,
    async () => {
      store.isSpinning = false;
      store.selectedResult = winner;
      
      store.addHistory(winner.label, winner.id);
      
      if (store.settings.enableSound && store.settings.enableVictorySound) {
        playVictory(store.settings.volume);
      }
      if (store.settings.enableConfetti) {
        fireConfetti();
      }

      store.isResultModalOpen = true;

      // Post spin analytics to server asynchronously
      if (store.id) {
        try {
          await wheelApi.recordSpin(store.id, winner.id, winner.label);
        } catch (e) {
          // Fail silently
        }
      }

      if (store.settings.autoRemoveWinner) {
        setTimeout(() => {
          if (store.selectedResult && store.selectedResult.id === winner.id) {
            removeWinnerFromList(false);
          }
        }, 1200);
      }
    }
  );
};

const closeResultModal = () => {
  store.isResultModalOpen = false;
};

const removeWinnerFromList = (showToastNotification = true) => {
  if (store.selectedResult) {
    const label = store.selectedResult.label;
    store.removeEntry(store.selectedResult.id);
    store.isResultModalOpen = false;
    
    if (showToastNotification) {
      store.showToast(`Đã xóa "${label}" khỏi danh sách.`, 'info');
    }
  }
};

const handleKeyDown = (e: KeyboardEvent) => {
  if (e.code === 'Space' && !isSpinning.value && !store.isResultModalOpen) {
    const active = document.activeElement;
    if (active && (active.tagName === 'INPUT' || active.tagName === 'TEXTAREA')) {
      return;
    }
    e.preventDefault();
    handleSpin();
  }
};

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown);
  const code = route.params.shareCode as string;
  if (code) {
    fetchWheelData(code);
  }
});

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown);
  stop();
});

// Watch for route shareCode change if visitor jumps to another shared URL
watch(
  () => route.params.shareCode,
  (newCode) => {
    if (newCode) {
      fetchWheelData(newCode as string);
    }
  }
);
</script>
