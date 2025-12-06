# مستندات پیشنهادی ساختار انبار و مدیریت موجودی

## 📋 نمای کلی

این سند، ساختار پیشنهادی برای مدیریت **انبار (Inventory/Storage)** در کنار مایکروسرویس کاتالوگ را تعریف می‌کند. هدف از این سند، ارائه یک معماری جامع برای مدیریت موجودی کالا، رزرو، تراکنش‌های انبار و جلوگیری از Oversell است.

> ⚠️ **توجه**: تمام موارد این سند پیشنهادی بوده و قبل از پیاده‌سازی نیاز به بررسی و تایید دارد.

---

## 🎯 تعریف نقش انبار (Inventory/Storage)

### مسئولیت‌های اصلی ماژول انبار:

| مسئولیت | توضیحات |
|---------|---------|
| **مدیریت موجودی** | نگهداری و بروزرسانی موجودی هر تنوع محصول (Variant) |
| **رزرو موجودی** | رزرو موقت موجودی هنگام ثبت سفارش |
| **تراکنش‌های انبار** | ثبت تمام ورودی‌ها و خروجی‌های انبار |
| **جلوگیری از Oversell** | اطمینان از فروش تنها محصولات موجود |
| **هشدار آستانه موجودی** | اعلام کمبود موجودی قبل از اتمام |
| **گزارش‌گیری** | تحلیل و گزارش وضعیت انبار |

### رابطه انبار با سایر ماژول‌ها:

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Catalog Service                              │
│  ┌─────────┐   ┌─────────┐   ┌─────────┐                            │
│  │ Product │──▶│ Variant │◀──│ Category│                            │
│  └─────────┘   └────┬────┘   └─────────┘                            │
│                     │                                                │
│                     │ stock_quantity                                 │
│                     ▼                                                │
├─────────────────────────────────────────────────────────────────────┤
│                       Inventory Service                              │
│  ┌───────────────┐   ┌─────────────────┐   ┌──────────────┐         │
│  │  Warehouse    │   │ StockTransaction │   │   Reserve    │         │
│  └───────────────┘   └─────────────────┘   └──────────────┘         │
│          │                    │                    │                 │
│          └────────────────────┼────────────────────┘                 │
│                               ▼                                      │
│                    ┌─────────────────┐                               │
│                    │   Inventory     │                               │
│                    │  (موجودی لحظه‌ای) │                               │
│                    └─────────────────┘                               │
├─────────────────────────────────────────────────────────────────────┤
│                         Order Service                                │
│  ┌─────────┐   ┌───────────┐   ┌─────────┐                          │
│  │  Order  │──▶│ OrderItem │◀──│  Cart   │                          │
│  └─────────┘   └───────────┘   └─────────┘                          │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 🔄 ارتباط موجودی با سایر موجودیت‌ها

### ارتباط با Variant (تنوع محصول)

هر تنوع محصول (Variant) دارای فیلد `stock_quantity` است که موجودی لحظه‌ای آن را نشان می‌دهد:

```
Variant
├── id
├── product_id
├── sku
├── price
├── sale_price
├── stock_quantity  ◀── موجودی قابل فروش
├── threshold       ◀── آستانه هشدار کمبود
└── stock_purchasable
```

### ارتباط با قیمت

```
┌─────────────┐      ┌─────────────────┐
│   Variant   │      │   StockTrans    │
├─────────────┤      ├─────────────────┤
│ price       │◀────▶│ purchase_price  │ (قیمت خرید)
│ sale_price  │      │ balance         │ (موجودی)
│ stock_qty   │◀─────│ amount          │ (تغییر)
└─────────────┘      └─────────────────┘
```

### ارتباط با فروش (Order)

```
Order Created
     │
     ▼
┌─────────────────┐
│  Reserve Stock  │ ◀── رزرو موجودی برای سفارش
└────────┬────────┘
         │
         ▼
   Order Completed?
    /          \
   Yes          No
   /              \
  ▼                ▼
┌──────────┐   ┌──────────────┐
│ Withdraw │   │ Release      │
│ Stock    │   │ Reserve      │
└──────────┘   └──────────────┘
```

---

## 📊 مدل‌های دامین انبار

### Warehouse (انبار)

