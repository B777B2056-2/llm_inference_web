CREATE TABLE `users` (
     `id` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
     `created_at` DATETIME NOT NULL,
     `updated_at` DATETIME NOT NULL,
     `deleted_at` DATETIME NULL,
     `name` VARCHAR(255) NOT NULL DEFAULT '',
     `password` VARCHAR(255) NOT NULL DEFAULT '',
     UNIQUE INDEX `idx_users_name` (`name`),
     INDEX `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;