# مستندات معماری مایکروسرویس go-catalog

## 📋 نمای کلی

این سند، ساختار کامل و معماری مایکروسرویس **go-catalog** را تعریف می‌کند. این مایکروسرویس معادل ماژول کاتالوگ در پروژه روچی (`packages/rochi/catalog`) است و با زبان **Go** و فریمورک **Gin** توسعه داده می‌شود.

---

## 🏗️ معماری پیشنهادی

معماری لایه‌ای (Layered Architecture) با الگوی Repository:

```
┌─────────────────────────────────────────────────────────┐
│                    HTTP Layer (Gin)                     │
│  ├── Handlers (Controllers)                             │
│  ├── Middlewares                                        │
│  └── Routes                                             │
├─────────────────────────────────────────────────────────┤
│                    Service Layer                        │
│  └── Business Logic                                     │
├─────────────────────────────────────────────────────────┤
│                   Repository Layer                      │
│  └── Data Access                                        │
├─────────────────────────────────────────────────────────┤
│                    Domain Layer                         │
│  ├── Models (Entities)                                  │
│  └── DTOs                                               │
├─────────────────────────────────────────────────────────┤
│                   Infrastructure                        │
│  ├── Database (MySQL/PostgreSQL)                        │
│  ├── Elasticsearch                                      │
│  ├── Redis (Cache)                                      │
│  └── External Services                                  │
└─────────────────────────────────────────────────────────┘
```

### جریان درخواست HTTP

```
Request → Router → Middleware → Handler → Service → Repository → DB
                                                              ↓
Response ← Handler ← Service ← Repository ← ─────────────────┘
```

---

## 📁 ساختار پوشه‌ها و فایل‌ها

