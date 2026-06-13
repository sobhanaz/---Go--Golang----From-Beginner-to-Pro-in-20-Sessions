<div dir="rtl">

> 🌐 **زبان / Language:** فارسی (همین فایل) · [English](session-17.md)

# جلسه‌ی ۱۷ — سرور HTTP 🌐

**هدف (۱ ساعت):** با کتابخانه‌ی استاندارد Go سرور وب بسازی — به فریم‌ورک نیازی
نیست. هندلر می‌نویسی، درخواست‌ها را مسیریابی می‌کنی، پارامترهای کوئری/مسیر و بدنه‌ی
JSON را می‌خوانی و پاسخ‌های JSON با کد وضعیت درست برمی‌گردانی. این دقیقاً پایه‌ای
است که پروژه‌ی نهایی **TaskFlow** رویش ساخته می‌شود.

> **مرور جلسه‌ی ۱۶:** می‌توانی کد را تست کنی. این جلسه *همچنین* می‌بینی چطور
> هندلرهای HTTP را با `httptest` تست کنی — بدون راه‌اندازی سرور واقعی.

---

## ۱. کوچک‌ترین سرور (۱۵ دقیقه)

پکیج `net/http` در Go یک سرور وب در سطح محصول است که در کتابخانه‌ی استاندارد قرار
دارد (سیستم‌های عظیمی را می‌راند). کوچک‌ترین سرور:

</div>

```go
func main() {
    mux := http.NewServeMux()   // یک روتر: الگوهای URL -> هندلرها

    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello from Go!")
    })

    http.ListenAndServe(":8080", mux)   // مسدود می‌شود، تا توقف سرویس می‌دهد
}
```

<div dir="rtl">

قلب توسعه‌ی وب Go، **امضای هندلر** است:

</div>

```go
func(w http.ResponseWriter, r *http.Request)
```

<div dir="rtl">

- `w http.ResponseWriter` — **پاسخ را** در این می‌نویسی (بدنه، هدر، وضعیت).
- `r *http.Request` — **درخواست ورودی** (متد، URL، هدرها، بدنه).

اجرایش کن، بعد URL را ببین یا curl کن:

</div>

```bash
go run examples/session17/hello/hello.go
# در ترمینال دیگر:
curl localhost:8080/hello?name=Sobhan
```

<div dir="rtl">

(برای توقف سرور Ctrl+C بزن.) فایل [`examples/session17/hello/hello.go`](../examples/session17/hello/hello.go) را اجرا کن.

---

## ۲. مسیریابی: متدها، مسیرها و پارامترها (۱۵ دقیقه)

از **Go 1.22**، روتر داخلی متدهای HTTP و wildcardهای مسیر را می‌فهمد — پس اغلب
اصلاً به روتر شخص‌ثالث نیازی نداری:

</div>

```go
mux.HandleFunc("GET /tasks", listTasks)        // متد + مسیر
mux.HandleFunc("POST /tasks", createTask)
mux.HandleFunc("GET /tasks/{id}", getTask)     // {id} یک wildcard مسیر است
```

<div dir="rtl">

خواندن انواع مختلف ورودی:

</div>

```go
// پارامتر مسیر:  GET /tasks/42  ->  r.PathValue("id") == "42"
id := r.PathValue("id")

// پارامتر کوئری:  GET /search?q=go  ->  r.URL.Query().Get("q")
q := r.URL.Query().Get("q")

// بدنه‌ی JSON درخواست (مثلاً از یک POST):
var input struct{ Title string `json:"title"` }
json.NewDecoder(r.Body).Decode(&input)
```

<div dir="rtl">

> 🔑 `r.PathValue("id")` یک **رشته** برمی‌گرداند — با `strconv.Atoi` (جلسه‌ی ۱۴)
> تبدیلش کن و اگر عدد معتبری نبود خطا را مدیریت کن.

---

## ۳. برگرداندن JSON درست (۱۵ دقیقه)

یک REST API با JSON و یک **کد وضعیت HTTP** معنادار پاسخ می‌دهد. الگوی تمیز، یک
کمک‌کننده‌ی کوچک است:

</div>

```go
func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")  // ۱. نوع محتوا را تنظیم کن
    w.WriteHeader(status)                                // ۲. کد وضعیت را تنظیم کن
    json.NewEncoder(w).Encode(v)                         // ۳. بدنه‌ی JSON را بنویس
}
```

<div dir="rtl">

> ⚠️ **ترتیب مهم است!** اول هدرها، بعد `WriteHeader(status)`، بعد بدنه را بنویس.
> به‌محض نوشتن بدنه، وضعیت قفل می‌شود. فراخوانی `WriteHeader` بعد از نوشتن همیشه به
> ۲۰۰ پیش‌فرض می‌شود.

کدهای وضعیت رایجی که استفاده می‌کنی:

| کد | ثابت | معنی |
|----|------|------|
| 200 | `http.StatusOK` | موفقیت (GET) |
| 201 | `http.StatusCreated` | ساخته شد (POST) |
| 400 | `http.StatusBadRequest` | ورودی بد از کلاینت |
| 404 | `http.StatusNotFound` | منبع وجود ندارد |
| 500 | `http.StatusInternalServerError` | خرابی سمت سرور |

از **ثابت‌های نام‌دار** استفاده کن، نه اعداد خام — واضح‌تر و خودتوضیح‌اند.

---

