-- MySQL dump 10.13  Distrib 8.0.43, for Linux (x86_64)
--
-- Host: localhost    Database: report_db_prod
-- ------------------------------------------------------
-- Server version	8.0.43

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `branches`
--

DROP TABLE IF EXISTS `branches`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `branches` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_by` int DEFAULT NULL,
  `updated_by` int DEFAULT NULL,
  `deleted_by` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `branches`
--

LOCK TABLES `branches` WRITE;
/*!40000 ALTER TABLE `branches` DISABLE KEYS */;
INSERT INTO `branches` VALUES (1,'สำนักงานใหญ่','2025-07-24 08:37:35','2025-07-24 08:37:35',NULL,NULL,NULL,NULL),(2,'สาขาสันกำแพง','2025-07-24 08:37:44','2025-07-24 08:37:44',NULL,NULL,NULL,NULL),(4,'สันผักหวาน','2025-09-13 06:19:44','2025-09-13 06:19:44',NULL,NULL,NULL,NULL),(5,'หอพักนพดล','2025-09-26 07:18:11','2025-09-26 07:18:11',NULL,NULL,NULL,NULL);
/*!40000 ALTER TABLE `branches` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `departments`
--

DROP TABLE IF EXISTS `departments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `departments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `branch_id` int DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `departments`
--

LOCK TABLES `departments` WRITE;
/*!40000 ALTER TABLE `departments` DISABLE KEYS */;
INSERT INTO `departments` VALUES (1,'จัดส่ง',2,'2025-07-24 08:49:55','2025-08-27 05:21:19',NULL),(2,'บัญชี',1,'2025-07-24 09:24:36','2025-07-24 09:24:36',NULL),(3,'บัญชีสินค้าคงคลัง',1,'2025-07-24 09:24:36','2025-07-24 09:24:36',NULL),(4,'พัฒนาธุรกิจ',1,'2025-07-24 09:24:36','2025-07-24 09:24:36',NULL),(5,'Information Technology',1,'2025-07-24 09:24:37','2025-07-24 09:24:37',NULL),(6,'การตลาด',1,'2025-07-24 09:24:37','2025-07-24 09:24:37',NULL),(7,'การเงิน',1,'2025-07-24 09:24:37','2025-07-24 09:24:37',NULL),(8,'สินเชื่อ',1,'2025-07-24 09:24:37','2025-07-24 09:24:37',NULL),(9,'บุคคล-ธุรการ',1,'2025-07-24 09:24:37','2025-07-24 09:24:37',NULL),(10,'บริหารสินค้า',1,'2025-07-24 09:24:37','2025-07-24 09:24:37',NULL),(11,'เมคคาเฟ่',1,'2025-07-24 09:24:37','2025-07-24 09:24:37',NULL),(12,'ขายโครงการ',1,'2025-07-24 09:24:37','2025-07-24 09:24:37',NULL),(13,'ลูกค้าสัมพันธ์',1,'2025-07-24 09:24:37','2025-08-04 09:47:59',NULL),(14,'แคชเชียร์',1,'2025-07-24 09:24:37','2025-08-04 09:48:04',NULL),(15,'Drive Thru Division',1,'2025-07-24 09:24:37','2025-08-04 09:48:10',NULL),(16,'ช่างเอ็กซ์เปิร์ท',1,'2025-07-24 09:24:37','2025-07-24 09:24:37',NULL),(17,'ตรวจรับสินค้า',1,'2025-07-24 09:24:37','2025-07-24 09:24:37',NULL),(18,'Store',1,'2025-07-24 09:24:37','2025-08-04 09:48:16',NULL),(19,'Stock9',1,'2025-07-24 09:24:37','2025-07-24 09:24:37',NULL),(20,'ลูกค้าสัมพันธ์',2,'2025-07-24 09:24:37','2025-08-04 09:48:25',NULL),(21,'แคชเชียร์',2,'2025-07-24 09:24:37','2025-08-04 09:48:31',NULL),(22,'Drive Thru Division',2,'2025-07-24 09:24:37','2025-08-04 09:48:38',NULL),(23,'ตรวจรับสินค้า',2,'2025-07-24 09:24:37','2025-08-04 09:48:44',NULL),(24,'Store',2,'2025-07-24 09:24:37','2025-08-04 09:48:51',NULL),(25,'Support Division',2,'2025-07-24 09:24:37','2025-08-04 09:48:57',NULL),(26,'ค้าส่ง',2,'2025-07-24 09:24:37','2025-07-24 09:24:37',NULL),(27,'เซรามิค',2,'2025-08-07 08:56:20','2025-08-07 08:56:20',NULL),(28,'เครื่องมือช่าง',1,'2025-08-07 08:59:58','2025-08-07 08:59:58',NULL),(29,'สีและเคมีภัณฑ์',1,'2025-08-07 08:59:58','2025-08-07 08:59:58',NULL),(30,'ประปา',1,'2025-08-07 08:59:58','2025-08-07 08:59:58',NULL),(31,'ตกแต่ง',1,'2025-08-07 08:59:58','2025-08-07 08:59:58',NULL),(32,'สุขภัณฑ์',1,'2025-08-07 08:59:58','2025-08-07 08:59:58',NULL),(33,'เขียนแบบ',1,'2025-08-07 09:03:51','2025-08-07 09:03:51',NULL),(35,'Sonoff R2',4,'2025-09-26 07:39:01','2025-09-26 07:39:01',NULL);
/*!40000 ALTER TABLE `departments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ip_phones`
--

DROP TABLE IF EXISTS `ip_phones`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ip_phones` (
  `id` int NOT NULL AUTO_INCREMENT,
  `number` int DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `department_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_by` int DEFAULT NULL,
  `updated_by` int DEFAULT NULL,
  `deleted_by` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=112 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ip_phones`
--

LOCK TABLES `ip_phones` WRITE;
/*!40000 ALTER TABLE `ip_phones` DISABLE KEYS */;
INSERT INTO `ip_phones` VALUES (7,100,'CR ',13,'2025-08-09 02:46:10','2025-08-09 02:46:10',NULL,NULL,NULL,NULL),(8,100,'CR',13,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(9,1000,'Make Kafe',11,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(10,101,'CR',13,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(11,102,'CR',13,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(12,103,'CR',13,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(13,104,'CR',13,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(14,105,'CR',13,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(15,107,'CR -Manager',13,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(16,118,'Operator',13,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(17,120,'Cashier Center Faham',14,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(18,121,'Cashier',14,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(19,125,'Cashier Drivethru',14,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(20,131,'Mechanic Tools',28,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(21,132,'Color department',29,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(22,133,'ประปา',30,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(23,134,'ตกแต่ง',31,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(24,136,'สุขภัณท์',32,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(25,140,'Drivethru cement',15,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(26,141,'Drivethru steel',15,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(27,150,'GR',17,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(28,151,'GR end control',17,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(29,152,'GR',17,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(30,153,'stock9',19,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(31,154,'Aui-stock9',19,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(32,209,'WS',12,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(33,210,'WS',12,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(34,211,'WS',12,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(35,212,'WS',12,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(36,213,'WS',12,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(37,214,'WS',12,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(38,220,'WS',12,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(39,221,'WS',12,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(40,222,'PJ',12,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(41,223,'PJ',12,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(42,224,'PJ',12,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(43,225,'PJ',12,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(44,226,'WS',12,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(45,310,'CAT1 -1',10,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(46,311,'CAT1 -2',10,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(47,320,'CAT2 -1',10,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(48,321,'CAT2 -2',10,'2025-08-09 02:52:05','2025-08-09 02:52:05',NULL,NULL,NULL,NULL),(49,322,'CAT2 -3',10,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(50,330,'CAT3 -1',10,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(51,332,'CAT3 -2',10,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(52,333,'CAT3 -3',10,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(53,340,'CAT4 -1',10,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(54,341,'CAT4 -2',10,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(55,350,'CAT5 -1',10,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(56,500,'CDM Credit & Finance',8,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(57,506,'AC',2,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(58,507,'AC',2,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(59,508,'ไม่มีข้อมูล',2,'2025-08-09 02:52:06','2025-08-22 02:47:14',NULL,NULL,NULL,NULL),(60,509,'AC',2,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(61,510,'CA Pat',7,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(62,512,'ไม่มีข้อมูล',8,'2025-08-09 02:52:06','2025-08-22 02:48:17',NULL,NULL,NULL,NULL),(63,513,'CD',8,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(64,514,'ไม่มีข้อมูล',7,'2025-08-09 02:52:06','2025-08-22 02:48:17',NULL,NULL,NULL,NULL),(65,515,'ไม่มีข้อมูล',8,'2025-08-09 02:52:06','2025-08-22 02:48:17',NULL,NULL,NULL,NULL),(66,520,'CA',7,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(67,521,'CA',7,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(68,530,'ICC',3,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(69,555,'ไม่มีข้อมูล',0,'2025-08-09 02:52:06','2025-08-22 02:48:17',NULL,NULL,NULL,NULL),(70,600,'HR',9,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(71,610,'HR',9,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(72,611,'HR',9,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(73,612,'HR',9,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(74,613,'HR',9,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(75,614,'HR',9,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(76,615,'LP',9,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(77,616,'BE',33,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(78,700,'IT -support',5,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(79,704,'IT -Champ',5,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(80,710,'IT -Moo',5,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(81,711,'IT Manager',5,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(82,712,'IT -Big',5,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(83,730,'MK',6,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(84,731,'MK',6,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(85,733,'MK',6,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(86,900,'ไม่มีข้อมูล',21,'2025-08-09 02:52:06','2025-08-22 02:48:17',NULL,NULL,NULL,NULL),(87,910,'HomeSolution CR',20,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(88,913,'Sale Division -S02',22,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(89,920,'Service Division -S02',22,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(90,921,'Branch-SCG Staff',27,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(91,922,'Exterior -SCG',27,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(92,924,'ไม่มีข้อมูล',20,'2025-08-09 02:52:06','2025-08-22 02:48:17',NULL,NULL,NULL,NULL),(93,925,'Ceramic warehouse',27,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(94,926,'Service -Division',22,'2025-08-09 02:52:06','2025-08-09 02:52:06',NULL,NULL,NULL,NULL),(95,927,'service -division',22,'2025-08-09 02:52:07','2025-08-09 02:52:07',NULL,NULL,NULL,NULL),(96,930,'Cashier Center HomeSolution',27,'2025-08-09 02:52:07','2025-08-09 02:52:07',NULL,NULL,NULL,NULL),(97,931,'Branch POS1',21,'2025-08-09 02:52:07','2025-08-09 02:52:07',NULL,NULL,NULL,NULL),(98,932,'Branch POS2',21,'2025-08-09 02:52:07','2025-08-09 02:52:07',NULL,NULL,NULL,NULL),(99,933,'ไม่มีข้อมูล',21,'2025-08-09 02:52:07','2025-08-22 02:48:17',NULL,NULL,NULL,NULL),(100,934,'ไม่มีข้อมูล',0,'2025-08-09 02:52:07','2025-08-22 02:48:17',NULL,NULL,NULL,NULL),(101,942,'ไม่มีข้อมูล',0,'2025-08-09 02:52:07','2025-08-22 02:48:17',NULL,NULL,NULL,NULL),(102,943,'Color -Tile',22,'2025-08-09 02:52:07','2025-08-09 02:52:07',NULL,NULL,NULL,NULL),(103,950,'Support Division',25,'2025-08-09 02:52:07','2025-08-09 02:52:07',NULL,NULL,NULL,NULL),(104,951,'Support Division',25,'2025-08-09 02:52:07','2025-08-09 02:52:07',NULL,NULL,NULL,NULL),(105,952,'Support Division',25,'2025-08-09 02:52:07','2025-08-09 02:52:07',NULL,NULL,NULL,NULL),(106,955,'Support Division',25,'2025-08-09 02:52:07','2025-08-09 02:52:07',NULL,NULL,NULL,NULL),(107,963,'ไม่มีข้อมูล',0,'2025-08-09 02:52:07','2025-08-22 02:48:17',NULL,NULL,NULL,NULL),(108,971,'Cement',22,'2025-08-09 02:52:07','2025-08-09 02:52:07',NULL,NULL,NULL,NULL),(109,972,'ไม่มีข้อมูล',22,'2025-08-09 02:52:07','2025-08-22 02:48:17',NULL,NULL,NULL,NULL),(111,701,'Soft Phone',5,'2025-09-13 06:20:25','2025-09-13 06:20:25',NULL,NULL,NULL,NULL);
/*!40000 ALTER TABLE `ip_phones` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `issue_types`
--

DROP TABLE IF EXISTS `issue_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `issue_types` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `issue_types`
--

LOCK TABLES `issue_types` WRITE;
/*!40000 ALTER TABLE `issue_types` DISABLE KEYS */;
INSERT INTO `issue_types` VALUES (1,'hardware','2025-08-16 05:13:30'),(2,'software','2025-08-16 05:13:30'),(3,'request','2025-08-16 05:13:30');
/*!40000 ALTER TABLE `issue_types` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `progress`
--

DROP TABLE IF EXISTS `progress`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `progress` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `task_id` bigint NOT NULL,
  `progress_text` text NOT NULL,
  `file_paths` json DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `progress`
--

LOCK TABLES `progress` WRITE;
/*!40000 ALTER TABLE `progress` DISABLE KEYS */;
/*!40000 ALTER TABLE `progress` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `resolutions`
--

DROP TABLE IF EXISTS `resolutions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `resolutions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `tasks_id` int DEFAULT NULL,
  `text` text,
  `file_paths` json DEFAULT NULL,
  `telegram_id` int DEFAULT NULL,
  `resolved_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `resolutions`
--

LOCK TABLES `resolutions` WRITE;
/*!40000 ALTER TABLE `resolutions` DISABLE KEYS */;
/*!40000 ALTER TABLE `resolutions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `responsibilities`
--

DROP TABLE IF EXISTS `responsibilities`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `responsibilities` (
  `id` int NOT NULL AUTO_INCREMENT,
  `telegram_username` varchar(50) DEFAULT NULL,
  `name` varchar(100) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `responsibilities`
--

LOCK TABLES `responsibilities` WRITE;
/*!40000 ALTER TABLE `responsibilities` DISABLE KEYS */;
INSERT INTO `responsibilities` VALUES (2,NULL,'สงกรานต์ ยิ้มยินดี','2025-08-23 02:26:44','2025-08-25 01:49:32'),(3,NULL,'กิตติศักดิ์ อินสองใจ','2025-08-23 02:26:58','2025-08-25 01:49:42'),(6,'@satit13','satit chomwattana','2025-09-13 04:44:13','2025-09-13 04:44:13'),(7,'@Panithan_T','ปณิธาน กันแก้ว','2025-09-13 04:44:46','2025-09-13 04:44:46'),(8,'@Weerawut','Weerawut Lukkanatorn','2025-09-13 04:45:24','2025-09-13 04:45:24');
/*!40000 ALTER TABLE `responsibilities` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `scores`
--

DROP TABLE IF EXISTS `scores`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `scores` (
  `id` int NOT NULL AUTO_INCREMENT,
  `department_id` int DEFAULT NULL,
  `year` int DEFAULT NULL,
  `month` int DEFAULT NULL,
  `score` int DEFAULT '100',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `scores`
--

LOCK TABLES `scores` WRITE;
/*!40000 ALTER TABLE `scores` DISABLE KEYS */;
/*!40000 ALTER TABLE `scores` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `systems_program`
--

DROP TABLE IF EXISTS `systems_program`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `systems_program` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `priority` tinyint DEFAULT '0',
  `type` int DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_by` int DEFAULT NULL,
  `updated_by` int DEFAULT NULL,
  `deleted_by` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `systems_program`
--

LOCK TABLES `systems_program` WRITE;
/*!40000 ALTER TABLE `systems_program` DISABLE KEYS */;
INSERT INTO `systems_program` VALUES (3,'ERP-WEB',0,2,'2025-08-05 03:24:10','2025-08-25 06:38:57',NULL,NULL,NULL,NULL),(4,'CHAMP-BC',0,2,'2025-08-05 03:24:15','2025-08-25 06:38:53',NULL,NULL,NULL,NULL),(5,'MERCHANT',0,2,'2025-08-05 03:24:20','2025-08-25 06:38:48',NULL,NULL,NULL,NULL),(6,'STOCK9-ADMIN',0,2,'2025-08-05 03:25:37','2025-08-25 06:38:43',NULL,NULL,NULL,NULL),(7,'STOCK9-APP',0,2,'2025-08-05 03:25:40','2025-08-25 06:38:39',NULL,NULL,NULL,NULL),(8,'MEMBER-SHOPPING',0,2,'2025-08-05 03:25:46','2025-08-25 06:38:35',NULL,NULL,NULL,NULL),(10,'Printer',2,1,'2025-09-13 06:25:42','2025-09-13 06:25:42',NULL,NULL,NULL,NULL);
/*!40000 ALTER TABLE `systems_program` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tasks`
--

DROP TABLE IF EXISTS `tasks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tasks` (
  `id` int NOT NULL AUTO_INCREMENT,
  `ticket_no` varchar(50) DEFAULT NULL,
  `phone_id` int DEFAULT NULL,
  `issue_type` int DEFAULT NULL,
  `system_id` int DEFAULT NULL,
  `issue_else` varchar(50) DEFAULT NULL,
  `department_id` int DEFAULT NULL,
  `text` mediumtext,
  `reported_by` varchar(50) DEFAULT NULL,
  `assignto_id` int DEFAULT NULL,
  `assignto` varchar(50) DEFAULT NULL,
  `solution_id` int DEFAULT NULL,
  `telegram_id` int DEFAULT NULL,
  `status` int DEFAULT NULL,
  `file_paths` json DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `resolved_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_by` int DEFAULT NULL,
  `updated_by` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tasks`
--

LOCK TABLES `tasks` WRITE;
/*!40000 ALTER TABLE `tasks` DISABLE KEYS */;
/*!40000 ALTER TABLE `tasks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `telegram_chat`
--

DROP TABLE IF EXISTS `telegram_chat`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `telegram_chat` (
  `id` int NOT NULL AUTO_INCREMENT,
  `chat_id` bigint DEFAULT NULL,
  `chat_name` varchar(50) DEFAULT NULL,
  `report_id` int DEFAULT NULL,
  `solution_id` int DEFAULT NULL,
  `assignto_id` int DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `telegram_chat`
--

LOCK TABLES `telegram_chat` WRITE;
/*!40000 ALTER TABLE `telegram_chat` DISABLE KEYS */;
/*!40000 ALTER TABLE `telegram_chat` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(100) DEFAULT NULL,
  `plain_password` varchar(255) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `role` varchar(50) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (3,'it',NULL,'$2a$10$RICikdc9QveLkuWMpH.4JuZOx4U9Y3XQlBZiZXsXU6aI34LNrmxqm','admin','2025-09-13 04:53:05','2025-09-13 04:53:05',NULL);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-09-27  6:52:06
