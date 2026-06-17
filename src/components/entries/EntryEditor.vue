<template>
  <div class="flex flex-col h-full bg-white dark:bg-slate-900 rounded-3xl border border-slate-100 dark:border-slate-800/60 shadow-sm overflow-hidden transition-all duration-300">
    <!-- Tab Headers -->
    <div class="flex border-b border-slate-100 dark:border-slate-800/60 bg-slate-50/50 dark:bg-slate-900/10">
      <button
        @click="activeTab = 'text'"
        class="flex-1 py-3 px-4 text-sm font-medium border-b-2 transition-all duration-200"
        :class="[
          activeTab === 'text'
            ? 'border-primary-600 text-primary-600 dark:border-primary-500 dark:text-primary-500'
            : 'border-transparent text-slate-500 hover:text-slate-700 dark:hover:text-slate-300'
        ]"
      >
        Nhập nhanh
      </button>
      <button
        @click="activeTab = 'list'"
        class="flex-1 py-3 px-4 text-sm font-medium border-b-2 transition-all duration-200"
        :class="[
          activeTab === 'list'
            ? 'border-primary-600 text-primary-600 dark:border-primary-500 dark:text-primary-500'
            : 'border-transparent text-slate-500 hover:text-slate-700 dark:hover:text-slate-300'
        ]"
      >
        Quản lý chi tiết
      </button>
    </div>

    <!-- Tab Content -->
    <div class="flex-1 p-5 overflow-hidden flex flex-col">
      <!-- 1. Text Mode (Textarea) -->
      <div v-if="activeTab === 'text'" class="flex-1 flex flex-col h-full">
        <label for="entries-textarea" class="sr-only">Danh sách lựa chọn</label>
        <textarea
          id="entries-textarea"
          ref="textareaRef"
          v-model="textareaContent"
          @input="handleTextareaInput"
          :disabled="!store.canEdit"
          placeholder="Mỗi dòng tương ứng với một lựa chọn...&#10;Ví dụ:&#10;Nguyễn Văn A&#10;Trần Thị B&#10;Lê Văn C"
          class="w-full flex-1 p-4 rounded-2xl border border-slate-200 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-950/40 focus:bg-white dark:focus:bg-slate-950 focus:ring-2 focus:ring-primary-500/20 focus:border-primary-500 outline-none resize-none font-medium text-sm transition-all duration-200 scrollbar"
        ></textarea>
      </div>

      <!-- 2. Detailed List Mode -->
      <div v-else class="flex-1 flex flex-col overflow-hidden">
        <!-- Add new entry input -->
        <div v-if="store.canEdit" class="flex gap-2 mb-4">
          <input
            v-model="newEntryLabel"
            @keyup.enter="addNewEntry"
            type="text"
            placeholder="Thêm lựa chọn mới..."
            class="flex-1 px-4 py-2 text-sm rounded-xl border border-slate-250 dark:border-slate-800 bg-slate-50/50 dark:bg-slate-950/40 focus:ring-2 focus:ring-primary-500/20 focus:border-primary-500 outline-none transition-all duration-200"
          />
          <button
            @click="addNewEntry"
            class="px-4 py-2 bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 text-sm font-semibold rounded-xl transition-all duration-200"
          >
            Thêm
          </button>
        </div>

        <!-- Scrollable list of items -->
        <div class="flex-1 overflow-y-auto pr-1 space-y-2.5">
          <div
            v-for="(entry, index) in store.entries"
            :key="entry.id"
            class="flex items-center gap-3 p-3 rounded-2xl border border-slate-100 dark:border-slate-800/80 bg-slate-50/30 dark:bg-slate-900/10 group transition-all duration-200 hover:shadow-sm"
          >
            <!-- Colored dot (or color picker) -->
            <div class="relative shrink-0">
              <input
                type="color"
                :value="store.getEntryColor(index, entry.color)"
                @change="(e) => handleColorChange(entry.id, (e.target as HTMLInputElement).value)"
                :disabled="!store.canEdit"
                class="w-6 h-6 rounded-lg cursor-pointer border-0 p-0 outline-none overflow-hidden"
                title="Đổi màu sắc phân đoạn"
              />
            </div>

            <!-- Input text label -->
            <input
              type="text"
              :value="entry.label"
              @change="(e) => handleLabelChange(entry.id, (e.target as HTMLInputElement).value)"
              :disabled="!store.canEdit"
              class="flex-1 bg-transparent border-0 focus:ring-1 focus:ring-primary-500 rounded px-1.5 py-0.5 text-sm font-semibold outline-none text-slate-800 dark:text-slate-100"
            />

            <!-- Weight input (shown if weights are enabled) -->
            <div v-if="store.settings.enableWeights" class="flex items-center gap-1">
              <span class="text-xs text-slate-400 dark:text-slate-500 font-medium">Hệ số:</span>
              <input
                type="number"
                step="0.5"
                min="0.1"
                max="10"
                :value="entry.weight"
                @change="(e) => handleWeightChange(entry.id, parseFloat((e.target as HTMLInputElement).value))"
                :disabled="!store.canEdit"
                class="w-12 px-1 py-0.5 text-center text-xs rounded border border-slate-200 dark:border-slate-700 bg-white dark:bg-slate-900 outline-none focus:ring-1 focus:ring-primary-500"
              />
            </div>

            <!-- Delete button -->
            <button
              v-if="store.canEdit"
              @click="store.removeEntry(entry.id)"
              class="p-1 rounded-lg text-slate-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-950/20 opacity-0 group-hover:opacity-100 focus:opacity-100 transition-all duration-200"
              aria-label="Xóa lựa chọn này"
            >
              <TrashIcon class="h-4 w-4" />
            </button>
          </div>

          <!-- Empty list state -->
          <div v-if="store.entries.length === 0" class="py-12 text-center text-slate-400 dark:text-slate-500">
            <ListCollapseIcon class="h-10 w-10 mx-auto opacity-40 mb-3" />
            <p class="text-sm">Danh sách trống. Nhập lựa chọn để bắt đầu!</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick action footer bar -->
    <div class="px-5 py-4 bg-slate-50 dark:bg-slate-900/40 border-t border-slate-100 dark:border-slate-800/60 flex items-center justify-between flex-wrap gap-2 text-xs font-semibold">
      <div class="text-slate-500 dark:text-slate-400">
        Tổng cộng: <span class="text-sm font-bold text-slate-700 dark:text-slate-200">{{ store.entries.length }}</span> dòng
      </div>

      <div v-if="store.canEdit" class="flex items-center gap-1">
        <button
          @click="shuffle"
          class="p-2 rounded-xl bg-white dark:bg-slate-800 border border-slate-200/60 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 flex items-center gap-1 transition-all duration-150"
          title="Xáo trộn vị trí"
        >
          <ShuffleIcon class="h-3.5 w-3.5" />
          <span>Trộn</span>
        </button>

        <button
          @click="sort"
          class="p-2 rounded-xl bg-white dark:bg-slate-800 border border-slate-200/60 dark:border-slate-700 hover:bg-slate-100 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 flex items-center gap-1 transition-all duration-150"
          title="Sắp xếp theo bảng chữ cái A-Z"
        >
          <SortAscIcon class="h-3.5 w-3.5" />
          <span>Xếp A-Z</span>
        </button>

        <button
          @click="confirmClear"
          class="p-2 rounded-xl bg-white dark:bg-slate-800 border border-slate-200/60 dark:border-slate-700 hover:bg-red-50 dark:hover:bg-red-950/20 hover:border-red-200 text-slate-700 dark:text-slate-200 hover:text-red-650 flex items-center gap-1 transition-all duration-150"
          title="Xóa toàn bộ danh sách"
        >
          <Trash2Icon class="h-3.5 w-3.5" />
          <span>Xóa hết</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import { useWheelStore } from '../../stores/wheelStore';
