INSERT INTO products (
    name, brand, operatingsystem, operatingsystemversion, 
    hdd, ssd, hddsize, ssdsize, ramsize, 
    cpumaker, cpugen, cpumodel, yom, imageurl, price, screensize,
    hasgpu, gpumake, gpumaker, hasigpu, isinstock, instock
) VALUES 
(
    'ThinkPad X1 Carbon Gen 11', 'Lenovo', 'Windows', '11 Pro',
    FALSE, TRUE, 0.0, 1024.0, 16.0,
    'Intel', '13th Gen', 'Core i7-1365U', '2023', 
    'https://images.lenovo.com/is/image/lenovo/thinkpad-x1-carbon-gen-11', 
    1899.99, 14.0,
    FALSE, NULL, NULL, TRUE, TRUE, 15
),
(
    'MacBook Pro 14"', 'Apple', 'macOS', 'Ventura 13.0',
    FALSE, TRUE, 0.0, 512.0, 16.0,
    'Apple', 'M2', 'M2 Pro', '2023',
    'https://store.storeimages.cdn-apple.com/macbook-pro-14-space-gray',
    2499.99, 14.2,
    TRUE, 'M2 Pro GPU', 'Apple', FALSE, TRUE, 8  -- Swapped: model, maker
),
(
    'ROG Strix G15', 'ASUS', 'Windows', '11 Home',
    TRUE, TRUE, 1000.0, 512.0, 32.0,
    'AMD', 'Ryzen 5000', 'Ryzen 7 5800H', '2022',
    'https://dlcdnwebimgs.asus.com/gain/rog-strix-g15-g513-black',
    1299.99, 15.6,
    TRUE, 'RTX 3070', 'NVIDIA', FALSE, TRUE, 12  -- Swapped: model, maker
),
(
    'XPS 13 Plus', 'Dell', 'Windows', '11 Home',
    FALSE, TRUE, 0.0, 256.0, 8.0,
    'Intel', '12th Gen', 'Core i5-1240P', '2022',
    'https://i.dell.com/is/image/DellContent/xps-13-plus-9320-laptop',
    999.99, 13.4,
    FALSE, NULL, NULL, TRUE, TRUE, 20
),
(
    'Surface Laptop 5', 'Microsoft', 'Windows', '11 Home',
    FALSE, TRUE, 0.0, 512.0, 16.0,
    'Intel', '12th Gen', 'Core i7-1255U', '2022',
    'https://img-prod-cms-rt-microsoft-com.surface-laptop-5-platinum',
    1599.99, 13.5,
    FALSE, NULL, NULL, TRUE, FALSE, 0
);