```go
// Warehouse نمایانگر یک انبار فیزیکی یا مجازی است
type Warehouse struct {
    ID          int64     `json:"id" db:"id"`
    Name        string    `json:"name" db:"name"`
    Code        string    `json:"code" db:"code"`
    Address     *string   `json:"address" db:"address"`
    IsDefault   bool      `json:"is_default" db:"is_default"`
    IsActive    bool      `json:"is_active" db:"is_active"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
    DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at"`
}
```

### Unit (واحد اندازه‌گیری)

```go
// Unit واحد اندازه‌گیری موجودی
type Unit struct {
    ID        int64           `json:"id" db:"id"`
    Name      string          `json:"name" db:"name"`      // عدد، کیلوگرم، متر
    Symbol    string          `json:"symbol" db:"symbol"`  // pcs, kg, m
    Type      int             `json:"type" db:"type"`      // 1=شمارشی, 2=وزنی, 3=طولی
    StepUnit  *float64        `json:"step_unit" db:"step_unit"`   // گام افزایش
    MinValue  *float64        `json:"min" db:"min"`
    SubUnit   json.RawMessage `json:"sub_unit" db:"sub_unit"`     // واحد فرعی (گرم برای کیلو)
    CreatedAt time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

// انواع واحد
const (
    UnitTypeCountable  = 1  // شمارشی (عدد)
    UnitTypeWeightable = 2  // وزنی (کیلوگرم)
    UnitTypeMeasurable = 3  // طولی (متر)
)
```

### Inventory (موجودی لحظه‌ای)

```go
// Inventory موجودی لحظه‌ای هر تنوع در هر انبار
type Inventory struct {
    ID              int64     `json:"id" db:"id"`
    WarehouseID     int64     `json:"warehouse_id" db:"warehouse_id"`
    VariantID       int64     `json:"variant_id" db:"variant_id"`
    UnitID          *int64    `json:"unit_id" db:"unit_id"`
    
    // موجودی‌ها
    Quantity        float64   `json:"quantity" db:"quantity"`           // موجودی کل
    ReservedQty     float64   `json:"reserved_qty" db:"reserved_qty"`   // موجودی رزرو شده
    AvailableQty    float64   `json:"available_qty" db:"available_qty"` // موجودی قابل فروش
    
    // آستانه‌ها
    MinThreshold    *float64  `json:"min_threshold" db:"min_threshold"` // حداقل موجودی
    MaxThreshold    *float64  `json:"max_threshold" db:"max_threshold"` // حداکثر موجودی
    ReorderPoint    *float64  `json:"reorder_point" db:"reorder_point"` // نقطه سفارش مجدد
    
    // قیمت‌ها
    AverageCost     *int64    `json:"average_cost" db:"average_cost"`   // میانگین قیمت خرید
    LastCost        *int64    `json:"last_cost" db:"last_cost"`         // آخرین قیمت خرید
    
    LastUpdated     time.Time `json:"last_updated" db:"last_updated"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
    
    // Relations
    Warehouse *Warehouse `json:"warehouse,omitempty"`
    Variant   *Variant   `json:"variant,omitempty"`
    Unit      *Unit      `json:"unit,omitempty"`
}

// محاسبه موجودی قابل فروش
func (i *Inventory) CalculateAvailable() float64 {
    return i.Quantity - i.ReservedQty
}

// بررسی نیاز به سفارش مجدد
func (i *Inventory) NeedsReorder() bool {
    if i.ReorderPoint == nil {
        return false
    }
    return i.AvailableQty <= *i.ReorderPoint
}
```

### StockTransaction (تراکنش انبار)

```go
// StockTransaction ثبت هر تغییر در موجودی
type StockTransaction struct {
    ID             int64           `json:"id" db:"id"`
    WarehouseID    int64           `json:"warehouse_id" db:"warehouse_id"`
    VariantID      int64           `json:"variant_id" db:"variant_id"`
    UnitID         *int64          `json:"unit_id" db:"unit_id"`
    
    // نوع و مقادیر
    Type           TransactionType `json:"type" db:"type"`
    Amount         float64         `json:"amount" db:"amount"`           // مقدار تغییر (+/-)
    Balance        float64         `json:"balance" db:"balance"`         // موجودی بعد از تراکنش
    ReservedAmount float64         `json:"reserved_amount" db:"reserved_amount"` // مقدار رزرو
    
    // قیمت‌ها
    UnitPrice      *int64          `json:"unit_price" db:"unit_price"`   // قیمت واحد
    TotalPrice     *int64          `json:"total_price" db:"total_price"` // قیمت کل
    
    // ارجاعات
    OrderItemID    *int64          `json:"order_item_id" db:"order_item_id"`
    ReferenceType  *string         `json:"reference_type" db:"reference_type"` // order, purchase, adjustment
    ReferenceID    *int64          `json:"reference_id" db:"reference_id"`
    
    // اطلاعات اپراتور
    OperatorID     int64           `json:"operator_id" db:"operator_id"`
    VendorID       *int64          `json:"vendor_id" db:"vendor_id"`     // فروشنده/تامین‌کننده
    
    // سایر
    Date           time.Time       `json:"date" db:"date"`
    Note           *string         `json:"note" db:"note"`
    Status         int             `json:"status" db:"status"`
    
    CreatedAt      time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
    DeletedAt      *time.Time      `json:"deleted_at" db:"deleted_at"`
    
    // Relations
    Warehouse *Warehouse `json:"warehouse,omitempty"`
    Variant   *Variant   `json:"variant,omitempty"`
    Operator  *User      `json:"operator,omitempty"`
}

// انواع تراکنش
type TransactionType string

const (
    TransactionTypeAdd        TransactionType = "add"        // ورود به انبار
    TransactionTypeSale       TransactionType = "sale"       // فروش
    TransactionTypeReserve    TransactionType = "reserve"    // رزرو
    TransactionTypeRelease    TransactionType = "release"    // آزادسازی رزرو
    TransactionTypeRefund     TransactionType = "refund"     // برگشت از فروش
    TransactionTypeAdjustment TransactionType = "adjustment" // تعدیل موجودی
    TransactionTypeTransfer   TransactionType = "transfer"   // انتقال بین انبار
    TransactionTypeDamage     TransactionType = "damage"     // ضایعات
    TransactionTypeExpiry     TransactionType = "expiry"     // انقضا
)

// وضعیت تراکنش
const (
    TransactionStatusPending   = 0
    TransactionStatusCompleted = 1
    TransactionStatusCancelled = 2
)
```

### StockReservation (رزرو موجودی)

```go
// StockReservation رزرو موقت موجودی برای سفارش
type StockReservation struct {
    ID          int64     `json:"id" db:"id"`
    WarehouseID int64     `json:"warehouse_id" db:"warehouse_id"`
    VariantID   int64     `json:"variant_id" db:"variant_id"`
    OrderID     int64     `json:"order_id" db:"order_id"`
    OrderItemID int64     `json:"order_item_id" db:"order_item_id"`
    
    Quantity    float64   `json:"quantity" db:"quantity"`
    Status      int       `json:"status" db:"status"`
    
    ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`   // زمان انقضای رزرو
    ReleasedAt  *time.Time `json:"released_at" db:"released_at"`
    ConfirmedAt *time.Time `json:"confirmed_at" db:"confirmed_at"`
    
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// وضعیت رزرو
const (
    ReservationStatusPending   = 0  // در انتظار
    ReservationStatusConfirmed = 1  // تایید شده
    ReservationStatusReleased  = 2  // آزاد شده
    ReservationStatusExpired   = 3  // منقضی شده
)
```

### InventoryLog (لاگ تغییرات)

```go
// InventoryLog لاگ کامل تغییرات موجودی
type InventoryLog struct {
    ID            int64           `json:"id" db:"id"`
    InventoryID   int64           `json:"inventory_id" db:"inventory_id"`
    VariantID     int64           `json:"variant_id" db:"variant_id"`
    
    Action        string          `json:"action" db:"action"`        // create, update, delete
    Field         string          `json:"field" db:"field"`          // فیلد تغییر یافته
    OldValue      json.RawMessage `json:"old_value" db:"old_value"`
    NewValue      json.RawMessage `json:"new_value" db:"new_value"`
    
    UserID        int64           `json:"user_id" db:"user_id"`
    IPAddress     *string         `json:"ip_address" db:"ip_address"`
    UserAgent     *string         `json:"user_agent" db:"user_agent"`
    
    CreatedAt     time.Time       `json:"created_at" db:"created_at"`
}
```

### QuantityAction (عملیات تغییر موجودی گروهی)

```go
// QuantityAction عملیات تغییر موجودی دسته‌ای
type QuantityAction struct {
    ID            int64           `json:"id" db:"id"`
    Request       json.RawMessage `json:"request" db:"request"`       // شرایط انتخاب محصولات
    Operation     string          `json:"operation" db:"operation"`   // add, subtract, set
    Amount        float64         `json:"amount" db:"amount"`
    TotalVariants int             `json:"total_variants" db:"total_variants"`
    Status        int             `json:"status" db:"status"`
    CreatedBy     int64           `json:"created_by" db:"created_by"`
    CreatedAt     time.Time       `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
    DeletedAt     *time.Time      `json:"deleted_at" db:"deleted_at"`
    
    // Relations
    Histories []QuantityHistory `json:"histories,omitempty"`
}

// عملیات‌های مجاز
const (
    QuantityOperationAdd      = "add"
    QuantityOperationSubtract = "subtract"
    QuantityOperationSet      = "set"
)
```

### QuantityHistory (تاریخچه تغییرات موجودی)

```go
// QuantityHistory تاریخچه تغییرات موجودی یک واریانت
type QuantityHistory struct {
    ID               int64     `json:"id" db:"id"`
    QuantityActionID int64     `json:"quantity_action_id" db:"quantity_action_id"`
    VariantID        int64     `json:"variant_id" db:"variant_id"`
    ProductID        int64     `json:"product_id" db:"product_id"`
    
    OldQuantity      float64   `json:"old_quantity" db:"old_quantity"`
    NewQuantity      float64   `json:"new_quantity" db:"new_quantity"`
    ChangeAmount     float64   `json:"change_amount" db:"change_amount"`
    
    CreatedAt        time.Time `json:"created_at" db:"created_at"`
}
```

---

## 📁 ساختار پوشه‌ها و فایل‌ها

### پیشنهاد ۱: ماژول انبار در کنار کاتالوگ (Monorepo)

```
go-catalog/
├── internal/
│   ├── domain/
│   │   ├── product.go
│   │   ├── variant.go
│   │   ├── category.go
│   │   │
│   │   └── inventory/              # ◀── ماژول انبار
│   │       ├── warehouse.go
│   │       ├── unit.go
│   │       ├── inventory.go
│   │       ├── stock_transaction.go
│   │       ├── stock_reservation.go
│   │       ├── inventory_log.go
│   │       ├── quantity_action.go
│   │       └── quantity_history.go
│   │
│   ├── repository/
│   │   ├── product_repo.go
│   │   ├── variant_repo.go
│   │   │
│   │   └── inventory/
│   │       ├── interfaces.go
│   │       ├── warehouse_repo.go
│   │       ├── inventory_repo.go
│   │       ├── transaction_repo.go
│   │       └── reservation_repo.go
│   │
│   ├── service/
│   │   ├── product_service.go
│   │   │
│   │   └── inventory/
│   │       ├── inventory_service.go
│   │       ├── stock_service.go
│   │       ├── reservation_service.go
│   │       └── report_service.go
│   │
│   ├── http/
│   │   ├── handlers/
│   │   │   ├── product_handler.go
│   │   │   │
│   │   │   └── inventory/
│   │   │       ├── warehouse_handler.go
│   │   │       ├── inventory_handler.go
│   │   │       ├── transaction_handler.go
│   │   │       └── report_handler.go
│   │   │
│   │   └── routes.go
│   │
│   └── events/
│       ├── stock_events.go
│       └── handlers/
│           ├── order_created_handler.go
│           ├── order_completed_handler.go
│           └── order_cancelled_handler.go
```

### پیشنهاد ۲: مایکروسرویس مستقل (go-inventory)

```
go-inventory/
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── app/
│   │   └── app.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── domain/
│   │   ├── warehouse.go
│   │   ├── unit.go
│   │   ├── inventory.go
│   │   ├── stock_transaction.go
│   │   ├── stock_reservation.go
│   │   └── quantity_action.go
│   │
│   ├── repository/
│   │   ├── interfaces.go
│   │   ├── warehouse_repo.go
│   │   ├── inventory_repo.go
│   │   ├── transaction_repo.go
│   │   └── reservation_repo.go
│   │
│   ├── service/
│   │   ├── inventory_service.go
│   │   ├── stock_service.go
│   │   ├── reservation_service.go
│   │   └── sync_service.go         # ◀── همگام‌سازی با کاتالوگ
│   │
│   ├── http/
│   │   ├── handlers/
│   │   │   ├── warehouse_handler.go
│   │   │   ├── inventory_handler.go
│   │   │   ├── transaction_handler.go
│   │   │   └── health_handler.go
│   │   │
│   │   ├── middlewares/
│   │   │   └── auth.go
│   │   │
│   │   └── routes.go
│   │
│   ├── grpc/                       # ◀── ارتباط RPC با سایر سرویس‌ها
│   │   ├── server.go
│   │   └── handlers/
│   │       └── inventory_handler.go
│   │
│   ├── events/                     # ◀── Event-Driven
│   │   ├── publisher.go
│   │   ├── consumer.go
│   │   └── handlers/
│   │       ├── order_events.go
│   │       └── catalog_events.go
│   │
│   └── pkg/
│       ├── logger/
│       ├── errors/
│       └── utils/
│
├── proto/
│   └── inventory.proto             # ◀── تعریف gRPC
│
├── database/
│   └── migrations/
│
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

---

## 🔌 API های پایه‌ای مدیریت موجودی

### Inventory APIs

| Method | Endpoint | توضیحات |
|--------|----------|---------|
| `GET` | `/api/inventory` | لیست موجودی‌ها |
| `GET` | `/api/inventory/:variant_id` | موجودی یک تنوع |
| `PUT` | `/api/inventory/:variant_id` | بروزرسانی موجودی |
| `POST` | `/api/inventory/adjust` | تعدیل موجودی |
| `POST` | `/api/inventory/bulk-update` | بروزرسانی گروهی |

### Stock Transaction APIs

| Method | Endpoint | توضیحات |
|--------|----------|---------|
| `GET` | `/api/transactions` | لیست تراکنش‌ها |
| `GET` | `/api/transactions/:id` | جزئیات تراکنش |
| `POST` | `/api/transactions/add` | ورود به انبار |
| `POST` | `/api/transactions/withdraw` | خروج از انبار |
| `POST` | `/api/transactions/transfer` | انتقال بین انبار |

### Reservation APIs

| Method | Endpoint | توضیحات |
|--------|----------|---------|
| `POST` | `/api/reservations/reserve` | رزرو موجودی |
| `POST` | `/api/reservations/release` | آزادسازی رزرو |
| `POST` | `/api/reservations/confirm` | تایید و کسر از موجودی |
| `GET` | `/api/reservations/order/:order_id` | رزروهای یک سفارش |

### Warehouse APIs

| Method | Endpoint | توضیحات |
|--------|----------|---------|
| `GET` | `/api/warehouses` | لیست انبارها |
| `POST` | `/api/warehouses` | ایجاد انبار جدید |
| `PUT` | `/api/warehouses/:id` | ویرایش انبار |
| `DELETE` | `/api/warehouses/:id` | حذف انبار |

### Report APIs

| Method | Endpoint | توضیحات |
|--------|----------|---------|
| `GET` | `/api/reports/stock-levels` | گزارش سطح موجودی |
| `GET` | `/api/reports/low-stock` | کالاهای کم‌موجود |
| `GET` | `/api/reports/movements` | گزارش گردش انبار |
| `GET` | `/api/reports/valuation` | ارزش‌گذاری انبار |

---

## 📥 نمونه Request/Response

### POST /api/transactions/add

**Request:**
```json
{
  "warehouse_id": 1,
  "variant_id": 123,
  "amount": 50,
  "unit_price": 150000,
  "vendor_id": 5,
  "note": "خرید از تامین‌کننده الف",
  "date": "2025-01-15"
}
```

**Response:**
```json
{
  "success": true,
  "transaction": {
    "id": 456,
    "warehouse_id": 1,
    "variant_id": 123,
    "type": "add",
    "amount": 50,
    "balance": 150,
    "unit_price": 150000,
    "total_price": 7500000,
    "created_at": "2025-01-15T10:30:00Z"
  },
  "inventory": {
    "variant_id": 123,
    "quantity": 150,
    "available_qty": 145,
    "reserved_qty": 5
  }
}
```

### POST /api/reservations/reserve

**Request:**
```json
{
  "order_id": 789,
  "items": [
    {
      "variant_id": 123,
      "quantity": 2
    },
    {
      "variant_id": 124,
      "quantity": 1
    }
  ],
  "expires_in_minutes": 30
}
```

**Response:**
```json
{
  "success": true,
  "reservations": [
    {
      "id": 1001,
      "variant_id": 123,
      "quantity": 2,
      "status": "pending",
      "expires_at": "2025-01-15T11:00:00Z"
    },
    {
      "id": 1002,
      "variant_id": 124,
      "quantity": 1,
      "status": "pending",
      "expires_at": "2025-01-15T11:00:00Z"
    }
  ]
}
```

---

## 🔄 جریان داده

### ۱. ورود به انبار (Stock In)

```
┌─────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  API Call   │────▶│ StockService    │────▶│ TransactionRepo │
│  /add       │     │ .AddStock()     │     │ .Create()       │
└─────────────┘     └────────┬────────┘     └────────┬────────┘
                             │                       │
                             ▼                       ▼
                    ┌─────────────────┐     ┌─────────────────┐
                    │ InventoryRepo   │     │ StockTransaction│
                    │ .UpdateQty()    │     │ (Record)        │
                    └────────┬────────┘     └─────────────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Publish Event   │
                    │ StockAdded      │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Catalog Service │
                    │ Update Variant  │
                    └─────────────────┘
```

### ۲. رزرو موجودی (Order Created)

```
┌─────────────┐     ┌─────────────────┐
│ Order Event │────▶│ Consumer        │
│ Created     │     │ Handler         │
└─────────────┘     └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Check Available │
                    │ Quantity        │
                    └────────┬────────┘
                             │
              ┌──────────────┴──────────────┐
              │                             │
         Available                    Not Available
              │                             │
              ▼                             ▼
     ┌─────────────────┐          ┌─────────────────┐
     │ Create Reserve  │          │ Publish Event   │
     │ Transaction     │          │ OutOfStock      │
     └────────┬────────┘          └─────────────────┘
              │
              ▼
     ┌─────────────────┐
     │ Update Inventory│
     │ reserved_qty++  │
     └────────┬────────┘
              │
              ▼
     ┌─────────────────┐
     │ Publish Event   │
     │ StockReserved   │
     └─────────────────┘
```

### ۳. تکمیل سفارش (Order Completed)

```
┌─────────────┐     ┌─────────────────┐
│ Order Event │────▶│ Consumer        │
│ Completed   │     │ Handler         │
└─────────────┘     └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Find Reservation│
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Confirm Reserve │
                    │ status=confirmed│
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Create Withdraw │
                    │ Transaction     │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Update Inventory│
                    │ quantity--      │
                    │ reserved_qty--  │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Sync with       │
                    │ Catalog Service │
                    └─────────────────┘
```

### ۴. لغو سفارش (Order Cancelled)

```
┌─────────────┐     ┌─────────────────┐
│ Order Event │────▶│ Consumer        │
│ Cancelled   │     │ Handler         │
└─────────────┘     └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Find Reservation│
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Release Reserve │
                    │ status=released │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Create Release  │
                    │ Transaction     │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Update Inventory│
                    │ reserved_qty--  │
                    │ available_qty++ │
                    └─────────────────┘
```

---

## 🔗 ارتباط مایکروسرویس انبار با کاتالوگ

### گزینه ۱: Event-Driven (پیشنهادی)

```
┌─────────────────┐                    ┌─────────────────┐
│ Catalog Service │                    │Inventory Service│
└────────┬────────┘                    └────────┬────────┘
         │                                      │
         │  ┌─────────────────────────────┐     │
         ├─▶│      Message Broker         │◀────┤
         │  │   (RabbitMQ / Kafka / NATS) │     │
         │  └─────────────────────────────┘     │
         │                                      │
         │  Events:                             │
         │  • product.created                   │
         │  • variant.created                   │
         │  • stock.updated ◀──────────────────┤
         │  • order.created                     │
         │  • order.completed                   │
         │  • order.cancelled                   │
         │                                      │
         ▼                                      ▼
   Update stock_quantity              Process stock changes
```

**مزایا:**
- Loose coupling
- Scalability
- Async processing
- Fault tolerance

### گزینه ۲: Sync REST API

```
┌─────────────────┐     HTTP/REST      ┌─────────────────┐
│ Catalog Service │◀──────────────────▶│Inventory Service│
└─────────────────┘                    └─────────────────┘

GET  /inventory/variant/{id}
POST /inventory/reserve
POST /inventory/release
POST /inventory/withdraw
```

### گزینه ۳: gRPC (برای عملیات Real-time)

```protobuf
// inventory.proto
service InventoryService {
    rpc CheckAvailability(CheckRequest) returns (CheckResponse);
    rpc ReserveStock(ReserveRequest) returns (ReserveResponse);
    rpc ReleaseStock(ReleaseRequest) returns (ReleaseResponse);
    rpc GetInventory(GetInventoryRequest) returns (InventoryResponse);
}

message CheckRequest {
    int64 variant_id = 1;
    double quantity = 2;
}

message CheckResponse {
    bool available = 1;
    double available_quantity = 2;
}
```

### پیشنهاد ترکیبی:

| عملیات | روش ارتباط | دلیل |
|--------|-----------|------|
| بررسی موجودی | gRPC | نیاز به پاسخ سریع |
| رزرو موجودی | gRPC | نیاز به تراکنش همزمان |
| بروزرسانی موجودی در کاتالوگ | Event | async، بدون نیاز به پاسخ |
| گزارشات | REST | کاربرد غیربحرانی |

---

## 🛡️ نکات مهندسی

### ۱. همگام‌سازی با دیتابیس (Sync)

```go
// استراتژی همگام‌سازی
type SyncStrategy interface {
    SyncInventoryToCatalog(ctx context.Context, variantID int64) error
    SyncAllInventories(ctx context.Context) error
}

// پیاده‌سازی با Eventual Consistency
type EventualSyncStrategy struct {
    publisher EventPublisher
}

func (s *EventualSyncStrategy) SyncInventoryToCatalog(ctx context.Context, variantID int64) error {
    inventory, _ := s.repo.GetByVariantID(ctx, variantID)
    
    event := StockUpdatedEvent{
        VariantID:    variantID,
        Quantity:     inventory.Quantity,
        AvailableQty: inventory.AvailableQty,
        UpdatedAt:    time.Now(),
    }
    
    return s.publisher.Publish("stock.updated", event)
}
```

### ۲. جلوگیری از Oversell

```go
// استفاده از Optimistic Locking
type Inventory struct {
    // ...
    Version int `json:"version" db:"version"`
}

func (r *InventoryRepo) ReserveWithLock(ctx context.Context, variantID int64, qty float64) error {
    tx, _ := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
    defer tx.Rollback()
    
    // قفل ردیف
    inv, err := r.getForUpdate(tx, variantID)
    if err != nil {
        return err
    }
    
    // بررسی موجودی
    if inv.AvailableQty < qty {
        return ErrInsufficientStock
    }
    
    // بروزرسانی
    inv.ReservedQty += qty
    inv.AvailableQty = inv.Quantity - inv.ReservedQty
    inv.Version++
    
    if err := r.updateWithVersion(tx, inv); err != nil {
        return err
    }
    
    return tx.Commit()
}
```

**روش‌های جلوگیری از Oversell:**

| روش | توضیحات | مزایا | معایب |
|-----|---------|-------|-------|
| **Pessimistic Lock** | قفل ردیف در DB | دقت بالا | کندی، bottleneck |
| **Optimistic Lock** | بررسی version | سریع‌تر | نیاز به retry |
| **Distributed Lock** | Redis SETNX | مناسب microservice | پیچیدگی |
| **Event Sourcing** | ثبت همه تراکنش‌ها | audit کامل | پیچیدگی بالا |

### ۳. ثبت تراکنش‌های انبار

```go
// Audit Trail برای هر تغییر
type StockAuditService struct {
    logRepo InventoryLogRepository
}

func (s *StockAuditService) LogChange(ctx context.Context, inv *Inventory, action string, userID int64) error {
    log := &InventoryLog{
        InventoryID: inv.ID,
        VariantID:   inv.VariantID,
        Action:      action,
        Field:       "quantity",
        OldValue:    json.RawMessage(`{"quantity": ` + strconv.FormatFloat(inv.OldQuantity, 'f', 2, 64) + `}`),
        NewValue:    json.RawMessage(`{"quantity": ` + strconv.FormatFloat(inv.Quantity, 'f', 2, 64) + `}`),
        UserID:      userID,
        IPAddress:   getIPFromContext(ctx),
        CreatedAt:   time.Now(),
    }
    
    return s.logRepo.Create(ctx, log)
}
```

### ۴. Event-Driven Architecture

```go
// تعریف Event ها
type StockEvent struct {
    Type      string    `json:"type"`
    VariantID int64     `json:"variant_id"`
    Quantity  float64   `json:"quantity"`
    Timestamp time.Time `json:"timestamp"`
}

const (
    EventStockAdded     = "stock.added"
    EventStockReserved  = "stock.reserved"
    EventStockReleased  = "stock.released"
    EventStockWithdrawn = "stock.withdrawn"
    EventLowStock       = "stock.low"
)

// Publisher
type EventPublisher interface {
    Publish(topic string, event interface{}) error
}

// Consumer Handler
type OrderEventHandler struct {
    stockService StockService
}

func (h *OrderEventHandler) HandleOrderCreated(event OrderCreatedEvent) error {
    for _, item := range event.Items {
        if err := h.stockService.Reserve(context.Background(), ReserveRequest{
            OrderID:   event.OrderID,
            VariantID: item.VariantID,
            Quantity:  item.Quantity,
        }); err != nil {
            return err
        }
    }
    return nil
}
```

### ۵. Cache Strategy

```go
// کش موجودی برای خواندن سریع
type CachedInventoryRepo struct {
    repo  InventoryRepository
    cache *redis.Client
    ttl   time.Duration
}

func (r *CachedInventoryRepo) GetAvailable(ctx context.Context, variantID int64) (float64, error) {
    key := fmt.Sprintf("inv:available:%d", variantID)
    
    // Try cache
    if val, err := r.cache.Get(ctx, key).Float64(); err == nil {
        return val, nil
    }
    
    // Fallback to DB
    inv, err := r.repo.GetByVariantID(ctx, variantID)
    if err != nil {
        return 0, err
    }
    
    // Update cache
    r.cache.Set(ctx, key, inv.AvailableQty, r.ttl)
    
    return inv.AvailableQty, nil
}

// Invalidate cache on update
func (r *CachedInventoryRepo) Update(ctx context.Context, inv *Inventory) error {
    if err := r.repo.Update(ctx, inv); err != nil {
        return err
    }
    
    key := fmt.Sprintf("inv:available:%d", inv.VariantID)
    r.cache.Del(ctx, key)
    
    return nil
}
```

---

## 🔧 Interfaces

```go
// InventoryRepository اینترفیس دسترسی به موجودی
type InventoryRepository interface {
    GetByVariantID(ctx context.Context, variantID int64) (*Inventory, error)
    GetByWarehouseAndVariant(ctx context.Context, warehouseID, variantID int64) (*Inventory, error)
    GetLowStock(ctx context.Context, threshold float64) ([]Inventory, error)
    Create(ctx context.Context, inv *Inventory) error
    Update(ctx context.Context, inv *Inventory) error
    UpdateQuantity(ctx context.Context, variantID int64, qty float64) error
    BulkUpdate(ctx context.Context, updates []InventoryUpdate) error
}

// StockTransactionRepository اینترفیس تراکنش‌های انبار
type StockTransactionRepository interface {
    Create(ctx context.Context, tx *StockTransaction) error
    GetByVariantID(ctx context.Context, variantID int64, filter TransactionFilter) ([]StockTransaction, error)
    GetBalance(ctx context.Context, variantID int64) (float64, error)
    SumByPeriod(ctx context.Context, variantID int64, from, to time.Time) (float64, error)
}

// ReservationRepository اینترفیس رزرو موجودی
type ReservationRepository interface {
    Create(ctx context.Context, res *StockReservation) error
    GetByOrderID(ctx context.Context, orderID int64) ([]StockReservation, error)
    GetPending(ctx context.Context, variantID int64) ([]StockReservation, error)
    GetExpired(ctx context.Context) ([]StockReservation, error)
    UpdateStatus(ctx context.Context, id int64, status int) error
    Release(ctx context.Context, id int64) error
    Confirm(ctx context.Context, id int64) error
}

// StockService اینترفیس سرویس موجودی
type StockService interface {
    // موجودی
    GetAvailable(ctx context.Context, variantID int64) (float64, error)
    CheckAvailability(ctx context.Context, variantID int64, qty float64) (bool, error)
    
    // تراکنش‌ها
    AddStock(ctx context.Context, req AddStockRequest) (*StockTransaction, error)
    WithdrawStock(ctx context.Context, req WithdrawRequest) (*StockTransaction, error)
    AdjustStock(ctx context.Context, req AdjustRequest) (*StockTransaction, error)
    
    // رزرو
    Reserve(ctx context.Context, req ReserveRequest) (*StockReservation, error)
    Release(ctx context.Context, reservationID int64) error
    Confirm(ctx context.Context, reservationID int64) error
    
    // همگام‌سازی
    SyncWithCatalog(ctx context.Context, variantID int64) error
}
```

---

## 📋 جداول دیتابیس

```sql
-- انبارها
CREATE TABLE warehouses (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE,
    address TEXT,
    is_default BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- واحدها
CREATE TABLE units (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    symbol VARCHAR(20),
    type TINYINT DEFAULT 1,
    step_unit DECIMAL(10,4),
    min DECIMAL(10,4),
    sub_unit JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- موجودی لحظه‌ای
CREATE TABLE inventories (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    warehouse_id BIGINT UNSIGNED NOT NULL,
    variant_id BIGINT UNSIGNED NOT NULL,
    unit_id BIGINT UNSIGNED,
    
    quantity DECIMAL(15,4) DEFAULT 0,
    reserved_qty DECIMAL(15,4) DEFAULT 0,
    available_qty DECIMAL(15,4) DEFAULT 0,
    
    min_threshold DECIMAL(15,4),
    max_threshold DECIMAL(15,4),
    reorder_point DECIMAL(15,4),
    
    average_cost BIGINT,
    last_cost BIGINT,
    
    version INT DEFAULT 1,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id),
    FOREIGN KEY (variant_id) REFERENCES variants(id),
    UNIQUE KEY unique_warehouse_variant (warehouse_id, variant_id)
);

-- تراکنش‌های انبار
CREATE TABLE stock_transactions (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    warehouse_id BIGINT UNSIGNED NOT NULL,
    variant_id BIGINT UNSIGNED NOT NULL,
    unit_id BIGINT UNSIGNED,
    
    type VARCHAR(20) NOT NULL,
    amount DECIMAL(15,4) NOT NULL,
    balance DECIMAL(15,4) NOT NULL,
    reserved_amount DECIMAL(15,4) DEFAULT 0,
    
    unit_price BIGINT,
    total_price BIGINT,
    
    order_item_id BIGINT UNSIGNED,
    reference_type VARCHAR(50),
    reference_id BIGINT UNSIGNED,
    
    operator_id BIGINT UNSIGNED NOT NULL,
    vendor_id BIGINT UNSIGNED,
    
    date DATE NOT NULL,
    note TEXT,
    status TINYINT DEFAULT 1,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id),
    FOREIGN KEY (variant_id) REFERENCES variants(id),
    INDEX idx_variant_date (variant_id, date),
    INDEX idx_type_status (type, status)
);

-- رزروهای موجودی
CREATE TABLE stock_reservations (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    warehouse_id BIGINT UNSIGNED NOT NULL,
    variant_id BIGINT UNSIGNED NOT NULL,
    order_id BIGINT UNSIGNED NOT NULL,
    order_item_id BIGINT UNSIGNED NOT NULL,
    
    quantity DECIMAL(15,4) NOT NULL,
    status TINYINT DEFAULT 0,
    
    expires_at TIMESTAMP NOT NULL,
    released_at TIMESTAMP NULL,
    confirmed_at TIMESTAMP NULL,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id),
    FOREIGN KEY (variant_id) REFERENCES variants(id),
    INDEX idx_order (order_id),
    INDEX idx_variant_status (variant_id, status),
    INDEX idx_expires (expires_at, status)
);

-- لاگ تغییرات
CREATE TABLE inventory_logs (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    inventory_id BIGINT UNSIGNED NOT NULL,
    variant_id BIGINT UNSIGNED NOT NULL,
    
    action VARCHAR(20) NOT NULL,
    field VARCHAR(50),
    old_value JSON,
    new_value JSON,
    
    user_id BIGINT UNSIGNED,
    ip_address VARCHAR(45),
    user_agent TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_inventory (inventory_id),
    INDEX idx_variant (variant_id),
    INDEX idx_created (created_at)
);
```

---

## 🔜 مراحل پیشنهادی پیاده‌سازی

1. [ ] تصمیم‌گیری: ماژول داخلی یا مایکروسرویس مستقل
2. [ ] ایجاد مدل‌های دامین و جداول
3. [ ] پیاده‌سازی Repository Layer
4. [ ] پیاده‌سازی Service Layer
5. [ ] پیاده‌سازی API Handlers
6. [ ] راه‌اندازی Event System
7. [ ] پیاده‌سازی همگام‌سازی با کاتالوگ
8. [ ] تست‌های Unit و Integration
9. [ ] پیاده‌سازی مکانیزم جلوگیری از Oversell
10. [ ] گزارشات و داشبورد انبار

---

## 📌 سوالات باز برای بررسی

1. آیا انبار باید مایکروسرویس مستقل باشد یا بخشی از کاتالوگ؟
2. آیا نیاز به پشتیبانی از چند انبار فیزیکی وجود دارد؟
3. آیا رزرو موجودی باید زمان انقضا داشته باشد؟
4. استراتژی محاسبه میانگین قیمت خرید چیست؟ (FIFO, LIFO, Weighted Average)
5. آیا نیاز به مدیریت سریال نامبر یا بچ نامبر وجود دارد؟
6. حداکثر زمان قابل قبول برای همگام‌سازی موجودی چقدر است؟

---

**تاریخ ایجاد:** ۱۴۰۴/۰۹/۱۶  
**نسخه:** 1.0.0 (پیشنهادی)  
**نویسنده:** تیم توسعه روچی

> ⚠️ این سند پیشنهادی است و قبل از پیاده‌سازی نیاز به بررسی و تایید تیم فنی دارد.

