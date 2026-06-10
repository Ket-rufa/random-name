import { ref } from 'vue';
import type { WheelEntry } from '../types/wheel';
import { useAudio } from './useAudio';

export function useWheelAnimation() {
  const rotationAngle = ref<number>(0);
  const isSpinning = ref<boolean>(false);
  const { playTick } = useAudio();

  let animationFrameId: number | null = null;

  const spin = (
    entries: WheelEntry[],
    useWeights: boolean,
    winner: WheelEntry,
    durationSeconds: number,
    volume: number,
    enableTickSound: boolean,
    onFinish: () => void
  ) => {
    if (isSpinning.value) return;
    isSpinning.value = true;

    // Calculate segment divisions based on weights
    const totalWeight = entries.reduce((sum, e) => sum + (useWeights ? e.weight : 1), 0);
    let accumulatedAngle = 0;
    const segments = entries.map(entry => {
      const weight = useWeights ? entry.weight : 1;
      const angleSize = (weight / totalWeight) * 360;
      const start = accumulatedAngle;
      const end = accumulatedAngle + angleSize;
      accumulatedAngle = end;
      return { entry, start, end };
    });

    const winnerSegment = segments.find(s => s.entry.id === winner.id);
    if (!winnerSegment) {
      isSpinning.value = false;
      return;
    }

    // Pick a point safely inside the winner segment (leave small padding near boundaries)
    const segmentWidth = winnerSegment.end - winnerSegment.start;
    const padding = Math.min(2.5, segmentWidth * 0.15); // max 2.5 degrees padding
    const minTarget = winnerSegment.start + padding;
    const maxTarget = winnerSegment.end - padding;
    const targetAngleOnWheel = minTarget + Math.random() * (maxTarget - minTarget);

    // Arrow is at the top (270 degrees)
    // To align targetAngleOnWheel directly with 270, the canvas rotation must be:
    const targetTheta = (270 - targetAngleOnWheel + 360) % 360;

    const startTheta = rotationAngle.value;
    const extraSpins = 6 + Math.floor(Math.random() * 4); // 6 to 9 full spins
    const diff = (targetTheta - (startTheta % 360) + 360) % 360;
    const totalRotation = startTheta + (extraSpins * 360) + diff;

    const durationMs = durationSeconds * 1000;
    const startTime = performance.now();
    let lastSegmentIdx = -1;

    const animate = (currentTime: number) => {
      const elapsed = currentTime - startTime;
      const t = Math.min(elapsed / durationMs, 1);

      // Ease out quartic: f(t) = 1 - (1 - t)^4 (very smooth slowdown)
      const ease = 1 - Math.pow(1 - t, 4);
      rotationAngle.value = startTheta + ease * (totalRotation - startTheta);

      // Play tick sound when boundaries cross the top pointer
      if (enableTickSound && entries.length > 1) {
        const currentArrowAngle = (270 - (rotationAngle.value % 360) + 360) % 360;
        const currentSegmentIdx = segments.findIndex(
          s => currentArrowAngle >= s.start && currentArrowAngle < s.end
        );

        if (currentSegmentIdx !== lastSegmentIdx) {
          if (lastSegmentIdx !== -1) {
            playTick(volume);
          }
          lastSegmentIdx = currentSegmentIdx;
        }
      }

      if (t < 1) {
        animationFrameId = requestAnimationFrame(animate);
      } else {
        isSpinning.value = false;
        rotationAngle.value = totalRotation % 360; // Normalize angle
        onFinish();
      }
    };

    animationFrameId = requestAnimationFrame(animate);
  };

  const stop = () => {
    if (animationFrameId !== null) {
      cancelAnimationFrame(animationFrameId);
      animationFrameId = null;
    }
    isSpinning.value = false;
  };

  return {
    rotationAngle,
    isSpinning,
    spin,
    stop
  };
}