```
go-catalog/
├── cmd/
│   └── server/
│       └── main.go                    # نقطه ورود برنامه
│
├── internal/
│   ├── app/
│   │   └── app.go                     # راه‌اندازی و اجرای سرور
│   │
│   ├── config/
│   │   └── config.go                  # تنظیمات و environment variables
│   │
│   ├── domain/
│   │   ├── product.go                 # موجودیت محصول
│   │   ├── variant.go                 # موجودیت تنوع محصول
│   │   ├── category.go                # موجودیت دسته‌بندی
│   │   ├── attribute.go               # موجودیت ویژگی
│   │   ├── attribute_value.go         # موجودیت مقادیر ویژگی
│   │   ├── tag.go                     # موجودیت برچسب
│   │   ├── comment.go                 # موجودیت نظر
│   │   ├── media.go                   # موجودیت رسانه
│   │   ├── brand.go                   # موجودیت برند
│   │   ├── badge.go                   # موجودیت نشان
│   │   └── wow_item.go                # موجودیت آیتم ویژه (شگفت‌انگیز)
│   │
│   ├── dto/
│   │   ├── request/
│   │   │   ├── product_request.go     # DTOهای ورودی محصول
│   │   │   ├── category_request.go    # DTOهای ورودی دسته‌بندی
│   │   │   ├── comment_request.go     # DTOهای ورودی نظر
│   │   │   └── filter_request.go      # DTOهای فیلترینگ
│   │   │
│   │   └── response/
│   │       ├── product_response.go    # DTOهای خروجی محصول
│   │       ├── category_response.go   # DTOهای خروجی دسته‌بندی
│   │       ├── variant_response.go    # DTOهای خروجی تنوع
│   │       ├── comment_response.go    # DTOهای خروجی نظر
│   │       └── paginate_response.go   # DTOهای صفحه‌بندی
│   │
│   ├── repository/
│   │   ├── interfaces.go              # تعریف اینترفیس‌های ریپوزیتوری
│   │   ├── product_repo.go            # ریپوزیتوری محصول
│   │   ├── variant_repo.go            # ریپوزیتوری تنوع
│   │   ├── category_repo.go           # ریپوزیتوری دسته‌بندی
│   │   ├── attribute_repo.go          # ریپوزیتوری ویژگی
│   │   ├── tag_repo.go                # ریپوزیتوری برچسب
│   │   ├── comment_repo.go            # ریپوزیتوری نظر
│   │   └── elastic_repo.go            # ریپوزیتوری Elasticsearch
│   │
│   ├── service/
│   │   ├── product_service.go         # سرویس محصول
│   │   ├── category_service.go        # سرویس دسته‌بندی
│   │   ├── comment_service.go         # سرویس نظرات
│   │   ├── home_service.go            # سرویس صفحه اصلی
│   │   ├── filter_service.go          # سرویس فیلترینگ
│   │   └── search_service.go          # سرویس جستجو
│   │
│   ├── http/
│   │   ├── handlers/
│   │   │   ├── product_handler.go     # هندلر API محصول
│   │   │   ├── category_handler.go    # هندلر API دسته‌بندی
│   │   │   ├── comment_handler.go     # هندلر API نظرات
│   │   │   ├── tag_handler.go         # هندلر API برچسب
│   │   │   ├── home_handler.go        # هندلر API صفحه اصلی
│   │   │   └── health_handler.go      # هندلر بررسی سلامت
│   │   │
│   │   ├── middlewares/
│   │   │   ├── auth.go                # میدل‌ور احراز هویت
│   │   │   ├── cors.go                # میدل‌ور CORS
│   │   │   ├── logging.go             # میدل‌ور لاگینگ
│   │   │   ├── ratelimit.go           # میدل‌ور محدودیت نرخ
│   │   │   └── recovery.go            # میدل‌ور بازیابی خطا
│   │   │
│   │   ├── validators/
│   │   │   └── validators.go          # اعتبارسنجی ورودی‌ها
│   │   │
│   │   └── routes.go                  # تعریف روت‌ها
│   │
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── mysql.go               # اتصال MySQL
│   │   │   └── migrations.go          # مایگریشن‌ها
│   │   ├── elasticsearch/
│   │   │   └── client.go              # کلاینت Elasticsearch
│   │   ├── redis/
│   │   │   └── client.go              # کلاینت Redis
│   │   └── cache/
│   │       └── cache.go               # لایه کش
│   │
│   └── pkg/
│       ├── logger/
│       │   └── logger.go              # سیستم لاگینگ
│       ├── errors/
│       │   └── errors.go              # مدیریت خطاها
│       └── utils/
│           ├── slug.go                # تولید slug
│           ├── pagination.go          # صفحه‌بندی
│           └── response.go            # پاسخ‌های استاندارد
│
├── database/
│   └── migrations/
│       └── *.sql                      # فایل‌های مایگریشن SQL
│
├── docs/
│   └── openapi.yaml                   # مستندات API (Swagger/OpenAPI)
│
├── .env.example                       # نمونه environment variables
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## 📊 مدل‌های دامین (Domain Models)

### Product (محصول)

```go
type Product struct {
    ID                int64           `json:"id" db:"id"`
    Title             string          `json:"title" db:"title"`
    Slug              string          `json:"slug" db:"slug"`
    Status            int             `json:"status" db:"status"`
    CategoryID        *int64          `json:"category_id" db:"category_id"`
    Featured          bool            `json:"featured" db:"featured"`
    SKU               *string         `json:"sku" db:"sku"`
    Attributes        json.RawMessage `json:"attributes" db:"attributes"`
    DefaultAttributes json.RawMessage `json:"default_attributes" db:"default_attributes"`
    Variations        json.RawMessage `json:"variations" db:"variations"`
    RelatedProducts   json.RawMessage `json:"related_products" db:"related_products"`
    UpSellProducts    json.RawMessage `json:"up_sell_products" db:"up_sell_products"`
    CrossSellProducts json.RawMessage `json:"cross_sell_products" db:"cross_sell_products"`
    Price             *int64          `json:"price" db:"price"`
    SalePrice         *int64          `json:"sale_price" db:"sale_price"`
    OnSale            bool            `json:"on_sale" db:"on_sale"`
    InStock           bool            `json:"in_stock" db:"in_stock"`
    StockQuantity     float64         `json:"stock_quantity" db:"stock_quantity"`
    Image             *string         `json:"image" db:"image"`
    Weight            *float64        `json:"weight" db:"weight"`
    AverageRating     float64         `json:"average_rating" db:"average_rating"`
    RatingCount       int             `json:"rating_count" db:"rating_count"`
    Dimensions        json.RawMessage `json:"dimensions" db:"dimensions"`
    ProductType       *int            `json:"product_type" db:"product_type"`
    PurchaseType      *int            `json:"purchase_type" db:"purchase_type"`
    ShortDescription  *string         `json:"short_description" db:"short_description"`
    Description       *string         `json:"description" db:"description"`
    Features          json.RawMessage `json:"features" db:"features"`
    Options           json.RawMessage `json:"options" db:"options"`
    Noindex           bool            `json:"noindex" db:"noindex"`
    FreeShipping      *time.Time      `json:"free_shipping" db:"free_shipping"`
    OrderingID        int             `json:"ordering_id" db:"ordering_id"`
    PublishedAt       *time.Time      `json:"published_at" db:"published_at"`
    CreatedAt         time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt         time.Time       `json:"updated_at" db:"updated_at"`
    DeletedAt         *time.Time      `json:"deleted_at" db:"deleted_at"`
    
    // Relations
    Category    *Category   `json:"category,omitempty"`
    Categories  []Category  `json:"categories,omitempty"`
    Variants    []Variant   `json:"variants,omitempty"`
    Tags        []Tag       `json:"tags,omitempty"`
    Media       []Media     `json:"media,omitempty"`
    Comments    []Comment   `json:"comments,omitempty"`
}

// ثابت‌های وضعیت محصول
const (
    ProductStatusOutOfStock = 0
    ProductStatusPublished  = 1
    ProductStatusDraft      = 3
    ProductStatusDeleted    = 6
    ProductStatusDisabled   = 7
)

// ثابت‌های نوع محصول
const (
    ProductTypeVariable = 1
    ProductTypeSimple   = 2
)

