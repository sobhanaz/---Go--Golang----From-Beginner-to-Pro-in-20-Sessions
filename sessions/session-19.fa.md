<div dir="rtl">

> 🌐 **زبان / Language:** فارسی (همین فایل) · [English](session-19.md)

# جلسه‌ی ۱۹ — احراز هویت، میان‌افزار و تنظیمات 🔐

**هدف (۱ ساعت):** TaskFlow را از یک اسباب‌بازی به یک بک‌اند واقعی تبدیل کنی.
**تنظیمات** را از محیط می‌خوانی، **میان‌افزار** (لاگ، بازیابی از panic، احراز
هویت) و **احراز هویت JWT** با هش رمز اضافه می‌کنی — و هر وظیفه را **محدود به مالکش**
می‌کنی تا کاربران فقط داده‌ی خودشان را ببینند. این جلسه پروژه‌ات را واقعاً
قابل‌ارائه در رزومه می‌کند.

> **مرور جلسه‌ی ۱۸:** TaskFlow یک API از نوع CRUD لایه‌ای و پشتیبانی‌شده با دیتابیس
> دارد (models ← store ← api). حالا امن و آماده‌ی محصولش می‌کنیم.

> 📂 همه‌ی کد در [`taskflow/`](../taskflow/) است. تست‌ها را با
> `cd taskflow && go test ./...` اجرا کن، یا سرور را با `go run .`.

---

## ۱. تنظیمات از محیط (۱۰ دقیقه)

هاردکد کردن پورت، مسیر دیتابیس و رازها اشتباه است — اپ‌های واقعی آن‌ها را از محیط
می‌خوانند تا همان فایل اجرایی بدون تغییر در dev، staging و production اجرا شود. یک
پکیج کوچک [`config`](../taskflow/internal/config/config.go) ساختیم:

</div>

```go
type Config struct {
    Addr        string // TASKFLOW_ADDR   (پیش‌فرض ":8080")
    DatabaseDSN string // TASKFLOW_DB     (پیش‌فرض "taskflow.db")
    JWTSecret   string // TASKFLOW_JWT_SECRET
}

func Load() Config { /* os.Getenv با مقادیر پیش‌فرض */ }
```

```bash
TASKFLOW_ADDR=:9000 TASKFLOW_JWT_SECRET=super-secret go run .
```

<div dir="rtl">

> 🔑 **هرگز رازهای واقعی را commit نکن.** راز پیش‌فرض JWT اینجا فقط برای dev محلی
> است. در production مقدار `TASKFLOW_JWT_SECRET` را از طریق محیط (یا یک مدیر راز)
> تزریق می‌کنی. این نکته‌ای است که در مصاحبه ارزش گفتن دارد.

---

## ۲. هش رمز و JWTها (۲۰ دقیقه)

دو پکیج شخص‌ثالث جدید (با `go get` اضافه شدند):

- `golang.org/x/crypto/bcrypt` — هش رمز.
- `github.com/golang-jwt/jwt/v5` — صدور و تأیید توکن.

### هش رمز — هرگز متن خام ذخیره نکن

[`auth/password.go`](../taskflow/internal/auth/password.go):

</div>

```go
func HashPassword(password string) (string, error) {
    b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(b), err
}

func CheckPassword(hash, password string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
```

<div dir="rtl">

> 🔑 **bcrypt عمداً کند است و خودکار salt می‌زند.** فقط هش را ذخیره می‌کنی؛ هرگز
> نمی‌توانی رمز را پس بگیری. هنگام ورود *مقایسه* می‌کنی، هرگز رمزگشایی نمی‌کنی.
> ذخیره‌ی رمز خام گناه امنیتی شماره‌ی ۱ است — bcrypt راه‌حل است.

### JSON Web Token — احراز هویت بدون‌حالت

یک **JWT** توکنی امضاشده است که کلاینت در هر درخواست می‌فرستد تا هویتش را اثبات
کند. [`auth/jwt.go`](../taskflow/internal/auth/jwt.go) شناسه‌ی کاربر را در
`Subject` توکن می‌گذارد و با راز امضایش می‌کند:

</div>

```go
func GenerateToken(secret string, userID int64) (string, error) {
    claims := jwt.RegisteredClaims{
        Subject:   strconv.FormatInt(userID, 10),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```

<div dir="rtl">

