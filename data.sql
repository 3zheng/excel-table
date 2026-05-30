-- --------------------------------------------------------
-- 主机:                           127.0.0.1
-- 服务器版本:                        11.8.7-MariaDB-ubu2404 - mariadb.org binary distribution
-- 服务器操作系统:                      debian-linux-gnu
-- HeidiSQL 版本:                  12.8.0.6908
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

-- 正在导出表  excel_table.users 的数据：~6 rows (大约)
INSERT IGNORE INTO `users` (`id`, `username`, `password_hash`, `user_role`, `can_edit_products`, `created_at`) VALUES
	(1, 'admin', 'changeme', 'admin', 1, '2026-05-27 09:42:25'),
	(2, 'dev', 'changeme', 'admin', 1, '2026-05-27 09:42:25'),
	(3, 'export_input', 'changeme', 'export_input', 0, '2026-05-27 09:42:25'),
	(4, 'export_review', 'changeme', 'export_review', 0, '2026-05-27 09:42:25'),
	(5, 'import_input_sa', 'changeme', 'import_input', 0, '2026-05-27 09:42:25'),
	(6, 'import_review_sa', 'changeme', 'import_review', 0, '2026-05-27 09:42:25');

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