// ثابت‌های نوع خرید
const (
    PurchaseTypeDownloadable = 1
    PurchaseTypeShippable    = 2
)
```

### Variant (تنوع محصول)

```go
type Variant struct {
    ID                int64           `json:"id" db:"id"`
    ProductID         *int64          `json:"product_id" db:"product_id"`
    Description       *string         `json:"description" db:"description"`
    SKU               *string         `json:"sku" db:"sku"`
    Attributes        json.RawMessage `json:"attributes" db:"attributes"`
    DefaultAttributes bool            `json:"default_attributes" db:"default_attributes"`
    Price             *int64          `json:"price" db:"price"`
    SalePrice         *int64          `json:"sale_price" db:"sale_price"`
    SinglePrice       *int64          `json:"single_price" db:"single_price"`
    OnSale            bool            `json:"on_sale" db:"on_sale"`
    PercentDiscount   *int            `json:"percent_discount" db:"percent_discount"`
    Weight            *float64        `json:"weight" db:"weight"`
    Dimensions        json.RawMessage `json:"dimensions" db:"dimensions"`
    PurchaseType      *int            `json:"purchase_type" db:"purchase_type"`
    StockQuantity     float64         `json:"stock_quantity" db:"stock_quantity"`
    Threshold         *float64        `json:"threshold" db:"threshold"`
    StockPurchasable  bool            `json:"stock_purchasable" db:"stock_purchasable"`
    PurchaseLimit     *int            `json:"purchase_limit" db:"purchase_limit"`
    Status            *int            `json:"status" db:"status"`
    Options           json.RawMessage `json:"options" db:"options"`
    SalePriceFromDate *time.Time      `json:"sale_price_from_date" db:"sale_price_from_date"`
    SalePriceToDate   *time.Time      `json:"sale_price_to_date" db:"sale_price_to_date"`
    ExtraName         *string         `json:"extra_name" db:"extra_name"`
    SupplierCode      *string         `json:"supplier_code" db:"supplier_code"`
    CreatedAt         time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt         time.Time       `json:"updated_at" db:"updated_at"`
    DeletedAt         *time.Time      `json:"deleted_at" db:"deleted_at"`
    
    // Relations
    Product          *Product          `json:"product,omitempty"`
    AttributesValues []AttributeValue  `json:"attributes_values,omitempty"`
}
```

### Category (دسته‌بندی)

```go
type Category struct {
    ID              int64           `json:"id" db:"id"`
    Title           string          `json:"title" db:"title"`
    Slug            string          `json:"slug" db:"slug"`
    Canonical       *string         `json:"canonical" db:"canonical"`
    Description     *string         `json:"description" db:"description"`
    MetaTitle       *string         `json:"meta_title" db:"meta_title"`
    MetaDescription *string         `json:"meta_description" db:"meta_description"`
    ParentID        *int64          `json:"parent_id" db:"parent_id"`
    Lft             int             `json:"_lft" db:"_lft"`               // Nested Set
    Rgt             int             `json:"_rgt" db:"_rgt"`               // Nested Set
    UnitID          *int64          `json:"unit_id" db:"unit_id"`
    AttributeID     *int64          `json:"attribute_id" db:"attribute_id"`
    StepUnit        *float64        `json:"step_unit" db:"step_unit"`
    Threshold       *float64        `json:"threshold" db:"threshold"`
    Commission      *int            `json:"commission" db:"commission"`
    Penalty         *int            `json:"penalty" db:"penalty"`
    Orderable       bool            `json:"orderable" db:"orderable"`
    Viewable        bool            `json:"viewable" db:"viewable"`
    Discountable    bool            `json:"discountable" db:"discountable"`
    ForceUnit       bool            `json:"force_unit" db:"force_unit"`
    Noindex         bool            `json:"noindex" db:"noindex"`
    CostPerView     *int            `json:"cost_per_view" db:"cost_per_view"`
    Suggest         bool            `json:"suggest" db:"suggest"`
    Meta            json.RawMessage `json:"meta" db:"meta"`
    DescendantsSelf json.RawMessage `json:"descendants_self" db:"descendants_self"`
    AttributesIDs   json.RawMessage `json:"attributes_ids" db:"attributes_ids"`
    ProductCount    int             `json:"product_count" db:"product_count"`
    CreatedAt       time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
    DeletedAt       *time.Time      `json:"deleted_at" db:"deleted_at"`
    
    // Relations
    Parent     *Category   `json:"parent,omitempty"`
    Children   []Category  `json:"children,omitempty"`
    Ancestors  []Category  `json:"ancestors,omitempty"`
    Products   []Product   `json:"products,omitempty"`
    Attributes []Attribute `json:"attributes,omitempty"`
}
```

### Attribute (ویژگی)

```go
type Attribute struct {
    ID          int64           `json:"id" db:"id"`
    Title       string          `json:"title" db:"title"`
    Slug        *string         `json:"slug" db:"slug"`
    Type        int             `json:"type" db:"type"`
    Featured    bool            `json:"featured" db:"featured"`
    Dedicated   bool            `json:"dedicated" db:"dedicated"`
    Status      *int            `json:"status" db:"status"`
    Placeholder *string         `json:"placeholder" db:"placeholder"`
    Meta        json.RawMessage `json:"meta" db:"meta"`
    CreatedAt   time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
    
    // Relations
    Values     []AttributeValue `json:"values,omitempty"`
    Categories []Category       `json:"categories,omitempty"`
}

