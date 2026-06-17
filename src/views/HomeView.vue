<template>
  <div class="flex-1 flex flex-col md:flex-row gap-6 p-6 overflow-hidden content-height">
    <!-- Left Column: Wheel Canvas & Pointer -->
    <div class="flex-1 flex flex-col items-center justify-center bg-white dark:bg-slate-900 rounded-3xl border border-slate-100 dark:border-slate-800/60 shadow-sm p-6 relative min-h-[380px] transition-all duration-300">
      <div v-if="store.entries.length === 1" class="absolute top-4 px-4 py-1.5 rounded-full bg-amber-50 dark:bg-amber-950/20 border border-amber-250 dark:border-amber-900 text-xs text-amber-700 dark:text-amber-300 font-semibold z-10 flex items-center gap-1.5">
        ⚠️ Thêm ít nhất 2 lựa chọn để vòng quay hoạt động tốt hơn.
      </div>
      
      <!-- The Canvas Component -->
      <WheelCanvas
        :entries="store.entries"
        :settings="store.settings"
        :rotation-angle="rotationAngle"
        :is-spinning="isSpinning"
        :get-entry-color="store.getEntryColor"
        @spin="handleSpin"
      />
    </div>

    <!-- Right Column: Entry Editor & History -->
    <div class="w-full md:w-[380px] flex flex-col gap-6 h-full min-h-[400px]">
      <!-- Entry Editor -->
      <EntryEditor class="flex-1 min-h-[250px]" />
      
      <!-- Spin History -->
      <SpinHistory />
    </div>

    <!-- Victory Result Modal -->
    <BaseModal :is-open="store.isResultModalOpen" @close="closeResultModal" size="sm" :header-color="winnerColor">
      <template #title>
        <div class="text-center w-full">🎉 KẾT QUẢ</div>
      </template>

      <div class="text-center py-6 flex flex-col items-center justify-center gap-4">
        <!-- Winner color swatch -->
        <div
          class="w-14 h-14 rounded-full shadow-lg flex items-center justify-center text-2xl"
          :style="{ backgroundColor: winnerColor, boxShadow: `0 0 24px ${winnerColor}66` }"
        >
          🎉
        </div>
        <span class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-widest">Lựa chọn trúng giải</span>
        <h2
          class="text-3xl font-black break-words max-w-full px-4"
          :style="{ color: winnerColor }"
        >
          {{ store.selectedResult?.label }}
        </h2>
        <div
          class="w-16 h-1.5 rounded-full my-1"
          :style="{ backgroundColor: winnerColor }"
        ></div>
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
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue';
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

const store = useWheelStore();
const { rotationAngle, isSpinning, spin, stop } = useWheelAnimation();
const { playVictory } = useAudio();
const { fireConfetti } = useConfetti();

// Compute the wheel segment color of the current winner
const winnerColor = computed(() => {
  const result = store.selectedResult;
  if (!result) return '#8b5cf6';
  const idx = store.entries.findIndex(e => e.id === result.id);
  return store.getEntryColor(idx >= 0 ? idx : 0, result.color);
});

const handleSpin = () => {
  if (isSpinning.value) return;
  if (store.entries.length === 0) {
    store.showToast('Vui lòng thêm lựa chọn trước khi quay!', 'error');
    return;
  }

  // 1. Pick winner cryptographically
  const winner = pickSecureRandom(store.entries, store.settings.enableWeights);
  
  // 2. Perform animation
  store.isSpinning = true;
  spin(
    store.entries,
    store.settings.enableWeights,
    winner,
    store.settings.spinDuration,
    store.settings.volume,
    store.settings.enableTickSound,
    () => {
      store.isSpinning = false;
      store.selectedResult = winner;
      
      // Append to local log history
      store.addHistory(winner.label, winner.id);
      
      // Celebration triggers
      if (store.settings.enableSound && store.settings.enableVictorySound) {
        playVictory(store.settings.volume);
      }
      if (store.settings.enableConfetti) {
        fireConfetti();
      }

      // Display Modal popup
      store.isResultModalOpen = true;

      // Handle auto-remove winner if toggled
      if (store.settings.autoRemoveWinner) {
        setTimeout(() => {
          if (store.selectedResult && store.selectedResult.id === winner.id) {
            removeWinnerFromList(false); // Quietly delete winner
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

// Keydown listener to trigger spin on spacebar presses
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
});

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown);
  stop();
});
</script>
