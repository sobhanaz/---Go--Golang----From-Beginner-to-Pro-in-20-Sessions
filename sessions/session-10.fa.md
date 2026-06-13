<div dir="rtl">

> 🌐 **زبان / Language:** فارسی (همین فایل) · [English](session-10.md)

# جلسه‌ی ۱۰ — اینترفیس‌ها 🔌

**هدف (۱ ساعت):** ویژگی شاخص Go را مسلط شوی. یک **اینترفیس** توصیف می‌کند *چیزی چه
کاری می‌تواند انجام دهد* (رفتارش) بدون این که بگوید *چیست*. اینترفیس‌ها کد
انعطاف‌پذیر و کم‌وابسته می‌دهند — و ترفند Go (ضمنی بودن) آن‌ها را دلپذیر می‌کند. این
جلسه مبتدی را به توسعه‌دهنده‌ی مطمئن Go تبدیل می‌کند.

> **مرور جلسه‌ی ۰۹:** استراکت + متد + اشاره‌گر برای مدل کردن داده و وصل رفتار.
> اینترفیس‌ها می‌گذارند نوع‌های مختلف بر اساس رفتار مشترک به‌جای هم استفاده شوند.

---

## ۱. اینترفیس چیست؟ (۱۰ دقیقه)

یک **اینترفیس** مجموعه‌ای از امضای متدهاست. هر نوعی که آن متدها را داشته باشد
*خودکار* اینترفیس را برآورده می‌کند — **به کلیدواژه‌ی `implements` نیازی نیست.**

</div>

```go
type Shape interface {
    Area() float64
    Perimeter() float64
}
```

<div dir="rtl">

این می‌گوید: «یک `Shape` هرچیزی است که متد `Area() float64` **و** متد
`Perimeter() float64` داشته باشد.» همین. هر نوعی با هر دو متد، خودکار یک `Shape` است.

</div>

```go
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

// Circle حالا Shape را برآورده می‌کند — هرگز اعلام نکردیم. فقط می‌کند.
var s Shape = Circle{Radius: 5}
fmt.Println(s.Area())
```

<div dir="rtl">

> 🔑 **برآورده‌سازی ضمنی ابرقدرت Go است.** در جاوا/سی‌شارپ باید `class Circle
> implements Shape` بنویسی. در Go اگر متدها بخوانند، فقط کار می‌کند. یعنی می‌توانی
> نوع‌های *موجود* (حتی از پکیج‌های دیگر) را به اینترفیس‌های *خودت* برآورنده کنی.

---

## ۲. چندریختی: یک تابع، چند نوع (۲۰ دقیقه)

چون هر نوع منطبق یک `Shape` است، می‌توانی توابعی بنویسی که روی *اینترفیس* کار کنند
و هر نوع مشخصی که برآورده‌اش می‌کند را بپذیرند:

</div>

```go
func describe(s Shape) {
    fmt.Printf("area=%.2f perimeter=%.2f\n", s.Area(), s.Perimeter())
}

describe(Circle{Radius: 5})
describe(Rectangle{Width: 3, Height: 4})
```

<div dir="rtl">

همان `describe` برای دایره، مستطیل، مثلث — هرچیزی که «یک Shape است» کار می‌کند.
حتی می‌توانی نوع‌های مختلط را در یک اسلایس نگه داری:

</div>

```go
shapes := []Shape{
    Circle{Radius: 2},
    Rectangle{Width: 3, Height: 4},
}
for _, s := range shapes {
    describe(s)
}
```

<div dir="rtl">

این **چندریختی** است — نوشتن کد در برابر رفتار، نه نوع‌های مشخص. کلید Go
انعطاف‌پذیر و قابل‌تست: توابعت به اینترفیس‌های کوچک وابسته‌اند و می‌توانی هر
پیاده‌سازی (واقعی، یا تقلبی برای تست) را جا بزنی.

فایل‌های [`examples/session10/basics/basics.go`](../examples/session10/basics/basics.go) و
[`examples/session10/polymorphism/polymorphism.go`](../examples/session10/polymorphism/polymorphism.go) را اجرا کن.

---

## ۳. اینترفیس `Stringer` — تا نوع‌هایت قشنگ چاپ شوند (۱۰ دقیقه)

کتابخانه‌ی استاندارد اینترفیس‌های کوچکی تعریف می‌کند که همیشه برآورده‌شان می‌کنی.
رایج‌ترین `fmt.Stringer` است:

</div>

```go
type Stringer interface {
    String() string
}
```

<div dir="rtl">

اگر نوعت یک متد `String() string` داشته باشد، `fmt` هنگام چاپ خودکار از آن استفاده
می‌کند — دقیقاً همان که با enum `Status` در جلسه‌ی ۰۳ دیدی:

</div>

```go
type Color struct{ R, G, B int }

func (c Color) String() string {
    return fmt.Sprintf("rgb(%d, %d, %d)", c.R, c.G, c.B)
}

fmt.Println(Color{255, 0, 0})   // rgb(255, 0, 0)  — از String() استفاده می‌کند!
```

<div dir="rtl">