// انواع ویژگی
const (
    AttributeTypeSelect    = 1  // سلکت باکس
    AttributeTypeRadio     = 2  // رادیو
    AttributeTypeInput     = 3  // ورودی
    AttributeTypeColor     = 4  // رنگ
    AttributeTypeCheckbox  = 5  // چک باکس
    AttributeTypeText      = 6  // تکست
)
```

### AttributeValue (مقدار ویژگی)

```go
type AttributeValue struct {
    ID          int64   `json:"id" db:"id"`
    AttributeID int64   `json:"attribute_id" db:"attribute_id"`
    Title       string  `json:"title" db:"title"`
    Slug        *string `json:"slug" db:"slug"`
    Color       *string `json:"color" db:"color"`
    Description *string `json:"description" db:"description"`
    Status      bool    `json:"status" db:"status"`
    ParentID    *int64  `json:"parent_id" db:"parent_id"`
    Lft         int     `json:"_lft" db:"_lft"`
    Rgt         int     `json:"_rgt" db:"_rgt"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
    
    // Relations
    Attribute *Attribute `json:"attribute,omitempty"`
}
```

### Tag (برچسب)

```go
type Tag struct {
    ID         int64     `json:"id" db:"id"`
    Title      string    `json:"title" db:"title"`
    Slug       string    `json:"slug" db:"slug"`
    Status     bool      `json:"status" db:"status"`
    IsBadge    bool      `json:"is_badge" db:"is_badge"`
    BadgeColor *string   `json:"badge_color" db:"badge_color"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
    UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
    
    // Relations
    Products []Product `json:"products,omitempty"`
}
```

### Comment (نظر)

```go
type Comment struct {
    ID              int64           `json:"id" db:"id"`
    Body            string          `json:"body" db:"body"`
    Reply           json.RawMessage `json:"reply" db:"reply"`
    ReplyBody       *string         `json:"reply_body" db:"reply_body"`
    ProductID       int64           `json:"product_id" db:"product_id"`
    UserID          *int64          `json:"user_id" db:"user_id"`
    Verified        bool            `json:"verified" db:"verified"`
    Rating          *int            `json:"rating" db:"rating"`
    Reviewer        *string         `json:"reviewer" db:"reviewer"`
    ReviewerEmail   *string         `json:"reviewer_email" db:"reviewer_email"`
    Status          string          `json:"status" db:"status"`
    Positive        json.RawMessage `json:"positive" db:"positive"`
    Negative        json.RawMessage `json:"negative" db:"negative"`
    Buyer           bool            `json:"buyer" db:"buyer"`
    ConsumerLike    int             `json:"consumer_like" db:"consumer_like"`
    ConsumerDislike int             `json:"consumer_dislike" db:"consumer_dislike"`
    CreatedAt       time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
    
    // Relations
    Product *Product `json:"product,omitempty"`
}

