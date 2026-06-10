<template>
  <div class="relative w-full max-w-[500px] aspect-square flex items-center justify-center select-none">
    <!-- Outer glow border ring -->
    <div class="absolute inset-0 rounded-full border-4 border-white dark:border-slate-800 shadow-[0_0_35px_rgba(133,28,255,0.15)] pointer-events-none z-0"></div>
    
    <!-- Canvas Element -->
    <canvas
      ref="canvasRef"
      class="w-full h-full cursor-pointer transition-transform duration-75"
      :style="{ transform: `rotate(${rotationAngle}deg)` }"
      @click="handleWheelClick"
    ></canvas>

    <!-- Top Pointer (Mũi tên chỉ góc 270 độ) -->
    <div class="absolute -top-1 left-1/2 -translate-x-1/2 z-25 flex flex-col items-center pointer-events-none">
      <div class="w-8 h-10 bg-gradient-to-b from-rose-500 to-rose-600 rounded-b-xl shadow-lg relative after:content-[''] after:absolute after:bottom-[-8px] after:left-1/2 after:-translate-x-1/2 after:border-t-[8px] after:border-t-rose-650 after:border-x-[12px] after:border-x-transparent"></div>
      <div class="w-3 h-3 rounded-full bg-white dark:bg-slate-950 border-2 border-rose-500 -mt-1 shadow-inner"></div>
    </div>

    <!-- Center SPIN Button -->
    <button
      @click="triggerSpin"
      :disabled="isSpinning || entries.length === 0"
      class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-20 w-20 h-20 rounded-full flex flex-col items-center justify-center border-4 border-white dark:border-slate-800 shadow-[0_10px_25px_rgba(0,0,0,0.25)] text-slate-800 dark:text-white transition-all duration-300 font-extrabold focus:outline-none"
      :class="[
        isSpinning 
          ? 'bg-slate-200 dark:bg-slate-700 opacity-80 cursor-not-allowed scale-95' 
          : 'bg-white/95 dark:bg-slate-900/95 hover:scale-110 active:scale-95 hover:shadow-[0_12px_30px_rgba(133,28,255,0.4)] cursor-pointer'
      ]"
      aria-label="Quay vòng quay"
    >
      <span class="text-xs uppercase tracking-widest text-slate-400 dark:text-slate-500 font-semibold mb-0.5">QUAY</span>
      <span class="text-lg font-black leading-none bg-gradient-to-r from-primary-600 to-violet-500 bg-clip-text text-transparent">SPIN</span>
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue';
import type { WheelEntry, WheelSettings } from '../../types/wheel';

const props = defineProps<{
  entries: WheelEntry[];
  settings: WheelSettings;
  rotationAngle: number;
  isSpinning: boolean;
  getEntryColor: (index: number, entryColor?: string) => string;
}>();

const emit = defineEmits<{
  (e: 'spin'): void;
}>();

const canvasRef = ref<HTMLCanvasElement | null>(null);
let resizeObserver: ResizeObserver | null = null;

const triggerSpin = () => {
  if (!props.isSpinning && props.entries.length > 0) {
    emit('spin');
  }
};

const handleWheelClick = () => {
  triggerSpin();
};

