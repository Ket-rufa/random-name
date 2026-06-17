<template>
  <BaseModal :is-open="isOpen" @close="close" size="lg">
    <template #title>
      <div class="flex items-center gap-2">
        <SettingsIcon class="h-5 w-5 text-primary-500" />
        <span>Tùy chỉnh vòng quay</span>
      </div>
    </template>

    <!-- Settings Content with Left Navigation Sidebar -->
    <div class="flex flex-col md:flex-row gap-6 min-h-[350px]">
      <!-- Navigation Menu -->
      <div class="flex md:flex-col gap-1 border-b md:border-b-0 md:border-r border-slate-100 dark:border-slate-800/60 pb-3 md:pb-0 md:pr-4 md:w-44 flex-shrink-0 overflow-x-auto">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeSubTab = tab.id"
          class="flex items-center gap-2 px-3 py-2 rounded-xl text-sm font-semibold transition-all duration-200 text-left whitespace-nowrap"
          :class="[
            activeSubTab === tab.id
              ? 'bg-primary-50 dark:bg-primary-950/30 text-primary-600 dark:text-primary-400'
              : 'text-slate-600 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800/40'
          ]"
        >
          <component :is="tab.icon" class="h-4 w-4" />
          <span>{{ tab.label }}</span>
        </button>
      </div>

      <!-- Settings Panels -->
      <div class="flex-1 min-w-0">
        <!-- 1. Appearance Settings -->
        <div v-if="activeSubTab === 'appearance'" class="space-y-5">
          <h4 class="text-sm font-bold text-slate-800 dark:text-slate-200">Giao diện vòng quay</h4>
          
          <!-- Themes -->
          <div class="space-y-2">
            <span class="text-xs font-bold text-slate-700 dark:text-slate-300">Chủ đề màu sắc</span>
            <div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
              <button
                v-for="theme in themes"
                :key="theme.id"
                @click="updateSetting('theme', theme.id)"
                class="flex flex-col gap-2 p-2.5 rounded-2xl border text-left transition-all duration-200"
                :class="[
                  store.settings.theme === theme.id
                    ? 'border-primary-500 bg-primary-50/20 dark:bg-primary-950/10'
                    : 'border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800/20'
                ]"
              >
                <span class="text-xs font-bold text-slate-800 dark:text-slate-200 capitalize">{{ theme.name }}</span>
                <!-- Color dots -->
                <div class="flex gap-1 overflow-hidden">
                  <div
                    v-for="c in theme.colors.slice(0, 5)"
                    :key="c"
                    :style="{ backgroundColor: c }"
                    class="w-3.5 h-3.5 rounded-full border border-white/20"
                  ></div>
                </div>
              </button>
            </div>
          </div>

          <!-- Font Size Slider -->
          <div class="space-y-2">
            <div class="flex justify-between items-center text-xs font-bold">
              <label for="font-size-range" class="text-slate-700 dark:text-slate-300">Kích thước chữ vòng quay</label>
              <span class="text-slate-800 dark:text-slate-200">{{ store.settings.fontSize }}px</span>
            </div>
            <input
              id="font-size-range"
              type="range"
              min="10"
              max="28"
              v-model.number="store.settings.fontSize"
              @input="(e) => updateSetting('fontSize', parseInt((e.target as HTMLInputElement).value))"
              class="w-full h-1.5 bg-slate-200 dark:bg-slate-800 rounded-lg appearance-none cursor-pointer accent-primary-600"
            />
          </div>

          <!-- Text on wheel switch -->
          <div class="flex items-center justify-between py-2">
            <div class="flex flex-col gap-0.5">
              <label for="show-text-switch" class="text-sm font-bold text-slate-800 dark:text-slate-200">Hiển thị chữ trên vòng quay</label>
              <span class="text-xs text-slate-500 dark:text-slate-400">Ẩn chữ nếu muốn hiển thị vòng quay trơn</span>
            </div>
            <input
              id="show-text-switch"
              type="checkbox"
              :checked="store.settings.showTextOnWheel"
              @change="(e) => updateSetting('showTextOnWheel', (e.target as HTMLInputElement).checked)"
              class="w-9 h-5 bg-slate-200 rounded-full appearance-none checked:bg-primary-600 relative transition-all duration-200 cursor-pointer before:content-[''] before:absolute before:w-4 before:h-4 before:rounded-full before:bg-white before:top-0.5 before:left-0.5 before:transition-all before:duration-200 checked:before:translate-x-4"
            />
          </div>
        </div>

        <!-- 2. Sound and Animation Settings -->
        <div v-if="activeSubTab === 'audio'" class="space-y-5">
          <h4 class="text-sm font-bold text-slate-800 dark:text-slate-200">Âm thanh & Hiệu ứng</h4>

          <!-- Enable sound toggle -->
          <div class="flex items-center justify-between py-2">
            <div class="flex flex-col gap-0.5">
              <label for="sound-master-switch" class="text-sm font-bold text-slate-800 dark:text-slate-200">Bật âm thanh</label>
              <span class="text-xs text-slate-400 dark:text-slate-500">Bật/tắt tất cả âm thanh trong ứng dụng</span>
            </div>
            <input
              id="sound-master-switch"
              type="checkbox"
              :checked="store.settings.enableSound"
              @change="(e) => updateSetting('enableSound', (e.target as HTMLInputElement).checked)"
              class="w-9 h-5 bg-slate-200 rounded-full appearance-none checked:bg-primary-600 relative transition-all duration-200 cursor-pointer before:content-[''] before:absolute before:w-4 before:h-4 before:rounded-full before:bg-white before:top-0.5 before:left-0.5 before:transition-all before:duration-200 checked:before:translate-x-4"
            />
          </div>

          <!-- Volume and sub-sounds -->
          <div v-if="store.settings.enableSound" class="space-y-4 pl-4 border-l-2 border-slate-100 dark:border-slate-800">
            <!-- Volume Slider -->
            <div class="space-y-2">
              <div class="flex justify-between items-center text-xs font-semibold">
                <label for="volume-range" class="text-slate-400 dark:text-slate-500">Âm lượng</label>
                <span class="text-slate-850 dark:text-slate-200">{{ store.settings.volume }}%</span>
              </div>
              <input
                id="volume-range"
                type="range"
                min="0"
                max="100"
                v-model.number="store.settings.volume"
                @input="(e) => updateSetting('volume', parseInt((e.target as HTMLInputElement).value))"
                class="w-full h-1.5 bg-slate-200 dark:bg-slate-800 rounded-lg appearance-none cursor-pointer accent-primary-600"
              />
            </div>

            <!-- Tick sound toggle -->
            <div class="flex items-center justify-between">
              <label for="sound-tick-switch" class="text-xs font-bold text-slate-700 dark:text-slate-300">Âm thanh click khi quay</label>
              <input
                id="sound-tick-switch"
                type="checkbox"
                :checked="store.settings.enableTickSound"
                @change="(e) => updateSetting('enableTickSound', (e.target as HTMLInputElement).checked)"
                class="w-8 h-4.5 bg-slate-200 rounded-full appearance-none checked:bg-primary-600 relative transition-all duration-200 cursor-pointer before:content-[''] before:absolute before:w-3.5 before:h-3.5 before:rounded-full before:bg-white before:top-0.5 before:left-0.5 before:transition-all before:duration-200 checked:before:translate-x-3.5"
              />
            </div>

            <!-- Victory chime toggle -->
            <div class="flex items-center justify-between">
              <label for="sound-victory-switch" class="text-xs font-bold text-slate-700 dark:text-slate-300">Âm thanh chúc mừng thắng giải</label>
              <input
                id="sound-victory-switch"
                type="checkbox"
                :checked="store.settings.enableVictorySound"
                @change="(e) => updateSetting('enableVictorySound', (e.target as HTMLInputElement).checked)"
                class="w-8 h-4.5 bg-slate-200 rounded-full appearance-none checked:bg-primary-600 relative transition-all duration-200 cursor-pointer before:content-[''] before:absolute before:w-3.5 before:h-3.5 before:rounded-full before:bg-white before:top-0.5 before:left-0.5 before:transition-all before:duration-200 checked:before:translate-x-3.5"
              />
            </div>
          </div>

          <!-- Spin Duration -->
          <div class="space-y-2">
            <div class="flex justify-between items-center text-xs font-semibold">
              <label for="spin-duration-range" class="text-slate-400 dark:text-slate-500">Thời gian quay</label>
              <span class="text-slate-850 dark:text-slate-200">{{ store.settings.spinDuration }} giây</span>
            </div>
            <input
              id="spin-duration-range"
              type="range"
              min="2"
              max="15"
              v-model.number="store.settings.spinDuration"
              @input="(e) => updateSetting('spinDuration', parseInt((e.target as HTMLInputElement).value))"
              class="w-full h-1.5 bg-slate-200 dark:bg-slate-800 rounded-lg appearance-none cursor-pointer accent-primary-600"
            />
          </div>

          <!-- Confetti switch -->
          <div class="flex items-center justify-between py-2">
            <div class="flex flex-col gap-0.5">
              <label for="confetti-switch" class="text-sm font-bold text-slate-800 dark:text-slate-200">Hiệu ứng pháo hoa giấy</label>
              <span class="text-xs text-slate-400 dark:text-slate-500">Bắn pháo hoa giấy rực rỡ khi hiển thị kết quả</span>
            </div>
            <input
              id="confetti-switch"
              type="checkbox"
              :checked="store.settings.enableConfetti"
              @change="(e) => updateSetting('enableConfetti', (e.target as HTMLInputElement).checked)"
              class="w-9 h-5 bg-slate-200 rounded-full appearance-none checked:bg-primary-600 relative transition-all duration-200 cursor-pointer before:content-[''] before:absolute before:w-4 before:h-4 before:rounded-full before:bg-white before:top-0.5 before:left-0.5 before:transition-all before:duration-200 checked:before:translate-x-4"
            />
          </div>

          <!-- Flashing lights switch -->
          <div class="flex items-center justify-between py-2">
            <div class="flex flex-col gap-0.5">
              <label for="flashing-lights-switch" class="text-sm font-bold text-slate-800 dark:text-slate-200">Hiệu ứng đèn nhấp nháy</label>
              <span class="text-xs text-slate-400 dark:text-slate-500">Bật/tắt nhấp nháy đèn LED viền ngoài vòng quay</span>
            </div>
            <input
              id="flashing-lights-switch"
              type="checkbox"
              :checked="store.settings.enableFlashingLights"
              @change="(e) => updateSetting('enableFlashingLights', (e.target as HTMLInputElement).checked)"
              class="w-9 h-5 bg-slate-200 rounded-full appearance-none checked:bg-primary-600 relative transition-all duration-200 cursor-pointer before:content-[''] before:absolute before:w-4 before:h-4 before:rounded-full before:bg-white before:top-0.5 before:left-0.5 before:transition-all before:duration-200 checked:before:translate-x-4"
            />
          </div>
        </div>

        <!-- 3. Rule Settings -->
        <div v-if="activeSubTab === 'rules'" class="space-y-5">
          <h4 class="text-sm font-bold text-slate-800 dark:text-slate-200">Quy tắc vòng quay</h4>

          <!-- Auto remove winner -->
          <div class="flex items-center justify-between py-2">
            <div class="flex flex-col gap-0.5">
              <label for="auto-remove-switch" class="text-sm font-bold text-slate-800 dark:text-slate-200">Tự động xóa người trúng</label>
              <span class="text-xs text-slate-400 dark:text-slate-500">Xóa lựa chọn vừa trúng khỏi danh sách sau 1 giây</span>
            </div>
            <input
              id="auto-remove-switch"
              type="checkbox"
              :checked="store.settings.autoRemoveWinner"
              @change="(e) => updateSetting('autoRemoveWinner', (e.target as HTMLInputElement).checked)"
              class="w-9 h-5 bg-slate-200 rounded-full appearance-none checked:bg-primary-600 relative transition-all duration-200 cursor-pointer before:content-[''] before:absolute before:w-4 before:h-4 before:rounded-full before:bg-white before:top-0.5 before:left-0.5 before:transition-all before:duration-200 checked:before:translate-x-4"
            />
          </div>

          <!-- Weights toggler -->
          <div class="flex items-center justify-between py-2">
            <div class="flex flex-col gap-0.5">
              <label for="enable-weights-switch" class="text-sm font-bold text-slate-800 dark:text-slate-200">Bật chế độ trọng số (Hệ số)</label>
              <span class="text-xs text-slate-400 dark:text-slate-500">Cho phép thiết lập tỉ lệ trúng khác nhau cho các dòng</span>
            </div>
            <input
              id="enable-weights-switch"
              type="checkbox"
              :checked="store.settings.enableWeights"
              @change="(e) => updateSetting('enableWeights', (e.target as HTMLInputElement).checked)"
              class="w-9 h-5 bg-slate-200 rounded-full appearance-none checked:bg-primary-600 relative transition-all duration-200 cursor-pointer before:content-[''] before:absolute before:w-4 before:h-4 before:rounded-full before:bg-white before:top-0.5 before:left-0.5 before:transition-all before:duration-200 checked:before:translate-x-4"
            />
          </div>

          <!-- Allow duplicate results -->
          <div class="flex items-center justify-between py-2">
            <div class="flex flex-col gap-0.5">
              <label for="allow-duplicates-switch" class="text-sm font-bold text-slate-800 dark:text-slate-200">Cho phép trùng kết quả</label>
              <span class="text-xs text-slate-400 dark:text-slate-500">Kết quả đã trúng ở vòng trước vẫn có thể trúng tiếp ở vòng sau</span>
            </div>
            <input
              id="allow-duplicates-switch"
              type="checkbox"
              :checked="store.settings.allowDuplicates"
              @change="(e) => updateSetting('allowDuplicates', (e.target as HTMLInputElement).checked)"
              class="w-9 h-5 bg-slate-200 rounded-full appearance-none checked:bg-primary-600 relative transition-all duration-200 cursor-pointer before:content-[''] before:absolute before:w-4 before:h-4 before:rounded-full before:bg-white before:top-0.5 before:left-0.5 before:transition-all before:duration-200 checked:before:translate-x-4"
            />
          </div>
        </div>

        <!-- 4. Contact Info -->
        <div v-if="activeSubTab === 'contact'" class="space-y-5">
          <div class="space-y-1">
            <h4 class="text-sm font-bold text-slate-800 dark:text-slate-200">Liên hệ tác giả</h4>
            <p class="text-xs text-slate-400 dark:text-slate-500">Góp ý, báo lỗi hoặc hợp tác – đừng ngại liên hệ!</p>
          </div>

          <div class="grid grid-cols-1 gap-3">
            <!-- Gmail -->
            <a
              href="mailto:nket865@gmail.com"
              target="_blank"
              rel="noopener"
              class="contact-card group flex items-center gap-3 p-3.5 rounded-2xl border border-slate-200 dark:border-slate-800 hover:border-red-400 dark:hover:border-red-500 transition-all duration-200 hover:shadow-md hover:shadow-red-500/10"
            >
              <div class="w-9 h-9 rounded-xl flex items-center justify-center shrink-0" style="background: linear-gradient(135deg, #ea4335, #fbbc04)">
                <MailIcon class="h-4 w-4 text-white" />
              </div>
              <div class="flex flex-col min-w-0">
                <span class="text-xs font-semibold text-slate-400 dark:text-slate-500">Gmail</span>
                <span class="text-sm font-bold text-slate-800 dark:text-slate-200 truncate group-hover:text-red-500 transition-colors">nket865@gmail.com</span>
              </div>
              <ExternalLinkIcon class="h-3.5 w-3.5 text-slate-300 dark:text-slate-600 ml-auto shrink-0" />
            </a>

            <!-- Facebook -->
            <a
              href="https://www.facebook.com/ket.nguyenvan.587268"
              target="_blank"
              rel="noopener"
              class="contact-card group flex items-center gap-3 p-3.5 rounded-2xl border border-slate-200 dark:border-slate-800 hover:border-blue-400 dark:hover:border-blue-500 transition-all duration-200 hover:shadow-md hover:shadow-blue-500/10"
            >
              <div class="w-9 h-9 rounded-xl flex items-center justify-center shrink-0" style="background: linear-gradient(135deg, #1877f2, #42b0ff)">
                <FacebookIcon class="h-4 w-4 text-white" />
              </div>
              <div class="flex flex-col min-w-0">
                <span class="text-xs font-semibold text-slate-400 dark:text-slate-500">Facebook</span>
                <span class="text-sm font-bold text-slate-800 dark:text-slate-200 truncate group-hover:text-blue-500 transition-colors">Nguyễn Văn Kết</span>
              </div>
              <ExternalLinkIcon class="h-3.5 w-3.5 text-slate-300 dark:text-slate-600 ml-auto shrink-0" />
            </a>

            <!-- Zalo -->
            <a
              href="https://zalo.me/0348644630"
              target="_blank"
              rel="noopener"
              class="contact-card group flex items-center gap-3 p-3.5 rounded-2xl border border-slate-200 dark:border-slate-800 hover:border-sky-400 dark:hover:border-sky-500 transition-all duration-200 hover:shadow-md hover:shadow-sky-500/10"
            >
              <div class="w-9 h-9 rounded-xl flex items-center justify-center shrink-0" style="background: linear-gradient(135deg, #0068ff, #00aaff)">
                <MessageCircleIcon class="h-4 w-4 text-white" />
              </div>
              <div class="flex flex-col min-w-0">
                <span class="text-xs font-semibold text-slate-400 dark:text-slate-500">Zalo</span>
                <span class="text-sm font-bold text-slate-800 dark:text-slate-200 group-hover:text-sky-500 transition-colors">0348 644 630</span>
              </div>
              <ExternalLinkIcon class="h-3.5 w-3.5 text-slate-300 dark:text-slate-600 ml-auto shrink-0" />
            </a>

            <!-- Telegram -->
            <a
              href="https://t.me/0348644630"
              target="_blank"
              rel="noopener"
              class="contact-card group flex items-center gap-3 p-3.5 rounded-2xl border border-slate-200 dark:border-slate-800 hover:border-cyan-400 dark:hover:border-cyan-500 transition-all duration-200 hover:shadow-md hover:shadow-cyan-500/10"
            >
              <div class="w-9 h-9 rounded-xl flex items-center justify-center shrink-0" style="background: linear-gradient(135deg, #26a5e4, #00b4d8)">
                <SendIcon class="h-4 w-4 text-white" />
              </div>
              <div class="flex flex-col min-w-0">
                <span class="text-xs font-semibold text-slate-400 dark:text-slate-500">Telegram</span>
                <span class="text-sm font-bold text-slate-800 dark:text-slate-200 group-hover:text-cyan-500 transition-colors">0348 644 630</span>
              </div>
              <ExternalLinkIcon class="h-3.5 w-3.5 text-slate-300 dark:text-slate-600 ml-auto shrink-0" />
            </a>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <BaseButton variant="primary" @click="close">
        Xong
      </BaseButton>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useWheelStore } from '../../stores/wheelStore';