// وضعیت‌های نظر
const (
    CommentStatusApproved   = "approved"
    CommentStatusHold       = "hold"
    CommentStatusProcessing = "processing"
    CommentStatusSpam       = "spam"
    CommentStatusTrash      = "trash"
)
```

### Media (رسانه)

```go
type Media struct {
    ID         int64     `json:"id" db:"id"`
    MediableID int64     `json:"mediable_id" db:"mediable_id"`
    MediableType string  `json:"mediable_type" db:"mediable_type"`
    Src        string    `json:"src" db:"src"`
    Source     *string   `json:"source" db:"source"`
    AWSPath    *string   `json:"aws_path" db:"aws_path"`
    IsS3       bool      `json:"is_s3" db:"is_s3"`
    Priority   int       `json:"priority" db:"priority"`
    CreatedAt  time.Time `json:"created_at" db:"created_at"`
    UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
    DeletedAt  *time.Time `json:"deleted_at" db:"deleted_at"`
}
```

### WowItem (آیتم شگفت‌انگیز)

```go
type WowItem struct {
    ID            int64           `json:"id" db:"id"`
    VariantID     int64           `json:"variant_id" db:"variant_id"`
    Price         int64           `json:"price" db:"price"`
    Quantity      int             `json:"quantity" db:"quantity"`
    SoldQuantity  int             `json:"sold_quantity" db:"sold_quantity"`
    Status        int             `json:"status" db:"status"`
    StartAt       *time.Time      `json:"start_at" db:"start_at"`
    EndAt         *time.Time      `json:"end_at" db:"end_at"`
    Options       json.RawMessage `json:"options" db:"options"`
    GroupNumber   *int            `json:"group_number" db:"group_number"`
    MarkupPercent *float64        `json:"markup_percent" db:"markup_percent"`
    CreatedAt     time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
    
    // Relations
    Variant *Variant `json:"variant,omitempty"`
}
```

---

## 🔌 API Endpoints

### Public APIs

| Method | Endpoint | توضیحات | پارامترها |
|--------|----------|---------|-----------|
| `GET` | `/health` | بررسی سلامت سرویس | - |
| `GET` | `/api/home` | صفحه اصلی | - |
| `GET` | `/api/products` | لیست محصولات | `category_id`, `price`, `attribute`, `tags`, `on_sale`, `in_stock`, `seller_id`, `q`, `sort`, `page`, `per_page` |
| `GET` | `/api/products/:id` | جزئیات محصول | `id` |
| `GET` | `/api/product-cat/:id` | محصولات یک دسته‌بندی | `id` |
| `GET` | `/api/categories/:id` | جزئیات دسته‌بندی | `id` |
| `GET` | `/api/tags` | لیست برچسب‌ها | - |
| `GET` | `/api/filters` | پارامترهای فیلتر | - |
| `GET` | `/api/comments/:id` | نظرات محصول | `id`, `page`, `per_page` |
| `POST` | `/api/comments` | ثبت نظر جدید | (نیاز به احراز هویت) |

### API Version 2

| Method | Endpoint | توضیحات |
|--------|----------|---------|
| `GET` | `/api/v2/home` | صفحه اصلی نسخه ۲ |
| `GET` | `/api/v2/products` | لیست محصولات نسخه ۲ |
| `GET` | `/api/v2/products/:id` | جزئیات محصول نسخه ۲ |
| `GET` | `/api/v2/categories/:id` | دسته‌بندی نسخه ۲ |
| `GET` | `/api/v2/menus` | منوها |
| `GET` | `/api/v2/product_list` | لیست ساده محصولات |
| `GET` | `/api/v2/wishlist` | لیست علاقه‌مندی‌ها (نیاز به auth) |
| `POST` | `/api/v2/wishlist` | افزودن به علاقه‌مندی‌ها (نیاز به auth) |
| `DELETE` | `/api/v2/wishlist` | حذف از علاقه‌مندی‌ها (نیاز به auth) |

### Admin APIs

| Method | Endpoint | توضیحات |
|--------|----------|---------|
| `POST` | `/api/admin/products/:id/replicate` | کپی محصول |
| `GET` | `/api/admin/variant/create/:id` | ایجاد تنوع جدید |

### Internal APIs (N8N / Automation)

| Method | Endpoint | توضیحات |
|--------|----------|---------|
| `POST` | `/api/update-stock` | بروزرسانی موجودی |
| `GET` | `/api/n8n/products/getAttributes/:id` | دریافت ویژگی‌های محصول |
| `GET` | `/api/n8n/categories/:id` | جزئیات دسته‌بندی برای n8n |

---

## 📥 نمونه ورودی/خروجی API

### GET /api/products

**Request:**
```http
GET /api/products?category_id=10&price=100000,500000&in_stock=1&sort=sale_price,asc&page=1&per_page=16
```

**Response:**
```json
{
  "attributes": [
    {
      "id": 1,
      "title": "رنگ",
      "values": [
        {"id": 1, "title": "مشکی", "color": "#000000"},
        {"id": 2, "title": "سفید", "color": "#FFFFFF"}
      ]
    }
  ],
  "products": {
    "data": [
      {
        "id": 123,
        "name": "عنوان محصول",
        "slug": "product-slug",
        "url": "/product/product-slug",
        "price": 500000,
        "sale_price": 450000,
        "percent_discount": 10,
        "on_sale": true,
        "in_stock": true,
        "image": "https://cdn.example.com/image.jpg",
        "category_title": "دسته‌بندی",
        "seller": {
          "id": 1,
          "store_name": "فروشگاه"
        }
      }
    ],
    "meta": {
      "current_page": 1,
      "per_page": 16,
      "total": 100,
      "last_page": 7
    }
  }
}
```

### GET /api/products/:id

**Response:**
```json
{
  "product": {
    "id": 123,
    "name": "عنوان محصول",
    "slug": "product-slug",
    "url": "/product/product-slug",
    "type": 1,
    "price": 500000,
    "sale_price": 450000,
    "percent_discount": 10,
    "on_sale": true,
    "in_stock": true,
    "sku": "SKU-123",
    "short_description": "توضیحات کوتاه",
    "description": "توضیحات کامل محصول",
    "category_id": 10,
    "category": {
      "id": 10,
      "title": "دسته‌بندی",
      "slug": "category-slug"
    },
    "category_ancestors": [
      {"id": 1, "title": "دسته اصلی", "slug": "main-category"}
    ],
    "images": [
      {"id": 1, "src": "https://cdn.example.com/image1.jpg"},
      {"id": 2, "src": "https://cdn.example.com/image2.jpg"}
    ],
    "tags": [
      {"id": 1, "title": "برچسب", "slug": "tag-slug"}
    ],
    "variants": {
      "1": [1, 2],
      "2": [3, 4]
    },
    "variants_details": [
      {
        "id": 1,
        "name": "مشکی - سایز L",
        "price": 500000,
        "sale_price": 450000,
        "stock_quantity": 10,
        "is_purchasable": true
      }
    ],
    "short_attributes": {
      "رنگ": ["مشکی", "سفید"],
      "سایز": ["L", "XL"]
    },
    "full_attributes": {
      "رنگ": ["مشکی", "سفید"],
      "سایز": ["L", "XL"],
      "جنس": ["پنبه"]
    },
    "rating_count": 25,
    "similar_products": [],
    "cross_products": []
  }
}
```

### POST /api/comments

**Request:**
```json
{
  "product_id": 123,
  "body": "متن نظر",
  "rating": 5,
  "positive": ["کیفیت خوب", "ارسال سریع"],
  "negative": []
}
```

**Response:**
```json
{
  "success": true,
  "message": "نظر شما با موفقیت ثبت شد و پس از تایید نمایش داده خواهد شد.",
  "comment": {
    "id": 456,
    "body": "متن نظر",
    "status": "hold",
    "created_at": "2025-01-15T10:30:00Z"
  }
}
```

---

## 🔧 اینترفیس‌ها (Interfaces)

### ProductRepository

```go
type ProductRepository interface {
    FindByID(ctx context.Context, id int64) (*domain.Product, error)
    FindBySlug(ctx context.Context, slug string) (*domain.Product, error)
    FindAll(ctx context.Context, filter ProductFilter) ([]domain.Product, *Pagination, error)
    FindByCategory(ctx context.Context, categoryID int64, filter ProductFilter) ([]domain.Product, *Pagination, error)
    Create(ctx context.Context, product *domain.Product) error
    Update(ctx context.Context, product *domain.Product) error
    Delete(ctx context.Context, id int64) error
    UpdateStock(ctx context.Context, variantID int64, quantity float64) error
    GetSimilarProducts(ctx context.Context, product *domain.Product) ([]domain.Product, error)
}
```

### CategoryRepository

```go
type CategoryRepository interface {
    FindByID(ctx context.Context, id int64) (*domain.Category, error)
    FindBySlug(ctx context.Context, slug string) (*domain.Category, error)
    FindAll(ctx context.Context) ([]domain.Category, error)
    FindWithAncestors(ctx context.Context, id int64) (*domain.Category, error)
    FindDescendants(ctx context.Context, id int64) ([]domain.Category, error)
    FindChildren(ctx context.Context, parentID *int64) ([]domain.Category, error)
    GetProductCount(ctx context.Context, id int64) (int, error)
}
```

### VariantRepository

```go
type VariantRepository interface {
    FindByID(ctx context.Context, id int64) (*domain.Variant, error)
    FindByProductID(ctx context.Context, productID int64) ([]domain.Variant, error)
    Create(ctx context.Context, variant *domain.Variant) error
    Update(ctx context.Context, variant *domain.Variant) error
    Delete(ctx context.Context, id int64) error
    UpdateStock(ctx context.Context, id int64, quantity float64) error
}
```

### CommentRepository

```go
type CommentRepository interface {
    FindByID(ctx context.Context, id int64) (*domain.Comment, error)
    FindByProductID(ctx context.Context, productID int64, filter CommentFilter) ([]domain.Comment, *Pagination, error)
    Create(ctx context.Context, comment *domain.Comment) error
    Update(ctx context.Context, comment *domain.Comment) error
    Delete(ctx context.Context, id int64) error
    Reaction(ctx context.Context, id int64, isLike bool) error
}
```

### SearchRepository

```go
type SearchRepository interface {
    SearchProducts(ctx context.Context, query string, filter ProductFilter) ([]domain.Product, *Pagination, error)
    IndexProduct(ctx context.Context, product *domain.Product) error
    DeleteProduct(ctx context.Context, id int64) error
    BulkIndex(ctx context.Context, products []domain.Product) error
}
```

---

## 🔄 پارامترهای فیلترینگ

```go
type ProductFilter struct {
    IDs         []int64  `form:"ids"`
    CategoryID  *int64   `form:"category_id"`
    Categories  []int64  `form:"categories"`
    SellerID    *int64   `form:"seller_id"`
    TagIDs      []int64  `form:"tags"`
    Attributes  []int64  `form:"attribute"`
    PriceMin    *int64   `form:"min_price"`
    PriceMax    *int64   `form:"max_price"`
    OnSale      *bool    `form:"on_sale"`
    InStock     *bool    `form:"in_stock"`
    Query       string   `form:"q"`
    Sort        string   `form:"sort"`       // id, price, sale_price, created_at
    Direction   string   `form:"direction"`  // asc, desc
    Page        int      `form:"page"`
    PerPage     int      `form:"per_page"`
}
```

---

## 📄 صفحه‌بندی (Pagination)

```go
type Pagination struct {
    CurrentPage int   `json:"current_page"`
    PerPage     int   `json:"per_page"`
    Total       int64 `json:"total"`
    LastPage    int   `json:"last_page"`
}