`ParseToken` این را برعکس می‌کند: امضا را تأیید می‌کند (الگوریتم‌های امضای
غیرمنتظره را رد می‌کند — یک بررسی امنیتی واقعی)، منقضی‌نبودن توکن را تأیید و شناسه‌ی
کاربر را برمی‌گرداند. چون توکن **امضاشده** است، سرور نیازی به ذخیره‌ی نشست ندارد —
می‌تواند به یک توکن معتبر اعتماد کند. این یعنی احراز هویت «بدون‌حالت».

---

## ۳. میان‌افزار — رفتار عرضی (۱۵ دقیقه)

یک **میان‌افزار** یک هندلر را می‌پیچد تا رفتاری اضافه کند، با این امضا:

</div>

```go
func(next http.Handler) http.Handler
```

<div dir="rtl">

یک هندلر جدید برمی‌گرداند که *قبل/بعد* از فراخوانی `next` کاری انجام می‌دهد. سه‌تا
در [`api/middleware.go`](../taskflow/internal/api/middleware.go) ساختیم:

- **`Logging`** — برای هر درخواست `method path -> status (duration)` را ثبت می‌کند.
  `ResponseWriter` را در یک `statusRecorder` می‌پیچد تا کد وضعیت را بگیرد.
- **`Recovery`** — یک `recover()` را `defer` می‌کند (جلسه‌ی ۱۱!) تا panic در هر
  هندلر یک ۵۰۰ تمیز برگرداند نه این که کل سرور را کرش کند.
- **`Auth`** — هدر `Authorization: Bearer <token>` را می‌خواند، JWT را تأیید و
  شناسه‌ی کاربر را در **context** درخواست می‌گذارد. بدون توکن ← `401`.

### پاس دادن داده از طریق `context`

`Auth` شناسه‌ی کاربر احرازشده را در context درخواست می‌گذارد تا هندلرها بخوانندش:

</div>

```go
ctx := context.WithValue(r.Context(), userIDKey, userID)
next.ServeHTTP(w, r.WithContext(ctx))
```

```go
func userIDFromContext(r *http.Request) int64 {
    id, _ := r.Context().Value(userIDKey).(int64)
    return id
}
```

<div dir="rtl">

> 🔑 **از یک نوع کلید export‌نشده** (`type ctxKey string`) برای کلیدهای context
> استفاده کن تا پکیج‌های دیگر اشتباهاً با مال تو تداخل نکنند. این الگوی اصطلاحی
> مقادیر دامنه‌ی-درخواست است.

### اعمال میان‌افزار

میان‌افزار سراسری کل mux را می‌پیچد؛ میان‌افزار هر-مسیر هندلرهای جداگانه را. در
[`api/server.go`](../taskflow/internal/api/server.go):

</div>

```go
// مسیرهای عمومی — بدون احراز هویت.
mux.HandleFunc("POST /auth/register", s.handleRegister)
mux.HandleFunc("POST /auth/login", s.handleLogin)

// مسیرهای محافظت‌شده — با Auth پیچیده شده.
mux.Handle("GET /tasks", s.Auth(http.HandlerFunc(s.handleListTasks)))
// ...

// زنجیره‌ی سراسری: Recovery بیرونی‌ترین است تا همه‌چیز را بگیرد.
return Recovery(Logging(mux))
```

<div dir="rtl">

---

## ۴. داده‌ی محدود به کاربر — بخش واقع‌گرایانه (۱۵ دقیقه)

یک اپ چندکاربره‌ی واقعی باید داده را ایزوله کند. کاری کردیم **هر وظیفه به یک کاربر
تعلق داشته باشد**:

- جدول `tasks` یک ستون `user_id` با کلید خارجی به `users` گرفت.
- هر متد [`TaskStore`](../taskflow/internal/store/task_store.go) حالا یک `userID`
  می‌گیرد و بر اساسش فیلتر می‌کند: `WHERE id = ? AND user_id = ?`.
- هندلرها شناسه را از context می‌خوانند (که `Auth` تنظیم کرده) و پایین پاس می‌دهند.

نتیجه ایزولاسیون واقعی است، که با یک تست اثبات شده: وقتی **باب** وظیفه‌ی **آلیس**
را با شناسه می‌خواهد، یک `404` می‌گیرد — نه چون وجود ندارد، بلکه چون *مالِ او* نیست.
عبارت SQL با `AND user_id = ?` نشت داده‌ی کاربر دیگر را غیرممکن می‌کند.

### هندلرهای احراز هویت همه را به هم می‌بافند

[`api/auth.go`](../taskflow/internal/api/auth.go):
- **`POST /auth/register`** — ورودی را اعتبارسنجی، رمز را هش، کاربر را بساز، یک JWT
  برگردان. ایمیل تکراری ← `409 Conflict`.
