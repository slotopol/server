/*
Адаптировано для PostgreSQL (Supabase)
*/

-- Вставляем клубы
INSERT INTO "club" ("cid", "name", "bank", "fund", "lock", "rate", "mrtp") VALUES
(1, 'virtual', 10000, 1000000, 0, 2.5, 95),
(2, 'global', 0, 0, 0, 2.5, 0);

-- Вставляем пользователей
INSERT INTO "user" ("uid", "email", "secret", "name", "status", "gal") VALUES
(1, 'admin@example.org', '0YBoaT', 'admin', 1, 31),
(2, 'dealer@example.org', 'LtpkAr', 'dealer', 1, 3),
(3, 'player@example.org', 'iVI05M', 'player', 1, 1);

-- Вставляем свойства (балансы и права в клубах)
-- Сумма 1+4+8 заменена на 13 (dealer/booker/master)
INSERT INTO "props" ("cid", "uid", "wallet", "access", "mrtp") VALUES
(1, 1, 10000, 1, 0),
(2, 1, 0, 0, 0),
(1, 2, 10000, 13, 0),
(2, 2, 0, 0, 0),
(1, 3, 1000, 1, 98),
(2, 3, 0, 0, 98);