type PaginatedResponse[T any] struct {
    Data []T         `json:"data"`
    Meta *Pagination `json:"meta"`
}
```

---

## ⚠️ مدیریت خطا (Error Handling)

```go
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

var (
    ErrNotFound          = &AppError{Code: 404, Message: "موردی یافت نشد"}
    ErrBadRequest        = &AppError{Code: 400, Message: "درخواست نامعتبر"}
    ErrUnauthorized      = &AppError{Code: 401, Message: "احراز هویت نشده"}
    ErrForbidden         = &AppError{Code: 403, Message: "دسترسی غیرمجاز"}
    ErrInternalServer    = &AppError{Code: 500, Message: "خطای داخلی سرور"}
    ErrValidation        = &AppError{Code: 422, Message: "خطای اعتبارسنجی"}
    ErrDatabaseError     = &AppError{Code: 500, Message: "خطای دیتابیس"}
)
```

---

## ✅ Validation

```go
type ProductRequest struct {
    Title            string  `json:"title" binding:"required,min=3,max=255"`
    Slug             string  `json:"slug" binding:"required,slug"`
    CategoryID       int64   `json:"category_id" binding:"required,gt=0"`
    Price            int64   `json:"price" binding:"omitempty,gt=0"`
    SalePrice        int64   `json:"sale_price" binding:"omitempty,gt=0,ltefield=Price"`
    ShortDescription string  `json:"short_description" binding:"max=500"`
    Description      string  `json:"description" binding:"max=10000"`
    Status           int     `json:"status" binding:"oneof=0 1 3 6 7"`
}

