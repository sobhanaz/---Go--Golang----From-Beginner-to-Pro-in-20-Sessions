<div dir="rtl">

> 🌐 **زبان / Language:** فارسی (همین فایل) · [English](session-11.md)

# جلسه‌ی ۱۱ — خطاها ⚠️

**هدف (۱ ساعت):** یاد بگیری Go چطور خرابی‌ها را مدیریت می‌کند. Go برای خرابی‌های
عادی **استثنا (exception) ندارد** — به‌جایش، خطاها **مقادیر** معمولی‌اند که
برمی‌گردانی و بررسی می‌کنی. الگوهای اصطلاحی را مسلط می‌شوی، خطای سفارشی می‌سازی، با
`errors.Is`/`errors.As` می‌پیچی و بازرسی می‌کنی، و یاد می‌گیری کِی `panic`/
`recover` مناسب‌اند. این چیزی است که کد آماتور را از Go حرفه‌ای جدا می‌کند.

> **مرور جلسه‌ی ۱۰:** اینترفیس‌ها رفتار را توصیف می‌کنند. در واقع، `error` *فقط* یک
> اینترفیس است — به همین خاطر مدیریت خطا انعطاف‌پذیر است.

---

## ۱. `error` فقط یک اینترفیس است (۱۰ دقیقه)

نوع داخلی `error` یک اینترفیس تک‌متدی است:

</div>

```go
type error interface {
    Error() string
}
```

<div dir="rtl">

هرچیزی با متد `Error() string` یک خطاست. تابعی که می‌تواند شکست بخورد یک `error`
به‌عنوان **آخرین** مقدار بازگشتی می‌دهد. `nil` یعنی «بدون خطا» (موفقیت):

</div>

```go
result, err := doThing()
if err != nil {
    // مدیریت خرابی
    return
}
// از result استفاده کن — امن، چون err برابر nil بود
```

<div dir="rtl">

> 🔑 **الگوی طلایی، دوباره:** `if err != nil { ... }`. خطاها را *فوراً* بررسی کن،
> مدیریتشان کن، و فقط وقتی `err == nil` شد از نتیجه استفاده کن.

### ساخت خطاهای ساده

</div>

```go
errors.New("something went wrong")          // یک خطای پایه
fmt.Errorf("user %d not found", id)          // یک خطای قالب‌بندی‌شده
```

<div dir="rtl">

`fmt.Errorf` همانی است که بیشتر استفاده می‌کنی — پیام خطا را با مقادیر جاسازی‌شده می‌سازد.

فایل [`examples/session11/basics/basics.go`](../examples/session11/basics/basics.go) را اجرا کن.

---

## ۲. مدیریت خوب خطا (۱۰ دقیقه)

چند اصطلاح که Go حرفه‌ای را مشخص می‌کنند:

**زود برگرد.** خطا را مدیریت و `return` کن، تا مسیر خوش (happy path) بدون تورفتگی بماند:

</div>

```go
f, err := os.Open("file.txt")
if err != nil {
    return err          // همین حالا بیرون بزن
}
defer f.Close()
// ... ادامه با f، با علم به این که معتبر است
```

<div dir="rtl">

**هنگام بالا رفتن خطا زمینه (context) اضافه کن.** با `%w` بپیچ تا فراخواننده‌ها
بفهمند *کجا* شکست خورد، بدون از دست رفتن اصلی:

</div>

```go
func loadConfig() error {
    data, err := os.ReadFile("config.json")
    if err != nil {
        return fmt.Errorf("loadConfig: %w", err)   // %w اصلی را می‌پیچد
    }
    ...
}
```

<div dir="rtl">

حالا پیام مثل یک رد می‌شود: `loadConfig: open config.json: no such file`.

**خطاها را نادیده نگیر.** `_ = doThing()` یک خرابی را دور می‌اندازد. فقط وقتی خطا
را نادیده بگیر که واقعاً اهمیت ندهی (و با کامنت بگو).

---

## ۳. نوع‌های خطای سفارشی و خطاهای sentinel (۲۰ دقیقه)

### خطاهای sentinel — مقادیر خطای از پیش‌تعریف‌شده و قابل‌مقایسه

وقتی فراخواننده‌ها باید یک خطای *مشخص* را بررسی کنند، یک بار به‌عنوان متغیر پکیج
تعریفش کن (قرارداد: نام با `Err` شروع می‌شود):

</div>

```go
var ErrNotFound = errors.New("not found")

func findUser(id int) (*User, error) {
    if id == 0 {
        return nil, ErrNotFound
    }
    ...
}

// فراخواننده با errors.Is آن خطای دقیق را بررسی می‌کند:
user, err := findUser(0)
if errors.Is(err, ErrNotFound) {
    fmt.Println("no such user")
}
```

<div dir="rtl">

> 🔑 **از `errors.Is` استفاده کن، نه `==`.** چون خطاها با `%w` *پیچیده* می‌شوند،
> `errors.Is` کل زنجیره‌ی پیچش را برای یافتن انطباق می‌پیماید. `==` ساده یک خطای
> پیچیده را از دست می‌دهد.

### استراکت‌های خطای سفارشی — خطاهایی که داده حمل می‌کنند

وقتی خطا باید اطلاعات اضافی حمل کند، یک استراکت با متد `Error()` بساز:

</div>

```go
type ValidationError struct {
    Field string
    Msg   string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Msg)
}
```

<div dir="rtl">

