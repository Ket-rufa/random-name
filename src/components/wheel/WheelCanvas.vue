<template>
  <div class="relative w-full max-w-[500px] aspect-square flex items-center justify-center select-none">
    <!-- Canvas Element -->
    <canvas
      ref="canvasRef"
      class="w-full h-full cursor-pointer"
      :style="{ transform: `rotate(${rotationAngle}deg)` }"
      @click="handleWheelClick"
    ></canvas>

    <!-- Top Pointer -->
    <div class="absolute -top-3 left-1/2 -translate-x-1/2 z-25 flex flex-col items-center pointer-events-none drop-shadow-lg">
      <svg width="32" height="40" viewBox="0 0 32 40" fill="none">
        <polygon points="16,38 0,2 32,2" fill="url(#pointerGrad)" filter="url(#shadow)"/>
        <defs>
          <linearGradient id="pointerGrad" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0%" stop-color="#ff4444"/>
            <stop offset="100%" stop-color="#cc0000"/>
          </linearGradient>
          <filter id="shadow" x="-20%" y="-20%" width="140%" height="140%">
            <feDropShadow dx="0" dy="2" stdDeviation="2" flood-opacity="0.4"/>
          </filter>
        </defs>
      </svg>
    </div>

    <!-- Center SPIN Button (golden hub) -->
    <button
      @click="triggerSpin"
      :disabled="isSpinning || entries.length === 0"
      class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-20 w-[17%] h-[17%] rounded-full flex flex-col items-center justify-center focus:outline-none transition-all duration-200"
      :class="[
        isSpinning
          ? 'opacity-80 cursor-not-allowed scale-95'
          : 'hover:scale-110 active:scale-95 cursor-pointer'
      ]"
      :style="{
        background: 'radial-gradient(circle at 38% 35%, #ffe566, #f5a623 55%, #c47e0a)',
        boxShadow: '0 4px 16px rgba(0,0,0,0.35), inset 0 2px 4px rgba(255,255,180,0.6), inset 0 -2px 4px rgba(100,50,0,0.3)'
      }"
      aria-label="Quay vòng quay"
    >
      <span class="text-[10px] uppercase tracking-widest font-semibold leading-none" style="color:#7a3f00; text-shadow: 0 1px 0 rgba(255,200,100,0.8)">QUAY</span>
      <span class="text-sm font-black leading-none" style="color:#5a2d00; text-shadow: 0 1px 0 rgba(255,200,100,0.8)">SPIN</span>
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

// Utility: check if hex color is dark
const isColorDark = (hex: string): boolean => {
  let color = hex.replace('#', '');
  if (color.length === 3) {
    color = color[0] + color[0] + color[1] + color[1] + color[2] + color[2];
  }
  const r = parseInt(color.substring(0, 2), 16);
  const g = parseInt(color.substring(2, 4), 16);
  const b = parseInt(color.substring(4, 6), 16);
  const yiq = (r * 299 + g * 587 + b * 114) / 1000;
  return yiq < 145;
};

// Lighten a hex color by amount (0-1)
const lightenColor = (hex: string, amount: number): string => {
  let color = hex.replace('#', '');
  if (color.length === 3) {
    color = color[0] + color[0] + color[1] + color[1] + color[2] + color[2];
  }
  const r = Math.min(255, parseInt(color.substring(0, 2), 16) + Math.round(255 * amount));
  const g = Math.min(255, parseInt(color.substring(2, 4), 16) + Math.round(255 * amount));
  const b = Math.min(255, parseInt(color.substring(4, 6), 16) + Math.round(255 * amount));
  return `#${r.toString(16).padStart(2, '0')}${g.toString(16).padStart(2, '0')}${b.toString(16).padStart(2, '0')}`;
};

