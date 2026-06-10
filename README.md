# Random Name Wheel (Vòng Quay May Mắn)

Một website vòng quay lựa chọn ngẫu nhiên hoàn chỉnh với thiết kế giao diện hiện đại, trải nghiệm mượt mà, hỗ trợ tùy chỉnh chủ đề, âm thanh, trọng số và đồng bộ hóa máy chủ thông tin.

## 🌟 Chức năng nổi bật

- **Vòng quay Canvas mượt mà**: Sử dụng HTML5 Canvas API với cơ chế chống vỡ hình trên màn hình Retina (High-DPI) và hiệu ứng giảm tốc tự nhiên (Ease-out).
- **Âm thanh chân thực**: Tự động tổng hợp âm thanh click cơ học và nhạc chiến thắng bằng **Web Audio API** mà không cần tải các file nhạc nặng.
- **Hiệu ứng chúc mừng**: Bắn pháo hoa giấy chúc mừng (Confetti) tưng bừng khi có kết quả.
- **Tùy chỉnh linh hoạt**: Điều chỉnh chủ đề màu sắc, cỡ chữ, thời gian quay, âm lượng, quy tắc tự động xóa và chế độ trọng số.
- **Quản lý danh sách tiện lợi**: Cho phép nhập nhanh dạng văn bản (mỗi dòng một lựa chọn) hoặc quản lý chi tiết trọng số và màu sắc của từng phân đoạn.
- **Đồng bộ hóa Server**: Lưu vòng quay lên PostgreSQL, sinh mã chia sẻ công khai và mã chỉnh sửa an toàn (lưu trong `localStorage` chủ sở hữu, bảo mật bằng mã hóa SHA-256 trên DB).

---

## 🛠️ Yêu cầu môi trường

- **Node.js** (v18 trở lên) & **npm**
- **Go** (v1.21 trở lên)
- **Docker Desktop** (để chạy PostgreSQL) hoặc PostgreSQL chạy local

---

## 🚀 Hướng dẫn khởi chạy nhanh

### Bước 1: Khởi chạy Database
Chạy container PostgreSQL 16 tại thư mục gốc bằng Docker Compose:
```bash
docker compose up -d
```
Cơ sở dữ liệu sẽ chạy trên cổng `5432` với tên DB `random_name_wheel`, tài khoản mặc định `postgres/postgres`.

### Bước 2: Khởi chạy Backend Go
1. Di chuyển vào thư mục backend:
   ```bash
   cd backend
   ```
2. Tạo file cấu hình môi trường `.env` từ file mẫu:
   - Trên **cmd/bash**: `cp .env.example .env`
   - Trên **PowerShell**: `Copy-Item .env.example .env`
3. Cài đặt các thư viện và chạy server:
   ```bash
   go mod tidy
   go run ./cmd/server
   ```
   Server backend sẽ chạy tại: `http://localhost:8080` (Tự động migrate cơ sở dữ liệu trên lần chạy đầu tiên).

### Bước 3: Khởi chạy Frontend Vue 3
1. Quay lại thư mục gốc dự án.
2. Khởi chạy môi trường phát triển:
   ```bash
   npm install
   npm run dev
   ```
   Mở trình duyệt truy cập: `http://localhost:5173`

---

## 📂 Cấu trúc thư mục dự án

```text
random-name/
├── backend/                   # Mã nguồn Go & Fiber API
│   ├── cmd/server/main.go     # Entrypoint khởi chạy server
│   ├── internal/
│   │   ├── config/            # Load cấu hình .env
│   │   ├── database/          # Kết nối database & Migration
│   │   ├── models/            # Schema GORM DB Models
│   │   ├── repositories/      # Các truy vấn SQL database
│   │   ├── services/          # Logic nghiệp vụ (Secure, Hash, Duplication)
│   │   ├── handlers/          # Controller tiếp nhận HTTP request
│   │   └── routes/            # Khai báo đường dẫn Router API
│   ├── go.mod
│   └── .env
├── src/                       # Frontend Vue 3 & TypeScript
│   ├── api/                   # Gọi API Axios client
│   ├── components/            # Giao diện components
│   │   ├── common/            # Component dùng chung (Button, Modal, Toast)
│   │   ├── entries/           # Bộ soạn thảo danh sách
│   │   ├── history/           # Danh sách trúng giải
│   │   ├── layout/            # Header điều khiển
│   │   └── wheel/             # Bản vẽ vòng quay Canvas
│   ├── composables/           # Reactivity hooks (Audio, Confetti, Animation)
│   ├── stores/                # Pinia store quản lý trạng thái
│   ├── types/                 # TypeScript interfaces
│   ├── views/                 # Trang Home và trang Shared View
│   ├── App.vue
│   └── main.ts
├── docker-compose.yml         # File khởi tạo Docker Postgres DB
├── package.json
└── README.md
```

---

## 🌐 Các API Endpoints (`/api/v1`)

| Phương thức | Đường dẫn | Chức năng | Quyền yêu cầu |
| :--- | :--- | :--- | :--- |
| **GET** | `/health` | Kiểm tra trạng thái server | Công khai |
| **POST** | `/wheels` | Lưu vòng quay mới lên server | Công khai (trả về `editToken`) |
| **GET** | `/wheels/:shareCode` | Tải vòng quay qua mã chia sẻ | Công khai |
| **PUT** | `/wheels/:id` | Cập nhật vòng quay | Chủ sở hữu (`X-Edit-Token`) |
| **DELETE** | `/wheels/:id` | Xóa vòng quay khỏi server | Chủ sở hữu (`X-Edit-Token`) |
| **POST** | `/wheels/:id/spin` | Ghi nhận phân tích lượt quay | Công khai |
| **GET** | `/wheels/:id/history` | Lấy lịch sử quay trên server | Công khai |
| **DELETE**| `/wheels/:id/history` | Xóa lịch sử quay của vòng quay | Chủ sở hữu (`X-Edit-Token`) |
| **POST** | `/wheels/:id/duplicate` | Nhân bản vòng quay thành bản mới| Công khai |

---

## 🏗️ Build Production

### Frontend
Build mã nguồn frontend tĩnh ra thư mục `dist/`:
```bash
npm run build
```
Bạn chỉ cần đưa thư mục `dist/` lên các dịch vụ hosting tĩnh (Vercel, Netlify, Nginx).

### Backend
Biên dịch dự án Go thành file binary chạy độc lập:
```bash
cd backend
go build -o server ./cmd/server
```

---

## 🔧 Xử lý lỗi thường gặp

1. **Lỗi kết nối cơ sở dữ liệu (`Failed to connect to database`)**:
   - Đảm bảo bạn đã bật Docker Desktop và container Postgres đang chạy (`docker ps`).
   - Kiểm tra xem cổng 5432 có bị chiếm dụng bởi service Postgres nào khác trên máy chủ hay không.
2. **Lỗi CORS trên trình duyệt**:
   - Hãy chắc chắn giá trị `FRONTEND_URL` trong file `backend/.env` khớp hoàn toàn với địa chỉ cổng frontend chạy thực tế (ví dụ: `http://localhost:5173`).
3. **Không nghe thấy tiếng click**:
   - Theo chính sách bảo mật của các trình duyệt hiện đại, bạn cần phải click tương tác lên trang web ít nhất một lần để kích hoạt `AudioContext` phát nhạc.