- **`POST /auth/login`** — کاربر را پیدا کن، `CheckPassword`، یک JWT برگردان. ایمیل
  *یا* رمز اشتباه ← **همان** پیام عمومی `401` (تا فاش نکنی کدام ایمیل‌ها ثبت‌شده‌اند
  — یک بهترین‌روش امنیتی ظریف).

توجه کن پاسخ ورود `User` را در خود دارد، اما هش رمز هرگز ظاهر نمی‌شود چون تگ
`json:"-"` دارد (جلسه‌ی ۱۵). رازها سمت سرور می‌مانند.

---

## ۵. ببین همه‌چیز کار می‌کند (۵ دقیقه)

</div>

```bash
cd taskflow
go test ./...                       # همه‌ی تست‌های احراز هویت + محدودسازی پاس می‌شوند

go run .                            # سرور را راه بینداز، بعد در ترمینال دیگر:
curl localhost:8080/health
curl -X POST localhost:8080/auth/register -d '{"email":"me@example.com","password":"secret123"}'
# مقدار "token" را از پاسخ کپی کن، بعد:
TOKEN=...
curl -X POST localhost:8080/tasks -H "Authorization: Bearer $TOKEN" -d '{"title":"Ship it"}'
curl localhost:8080/tasks -H "Authorization: Bearer $TOKEN"
```

<div dir="rtl">

مجموعه‌ی تست ([`api/tasks_test.go`](../taskflow/internal/api/tasks_test.go)) کل
داستان را پوشش می‌دهد: health عمومی، ۴۰۱ بدون توکن، CRUD کامل با احراز هویت،
**ایزولاسیون به‌ازای-کاربر**، جریان ورود و رد ایمیل تکراری — همه با `httptest` و یک
دیتابیس واقعی (موقت).

---

## 🎯 تمرین‌ها (قبل از جلسه‌ی ۲۰ انجام بده!)

در پروژه‌ی `taskflow/` کار کن:

۱. **`GET /me`:** یک endpoint محافظت‌شده که اطلاعات کاربر فعلی (شناسه، ایمیل) را برمی‌گرداند اضافه کن. به یک متد `UserStore.GetByID` نیاز داری. هش را نشت نده.
۲. **اعتبارسنجی قوی‌تر:** ثبت‌نام را وقتی ایمیل `@` ندارد رد کن، و رمز ≥ ۸ کاراکتر لازم کن. پیام‌های `400` واضح برگردان.
۳. **میان‌افزار CORS:** یک میان‌افزار `CORS` که `Access-Control-Allow-Origin: *` را تنظیم و درخواست‌های preflight از نوع `OPTIONS` را مدیریت می‌کند اضافه کن.
۴. **تست انقضای توکن:** تستی بنویس که توکنی با انقضای گذشته تولید می‌کند و تأیید می‌کند API مقدار `401` برمی‌گرداند.
۵. **تست‌های واحد auth:** یک `password_test.go` و `jwt_test.go` در پکیج `auth` اضافه کن: رفت‌وبرگشت هش←بررسی، و رفت‌وبرگشت تولید←تجزیه به‌علاوه‌ی شکست توکن دستکاری‌شده.

---

## ✅ چک‌لیست جلسه‌ی ۱۹

- [ ] تنظیمات را از متغیرهای محیطی با پیش‌فرض‌های منطقی بارگذاری می‌کنم
- [ ] رمزها را با bcrypt هش می‌کنم و هرگز متن خام ذخیره نمی‌کنم
- [ ] می‌توانم JWT تولید و تأیید کنم، و توکن منقضی/دستکاری‌شده را رد کنم
- [ ] شکل میان‌افزار `func(next http.Handler) http.Handler` را می‌فهمم
- [ ] میان‌افزار لاگ، بازیابی و احراز هویت ساختم
- [ ] شناسه‌ی کاربر را با یک نوع کلید export‌نشده از `context` پاس می‌دهم
- [ ] داده‌ی وظایفم محدود به هر کاربر است (`WHERE user_id = ?`)
- [ ] ورود برای هر دو حالت ایمیل-اشتباه و رمز-اشتباه از خطای عمومی استفاده می‌کند
- [ ] هر ۵ تمرین را انجام دادم

**قبلی:** [→ جلسه‌ی ۱۸](session-18.fa.md) · **بعدی:** [جلسه‌ی ۲۰ — جلا و انتشار ←](session-20.fa.md)

</div>