const drawWheel = () => {
  const canvas = canvasRef.value;
  if (!canvas) return;
  const ctx = canvas.getContext('2d');
  if (!ctx) return;

  const dpr = window.devicePixelRatio || 1;
  const rect = canvas.getBoundingClientRect();
  const width = rect.width;
  const height = rect.height;

  canvas.width = width * dpr;
  canvas.height = height * dpr;
  ctx.scale(dpr, dpr);

  const cx = width / 2;
  const cy = height / 2;

  // Ring dimensions
  const outerRingRadius = Math.min(cx, cy) - 4;
  const ringThickness = outerRingRadius * 0.12;
  const wheelRadius = outerRingRadius - ringThickness;
  const ballRadius = ringThickness * 0.42;
  const ballCount = Math.max(16, Math.round(2 * Math.PI * outerRingRadius / (ballRadius * 3)));

  ctx.clearRect(0, 0, width, height);

  const totalEntries = props.entries.length;

  // === EMPTY STATE ===
  if (totalEntries === 0) {
    // Draw golden ring
    drawGoldenRing(ctx, cx, cy, outerRingRadius, ringThickness, ballRadius, ballCount, 0);
    // Empty wheel
    ctx.beginPath();
    ctx.arc(cx, cy, wheelRadius, 0, 2 * Math.PI);
    ctx.fillStyle = '#f1f5f9';
    ctx.fill();
    ctx.fillStyle = '#94a3b8';
    ctx.font = `500 ${wheelRadius / 8}px Outfit, sans-serif`;
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText('Nhập danh sách lựa chọn', cx, cy);
    return;
  }

  // === DRAW SEGMENTS ===
  const totalWeight = props.entries.reduce((sum, e) => sum + (props.settings.enableWeights ? e.weight : 1), 0);
  let currentAngle = 0;
  const sliceAngles: { start: number; end: number; mid: number; color: string }[] = [];

  props.entries.forEach((entry, index) => {
    const weight = props.settings.enableWeights ? entry.weight : 1;
    const sliceAngle = (weight / totalWeight) * 2 * Math.PI;
    const startAngle = currentAngle;
    const endAngle = currentAngle + sliceAngle;
    const midAngle = startAngle + sliceAngle / 2;
    const baseColor = props.getEntryColor(index, entry.color);
    currentAngle = endAngle;

    sliceAngles.push({ start: startAngle, end: endAngle, mid: midAngle, color: baseColor });

    // Draw segment with radial gradient for 3D look
    const grad = ctx.createRadialGradient(
      cx + Math.cos(midAngle) * wheelRadius * 0.3,
      cy + Math.sin(midAngle) * wheelRadius * 0.3,
      0,
      cx, cy, wheelRadius
    );
    grad.addColorStop(0, lightenColor(baseColor, 0.18));
    grad.addColorStop(0.6, baseColor);
    grad.addColorStop(1, baseColor);

    ctx.beginPath();
    ctx.moveTo(cx, cy);
    ctx.arc(cx, cy, wheelRadius, startAngle, endAngle);
    ctx.closePath();
    ctx.fillStyle = grad;
    ctx.fill();
  });

  // === SEGMENT DIVIDER LINES ===
  sliceAngles.forEach(({ start }) => {
    ctx.beginPath();
    ctx.moveTo(cx, cy);
    ctx.lineTo(
      cx + Math.cos(start) * wheelRadius,
      cy + Math.sin(start) * wheelRadius
    );
    ctx.strokeStyle = 'rgba(255,255,255,0.75)';
    ctx.lineWidth = 2;
    ctx.stroke();
  });

  // === OUTER WHEEL BORDER ===
  ctx.beginPath();
  ctx.arc(cx, cy, wheelRadius, 0, 2 * Math.PI);
  ctx.strokeStyle = 'rgba(255,255,255,0.85)';
  ctx.lineWidth = 3;
  ctx.stroke();

  // === TEXT LABELS ===
  if (props.settings.showTextOnWheel) {
    sliceAngles.forEach(({ start, end, color }, i) => {
      const sliceAngle = end - start;
      const midAngle = start + sliceAngle / 2;
      const isBgDark = isColorDark(color);

      ctx.save();
      ctx.translate(cx, cy);
      ctx.rotate(midAngle);

      ctx.fillStyle = isBgDark ? '#ffffff' : '#1e293b';

      // Text shadow for readability
      ctx.shadowColor = isBgDark ? 'rgba(0,0,0,0.5)' : 'rgba(255,255,255,0.5)';
      ctx.shadowBlur = 3;

      const sizeMultiplier = props.settings.fontSize / 18;
      const calculatedFontSize = Math.max(9, Math.min(20, (wheelRadius / 14) * sizeMultiplier));
      ctx.font = `700 ${calculatedFontSize}px Outfit, sans-serif`;
      ctx.textAlign = 'right';
      ctx.textBaseline = 'middle';

      let text = props.entries[i]?.label || '';
      const maxTextLength = Math.max(5, Math.floor(wheelRadius / (calculatedFontSize * 0.75)));
      if (text.length > maxTextLength) {
        text = text.substring(0, maxTextLength - 2) + '..';
      }

      const textOffset = wheelRadius - 14;
      ctx.fillText(text, textOffset, 0);
      ctx.restore();
    });
  }

  // === GOLDEN OUTER RING ===
  drawGoldenRing(ctx, cx, cy, outerRingRadius, ringThickness, ballRadius, ballCount, currentAngle);
};

