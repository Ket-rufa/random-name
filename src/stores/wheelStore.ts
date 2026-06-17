import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { WheelEntry, WheelSettings, SpinHistory, Wheel } from '../types/wheel';
import { wheelApi } from '../api/wheelApi';

const THEME_PALETTES: Record<string, string[]> = {
  default: ['#f5f5f0', '#00b4d8', '#ffd600', '#9b5de5', '#f72585', '#ff6d00', '#06d6a0', '#ff006e'],
  neon: ['#ff007f', '#39ff14', '#00ffff', '#ff00ff', '#ffff00', '#9d00ff'],
  candy: ['#ff9eb5', '#ffbeb2', '#e2f0cb', '#b5ead7', '#c7ceea', '#ffdac1'],
  pastel: ['#f4a261', '#e76f51', '#2a9d8f', '#e9c46a', '#457b9d', '#a8dadc'],
  monochrome: ['#1e293b', '#334155', '#475569', '#64748b', '#94a3b8', '#cbd5e1'],
  ocean: ['#0284c7', '#0369a1', '#0f766e', '#0d9488', '#06b6d4', '#0891b2']
};

const DEFAULT_SETTINGS: WheelSettings = {
  theme: 'default',
  customColors: [],
  spinDuration: 5,
  volume: 50,
  enableSound: true,
  enableTickSound: true,
  enableVictorySound: true,
  enableConfetti: true,
  allowDuplicates: true,
  autoRemoveWinner: false,
  enableWeights: false,
  showTextOnWheel: true,
  fontSize: 18
};

const DEFAULT_ENTRIES: string[] = [
  'Nguyễn Văn A',
  'Trần Thị B',
  'Lê Văn C',
  'Phạm Thị D',
  'Hoàng Văn E',
  'Huỳnh Thị F'
];

