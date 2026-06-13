<div dir="rtl">

> 🌐 **زبان / Language:** فارسی (همین فایل) · [English](session-18.md)

# جلسه‌ی ۱۸ — REST API + پایگاه‌داده 🗄️

**هدف (۱ ساعت):** ساخت **TaskFlow**، پروژه‌ی نمونه‌کارت، را شروع کنی. یک چیدمان
درست پروژه‌ی Go راه می‌اندازی، به یک **پایگاه‌داده‌ی SQLite** واقعی وصل می‌شوی و
endpointهای کامل **CRUD** (ساخت، خواندن، به‌روزرسانی، حذف) را با معماری تمیز و
لایه‌ای می‌سازی. در پایان یک REST API کارکننده و پشتیبانی‌شده با دیتابیس داری.

> **مرور جلسه‌ی ۱۷:** می‌توانی هندلرهای JSON بسازی و تست کنی. حالا مپ حافظه‌ای را
> با یک دیتابیس واقعی عوض و کد را حرفه‌ای سازمان می‌دهیم.

> 📁 **پروژه در [`taskflow/`](../taskflow/)** در ریشه‌ی مخزن قرار دارد، به‌عنوان
> ماژول Go مستقل خودش — پس می‌توانی آن پوشه را مستقیم به گیت‌هاب به‌عنوان یک
> پروژه‌ی مستقل برای رزومه کپی کنی.

---

## ۱. چیدمان پروژه و قرارداد `internal/` (۱۰ دقیقه)

پروژه‌های واقعی Go دغدغه‌ها را به پکیج‌ها جدا می‌کنند. چیدمان TaskFlow:

</div>

```
taskflow/
├── go.mod                       ماژول مستقل خودش: «taskflow»
├── main.go                      نقطه‌ی ورود: باز کردن DB، اتصال لایه‌ها، سرویس
└── internal/
    ├── models/   task.go        انواع دامنه
    ├── store/    store.go       اتصال DB + مهاجرت‌ها
    │             task_store.go  همه‌ی SQL برای taskها («مخزن/repository»)
    └── api/      server.go      روتر، کمک‌کننده‌های JSON، اینترفیس مخزن
                  tasks.go       هندلرهای HTTP
                  tasks_test.go  تست‌ها
```

<div dir="rtl">

> 🔑 **`internal/` در Go ویژه است:** پکیج‌های زیر `internal/` فقط توسط کد *همان
> ماژول* قابل import‌اند. این روش زبان‌محورِ خصوصی نگه‌داشتن پیاده‌سازی توست.

**معماری لایه‌ای** (ایده‌ی کلیدی):

</div>

```
درخواست HTTP → api (هندلرها) → store (مخزن) → پایگاه‌داده
```

<div dir="rtl">

هر لایه فقط با لایه‌ی زیرش حرف می‌زند. هندلرها SQL نمی‌نویسند؛ store از HTTP خبر
ندارد. این جداسازی هر تکه را ساده، قابل‌تعویض و قابل‌تست می‌کند.

---

## ۲. اتصال به دیتابیس با `database/sql` (۱۵ دقیقه)

پکیج استاندارد `database/sql` یک رابط عمومی به دیتابیس‌های SQL است. آن را با یک
**درایور** برای دیتابیس مشخصت جفت می‌کنی. ما از درایور SQLite کاملاً-Go استفاده
می‌کنیم (به کامپایلر C نیازی نیست):

</div>

```go
import (
    "database/sql"
    _ "modernc.org/sqlite"   // import خالی: درایور «sqlite» را ثبت می‌کند
)

db, err := sql.Open("sqlite", "taskflow.db")
err = db.Ping()              // تأیید کن اتصال واقعاً کار می‌کند
```

<div dir="rtl">

> 🔑 **import خالی `_ "modernc.org/sqlite"`** تابع `init()` درایور را اجرا می‌کند
> تا با `database/sql` ثبت شود، اما خودِ پکیج را مستقیم ارجاع نمی‌دهی. این الگوی
> ثبت-درایور برای SQL در Go استاندارد است.

