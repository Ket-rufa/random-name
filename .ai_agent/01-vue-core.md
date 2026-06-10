# VUE 3 CORE SKILLS & RULES

1. **Strict Composition API:** MỌI component đều phải sử dụng `<script setup>` và Composition API. TUYỆT ĐỐI KHÔNG sử dụng Options API (`data`, `methods`, `created`).
2. **TypeScript by Default:** Luôn sử dụng `<script setup lang="ts">`. Phải define Types/Interfaces rõ ràng cho Props, Emits và Refs.
3. **Reactivity:** 
   - Ưu tiên sử dụng `ref` cho các kiểu dữ liệu nguyên thủy (string, boolean, number).
   - Chỉ dùng `reactive` cho các object phức tạp có tính liên kết chặt chẽ.
4. **Macros:** Sử dụng `defineProps`, `defineEmits`, và `defineExpose` chuẩn xác. Không cần import các macro này.
5. **Lifecycle Hooks:** Sử dụng `onMounted`, `onUnmounted`, v.v., thay vì các hook cũ.