import BaseModal from '../common/BaseModal.vue';
import BaseButton from '../common/BaseButton.vue';
import { SettingsIcon, PaletteIcon, MusicIcon, FileQuestionIcon, MailIcon, FacebookIcon, MessageCircleIcon, SendIcon, ExternalLinkIcon, PhoneIcon } from 'lucide-vue-next';
import type { WheelSettings } from '../../types/wheel';

defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
}>();

const store = useWheelStore();
const activeSubTab = ref<string>('appearance');

const tabs = [
  { id: 'appearance', label: 'Thiết kế', icon: PaletteIcon },
  { id: 'audio', label: 'Âm thanh', icon: MusicIcon },
  { id: 'rules', label: 'Quy tắc', icon: FileQuestionIcon },
  { id: 'contact', label: 'Liên hệ', icon: PhoneIcon }
];

const themes: { id: WheelSettings['theme']; name: string; colors: string[] }[] = [
  { id: 'default', name: 'Mặc định', colors: ['#f43f5e', '#f97316', '#eab308', '#22c55e', '#3b82f6', '#8b5cf6', '#ec4899'] },
  { id: 'neon', name: 'Neon sáng', colors: ['#ff007f', '#39ff14', '#00ffff', '#ff00ff', '#ffff00', '#9d00ff'] },
  { id: 'candy', name: 'Kẹo ngọt', colors: ['#ff9eb5', '#ffbeb2', '#e2f0cb', '#b5ead7', '#c7ceea', '#ffdac1'] },
  { id: 'pastel', name: 'Pastel dịu', colors: ['#f4a261', '#e76f51', '#2a9d8f', '#e9c46a', '#457b9d', '#a8dadc'] },
  { id: 'monochrome', name: 'Đơn sắc', colors: ['#1e293b', '#334155', '#475569', '#64748b', '#94a3b8', '#cbd5e1'] },
  { id: 'ocean', name: 'Đại dương', colors: ['#0284c7', '#0369a1', '#0f766e', '#0d9488', '#06b6d4', '#0891b2'] }
];

const close = () => {
  emit('close');
};

const updateSetting = <K extends keyof WheelSettings>(key: K, value: WheelSettings[K]) => {
  store.settings[key] = value;
  store.saveToLocalStorage();
};
</script>