type CommentRequest struct {
    ProductID int64    `json:"product_id" binding:"required,gt=0"`
    Body      string   `json:"body" binding:"required,min=10,max=1000"`
    Rating    int      `json:"rating" binding:"required,min=1,max=5"`
    Positive  []string `json:"positive" binding:"max=5,dive,max=100"`
    Negative  []string `json:"negative" binding:"max=5,dive,max=100"`
}
```

---

## 🔗 ارتباط با سایر سرویس‌ها

### سرویس‌های وابسته

| سرویس | نوع ارتباط | توضیحات |
|-------|-----------|---------|
| **go-login** | REST API | احراز هویت کاربران و JWT validation |
| **Seller Service** | Database | اطلاعات فروشندگان |
| **Order Service** | Database | ارتباط با سفارشات و سبد خرید |
| **Warehouse Service** | REST/Event | مدیریت موجودی |
| **Elasticsearch** | Client | جستجوی محصولات |
| **Redis** | Client | کش و Rate Limiting |

### Event-Driven Communication (آینده)

```go
// رویدادهای منتشر شده
type ProductPublishedEvent struct {
    ProductID int64     `json:"product_id"`
    Timestamp time.Time `json:"timestamp"`
}

type StockUpdatedEvent struct {
    VariantID     int64   `json:"variant_id"`
    OldQuantity   float64 `json:"old_quantity"`
    NewQuantity   float64 `json:"new_quantity"`
    Timestamp     time.Time `json:"timestamp"`
}

// رویدادهای دریافتی
type OrderCompletedEvent struct {
    OrderID   int64 `json:"order_id"`
    VariantID int64 `json:"variant_id"`
    Quantity  int   `json:"quantity"`
}
```

---

## 📦 وابستگی‌ها و پکیج‌ها

```go
// go.mod
require (
    github.com/gin-gonic/gin v1.11.0
    github.com/go-sql-driver/mysql v1.8.1
    github.com/jmoiron/sqlx v1.4.0
    github.com/redis/go-redis/v9 v9.7.0
    github.com/elastic/go-elasticsearch/v8 v8.16.0
    github.com/golang-jwt/jwt/v5 v5.2.1
    github.com/go-playground/validator/v10 v10.22.0
    github.com/spf13/viper v1.19.0
    go.uber.org/zap v1.27.0
    github.com/swaggo/swag v1.16.0
    github.com/swaggo/gin-swagger v1.6.0
)
```

---

## 🪵 لاگینگ (Logging)

```go
// سطوح لاگ
const (
    LogLevelDebug = "debug"
    LogLevelInfo  = "info"
    LogLevelWarn  = "warn"
    LogLevelError = "error"
)