import { TrashIcon, ShuffleIcon, SortAscIcon, Trash2Icon, ListCollapseIcon } from 'lucide-vue-next';

const store = useWheelStore();
const activeTab = ref<'text' | 'list'>('text');
const textareaContent = ref<string>('');
const newEntryLabel = ref<string>('');

// Sync local textarea content when entries change in Pinia
watch(
  () => store.entries,
  () => {
    // Only update textarea content if not focused or active tab isn't text, to prevent cursor jumping
    if (activeTab.value !== 'text' || document.activeElement !== document.getElementById('entries-textarea')) {
      textareaContent.value = store.getEntriesText();
    }
  },
  { deep: true, immediate: true }
);

const handleTextareaInput = () => {
  store.setEntriesFromText(textareaContent.value);
};

const addNewEntry = () => {
  if (newEntryLabel.value.trim()) {
    store.addEntry(newEntryLabel.value);
    newEntryLabel.value = '';
    // Scroll detailed list down
  }
};

const handleLabelChange = (id: string, newLabel: string) => {
  store.updateEntry(id, { label: newLabel });
};

const handleWeightChange = (id: string, newWeight: number) => {
  if (isNaN(newWeight) || newWeight <= 0) return;
  store.updateEntry(id, { weight: newWeight });
};

const handleColorChange = (id: string, newColor: string) => {
  store.updateEntry(id, { color: newColor });
};

const shuffle = () => {
  store.shuffleEntries();
  textareaContent.value = store.getEntriesText();
};

const sort = () => {
  store.sortEntries();
  textareaContent.value = store.getEntriesText();
};

const confirmClear = () => {
  if (confirm('Bạn có chắc chắn muốn xóa toàn bộ danh sách lựa chọn không?')) {
    store.clearEntries();
    textareaContent.value = '';
    store.showToast('Đã xóa toàn bộ danh sách', 'info');
  }
};

onMounted(() => {
  textareaContent.value = store.getEntriesText();
});
</script>