const drawGoldenRing = (
  ctx: CanvasRenderingContext2D,
  cx: number, cy: number,
  outerR: number,
  thickness: number,
  ballR: number,
  ballCount: number,
  _angle: number
) => {
  const innerR = outerR - thickness;

  // === Outer shadow ===
  ctx.save();
  ctx.shadowColor = 'rgba(0,0,0,0.35)';
  ctx.shadowBlur = 12;
  ctx.shadowOffsetY = 4;

  // Gold ring base
  const ringGrad = ctx.createRadialGradient(cx - outerR * 0.2, cy - outerR * 0.2, innerR * 0.5, cx, cy, outerR);
  ringGrad.addColorStop(0,   '#fffde0');
  ringGrad.addColorStop(0.2, '#ffd84d');
  ringGrad.addColorStop(0.5, '#e8900a');
  ringGrad.addColorStop(0.75,'#f5c842');
  ringGrad.addColorStop(1,   '#c47a00');

  ctx.beginPath();
  ctx.arc(cx, cy, outerR, 0, 2 * Math.PI);
  ctx.arc(cx, cy, innerR, 0, 2 * Math.PI, true);
  ctx.fillStyle = ringGrad;
  ctx.fill('evenodd');
  ctx.restore();

  // === Ring highlight stripe ===
  ctx.beginPath();
  ctx.arc(cx, cy, outerR - 1, 0, 2 * Math.PI);
  ctx.strokeStyle = 'rgba(255,255,220,0.45)';
  ctx.lineWidth = 2;
  ctx.stroke();

  ctx.beginPath();
  ctx.arc(cx, cy, innerR + 1, 0, 2 * Math.PI);
  ctx.strokeStyle = 'rgba(120,60,0,0.35)';
  ctx.lineWidth = 1.5;
  ctx.stroke();

  // === Pearl balls ===
  const ballOrbitR = outerR - thickness / 2;
  for (let i = 0; i < ballCount; i++) {
    const angle = (i / ballCount) * 2 * Math.PI - Math.PI / 2;
    const bx = cx + Math.cos(angle) * ballOrbitR;
    const by = cy + Math.sin(angle) * ballOrbitR;

    // Ball shadow
    ctx.save();
    ctx.shadowColor = 'rgba(0,0,0,0.3)';
    ctx.shadowBlur = 4;
    ctx.shadowOffsetY = 2;

    // Pearl gradient: white highlight top-left
    const ballGrad = ctx.createRadialGradient(
      bx - ballR * 0.35, by - ballR * 0.35, 0,
      bx, by, ballR
    );
    ballGrad.addColorStop(0,   '#ffffff');
    ballGrad.addColorStop(0.4, '#f8f0d0');
    ballGrad.addColorStop(0.75,'#deba6a');
    ballGrad.addColorStop(1,   '#b07d10');

    ctx.beginPath();
    ctx.arc(bx, by, ballR, 0, 2 * Math.PI);
    ctx.fillStyle = ballGrad;
    ctx.fill();
    ctx.restore();

    // Ball rim
    ctx.beginPath();
    ctx.arc(bx, by, ballR, 0, 2 * Math.PI);
    ctx.strokeStyle = 'rgba(180,120,20,0.6)';
    ctx.lineWidth = 0.8;
    ctx.stroke();
  }

  // === Center hub (drawn over segments) ===
  const hubR = outerR * 0.16;

  // Hub shadow
  ctx.save();
  ctx.shadowColor = 'rgba(0,0,0,0.45)';
  ctx.shadowBlur = 10;
  ctx.shadowOffsetY = 3;

  const hubGrad = ctx.createRadialGradient(
    cx - hubR * 0.35, cy - hubR * 0.35, 0,
    cx, cy, hubR
  );
  hubGrad.addColorStop(0,   '#fffde0');
  hubGrad.addColorStop(0.3, '#ffd84d');
  hubGrad.addColorStop(0.65,'#f5a623');
  hubGrad.addColorStop(1,   '#c47e0a');

  ctx.beginPath();
  ctx.arc(cx, cy, hubR, 0, 2 * Math.PI);
  ctx.fillStyle = hubGrad;
  ctx.fill();
  ctx.restore();

  // Hub rim
  ctx.beginPath();
  ctx.arc(cx, cy, hubR, 0, 2 * Math.PI);
  ctx.strokeStyle = 'rgba(120,60,0,0.5)';
  ctx.lineWidth = 2;
  ctx.stroke();

  // Hub shine dot
  ctx.beginPath();
  ctx.arc(cx - hubR * 0.3, cy - hubR * 0.3, hubR * 0.22, 0, 2 * Math.PI);
  ctx.fillStyle = 'rgba(255,255,255,0.65)';
  ctx.fill();
};

// Lifecycle
onMounted(() => {
  drawWheel();
  if (canvasRef.value) {
    resizeObserver = new ResizeObserver(() => { drawWheel(); });
    resizeObserver.observe(canvasRef.value);
  }
});

onUnmounted(() => {
  if (resizeObserver) resizeObserver.disconnect();
});

watch(() => [props.entries, props.settings], () => { drawWheel(); }, { deep: true });

const observer = new MutationObserver(() => { drawWheel(); });
onMounted(() => { observer.observe(document.body, { attributes: true, attributeFilter: ['class'] }); });
onUnmounted(() => { observer.disconnect(); });
</script>

<style scoped>
canvas {
  image-rendering: -webkit-optimize-contrast;
  image-rendering: crisp-edges;
}
</style>
