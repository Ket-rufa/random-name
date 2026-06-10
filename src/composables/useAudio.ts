import { ref } from 'vue';

const audioCtx = ref<AudioContext | null>(null);

export function useAudio() {
  const initAudio = () => {
    if (!audioCtx.value) {
      audioCtx.value = new (window.AudioContext || (window as any).webkitAudioContext)();
    }
    if (audioCtx.value.state === 'suspended') {
      audioCtx.value.resume();
    }
  };

  const playTick = (volumePercent: number) => {
    try {
      initAudio();
      if (!audioCtx.value) return;

      const ctx = audioCtx.value;
      const vol = volumePercent / 100;
      if (vol <= 0) return;

      const osc = ctx.createOscillator();
      const gainNode = ctx.createGain();

      osc.connect(gainNode);
      gainNode.connect(ctx.destination);

      // Triangle wave has a nice mechanical "click" timbre
      osc.type = 'triangle';
      
      const now = ctx.currentTime;
      
      // Fast drop in frequency creates a nice tick sound
      osc.frequency.setValueAtTime(600, now);
      osc.frequency.exponentialRampToValueAtTime(100, now + 0.04);

      // Exponential decay of volume
      gainNode.gain.setValueAtTime(vol * 0.7, now);
      gainNode.gain.exponentialRampToValueAtTime(0.001, now + 0.045);

      osc.start(now);
      osc.stop(now + 0.05);
    } catch (e) {
      console.warn('Failed to play tick sound:', e);
    }
  };

  const playVictory = (volumePercent: number) => {
    try {
      initAudio();
      if (!audioCtx.value) return;

      const ctx = audioCtx.value;
      const vol = volumePercent / 100;
      if (vol <= 0) return;

      const now = ctx.currentTime;
      // Arpeggio of notes: C4 -> E4 -> G4 -> C5
      const notes = [261.63, 329.63, 392.00, 523.25];
      
      notes.forEach((freq, index) => {
        const osc = ctx.createOscillator();
        const gainNode = ctx.createGain();

        osc.connect(gainNode);
        gainNode.connect(ctx.destination);

        osc.type = 'sine';
        osc.frequency.setValueAtTime(freq, now + index * 0.12);

        // Soft onset and decay
        gainNode.gain.setValueAtTime(0, now + index * 0.12);
        gainNode.gain.linearRampToValueAtTime(vol * 0.4, now + index * 0.12 + 0.02);
        gainNode.gain.exponentialRampToValueAtTime(0.001, now + index * 0.12 + 0.45);

        osc.start(now + index * 0.12);
        osc.stop(now + index * 0.12 + 0.5);
      });
    } catch (e) {
      console.warn('Failed to play victory melody:', e);
    }
  };

  return {
    initAudio,
    playTick,
    playVictory,
  };
}