// Draw logic
const drawWheel = () => {
  const canvas = canvasRef.value;
  if (!canvas) return;

  const ctx = canvas.getContext('2d');
  if (!ctx) return;

  const dpr = window.devicePixelRatio || 1;
  const rect = canvas.getBoundingClientRect();
  const width = rect.width;
  const height = rect.height;

  // Set physical drawing buffer size
  canvas.width = width * dpr;
  canvas.height = height * dpr;
  ctx.scale(dpr, dpr);

  const cx = width / 2;
  const cy = height / 2;
  const radius = Math.min(cx, cy) - 8; // margin for pointer shadows

  ctx.clearRect(0, 0, width, height);

  const totalEntries = props.entries.length;
  if (totalEntries === 0) {
    // Draw empty state circle
    ctx.beginPath();
    ctx.arc(cx, cy, radius, 0, 2 * Math.PI);
    ctx.fillStyle = '#f1f5f9';
    ctx.fill();
    ctx.lineWidth = 4;
    ctx.strokeStyle = '#cbd5e1';
    ctx.stroke();
    
    ctx.fillStyle = '#64748b';
    ctx.font = '16px Outfit, sans-serif';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText('Nhập danh sách lựa chọn', cx, cy);
    return;
  }

  // Calculate slice size angles (weighted or uniform)
  const totalWeight = props.entries.reduce((sum, e) => sum + (props.settings.enableWeights ? e.weight : 1), 0);
  let currentAngle = 0;

  props.entries.forEach((entry, index) => {
    const weight = props.settings.enableWeights ? entry.weight : 1;
    const sliceAngle = (weight / totalWeight) * 2 * Math.PI;
    const startAngle = currentAngle;
    const endAngle = currentAngle + sliceAngle;
    currentAngle = endAngle;

    // 1. Draw wedge segment
    ctx.beginPath();
    ctx.moveTo(cx, cy);
    ctx.arc(cx, cy, radius, startAngle, endAngle);
    ctx.closePath();
    ctx.fillStyle = props.getEntryColor(index, entry.color);
    ctx.fill();

    // 2. Draw border line
    ctx.lineWidth = 1.5;
    ctx.strokeStyle = 'rgba(255, 255, 255, 0.25)';
    ctx.stroke();

    // 3. Draw text label
    if (props.settings.showTextOnWheel) {
      ctx.save();
      // Move origin to center
      ctx.translate(cx, cy);
      // Rotate to the center of the slice
      const textAngle = startAngle + sliceAngle / 2;
      ctx.rotate(textAngle);

      // Determine text color based on background luminance for high contrast readability
      const bgColor = props.getEntryColor(index, entry.color);
      const isBgDark = isColorDark(bgColor);
      ctx.fillStyle = isBgDark ? '#ffffff' : '#1e293b';

      // Font size configuration
      const sizeMultiplier = props.settings.fontSize / 18;
      const calculatedFontSize = Math.max(9, Math.min(22, (radius / 15) * sizeMultiplier));
      ctx.font = `600 ${calculatedFontSize}px Outfit, sans-serif`;
      ctx.textAlign = 'right';
      ctx.textBaseline = 'middle';

      // Truncate text if too long
      let text = entry.label;
      const maxTextLength = Math.max(5, Math.floor(radius / (calculatedFontSize * 0.7)));
      if (text.length > maxTextLength) {
        text = text.substring(0, maxTextLength - 2) + '..';
      }

      // Draw text starting from near the outer edge inwards
      const textOffset = radius - 15;
      ctx.fillText(text, textOffset, 0);
      ctx.restore();
    }
  });

  // 4. Draw outer border ring on top
  ctx.beginPath();
  ctx.arc(cx, cy, radius, 0, 2 * Math.PI);
  ctx.lineWidth = 5;
  ctx.strokeStyle = '#ffffff';
  ctx.stroke();

  // Draw a darker ring in dark mode for contrast
  const isDarkMode = document.body.classList.contains('dark');
  if (isDarkMode) {
    ctx.beginPath();
    ctx.arc(cx, cy, radius, 0, 2 * Math.PI);
    ctx.lineWidth = 2;
    ctx.strokeStyle = '#334155';
    ctx.stroke();
  }
};

// Utility: check if hex color is dark
const isColorDark = (hex: string): boolean => {
  let color = hex.replace('#', '');
  if (color.length === 3) {
    color = color[0] + color[0] + color[1] + color[1] + color[2] + color[2];
  }
  const r = parseInt(color.substring(0, 2), 16);
  const g = parseInt(color.substring(2, 4), 16);
  const b = parseInt(color.substring(4, 6), 16);
  // YIQ formula
  const yiq = (r * 299 + g * 587 + b * 114) / 1000;
  return yiq < 145;
};

// Lifecycle
onMounted(() => {
  drawWheel();

  // Watch for element size changes to redraw sharp
  if (canvasRef.value) {
    resizeObserver = new ResizeObserver(() => {
      drawWheel();
    });
    resizeObserver.observe(canvasRef.value);
  }
});

onUnmounted(() => {
  if (resizeObserver) {
    resizeObserver.disconnect();
  }
});

// Watch reactive props to trigger redraw
watch(() => [props.entries, props.settings], () => {
  drawWheel();
}, { deep: true });

// Listen to dark mode mutations
const observer = new MutationObserver(() => {
  drawWheel();
});

onMounted(() => {
  observer.observe(document.body, { attributes: true, attributeFilter: ['class'] });
});

onUnmounted(() => {
  observer.disconnect();
});
</script>

<style scoped>
canvas {
  image-rendering: -webkit-optimize-contrast;
  image-rendering: crisp-edges;
}
</style>
