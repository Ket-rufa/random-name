<template>
  <div class="flex flex-col bg-white dark:bg-slate-900 rounded-3xl border border-slate-100 dark:border-slate-800/60 shadow-sm overflow-hidden p-5 max-h-[300px] shrink-0 transition-all duration-300">
    <div class="flex items-center justify-between mb-3.5">
      <h3 class="text-sm font-bold text-slate-850 dark:text-slate-200 flex items-center gap-2">
        <HistoryIcon class="h-4.5 w-4.5 text-primary-500" />
        Lịch sử kết quả
      </h3>
      
      <div class="flex items-center gap-3" v-if="store.spinHistory.length > 0">
        <button
          @click="exportHistoryCSV"
          class="text-xs font-semibold text-slate-500 hover:text-primary-600 transition-colors cursor-pointer flex items-center gap-1"
          title="Xuất lịch sử ra file CSV"
        >
          <DownloadIcon class="h-3 w-3" />
          <span>Xuất CSV</span>
        </button>
        <button
          @click="clearAllHistory"
          class="text-xs font-semibold text-slate-400 hover:text-red-500 transition-colors cursor-pointer"
        >
          Xóa hết
        </button>
      </div>
    </div>

    <!-- History log list -->
    <div class="flex-1 overflow-y-auto pr-1 space-y-2 scrollbar text-xs">
      <div
        v-for="(log, index) in store.spinHistory"
        :key="log.id"
        class="flex items-center justify-between p-2.5 rounded-xl bg-slate-50 dark:bg-slate-950/30 border border-slate-100/50 dark:border-slate-800/40"
      >
        <div class="flex items-center gap-2 overflow-hidden">
          <span class="font-bold text-slate-400">#{{ store.spinHistory.length - index }}</span>
          <span class="font-semibold text-slate-850 dark:text-slate-250 truncate" :title="log.resultLabel">
            {{ log.resultLabel }}
          </span>
        </div>
        <div class="flex items-center gap-2 text-slate-400 shrink-0">
          <span>{{ formatTime(log.spunAt) }}</span>
          <button
            @click="store.removeHistoryItem(log.id)"
            class="text-slate-350 hover:text-red-500 p-0.5 rounded transition-colors cursor-pointer"
            title="Xóa log này"
          >
            <TrashIcon class="h-3.5 w-3.5" />
          </button>
        </div>
      </div>

      <div v-if="store.spinHistory.length === 0" class="py-8 text-center text-slate-400 dark:text-slate-500">
        <p>Chưa có lượt quay nào</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useWheelStore } from '../../stores/wheelStore';
import { HistoryIcon, TrashIcon, DownloadIcon } from 'lucide-vue-next';

const store = useWheelStore();

const formatTime = (isoString: string) => {
  try {
    const d = new Date(isoString);
    return d.toLocaleTimeString('vi-VN', { hour: '2-digit', minute: '2-digit', second: '2-digit' });
  } catch (e) {
    return '';
  }
};

const clearAllHistory = () => {
  if (confirm('Bạn có chắc chắn muốn xóa toàn bộ lịch sử quay không?')) {
    store.clearHistory();
    store.showToast('Đã xóa lịch sử quay', 'info');
  }
};

const exportHistoryCSV = () => {
  if (store.spinHistory.length === 0) {
    store.showToast('Lịch sử quay rỗng!', 'info');
    return;
  }
  
  let csvContent = "\uFEFF"; // UTF-8 BOM
  csvContent += "STT,Kết quả,Thời gian quay\n";
  
  store.spinHistory.forEach((log, index) => {
    const time = formatTime(log.spunAt);
    const label = log.resultLabel.replace(/"/g, '""');
    csvContent += `${store.spinHistory.length - index},"${label}","${time}"\n`;
  });
  
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.setAttribute("href", url);
  link.setAttribute("download", `lich_su_quay_${Date.now()}.csv`);
  link.style.visibility = 'hidden';
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  
  store.showToast('Đã xuất lịch sử quay dạng CSV!', 'success');
};
</script>
