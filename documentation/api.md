# Endpoint Structure
Organize endpoints by resource type, following RESTful conventions:

## **1. Users Endpoints**
```
GET    /api/v1/users/                    # List all users (paginated)
POST   /api/v1/users/                    # Create new user
GET    /api/v1/users/{id}/               # Get user by ID
PUT    /api/v1/users/{id}/               # Update user
DELETE /api/v1/users/{id}/               # Delete user
```

## **2. Categories Endpoints**
```
GET    /api/v1/categories/               # List all categories (paginated)
POST   /api/v1/categories/               # Create new category
GET    /api/v1/categories/{id}/          # Get category by ID
PUT    /api/v1/categories/{id}/          # Update category
DELETE /api/v1/categories/{id}/          # Delete category
GET    /api/v1/categories/{id}/products/ # Get products under this category
GET    /api/v1/categories/{id}/children/ # Get all children of this category
```

## **3. Attributes Endpoints**
```
GET    /api/v1/attributes/               # List all attributes (paginated)
POST   /api/v1/attributes/               # Create new attribute
GET    /api/v1/attributes/{id}/          # Get attribute by ID
PUT    /api/v1/attributes/{id}/          # Update attribute
DELETE /api/v1/attributes/{id}/          # Delete attribute
```

## **4. Products Endpoints**
```
GET    /api/v1/products/                 # List all products (paginated)
POST   /api/v1/products/                 # Create new product
GET    /api/v1/products/{id}/            # Get product by ID
PUT    /api/v1/products/{id}/            # Update product
DELETE /api/v1/products/{id}/            # Delete product
GET    /api/v1/products/{id}/skus/       # Get SKUs for this product
```

## **5. SKUs Endpoints**
```
GET    /api/v1/skus/                     # List all SKUs (paginated)
POST   /api/v1/skus/                     # Create new SKU
GET    /api/v1/skus/{id}/                # Get SKU by ID
PUT    /api/v1/skus/{id}/                # Update SKU
DELETE /api/v1/skus/{id}/                # Delete SKU
GET    /api/v1/skus/{id}/attributes/     # Get all attributes for a SKU
POST   /api/v1/skus/{id}/attributes/     # Add/Update attributes to SKU (bulk)
DELETE /api/v1/skus/{id}/attributes/{attribute_id}/  # Remove specific attribute from SKU
```

## **6. Images Endpoints**
```
# Product Images
GET    /api/v1/products/{id}/images/     # Get all images for a product
POST   /api/v1/products/{id}/images/     # Upload image(s) to product
PUT    /api/v1/products/{id}/images/{image_id}/  # Update image metadata (title, is_primary)
DELETE /api/v1/products/{id}/images/{image_id}/  # Delete product image

# SKU Images
GET    /api/v1/skus/{id}/images/         # Get all images for a SKU
POST   /api/v1/skus/{id}/images/         # Upload image(s) to SKU
PUT    /api/v1/skus/{id}/images/{image_id}/  # Update image metadata (title, is_primary)
DELETE /api/v1/skus/{id}/images/{image_id}/  # Delete SKU image
```

## **7. Public Catalog Endpoints** (Customer-facing)
```
GET    /api/v1/catalog/categories/       # Get public category tree
GET    /api/v1/catalog/categories/{id}/  # Get category detail with products
GET    /api/v1/catalog/products/         # Browse products (with filters & search)
GET    /api/v1/catalog/products/{slug}/  # Get product detail by slug (with SKUs)
GET    /api/v1/catalog/skus/{sku_number}/  # Get SKU detail by sku_number
```

## **8. Authentication Endpoints**
```
POST   /api/v1/auth/login/               # User login
POST   /api/v1/auth/logout/              # User logout
POST   /api/v1/auth/refresh/             # Refresh access token
```

## **9. Profile Endpoints**
```
GET    /api/v1/profile/me/               # Get current user info
PUT    /api/v1/profile/me/               # Update current user profile
POST   /api/v1/profile/change-password/  # Change password
```


# Response Format
```python
# Success response
{
    "success": true,
    "data": {...},
    "error": None
}
```

# Pagination
Implement pagination for list endpoints with these parameters:
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20, max: 100)
- `sort_field`: Field to sort by
- `order_rule`: Sort order (asc/desc, default: asc)

Include pagination metadata in responses:
```
{
  "success": true,
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 157,
    "pages": 8
  },
  "error": null
}
```

