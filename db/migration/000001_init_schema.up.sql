CREATE TABLE IF NOT EXISTS `Todo` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `title` varchar(255) DEFAULT NULL,
  `description` text,
  `create_date` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_date` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;