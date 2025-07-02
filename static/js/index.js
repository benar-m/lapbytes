    const API_BASE = '/api';
    const productsGrid = document.querySelector('.products-grid');
    let loadMoreBtn = document.querySelector('.load-more-btn') || document.querySelector('.btn-outline');
    let currentPage = 1;
    const limit = 6;
    let isLoading = false;
    let hasMoreItems = true;
    
    async function fetchLaptops(page, limit) {
        try {
            isLoading = true;
            updateLoadingState();
            const response = await fetch(`${API_BASE}/products/${limit}/${page}`);
            if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
            return await response.json();
        } catch (error) {
            showError('Failed to load laptops. Please try again.');
            return null;
        } finally {
            isLoading = false;
            updateLoadingState();
        }
    }
    
    function updateLoadingState() {
        if (loadMoreBtn) {
            loadMoreBtn.disabled = isLoading;
            loadMoreBtn.textContent = isLoading ? 'Loading...' : 'Load More';
        }
    }
    
    function createLaptopCard(laptop) {
        const name = laptop.name || 'Unnamed Laptop';
        const brand = laptop.brand || '';
        const price = laptop.price || 0;
        const stockBadge = laptop.is_in_stock ? 
            `<span class="badge in-stock">In Stock</span>` : 
            `<span class="badge out-of-stock">Out of Stock</span>`;
        const imageUrl = laptop.image_url || '';
        
        return `
            <div class="product-card" onclick="window.location.href='/product/${laptop.id}'" style="cursor: pointer;">
                <div class="product-image">
                    ${imageUrl ? `<img src="${imageUrl}" alt="${name}" onerror="this.style.display='none'">` : ''}
                    <div class="product-badges">${stockBadge}</div>
                </div>
                <div class="product-info">
                    <div class="product-header">
                        <h3 class="product-name">${name}</h3>
                        <span class="product-brand">${brand}</span>
                    </div>
                    <div class="product-price">KSH ${price.toLocaleString()}</div>
                    <div class="product-specs">
                        <div class="spec-group">
                            <h4>Quick Specs</h4>
                            <div class="spec-item">
                                <span class="spec-label">CPU:</span>
                                <span class="spec-value">${laptop.cpu_maker || ''} ${laptop.cpu_generation || ''}</span>
                            </div>
                            <div class="spec-item">
                                <span class="spec-label">RAM:</span>
                                <span class="spec-value">${laptop.ram_size || '0'}GB</span>
                            </div>
                            <div class="spec-item">
                                <span class="spec-label">Storage:</span>
                                <span class="spec-value">${laptop.ssd ? laptop.ssd_size + 'GB SSD' : (laptop.hdd ? laptop.hdd_size + 'GB HDD' : 'N/A')}</span>
                            </div>
                        </div>
                    </div>
                    <div class="product-actions">
                        <button class="btn btn-primary" onclick="event.stopPropagation(); addToCart(${laptop.id})" ${!laptop.is_in_stock ? 'disabled' : ''}>
                            ${!laptop.is_in_stock ? 'Out of Stock' : 'Add to Cart'}
                        </button>
                        <button class="btn btn-secondary" onclick="event.stopPropagation(); toggleWishlist(${laptop.id})">
                            <i class="fas fa-heart"></i>
                        </button>
                    </div>
                </div>
            </div>
        `;
    }
    
    function renderLaptops(laptops, append = false) {
        if (!laptops || laptops.length === 0) {
            if (!append) {
                productsGrid.innerHTML = `
                    <div class="no-products" style="grid-column: 1/-1; text-align: center; padding: 2rem;">
                        <h3>No laptops found</h3>
                        <p>Check back later for new arrivals!</p>
                    </div>
                `;
            }
            return;
        }
        
        const cardsHTML = laptops.map(laptop => createLaptopCard(laptop)).join('');
        if (append) {
            productsGrid.insertAdjacentHTML('beforeend', cardsHTML);
        } else {
            productsGrid.innerHTML = cardsHTML;
        }
        hasMoreItems = laptops.length >= limit;
        updateLoadMoreButton();
    }
    
    function updateLoadMoreButton() {
        if (loadMoreBtn) {
            loadMoreBtn.style.display = hasMoreItems ? 'block' : 'none';
        }
    }
    
    function showError(message) {
        productsGrid.innerHTML = `
            <div class="error-message" style="grid-column: 1/-1; text-align: center; padding: 2rem;">
                <h3 style="color: #dc3545;">${message}</h3>
                <button onclick="loadInitialLaptops()" class="btn btn-primary" style="margin-top: 1rem;">Try Again</button>
            </div>
        `;
    }
    
    async function loadInitialLaptops() {
        currentPage = 1;
        const data = await fetchLaptops(currentPage, limit);
        if (data && data.products) {
            renderLaptops(data.products, false);
        }
    }
    
    async function loadMoreLaptops() {
        if (isLoading || !hasMoreItems) return;
        currentPage++;
        const data = await fetchLaptops(currentPage, limit);
        if (data && data.products) {
            renderLaptops(data.products, true);
        }
    }
    
    function addToCart(laptopId) {
        const cartCount = document.querySelector('.cart-count');
        if (cartCount) {
            const currentCount = parseInt(cartCount.textContent) || 0;
            cartCount.textContent = currentCount + 1;
        }
    }
    
    function toggleWishlist(laptopId) {
    }
    
    document.addEventListener('DOMContentLoaded', function() {
        if (!loadMoreBtn) {
            const loadMoreContainer = document.createElement('div');
            loadMoreContainer.className = 'load-more-container';
            loadMoreContainer.style.cssText = 'text-align: center; margin-top: 2rem;';
            const newLoadMoreBtn = document.createElement('button');
            newLoadMoreBtn.className = 'btn btn-outline load-more-btn';
            newLoadMoreBtn.textContent = 'Load More';
            newLoadMoreBtn.addEventListener('click', loadMoreLaptops);
            loadMoreContainer.appendChild(newLoadMoreBtn);
            if (productsGrid && productsGrid.parentNode) {
                productsGrid.parentNode.insertBefore(loadMoreContainer, productsGrid.nextSibling);
            }
            loadMoreBtn = newLoadMoreBtn;
        }
        loadInitialLaptops();
    });