## ۴. ساختار برای قابلیت تست — استراکت `Server` (۱۵ دقیقه)

یک الگوی حرفه‌ای (و همانی که در TaskFlow استفاده می‌کنیم): وابستگی‌هایت (دیتابیس،
تنظیمات، انبار حافظه‌ای) را در یک استراکت `Server` بگذار، هندلرها را **متد** روی آن
کن، و یک متد `routes()` که روتر را برمی‌گرداند ارائه بده.

</div>

```go
type Server struct {
    mu    sync.Mutex
    tasks map[int]Task   // بعداً: یک دیتابیس واقعی
}

func (s *Server) handleListTasks(w http.ResponseWriter, r *http.Request) { ... }

func (s *Server) routes() http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /tasks", s.handleListTasks)
    mux.HandleFunc("POST /tasks", s.handleCreateTask)
    return mux
}
```

<div dir="rtl">

چرا مهم است: هندلرها را **به‌سادگی قابل‌تست** می‌کند. پکیج `net/http/httptest`
می‌گذارد یک درخواست تقلبی به `routes()` بفرستی و پاسخ را بازرسی کنی — **بدون سرور
واقعی، بدون پورت واقعی:**

</div>

```go
srv := NewServer()
req := httptest.NewRequest("GET", "/health", nil)
rec := httptest.NewRecorder()        // پاسخ را ضبط می‌کند
srv.routes().ServeHTTP(rec, req)

if rec.Code != http.StatusOK { t.Fatalf("got %d", rec.Code) }
```

<div dir="rtl">

این جلسه‌ی ۱۶ (تست) را با این جلسه ترکیب می‌کند، و روشی است که APIهای واقعی Go در
CI تست می‌شوند. مثال کامل و تست‌هایش را مطالعه کن:

</div>

```bash
go run examples/session17/jsonapi/jsonapi.go    # سرور واقعی را اجرا کن
go test -v ./examples/session17/jsonapi/        # هندلرها را تست کن (بدون سرور)
```

<div dir="rtl">

فایل‌های [`examples/session17/jsonapi/jsonapi.go`](../examples/session17/jsonapi/jsonapi.go) و
[`jsonapi_test.go`](../examples/session17/jsonapi/jsonapi_test.go) را ببین.

> 📦 **پیش‌نمایش میان‌افزار:** یک میان‌افزار تابعی است که یک هندلر را می‌پیچد تا رفتار
> عرضی (لاگ، احراز هویت، بازیابی) اضافه کند. این‌ها را در جلسه‌ی ۱۹ می‌سازیم.
> شکلشان: `func(next http.Handler) http.Handler`.

---

## 🎯 تمرین‌ها (قبل از جلسه‌ی ۱۸ انجام بده!)

در `examples/session17/practice/` کار کن:

۱. **سرور اکو:** یک endpoint با `GET /echo?msg=hi` که JSON `{"echo":"hi"}` برمی‌گرداند. وقتی `msg` نیست به `"nothing"` پیش‌فرض شو.
۲. **سلام با پارامتر مسیر:** `GET /greet/{name}` که `{"greeting":"Hello, <name>"}` برمی‌گرداند.
۳. **API یادداشت حافظه‌ای:** `POST /notes` (ساخت از بدنه‌ی JSON) و `GET /notes` (فهرست همه). از یک استراکت `Server` با اسلایس یا مپ استفاده کن.
۴. **کدهای وضعیت:** کاری کن `POST /notes` اگر بدنه فیلد `text` نداشت ۴۰۰ و در موفقیت ۲۰۱ برگرداند.
۵. **تستش کن:** تست‌های مبتنی بر `httptest` برای API یادداشت بنویس که ساخت، فهرست و حالت اعتبارسنجی ۴۰۰ را پوشش دهد. `go test` را اجرا کن.

---

## ✅ چک‌لیست جلسه‌ی ۱۷

- [ ] می‌توانم سرور را با `http.NewServeMux` + `http.ListenAndServe` راه بیندازم
- [ ] امضای هندلر `func(w, r)` را می‌فهمم
- [ ] می‌توانم با متد و مسیر، از جمله wildcardهای `{id}`، مسیریابی کنم
- [ ] می‌توانم پارامتر مسیر، پارامتر کوئری و بدنه‌ی JSON را بخوانم
- [ ] می‌توانم JSON را با `Content-Type` و کد وضعیت درست برگردانم
- [ ] می‌دانم باید هدر/وضعیت را قبل از نوشتن بدنه تنظیم کنم
- [ ] هندلرها را به‌صورت متد روی `Server` برای قابلیت تست ساختار می‌دهم
- [ ] می‌توانم هندلرها را با `httptest` تست کنم (بدون سرور واقعی)
- [ ] هر ۵ تمرین را انجام دادم

**قبلی:** [→ جلسه‌ی ۱۶](session-16.fa.md) · **بعدی:** [جلسه‌ی ۱۸ — REST API + پایگاه‌داده ←](session-18.fa.md)

---

🎉 **نقطه‌ی عطف:** بخش ۴ (Go در دنیای واقعی) تمام شد! حالا می‌توانی سرویس‌های وب
JSON تست‌شده بسازی. بعدی **بخش ۵** را شروع می‌کنیم — ساخت پروژه‌ی کامل نمونه‌کار
TaskFlow، با شروع از ذخیره‌سازی پشتیبانی‌شده با دیتابیس.

</div>
