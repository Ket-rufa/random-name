# COMPONENT ARCHITECTURE & UI

1. **Smart & Dumb Components:**
   - Phân biệt rõ Container Component (chứa logic gọi API, Pinia) và Presentational Component (chỉ nhận Props và hiển thị UI).
2. **Reusability:** Nếu một đoạn UI được dùng ở 2 nơi trở lên, BẮT BUỘC phải tách nó ra thành Component riêng đặt trong `src/components/common`.
3. **Slots:** Sử dụng Slots (`<slot />`) và Scoped Slots để làm cho các UI component (Card, Modal, Table) linh hoạt nhất có thể.
4. **Styling:** (Bạn có thể quy định thêm AI dùng TailwindCSS, SCSS hay CSS Modules ở đây). Không viết CSS global trừ khi thật sự cần thiết. Dùng `<style scoped>`.