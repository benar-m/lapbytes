# LapBytes API Endpoints

##  Authentication
- `POST /api/auth/register` — Register a new user  
- `POST /api/auth/login` — Log in and get JWT  
- `GET /api/auth/profile` — Get current user info  
- `POST /api/auth/logout` — (Optional) Invalidate session

---

##  Products
- `GET /api/products` — List all laptops  
- `GET /api/products/{id}` — Get product details  
- `GET /api/products/search?q=` — Search or filter products

---

##  Cart
- `GET /api/cart` — View cart items  
- `POST /api/cart` — Add item to cart  
- `PUT /api/cart/{item_id}` — Update quantity  
- `DELETE /api/cart/{item_id}` — Remove item from cart

---

##  Checkout / Orders
- `POST /api/checkout` — Create new order from cart  
- `GET /api/orders` — List user's orders  
- `GET /api/orders/{id}` — Get single order details

---

## Admin (Protected)
- `POST /api/admin/products` — Add new laptop  
- `PUT /api/admin/products/{id}` — Update laptop info  
- `DELETE /api/admin/products/{id}` — Delete laptop  
- `GET /api/admin/orders` — View all orders  
- `GET /api/admin/users` — View all registered users  
- `GET /api/admin/users/{id}` — View specific user details

---

- `POST /api/products/{id}/reviews` — Add product review  
- `GET /api/products/{id}/reviews` — Fetch reviews  
- `GET /api/wishlist` — View wishlist  
- `POST /api/wishlist` — Add to wishlist  
- `DELETE /api/wishlist/{id}` — Remove from wishlist


---

## Render Endpoints 
- `GET /` — Homepage  
- `GET /login` — Login page  
- `GET /register` — Registration page  
- `GET /products` — Product listing page  
- `GET /products/{id}` — Product details page  
- `GET /cart` — View cart  
- `GET /checkout` — Checkout page  
- `GET /orders` — User’s order history  
- `GET /admin` — Admin dashboard  
- `GET /admin/products` — Admin laptop listing  
- `GET /admin/orders` — Admin order view  
- `GET /admin/users` — Admin user management  
