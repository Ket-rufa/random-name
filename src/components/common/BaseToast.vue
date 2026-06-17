<template>
  <Transition name="slide-up">
    <div
      v-if="message"
      :class="[
        'fixed bottom-6 right-6 z-9999 flex items-center gap-3 px-5 py-3.5 rounded-2xl shadow-xl border backdrop-blur-md transition-all duration-300 max-w-sm',
        toastClass
      ]"
    >
      <component :is="icon" class="h-5 w-5 shrink-0" />
      <span class="text-sm font-medium">{{ message }}</span>
      <button @click="close" class="ml-auto p-0.5 rounded-md hover:bg-black/5 dark:hover:bg-white/5 transition-colors">
        <XIcon class="h-4 w-4 opacity-60 hover:opacity-100" />
      </button>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { CheckCircleIcon, AlertCircleIcon, InfoIcon, XIcon } from 'lucide-vue-next';

const props = withDefaults(
  defineProps<{
    message: string | null;
    type?: 'success' | 'error' | 'info';
  }>(),
  {
    type: 'info'
  }
);

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const close = () => {
  emit('close');
};

const toastClass = computed(() => {
  switch (props.type) {
    case 'success':
      return 'bg-emerald-50/90 border-emerald-250 text-emerald-800 dark:bg-emerald-950/90 dark:border-emerald-800 dark:text-emerald-200';
    case 'error':
      return 'bg-rose-50/90 border-rose-250 text-rose-800 dark:bg-rose-950/90 dark:border-rose-800 dark:text-rose-200';
    case 'info':
    default:
      return 'bg-slate-50/90 border-slate-250 text-slate-800 dark:bg-slate-900/90 dark:border-slate-800 dark:text-slate-200';
  }
});

const icon = computed(() => {
  switch (props.type) {
    case 'success':
      return CheckCircleIcon;
    case 'error':
      return AlertCircleIcon;
    case 'info':
    default:
      return InfoIcon;
  }
});
</script>

<style scoped>
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.35s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}
.slide-up-enter-from {
  opacity: 0;
  transform: translateY(20px) scale(0.95);
}
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(-20px) scale(0.95);
}
</style>