درایور را یک بار به ماژولت اضافه کن:

</div>

```bash
cd taskflow
go get modernc.org/sqlite
```

<div dir="rtl">

### مهاجرت‌ها — ساخت اسکیما

هنگام راه‌اندازی مطمئن می‌شویم جدول وجود دارد:

</div>

```go
const schema = `
CREATE TABLE IF NOT EXISTS tasks (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    title      TEXT    NOT NULL,
    done       INTEGER NOT NULL DEFAULT 0,
    created_at TEXT    NOT NULL
);`
db.Exec(schema)
```

<div dir="rtl">

فایل [`taskflow/internal/store/store.go`](../taskflow/internal/store/store.go) را ببین.

---

## ۳. مخزن: همه‌ی SQL در یک جا (۲۰ دقیقه)

**الگوی مخزن (repository)** هر کوئری دیتابیس برای یک موجودیت را در یک نوع می‌گذارد.
بقیه‌ی اپ متدهایش را صدا می‌زند و هرگز SQL نمی‌بیند.

</div>

```go
type TaskStore struct { db *sql.DB }

func (s *TaskStore) Create(title string) (models.Task, error) {
    res, err := s.db.Exec(
        `INSERT INTO tasks (title, done, created_at) VALUES (?, 0, ?)`,
        title, time.Now().UTC().Format(time.RFC3339),
    )
    id, _ := res.LastInsertId()
    return models.Task{ID: id, Title: title, CreatedAt: ...}, nil
}
```

<div dir="rtl">

دو نکته‌ی ضروری `database/sql`:

- **`db.Exec`** — برای `INSERT`/`UPDATE`/`DELETE` (ردیفی برنمی‌گردد).
  `LastInsertId()` و `RowsAffected()` می‌دهد.
- **`db.Query`** (چند ردیف) و **`db.QueryRow`** (یک ردیف) — برای `SELECT`. ستون‌ها
  را با `rows.Scan(&a, &b, ...)` می‌خوانی.

> ⚠️ **همیشه از جای‌نگه‌دار `?` استفاده کن** برای مقادیر — *هرگز* SQL را با الحاق
> رشته نساز. جای‌نگه‌دارها از **تزریق SQL (SQL injection)** جلوگیری می‌کنند، حفره‌ی
> امنیتی کلاسیک. `Exec("... WHERE id = ?", id)` امن است؛
> `Exec("... WHERE id = " + id)` خطرناک. این عادتی غیرقابل‌مذاکره است.

وقتی ردیفی پیدا نشود، `QueryRow(...).Scan(...)` مقدار `sql.ErrNoRows` برمی‌گرداند.
ما آن را به sentinel خودمان `ErrNotFound` (جلسه‌ی ۱۱) ترجمه می‌کنیم تا لایه‌ی API
بتواند به `404` نگاشتش کند:

</div>

```go
if errors.Is(err, sql.ErrNoRows) {
    return models.Task{}, ErrNotFound
}
```

<div dir="rtl">

مخزن کامل را مطالعه کن: [`taskflow/internal/store/task_store.go`](../taskflow/internal/store/task_store.go).

---

## ۴. اتصال هندلرها به مخزن از طریق یک اینترفیس (۱۵ دقیقه)

لایه‌ی API رفتاری را که نیاز دارد به‌صورت **اینترفیس** تعریف می‌کند، بعد به آن
وابسته می‌شود — نه به `TaskStore` مشخص. این وارونگی وابستگی است، و دلیل قابل‌تست
بودن هندلرهاست (جلسه‌ی ۱۰ در عمل):

</div>

```go
// در پکیج api
type TaskRepository interface {
    Create(title string) (models.Task, error)
    List() ([]models.Task, error)
    Get(id int64) (models.Task, error)
    Update(id int64, title string, done bool) (models.Task, error)
    Delete(id int64) error
}

type Server struct { tasks TaskRepository }
```

<div dir="rtl">

`*store.TaskStore` این اینترفیس را خودکار برآورده می‌کند (برآورده‌سازی ضمنی!).
هندلرها نتایج و خطاهای مخزن را به HTTP نگاشت می‌کنند:

