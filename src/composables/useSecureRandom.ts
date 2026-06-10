import type { WheelEntry } from '../types/wheel';

/**
 * Generates a cryptographically secure random float in the range [0, 1)
 */
export function secureRandomFloat(): number {
  const array = new Uint32Array(1);
  window.crypto.getRandomValues(array);
  return array[0] / 4294967296; // 2^32
}

/**
 * Generates a cryptographically secure random integer in the range [0, max - 1]
 */
export function secureRandomInt(max: number): number {
  if (max <= 0) return 0;
  
  const array = new Uint32Array(1);
  window.crypto.getRandomValues(array);
  
  // Guard against modulo bias
  const maxSafe = Math.floor(4294967296 / max) * max;
  let val = array[0];
  while (val >= maxSafe) {
    window.crypto.getRandomValues(array);
    val = array[0];
  }
  
  return val % max;
}

/**
 * Selects a WheelEntry securely. Supports both uniform distribution and weighted options.
 */
export function pickSecureRandom(entries: WheelEntry[], useWeights: boolean): WheelEntry {
  if (entries.length === 0) {
    throw new Error('Selection list is empty');
  }

  if (!useWeights) {
    const idx = secureRandomInt(entries.length);
    return entries[idx];
  }

  const totalWeight = entries.reduce((sum, entry) => sum + (entry.weight || 1), 0);
  if (totalWeight <= 0) {
    const idx = secureRandomInt(entries.length);
    return entries[idx];
  }

  const randVal = secureRandomFloat() * totalWeight;
  let cumulativeWeight = 0;
  for (const entry of entries) {
    cumulativeWeight += entry.weight || 1;
    if (randVal <= cumulativeWeight) {
      return entry;
    }
  }

  return entries[entries.length - 1];
}
