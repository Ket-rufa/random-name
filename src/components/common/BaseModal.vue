<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="isOpen"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 dark:bg-slate-950/80 backdrop-blur-sm"
        @click.self="close"
        @keydown.esc="close"
        tabindex="-1"
      >
        <div
          :class="[
            'w-full overflow-hidden glass rounded-3xl shadow-2xl flex flex-col max-h-[90vh] transition-all duration-300 scale-transition',
            sizeClass
          ]"
          role="dialog"
          aria-modal="true"
        >
          <!-- Header -->
          <div class="flex items-center justify-between px-6 py-4 border-b border-slate-100 dark:border-slate-800/60">
            <h3 class="text-lg font-bold text-slate-900 dark:text-white">
              <slot name="title"></slot>
            </h3>
            <button
              @click="close"
              class="p-1.5 rounded-lg text-slate-400 hover:text-slate-650 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
              aria-label="Close modal"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </button>
          </div>

          <!-- Body -->
          <div class="flex-1 px-6 py-5 overflow-y-auto">
            <slot></slot>
          </div>

          <!-- Footer -->
          <div v-if="$slots.footer" class="px-6 py-4 border-t border-slate-100 dark:border-slate-800/60 flex items-center justify-end gap-3 bg-slate-50/50 dark:bg-slate-900/10">
            <slot name="footer"></slot>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, watch, onUnmounted } from 'vue';

const props = withDefaults(
  defineProps<{
    isOpen: boolean;
    size?: 'sm' | 'md' | 'lg' | 'xl';
  }>(),
  {
    size: 'md'
  }
);

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const close = () => {
  emit('close');
};

const handleKeyDown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && props.isOpen) {
    close();
  }
};

watch(() => props.isOpen, (newVal) => {
  if (newVal) {
    document.body.classList.add('overflow-hidden');
    window.addEventListener('keydown', handleKeyDown);
  } else {
    document.body.classList.remove('overflow-hidden');
    window.removeEventListener('keydown', handleKeyDown);
  }
}, { immediate: true });

onUnmounted(() => {
  document.body.classList.remove('overflow-hidden');
  window.removeEventListener('keydown', handleKeyDown);
});

const sizeClass = computed(() => {
  switch (props.size) {
    case 'sm':
      return 'max-w-md';
    case 'md':
      return 'max-w-lg';
    case 'lg':
      return 'max-w-2xl';
    case 'xl':
      return 'max-w-4xl';
    default:
      return 'max-w-lg';
  }
});
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.fade-enter-active .scale-transition {
  animation: scale-up 0.25s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.fade-leave-active .scale-transition {
  animation: scale-up 0.2s cubic-bezier(0.34, 1.56, 0.64, 1) reverse;
}

@keyframes scale-up {
  0% {
    opacity: 0;
    transform: scale(0.95) translateY(10px);
  }
  100% {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}
</style>
