INSERT INTO users (username, email, passwordhash, isadmin, accesslevel, createdat, updatedat) VALUES
('benar',     'benar@lapbytes.com',     'hash_ben123',   true,  10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('mike.dev',  'mike@lapbytes.com',      'hash_mike456',  false, 6,  CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('shiko.ui',  'shiko@lapbytes.com',     'hash_shiko789', false, 4,  CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('kevo.ops',  'kevo@lapbytes.com',      'hash_kevo321',  false, 3,  CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('zack.qa',   'zack.qa@lapbytes.com',   'hash_zack654',  false, 2,  CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
