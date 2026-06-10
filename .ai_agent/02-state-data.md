# STATE MANAGEMENT & DATA FETCHING

1. **Global State:** CHỈ sử dụng Pinia. KHÔNG sử dụng Vuex.
2. **Pinia Setup Syntax:** Viết Pinia store dưới dạng Setup Store (tương tự Composition API), không dùng Option Store. 
   - State = `ref()`
   - Getters = `computed()`
   - Actions = `function()`
3. **Data Fetching:** 
   - Tách logic gọi API ra các file "Composables" (ví dụ: `useUserFetch.ts`) thay vì nhét thẳng vào Component.
   - Luôn xử lý các trạng thái: `isLoading`, `isError`, và `data`.
4. **Prop Drilling:** Nếu component lồng nhau quá 3 tầng, hãy sử dụng `Provide / Inject` thay vì truyền props liên tục.