برای بیرون کشیدن خطای نوع‌دار (تا فیلدهایش را بخوانی) از `errors.As` استفاده کن:

</div>

```go
err := validate(-5)
var ve *ValidationError
if errors.As(err, &ve) {
    fmt.Println("bad field:", ve.Field)   // "age"
}
```

<div dir="rtl">

> 🔑 **`errors.Is` در مقابل `errors.As`:**
> - `errors.Is(err, target)` — «آیا این (هرجای زنجیره) *همان مقدار خطای مشخص*
>   است؟» برای sentinelها.
> - `errors.As(err, &target)` — «آیا خطایی از *این نوع* در زنجیره هست؟ اگر بله،
>   بیرونش بکش.» برای استراکت‌های خطای سفارشی که فیلدشان را لازم داری.

فایل [`examples/session11/custom/custom.go`](../examples/session11/custom/custom.go) را اجرا کن.

---

## ۴. `panic` و `recover` — خروج اضطراری (۱۵ دقیقه)

یک **panic** برای مشکلات *غیرقابل‌بازیابی* و سطح-برنامه‌نویس است — نه برای خطاهای
عادی. وقتی panic رخ می‌دهد، برنامه پشته را باز و کرش می‌کند (مگر بازیابی شود).

</div>

```go
panic("this should never happen")
```

<div dir="rtl">

چیزهایی که خودکار panic می‌کنند: اندیس اسلایس خارج از محدوده، نوشتن روی مپ nil،
دی‌رفرنس اشاره‌گر nil، تقسیم صحیح بر صفر.

> ⚠️ **برای خطاهای عادی panic نکن** مثل «فایل پیدا نشد» یا «ورودی نامعتبر». به‌جایش
> یک `error` برگردان. panic برای موقعیت‌های «دنیا خراب شده» یا حالت‌های واقعاً
> غیرممکن است.

### `recover` — گرفتن یک panic

`recover` (فقط درون یک تابع `defer`شده مفید است) جلوی کرش کردن برنامه با panic را
می‌گیرد. کاربرد مشروع اصلی، در یک مرز است که نمی‌خواهی کل پروسه را پایین بیاورد —
مثلاً یک وب‌سرور که از panic در یک هندلر درخواست بازیابی می‌کند تا *سرور* روشن بماند:

</div>

```go
func safeRun(task func()) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("recovered from panic:", r)
        }
    }()
    task()   // حتی اگر این panic کند، بازیابی و ادامه می‌دهیم
}
```

<div dir="rtl">

دقیقاً همین الگو را به‌عنوان **میان‌افزار بازیابی** در پروژه‌ی نهایی استفاده می‌کنی.

فایل [`examples/session11/panic/panic.go`](../examples/session11/panic/panic.go) را اجرا کن.

> 💡 **قاعده‌ی سرانگشتی:** ۹۹٪ مواقع، یک `error` برگردان. panic را برای حالت‌های
> غیرممکن و recover را برای محافظت از یک مرز (مثل سرور) نگه دار.

---

## 🎯 تمرین‌ها (قبل از جلسه‌ی ۱۲ انجام بده!)

فایل `examples/session11/practice/practice.go` را بساز:

۱. **تقسیم با خطا:** `func divide(a, b float64) (float64, error)` که وقتی `b == 0` خطا برگرداند. هر دو حالت را با `if err != nil` مدیریت کن.
۲. **پیچیدن با زمینه:** `func parseAge(s string) (int, error)` که از `strconv.Atoi` استفاده و در صورت شکست `fmt.Errorf("parseAge: %w", err)` برگرداند. پیام پیچیده را چاپ کن.
۳. **خطای sentinel:** `var ErrEmptyName = errors.New("name is empty")`. `func greet(name string) (string, error)` که برای `""` آن را برگرداند. در فراخواننده با `errors.Is` تشخیصش بده.
۴. **نوع خطای سفارشی:** `RangeError struct { Value, Min, Max int }` با متد `Error()`. `func check(n int) error` که خارج از محدوده آن را برگرداند، و با `errors.As` فیلدهای `.Min`/`.Max` را در فراخواننده بخوان.
۵. **recover:** `func safeDivideInts(a, b int) (result int, err error)` که `a / b` را انجام می‌دهد اما با `defer`+`recover` panic تقسیم بر صفر را به خطای بازگشتی تبدیل می‌کند نه کرش.

---

## ✅ چک‌لیست جلسه‌ی ۱۱

- [ ] می‌دانم `error` یک اینترفیس تک‌متدی است و `nil` یعنی موفقیت
- [ ] از `errors.New` و `fmt.Errorf` برای ساخت خطا استفاده می‌کنم
- [ ] خطاها را فوراً بررسی و زود برمی‌گردم
- [ ] خطاها را با `%w` برای افزودن زمینه می‌پیچم
- [ ] خطای sentinel تعریف و با `errors.Is` تشخیصش می‌دهم
- [ ] استراکت خطای سفارشی می‌سازم و با `errors.As` بیرونش می‌کشم
- [ ] می‌دانم panic برای موقعیت‌های غیرقابل‌بازیابی است، نه خطای عادی
- [ ] می‌توانم از یک panic درون تابع defer‌شده بازیابی کنم
- [ ] هر ۵ تمرین را انجام دادم

**قبلی:** [→ جلسه‌ی ۱۰](session-10.fa.md) · **بعدی:** [جلسه‌ی ۱۲ — همروندی ۱ ←](session-12.fa.md)

</div>