</div>

```go
task, err := s.tasks.Get(id)
if errors.Is(err, store.ErrNotFound) {
    writeError(w, http.StatusNotFound, "task not found")
    return
}
```

<div dir="rtl">

فایل `main.go` سه لایه را به هم وصل می‌کند:

</div>

```go
db, _ := store.Open("taskflow.db")
taskStore := store.NewTaskStore(db)        // لایه‌ی store
server := api.NewServer(taskStore)         // لایه‌ی api (مخزن را می‌گیرد)
http.ListenAndServe(":8080", server.Routes())
```

<div dir="rtl">

### اجرا و تست

</div>

```bash
cd taskflow
go run .                       # سرور واقعی + DB واقعی
go test ./...                  # تست‌های یکپارچگی روی DB موقت
```

<div dir="rtl">

تست‌ها از یک **DB واقعی SQLite در `t.TempDir()`** استفاده می‌کنند — اثبات می‌کند
کل پشته (هندلر → SQL → دیتابیس) کار می‌کند، بعد خودکار پاک‌سازی می‌شود. فایل
[`taskflow/internal/api/tasks_test.go`](../taskflow/internal/api/tasks_test.go) را ببین.

> 💡 **چرا DB واقعی در تست، نه تقلبی؟** برای پروژه‌ی کوچک، تمرین SQL واقعی باگ
> بیشتری می‌گیرد و با پشتیبانی temp/in-memory در SQLite ساده است. اینترفیس
> `TaskRepository` همچنان می‌گذارد اگر لازم شد یک تقلبی جا بزنی.

---

## 🎯 تمرین‌ها (قبل از جلسه‌ی ۱۹ انجام بده!)

داخل `taskflow/` کار کن:

۱. **اجرای کل جریان:** سرور را راه بینداز و با `curl` یک task را بساز، فهرست، به‌روز، دریافت و حذف کن. ظاهرشدن `taskflow.db` را تماشا کن.
۲. **افزودن فیلد:** یک ستون `Priority int` اضافه کن (مهاجرت + مدل + SQL در Create/Update/Scan). تأیید کن از طریق API رفت‌وبرگشت می‌کند.
۳. **endpoint فیلتر:** `GET /tasks?done=true` اضافه کن که فقط taskهای تکمیل‌شده را فهرست کند (پارامتر کوئری را بخوان، متد کوئری با `WHERE done = ?` اضافه کن).
۴. **اعتبارسنجی:** عنوان‌های بلندتر از ۲۰۰ کاراکتر را با `400` رد کن.
۵. **تست بیشتر:** تستی برای حالت به‌روزرسانی-پیدانشد (`PUT /tasks/999` → 404) و برای endpoint فیلتر جدیدت اضافه کن.

---

## ✅ چک‌لیست جلسه‌ی ۱۸

- [ ] چیدمان لایه‌ای (api → store → db) و `internal/` را می‌فهمم
- [ ] می‌توانم DB SQLite را با `database/sql` + import خالی درایور باز کنم
- [ ] هنگام راه‌اندازی یک مهاجرت برای ساخت اسکیما اجرا می‌کنم
- [ ] می‌توانم `Exec`، `Query` و `QueryRow` را با جای‌نگه‌دار `?` بنویسم
- [ ] می‌دانم چرا جای‌نگه‌دارها از تزریق SQL جلوگیری می‌کنند
- [ ] `sql.ErrNoRows` را به `ErrNotFound` خودم ترجمه می‌کنم
- [ ] در لایه‌ی API به اینترفیس `TaskRepository` وابسته‌ام
- [ ] می‌توانم سرور و تست‌های یکپارچگی را اجرا کنم
- [ ] هر ۵ تمرین را انجام دادم

**قبلی:** [→ جلسه‌ی ۱۷](session-17.fa.md) · **بعدی:** [جلسه‌ی ۱۹ — احراز هویت، میان‌افزار و تنظیمات ←](session-19.fa.md)

</div>
