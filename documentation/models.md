# Entity Relationships
Base model:
```
┌─────────────────┐
│      Base       │
│─────────────────│
│ id (PK)         │
│ created_at      │
│ updated_at      │
│ created_by      │
│ updated_by      │
│ is_active       │
│ sequence        │
└─────────────────┘
```
This model will be inherited to all models in the ERD structure except Users

ERD structure:

```
┌─────────────────┐
│      Users      │
│─────────────────│
│ id (PK)         │
│ username        │
│ password        │
│ name            │
│ role            │
│ created_at      │
│ updated_at      │
└─────────────────┘

┌─────────────────┐    1:N    ┌─────────────────┐
│   Categories    │◄──────────│    Products     │
│─────────────────│           │─────────────────│
│ id (PK)         │           │ id (PK)         │
│ name            │           │ name            │
│ slug            │           │ slug            │
│ description     │           │ description     │
│ parent_id (FK)  │           │ category_id (FK)│
│ created_at      │           │ created_at      │
│ updated_at      │           │ updated_at      │
│ created_by (FK) │           │ created_by (FK) │
│ updated_by (FK) │           │ updated_by (FK) │
│ is_active       │           │ is_active       │
│ sequence        │           │ sequence        │
└─────────────────┘           └─────────────────┘
        │                              │ 1:N
        │ 1:N (self-referencing)       ▼
        ▼                         ┌─────────────────┐
┌─────────────────┐               │      Skus       │
│   Categories    │               │─────────────────│
│   (children)    │               │ id (PK)         │
│─────────────────│               │ name            │
│ parent_id (FK)  │               │ slug            │
└─────────────────┘               │ description     │
                                  │ sku_number      │
                                  │ price           │
                                  │ product_id (FK) │
                                  │ created_at      │
                                  │ updated_at      │
                                  │ created_by (FK) │
                                  │ updated_by (FK) │
                                  │ is_active       │
                                  │ sequence        │
                                  └─────────────────┘


# Attribute Management System

┌─────────────────┐    1:N    ┌─────────────────┐    N:1    ┌─────────────────┐
│   Attributes    │◄──────────│SkuAttributeValue│──────────►│      Skus       │
│─────────────────│           │─────────────────│           │─────────────────│
│ id (PK)         │           │ id (PK)         │           │ id (PK)         │
│ name            │           │ sku_id (FK)     │           │ name            │
│ code            │           │ attribute_id(FK)│           │ slug            │
│ data_type       │           │ value           │           │ description     │
│ uom             │           │ created_at      │           │ sku_number      │
│ created_at      │           │ updated_at      │           │ price           │
│ updated_at      │           │ created_by (FK) │           │ product_id (FK) │
│ created_by (FK) │           │ updated_by (FK) │           │ created_at      │
│ updated_by (FK) │           │ is_active       │           │ updated_at      │
│ is_active       │           │ sequence        │           │ created_by (FK) │
│ sequence        │           └─────────────────┘           │ updated_by (FK) │
└─────────────────┘                                         │ is_active       │
                                                            │ sequence        │
                                                            └─────────────────┘

┌─────────────────┐    Generic FK    ┌─────────────────┐
│     Images      │◄─────────────────│   Products      │
│─────────────────│                  │      Skus       │
│ id (PK)         │                  │─────────────────│
│ file            │                  │ content_type    │
│ title           │                  │ object_id       │
│ is_primary      │                  └─────────────────┘
│ content_type    │
│ object_id       │
│ created_at      │
│ updated_at      │
│ created_by (FK) │
│ updated_by (FK) │
│ is_active       │
│ sequence        │
└─────────────────┘
```

# Key Relationships Summary:
- **User** → **All Models** (created_by, updated_by)
- **Categories** → **Categories** (1:N self-referencing via parent_id)
- **Categories** → **Products** (1:N)
- **Products** → **Skus** (1:N)
- **Skus** ↔ **Attributes** (M:N via SkuAttributeValue with actual values)
- **Images** → **Products/Skus** (Generic Foreign Key)

# Attribute System Explanation:

## Design Approach: Simple & Flexible
Attributes use  **simple Many-to-Many relationship** between SKUs and Attributes through SkuAttributeValue. This approach gives **maximum flexibility** to handle various unexpected product variant.

## Flow:

```
1. Define Attributes (Global Master Data)
   ↓
2. SKU → Add any attributes as needed (SkuAttributeValue)
```

## Models & Relationships:

### 1. **Attributes** (Master Data)
Global attribute definitions yang bisa digunakan across products:
- `name`: "RAM", "Storage", "Processor", "Color", "Warranty"
- `code`: "ram", "storage", "processor", "color", "warranty"
- `data_type`: "string", "number", "decimal", "boolean"
- `uom`: Unit of measurement (GB, inch, GHz, years, etc.)

### 2. **SkuAttributeValue** (Actual Data)
Stores the actual value of each attribute for a specific SKU:
- `sku_id`: Which SKU this value belongs to
- `attribute_id`: Which attribute from master data
- `value`: The actual value (e.g., "16GB", "512GB SSD", "Intel i7", "Black")

## Example Scenario:

```
Product: "ASUS ROG Strix G15"

SKU 1: "ASUS ROG Strix G15 - Standard"
  SkuAttributeValue records:
  - attribute: "RAM", value: "16GB"
  - attribute: "Storage", value: "512GB SSD"
  - attribute: "Processor", value: "AMD Ryzen 7 6800H"
  - attribute: "Color", value: "Black"

SKU 2: "ASUS ROG Strix G15 - Gaming Pro"
  SkuAttributeValue records:
  - attribute: "RAM", value: "32GB"
  - attribute: "Storage", value: "1TB SSD"
  - attribute: "Processor", value: "AMD Ryzen 9 6900H"    ← Different processor
  - attribute: "Color", value: "RGB Backlit"
  - attribute: "Gaming_Mouse_Included", value: "Yes"      ← Unique attribute
  - attribute: "Extended_Warranty", value: "3 Years"     ← Unique attribute

SKU 3: "ASUS ROG Strix G15 - Limited Edition"
  SkuAttributeValue records:
  - attribute: "RAM", value: "64GB"
  - attribute: "Storage", value: "2TB NVMe"
  - attribute: "Color", value: "Gold"
  - attribute: "Limited_Edition_Number", value: "#001"    ← Completely unique
  - attribute: "Signed_Certificate", value: "Yes"        ← Completely unique
```

## Benefits:

✅ **Maximum Flexibility**: Each SKU can have any combination of attributes  
✅ **Future-Proof**: Easy to add new attributes without schema changes  
✅ **Real-World Friendly**: Handles unpredictable product variations  
✅ **Fast Development**: No need to setup templates first  
✅ **Industry Standard**: Used by most e-commerce platforms  

## Considerations:

⚠️ **Data Consistency**: Requires application-level validation to prevent typos  
⚠️ **UI Complexity**: Frontend needs to handle dynamic attribute sets  
⚠️ **Comparison**: May need additional logic for product comparison features  

## Recommended Best Practices:

1. **Standardize Common Attributes**: Use consistent naming (e.g., always "RAM" not "Memory")
2. **UI Guidelines**: Provide dropdown suggestions for common attributes
3. **Soft Validation**: Warn (don't block) if important attributes are missing
4. **Search Optimization**: Index commonly searched attributes