<template>
  <button
    :class="[
      'inline-flex items-center justify-center font-medium rounded-xl transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed active:scale-95',
      variantClass,
      sizeClass
    ]"
    :disabled="disabled || loading"
    v-bind="$attrs"
  >
    <svg v-if="loading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-current" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>
    <span v-else class="mr-2 inline-flex items-center justify-center">
      <slot name="icon"></slot>
    </span>
    <slot></slot>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'danger' | 'ghost' | 'glass';
    size?: 'sm' | 'md' | 'lg';
    disabled?: boolean;
    loading?: boolean;
  }>(),
  {
    variant: 'primary',
    size: 'md',
    disabled: false,
    loading: false
  }
);

const variantClass = computed(() => {
  switch (props.variant) {
    case 'primary':
      return 'bg-primary-600 hover:bg-primary-700 text-white shadow-lg shadow-primary-500/20 focus:ring-primary-500';
    case 'secondary':
      return 'bg-slate-200 hover:bg-slate-300 text-slate-800 focus:ring-slate-500 dark:bg-slate-800 dark:hover:bg-slate-700 dark:text-slate-200';
    case 'danger':
      return 'bg-red-600 hover:bg-red-700 text-white shadow-lg shadow-red-500/20 focus:ring-red-500';
    case 'ghost':
      return 'hover:bg-slate-100 dark:hover:bg-slate-900 text-slate-700 dark:text-slate-300 focus:ring-slate-500';
    case 'glass':
      return 'glass text-slate-850 dark:text-slate-100 hover:bg-white/80 dark:hover:bg-slate-900/80 shadow-sm focus:ring-primary-500';
    default:
      return '';
  }
});

const sizeClass = computed(() => {
  switch (props.size) {
    case 'sm':
      return 'px-3 py-1.5 text-xs';
    case 'md':
      return 'px-4 py-2 text-sm';
    case 'lg':
      return 'px-6 py-3 text-base';
    default:
      return '';
  }
});
</script>