> 💡 **اصل طراحی:** Go **اینترفیس‌های کوچک** را ترجیح می‌دهد. معروف‌ترین‌ها،
> `io.Writer` و `io.Reader`، هرکدام *یک* متد دارند. «هرچه اینترفیس بزرگ‌تر،
> انتزاع ضعیف‌تر» — یک ضرب‌المثل Go.

---

## ۴. اینترفیس خالی و type assertion (۱۵ دقیقه)

اینترفیس خالی `interface{}` (یا نام مستعارش **`any`** از Go 1.18) *هیچ* متدی ندارد
— پس **هر** نوعی برآورده‌اش می‌کند. یعنی «هر مقداری»:

</div>

```go
var x any        // می‌تواند هرچیزی نگه دارد
x = 42
x = "hello"
x = []int{1, 2, 3}
```

<div dir="rtl">

`any` را در توابعی می‌بینی که مقادیر دلخواه می‌پذیرند (مثل `fmt.Println`). اما وقتی
مقداری به‌صورت `any` ذخیره شد، نوع مشخصش را از دست داده‌ای. برای پس گرفتنش از یک
**type assertion** استفاده کن:

</div>

```go
var x any = "hello"

s := x.(string)       // ادعا کن x یک string است -> "hello"

if s, ok := x.(string); ok {   // فرم امن comma-ok (بدون panic در صورت اشتباه)
    fmt.Println("it's a string:", s)
}
```

<div dir="rtl">

### type switch — مدیریت چند نوع ممکن

</div>

```go
func describe(i any) {
    switch v := i.(type) {
    case int:
        fmt.Println("int:", v)
    case string:
        fmt.Println("string of length", len(v))
    case bool:
        fmt.Println("bool:", v)
    default:
        fmt.Printf("unknown type %T\n", v)
    }
}
```

<div dir="rtl">

> ⚠️ **`any` را کم استفاده کن.** دست بردن به `any` اغلب یعنی داری امنیت نوعِ Go را
> دور می‌اندازی. نوع‌های واقعی و اینترفیس‌های کوچک را ترجیح بده؛ `any` را فقط وقتی
> واقعاً باید «هرچیزی» بپذیری استفاده کن (مثل چاپگرهای عمومی یا JSON).

فایل [`examples/session10/empty/empty.go`](../examples/session10/empty/empty.go) را اجرا کن.

---

## 🎯 تمرین‌ها (قبل از جلسه‌ی ۱۱ انجام بده!)

فایل `examples/session10/practice/practice.go` را بساز:

۱. **اینترفیس Shape:** `Shape` با `Area() float64` تعریف کن. برای `Circle` و `Square` پیاده کن. `func totalArea(shapes []Shape) float64` بنویس که مساحت اسلایس مختلط را جمع کند.
۲. **Stringer:** `Money struct { Amount float64; Currency string }` با متد `String()` که مثل `"$19.99"` چاپ کند. یک `Money` چاپ کن و تأیید کن `fmt` از متدت استفاده می‌کند.
۳. **Speaker:** `type Speaker interface { Speak() string }`. برای `Dog` و `Cat` پیاده کن. تابعی که `[]Speaker` می‌گیرد و گفته‌ی هرکدام را چاپ می‌کند بنویس.
۴. **type switch:** `func printType(values ...any)` بنویس که با type switch نوع و مقدار هر آرگومان را گزارش کند. با یک int، string، float و bool در یک فراخوانی صدایش بزن.
۵. **Mock برای تست (پیش‌نمایش):** `type Notifier interface { Send(msg string) error }`. یک `EmailNotifier` واقعی و یک `MockNotifier` که فقط پیام‌ها را ثبت می‌کند بساز. تابعی که `Notifier` می‌گیرد بنویس — ببین چطور می‌توانی بدون ارسال ایمیل واقعی تستش کنی. *(دقیقاً همین‌طور اینترفیس‌ها تست را در پروژه‌ی نهایی ممکن می‌کنند!)*

---

## ✅ چک‌لیست جلسه‌ی ۱۰

- [ ] می‌توانم اینترفیس را به‌صورت مجموعه‌ای از امضای متد تعریف کنم
- [ ] می‌فهمم برآورده‌سازی ضمنی است — بدون کلیدواژه‌ی `implements`
- [ ] می‌توانم تابعی بنویسم که اینترفیس بپذیرد و روی چند نوع کار کند
- [ ] می‌توانم نوع‌های مختلف را در یک `[]Interface` نگه دارم
- [ ] می‌توانم `String()` را پیاده کنم تا نحوه‌ی چاپ نوع را کنترل کنم
- [ ] می‌دانم `any` = `interface{}` و چطور با assertion نوع را پس بگیرم
- [ ] می‌توانم type switch بنویسم
- [ ] می‌دانم اینترفیس کوچک را ترجیح دهم و `any` را کم استفاده کنم
- [ ] هر ۵ تمرین را انجام دادم

**قبلی:** [→ جلسه‌ی ۰۹](session-09.fa.md) · **بعدی:** [جلسه‌ی ۱۱ — خطاها ←](session-11.fa.md)

</div>
