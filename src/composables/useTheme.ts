import { ref, readonly } from 'vue';

const THEME_KEY = 'vqmm_theme';

// Singleton state - shared across all composable instances
const isDark = ref(false);

function applyTheme(dark: boolean) {
  if (dark) {
    document.documentElement.classList.add('dark');
  } else {
    document.documentElement.classList.remove('dark');
  }
}

// Initialize once on first import
function initTheme() {
  const saved = localStorage.getItem(THEME_KEY);
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
  const shouldBeDark = saved === 'dark' || (saved === null && prefersDark);
  isDark.value = shouldBeDark;
  applyTheme(shouldBeDark);
}

// Watch for system preference changes
let systemListener: MediaQueryList | null = null;
function attachSystemListener() {
  if (systemListener) return;
  systemListener = window.matchMedia('(prefers-color-scheme: dark)');
  systemListener.addEventListener('change', (e) => {
    // Only follow system if user hasn't explicitly set a preference
    if (!localStorage.getItem(THEME_KEY)) {
      isDark.value = e.matches;
      applyTheme(e.matches);
    }
  });
}

export function useTheme() {
  const toggle = () => {
    isDark.value = !isDark.value;
    applyTheme(isDark.value);
    localStorage.setItem(THEME_KEY, isDark.value ? 'dark' : 'light');
  };

  const setDark = (dark: boolean) => {
    isDark.value = dark;
    applyTheme(dark);
    localStorage.setItem(THEME_KEY, dark ? 'dark' : 'light');
  };

  return {
    isDark: readonly(isDark),
    toggle,
    setDark,
    initTheme,
    attachSystemListener,
  };
}
