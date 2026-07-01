CREATE DATABASE IF NOT EXISTS `order`;
CREATE DATABASE IF NOT EXISTS `payment`;
CREATE DATABASE IF NOT EXISTS `shipping`;

USE `order`;

CREATE TABLE IF NOT EXISTS `stock_items` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `product_code` varchar(191) NOT NULL,
  `description` varchar(191) DEFAULT NULL,
  `quantity` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_stock_items_product_code` (`product_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `stock_items` (`product_code`, `description`, `quantity`) VALUES
  ('prod_1', 'Produto 1 - Camiseta', 100),
  ('prod_2', 'Produto 2 - Calça Jeans', 50),
  ('prod_3', 'Produto 3 - Tênis Esportivo', 75),
  ('prod_4', 'Produto 4 - Mochila', 30),
  ('prod_5', 'Produto 5 - Relógio', 20)
ON DUPLICATE KEY UPDATE `product_code` = `product_code`;