# Error response
```
{
    "success": false,
    "data": None
    "error": {
        "code": "VALIDATION_ERROR",
        "message": "Invalid input data",
        "details": {...}
    }
}
```
# Filtering
- Simple filter: `?filter[field]=value`
- Operator filter: `?filter[field][operator]=value`
- Supported operators: `gt`, `lt`, `ge`, `le`, `ne`, `like`

# Request/Response Examples

## SKU Attributes Management

### Add/Update SKU Attributes (Bulk)
**Request:** `POST /api/v1/skus/{id}/attributes/`
```json
{
  "attributes": [
    {
      "attribute_id": 1,
      "value": "16GB"
    },
    {
      "attribute_id": 2,
      "value": "512GB SSD"
    },
    {
      "attribute_id": 3,
      "value": "Intel i7"
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "sku_id": 123,
    "attributes": [
      {
        "id": 1001,
        "attribute_id": 1,
        "attribute_name": "RAM",
        "value": "16GB"
      },
      {
        "id": 1002,
        "attribute_id": 2,
        "attribute_name": "Storage",
        "value": "512GB SSD"
      },
      {
        "id": 1003,
        "attribute_id": 3,
        "attribute_name": "Processor",
        "value": "Intel i7"
      }
    ]
  },
  "error": null
}
```

## Image Upload

### Upload Product Image
**Request:** `POST /api/v1/products/{id}/images/`
- Content-Type: `multipart/form-data`
- Fields:
  - `file`: Image file (required)
  - `title`: Image title (optional)
  - `is_primary`: Set as primary image (boolean, optional)

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 501,
    "file": "/uploads/products/product-123-image-501.jpg",
    "title": "Front View",
    "is_primary": true,
    "content_type": "product",
    "object_id": 123,
    "created_at": "2025-10-17T10:30:00Z"
  },
  "error": null
}
```

## Public Catalog

### Get Product Detail by Slug
**Request:** `GET /api/v1/catalog/products/asus-rog-strix-g15/`

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 123,
    "name": "ASUS ROG Strix G15",
    "slug": "asus-rog-strix-g15",
    "description": "Gaming laptop with powerful specs",
    "category": {
      "id": 5,
      "name": "Laptops",
      "slug": "laptops"
    },
    "images": [
      {
        "id": 501,
        "file": "/uploads/products/product-123-image-501.jpg",
        "title": "Front View",
        "is_primary": true
      }
    ],
    "skus": [
      {
        "id": 1001,
        "name": "ASUS ROG Strix G15 - 16GB/512GB",
        "sku_number": "ASUS-ROG-G15-001",
        "price": 15000000,
        "attributes": [
          {"name": "RAM", "value": "16GB"},
          {"name": "Storage", "value": "512GB SSD"},
          {"name": "Processor", "value": "AMD Ryzen 7 6800H"}
        ],
        "images": []
      },
      {
        "id": 1002,
        "name": "ASUS ROG Strix G15 - 32GB/1TB",
        "sku_number": "ASUS-ROG-G15-002",
        "price": 20000000,
        "attributes": [
          {"name": "RAM", "value": "32GB"},
          {"name": "Storage", "value": "1TB SSD"},
          {"name": "Processor", "value": "AMD Ryzen 9 6900H"}
        ],
        "images": []
      }
    ]
  },
  "error": null
}
```

# Implementation Notes

## Authentication & Authorization
- Admin endpoints (1-6, 9): Require authentication token
- Public catalog endpoints (7): No authentication required, but respect `is_active` flag
- Use JWT tokens with refresh token mechanism

## Image Handling
- Accept formats: JPEG, PNG, WebP
- Max file size: 5MB per image
- Auto-generate thumbnails (optional but recommended)
- Store file path/URL in database, actual file in storage (local/cloud)

## SKU Attributes
- POST `/skus/{id}/attributes/` should be **upsert** operation:
  - If attribute exists for SKU: update value
  - If attribute doesn't exist: create new
- Validate `attribute_id` exists in Attributes master data
- Return joined data with attribute name for better frontend UX

## Slug & SKU Number
- Slugs and SKU numbers must be unique
- Auto-generate if not provided:
  - Slug: from name (e.g., "ASUS ROG" → "asus-rog")
  - SKU Number: with prefix/pattern (e.g., "SKUName-{specifications}-{attributevalues}")

## Soft Delete (Recommended)
- Use `is_active` flag instead of hard delete
- DELETE endpoints set `is_active = false`
- Public catalog only shows `is_active = true` items