export const useWheelStore = defineStore('wheel', () => {
  // State
  const id = ref<string | null>(null);
  const title = ref<string>('Vòng Quay May Mắn');
  const shareCode = ref<string | null>(null);
  const editToken = ref<string | null>(null);
  const permission = ref<'view' | 'spin' | 'edit'>('edit');
  
  const entries = ref<WheelEntry[]>([]);
  const settings = ref<WheelSettings>({ ...DEFAULT_SETTINGS });
  const spinHistory = ref<SpinHistory[]>([]);
  
  // UI states
  const selectedResult = ref<WheelEntry | null>(null);
  const isSpinning = ref<boolean>(false);
  const isResultModalOpen = ref<boolean>(false);
  const isLoading = ref<boolean>(false);
  const isError = ref<boolean>(false);
  const errorMessage = ref<string>('');
  
  // Toast notifications
  const toastMessage = ref<string | null>(null);
  const toastType = ref<'success' | 'error' | 'info'>('info');
  
  let toastTimeout: any = null;
  const showToast = (message: string, type: 'success' | 'error' | 'info' = 'info', duration = 3000) => {
    if (toastTimeout) clearTimeout(toastTimeout);
    toastMessage.value = message;
    toastType.value = type;
    toastTimeout = setTimeout(() => {
      toastMessage.value = null;
    }, duration);
  };

  // Last removed entry for Undo feature
  const lastRemovedEntry = ref<{ entry: WheelEntry; index: number } | null>(null);

  // Getters
  const currentPalette = computed(() => {
    if (settings.value.theme === 'custom' && settings.value.customColors.length > 0) {
      return settings.value.customColors;
    }
    return THEME_PALETTES[settings.value.theme] || THEME_PALETTES.default;
  });

  const canEdit = computed(() => permission.value === 'edit');
  const canSpin = computed(() => permission.value === 'edit' || permission.value === 'spin');

  // Helper: map index to palette color
  const getEntryColor = (index: number, entryColor?: string) => {
    if (entryColor) return entryColor;
    const palette = currentPalette.value;
    return palette[index % palette.length];
  };

  // Actions
  const addEntry = (label: string, weight = 1, color?: string) => {
    if (!label.trim()) return;
    const newEntry: WheelEntry = {
      id: crypto.randomUUID(),
      label: label.trim(),
      weight: Math.max(0.1, weight),
      color,
      position: entries.value.length
    };
    entries.value.push(newEntry);
    saveToLocalStorage();
  };

  const removeEntry = (entryId: string) => {
    const idx = entries.value.findIndex(e => e.id === entryId);
    if (idx !== -1) {
      const [removed] = entries.value.splice(idx, 1);
      lastRemovedEntry.value = { entry: removed, index: idx };
      // Update positions
      entries.value.forEach((entry, i) => {
        entry.position = i;
      });
      saveToLocalStorage();
    }
  };

  const restoreLastRemoved = () => {
    if (lastRemovedEntry.value) {
      const { entry, index } = lastRemovedEntry.value;
      entries.value.splice(index, 0, entry);
      entries.value.forEach((e, i) => {
        e.position = i;
      });
      lastRemovedEntry.value = null;
      saveToLocalStorage();
      return true;
    }
    return false;
  };

  const updateEntry = (entryId: string, data: Partial<Omit<WheelEntry, 'id'>>) => {
    const entry = entries.value.find(e => e.id === entryId);
    if (entry) {
      Object.assign(entry, data);
      saveToLocalStorage();
    }
  };

  const shuffleEntries = () => {
    // Cryptographically secure shuffle
    const array = [...entries.value];
    for (let i = array.length - 1; i > 0; i--) {
      const randArr = new Uint32Array(1);
      window.crypto.getRandomValues(randArr);
      const j = randArr[0] % (i + 1);
      [array[i], array[j]] = [array[j], array[i]];
    }
    array.forEach((entry, idx) => {
      entry.position = idx;
    });
    entries.value = array;
    saveToLocalStorage();
  };

  const sortEntries = () => {
    entries.value.sort((a, b) => a.label.localeCompare(b.label, 'vi'));
    entries.value.forEach((entry, idx) => {
      entry.position = idx;
    });
    saveToLocalStorage();
  };

  const clearEntries = () => {
    entries.value = [];
    saveToLocalStorage();
  };

  const setEntriesFromText = (text: string) => {
    const lines = text.split('\n').map(line => line.trim()).filter(line => line.length > 0);
    entries.value = lines.map((label, index) => ({
      id: crypto.randomUUID(),
      label,
      weight: 1,
      position: index
    }));
    saveToLocalStorage();
  };

  const getEntriesText = () => {
    return entries.value.map(e => e.label).join('\n');
  };

  const clearHistory = () => {
    spinHistory.value = [];
    saveToLocalStorage();
  };

  const addHistory = (label: string, entryId?: string) => {
    spinHistory.value.unshift({
      id: crypto.randomUUID(),
      resultLabel: label,
      spunAt: new Date().toISOString(),
      entryId
    });
    saveToLocalStorage();
  };

  const removeHistoryItem = (historyId: string) => {
    spinHistory.value = spinHistory.value.filter(h => h.id !== historyId);
    saveToLocalStorage();
  };

  // LocalStorage logic
  const saveToLocalStorage = () => {
    if (shareCode.value && permission.value !== 'edit') {
      // Don't save over owner settings if we are in view/spin mode
      return;
    }
    try {
      const dataToSave = {
        title: title.value,
        entries: entries.value,
        settings: settings.value,
        spinHistory: spinHistory.value,
        editToken: editToken.value,
        id: id.value,
        shareCode: shareCode.value
      };
      localStorage.setItem('random_name_wheel_data', JSON.stringify(dataToSave));
    } catch (e) {
      console.error('Failed to save state to localStorage:', e);
    }
  };

  const loadFromLocalStorage = () => {
    try {
      const raw = localStorage.getItem('random_name_wheel_data');
      if (raw) {
        const parsed = JSON.parse(raw);
        if (parsed) {
          title.value = parsed.title || 'Vòng Quay May Mắn';
          entries.value = Array.isArray(parsed.entries) ? parsed.entries : [];
          settings.value = { ...DEFAULT_SETTINGS, ...parsed.settings };
          spinHistory.value = Array.isArray(parsed.spinHistory) ? parsed.spinHistory : [];
          editToken.value = parsed.editToken || null;
          id.value = parsed.id || null;
          shareCode.value = parsed.shareCode || null;
          permission.value = 'edit'; // Local storage owner is always editor
          return true;
        }
      }
    } catch (e) {
      console.warn('Failed to parse localStorage data, loading defaults.');
    }
    loadDefaultData();
    return false;
  };

  const loadDefaultData = () => {
    title.value = 'Vòng Quay May Mắn';
    entries.value = DEFAULT_ENTRIES.map((label, index) => ({
      id: crypto.randomUUID(),
      label,
      weight: 1,
      position: index
    }));
    settings.value = { ...DEFAULT_SETTINGS };
    spinHistory.value = [];
    editToken.value = null;
    id.value = null;
    shareCode.value = null;
    permission.value = 'edit';
  };

  const resetWheel = () => {
    loadDefaultData();
    localStorage.removeItem('random_name_wheel_data');
  };

  // Set store data from Server response
  const setWheelDataFromServer = (wheelData: Wheel, remoteEditToken?: string) => {
    id.value = wheelData.id;
    title.value = wheelData.title;
    shareCode.value = wheelData.shareCode;
    settings.value = { ...DEFAULT_SETTINGS, ...wheelData.settings };
    entries.value = Array.isArray(wheelData.entries)
      ? wheelData.entries.sort((a, b) => a.position - b.position)
      : [];
    permission.value = wheelData.permission;
    
    if (remoteEditToken) {
      editToken.value = remoteEditToken;
    }
    
    // Save to local storage if we are the owner
    if (permission.value === 'edit') {
      saveToLocalStorage();
    }
  };

  const saveWheelToServer = async () => {
    isLoading.value = true;
    isError.value = false;
    errorMessage.value = '';
    try {
      if (id.value && editToken.value && permission.value === 'edit') {
        const payload = {
          title: title.value,
          entries: entries.value.map(e => ({
            label: e.label,
            weight: e.weight,
            color: e.color,
            position: e.position
          })),
          settings: settings.value
        };
        const res = await wheelApi.updateWheel(id.value, payload);
        if (res.data && res.data.success) {
          showToast('Đã lưu thay đổi lên server thành công!', 'success');
          saveToLocalStorage();
          return true;
        }
      } else {
        const payload = {
          title: title.value,
          entries: entries.value.map(e => ({
            label: e.label,
            weight: e.weight,
            color: e.color
          })),
          settings: settings.value
        };
        const res = await wheelApi.createWheel(payload);
        if (res.data && res.data.success) {
          const { wheel, editToken: localEditToken } = res.data.data;
          setWheelDataFromServer(wheel, localEditToken);
          showToast('Đã tạo và lưu vòng quay lên server thành công!', 'success');
          return true;
        }
      }
    } catch (err: any) {
      isError.value = true;
      errorMessage.value = err.response?.data?.message || 'Không thể lưu vòng quay lên server.';
      showToast(errorMessage.value, 'error');
    } finally {
      isLoading.value = false;
    }
    return false;
  };

  const deleteWheelOnServer = async () => {
    if (!id.value || !editToken.value) return false;
    isLoading.value = true;
    try {
      const res = await wheelApi.deleteWheel(id.value);
      if (res.data && res.data.success) {
        showToast('Đã xóa vòng quay khỏi server!', 'success');
        resetWheel();
        return true;
      }
    } catch (err: any) {
      showToast(err.response?.data?.message || 'Không thể xóa vòng quay.', 'error');
    } finally {
      isLoading.value = false;
    }
    return false;
  };

  return {
    // State
    id,
    title,
    shareCode,
    editToken,
    permission,
    entries,
    settings,
    spinHistory,
    selectedResult,
    isSpinning,
    isResultModalOpen,
    isLoading,
    isError,
    errorMessage,
    lastRemovedEntry,
    toastMessage,
    toastType,

    // Getters
    currentPalette,
    canEdit,
    canSpin,
    getEntryColor,

    // Actions
    addEntry,
    removeEntry,
    restoreLastRemoved,
    updateEntry,
    shuffleEntries,
    sortEntries,
    clearEntries,
    setEntriesFromText,
    getEntriesText,
    clearHistory,
    addHistory,
    removeHistoryItem,
    saveToLocalStorage,
    loadFromLocalStorage,
    resetWheel,
    setWheelDataFromServer,
    showToast,
    saveWheelToServer,
    deleteWheelOnServer
  };
});