// فرمت لاگ
type LogEntry struct {
    Level     string    `json:"level"`
    Timestamp time.Time `json:"timestamp"`
    Message   string    `json:"message"`
    RequestID string    `json:"request_id,omitempty"`
    UserID    *int64    `json:"user_id,omitempty"`
    Method    string    `json:"method,omitempty"`
    Path      string    `json:"path,omitempty"`
    Status    int       `json:"status,omitempty"`
    Duration  float64   `json:"duration_ms,omitempty"`
    Error     string    `json:"error,omitempty"`
}
```

---

## ⚙️ تنظیمات (Configuration)

```go
type Config struct {
    // Server
    ServerAddr string `env:"SERVER_ADDR" default:":8080"`
    AppEnv     string `env:"APP_ENV" default:"production"`
    
    // Database
    DBHost     string `env:"DB_HOST" default:"localhost"`
    DBPort     string `env:"DB_PORT" default:"3306"`
    DBUser     string `env:"DB_USER" default:"root"`
    DBPassword string `env:"DB_PASSWORD"`
    DBName     string `env:"DB_NAME" default:"rochi"`
    
    // Redis
    RedisAddr     string `env:"REDIS_ADDR"`
    RedisPassword string `env:"REDIS_PASSWORD"`
    RedisDB       int    `env:"REDIS_DB" default:"0"`
    
    // Elasticsearch
    ElasticURL      string `env:"ELASTIC_URL"`
    ElasticUsername string `env:"ELASTIC_USERNAME"`
    ElasticPassword string `env:"ELASTIC_PASSWORD"`
    
    // JWT
    JWTSecret string `env:"JWT_SECRET"`
    
    // Logging
    LogLevel string `env:"LOG_LEVEL" default:"info"`
    LogFile  string `env:"LOG_FILE" default:"logs/app.log"`
    
    // Cache
    CacheTTL time.Duration `env:"CACHE_TTL" default:"5m"`
}
```

### نمونه .env

```env
# Server
SERVER_ADDR=:8080
APP_ENV=production

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=secret
DB_NAME=rochi

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Elasticsearch
ELASTIC_URL=http://localhost:9200
ELASTIC_USERNAME=elastic
ELASTIC_PASSWORD=changeme

# JWT
JWT_SECRET=your-super-secret-key

# Logging
LOG_LEVEL=info
LOG_FILE=logs/app.log

# Cache
CACHE_TTL=5m
```

---

## 🧪 تست (Testing)

### ساختار تست

```
tests/
├── unit/
│   ├── service/
│   │   ├── product_service_test.go
│   │   └── category_service_test.go
│   └── repository/
│       └── product_repo_test.go
├── integration/
│   ├── api/
│   │   ├── product_api_test.go
│   │   └── category_api_test.go
│   └── database/
│       └── migrations_test.go
└── fixtures/
    ├── products.json
    └── categories.json
```

### دستورات تست

```bash
# اجرای تمام تست‌ها
make test

# اجرای تست‌های unit
make test-unit

# اجرای تست‌های integration
make test-integration

# گزارش coverage
make test-coverage
```

---

## 🚀 دستورات Makefile

```makefile
.PHONY: build run test migrate docker

build:
	go build -o bin/server cmd/server/main.go

run:
	go run cmd/server/main.go

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

docker-build:
	docker build -t go-catalog .

docker-run:
	docker-compose up -d

lint:
	golangci-lint run

swagger:
	swag init -g cmd/server/main.go -o docs
```

---

## 📋 جدول مایگریشن‌ها

جداول اصلی مورد نیاز:

| جدول | توضیحات |
|------|---------|
| `products` | محصولات |
| `variants` | تنوع‌های محصول |
| `categories` | دسته‌بندی‌ها (Nested Set) |
| `attributes` | ویژگی‌ها |
| `attribute_values` | مقادیر ویژگی |
| `attribute_category` | رابطه ویژگی-دسته‌بندی |
| `category_product` | رابطه محصول-دسته‌بندی |
| `tags` | برچسب‌ها |
| `product_tag` | رابطه محصول-برچسب |
| `comments` | نظرات |
| `media` | رسانه‌ها |
| `wow_items` | آیتم‌های شگفت‌انگیز |
| `badges` | نشان‌ها |
| `product_badges` | رابطه محصول-نشان |

---

## 📌 نکات مهم توسعه

1. **Soft Delete**: تمام موجودیت‌های اصلی از soft delete پشتیبانی می‌کنند
2. **Nested Set**: دسته‌بندی‌ها از الگوریتم Nested Set برای درختی بودن استفاده می‌کنند
3. **JSON Fields**: فیلدهای `attributes`, `options`, `meta` به صورت JSON ذخیره می‌شوند
4. **Image CDN**: تصاویر از CDN سرو می‌شوند
5. **Elasticsearch**: برای جستجوی سریع محصولات از Elasticsearch استفاده می‌شود
6. **Cache Strategy**: استفاده از Redis برای کش کردن لیست محصولات و دسته‌بندی‌ها
7. **Rate Limiting**: محدودیت نرخ درخواست برای API عمومی

---

## 🔜 مراحل بعدی

1. [ ] پیاده‌سازی لایه Domain و DTOs
2. [ ] پیاده‌سازی Repository Layer با MySQL
3. [ ] پیاده‌سازی Service Layer
4. [ ] پیاده‌سازی HTTP Handlers
5. [ ] افزودن Elasticsearch برای جستجو
6. [ ] افزودن Redis برای کش
7. [ ] پیاده‌سازی احراز هویت JWT
8. [ ] نوشتن تست‌های Unit و Integration
9. [ ] مستندات OpenAPI/Swagger
10. [ ] Dockerize کردن پروژه

---

**تاریخ ایجاد:** ۱۴۰۴/۰۹/۱۱  
**نسخه:** 1.0.0  
**نویسنده:** تیم توسعه روچی

