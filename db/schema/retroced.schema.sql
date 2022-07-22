-- MySQL dump 10.13  Distrib 8.0.22, for Linux (x86_64)
--
-- Host: 35.193.223.28    Database: retroced
-- ------------------------------------------------------
-- Server version	5.7.34-google

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `Clicks`
--

DROP TABLE IF EXISTS `Clicks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Clicks` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `link` mediumtext COLLATE utf8_unicode_ci NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=693 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Favourites`
--

DROP TABLE IF EXISTS `Favourites`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Favourites` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int(11) unsigned NOT NULL,
  `user_id` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `Favourites_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `Users` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `Favourites_ibfk_2` FOREIGN KEY (`product_id`) REFERENCES `Products` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=290 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Logs`
--

DROP TABLE IF EXISTS `Logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Logs` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `site_id` int(11) unsigned NOT NULL,
  `url_id` int(11) unsigned NOT NULL,
  `level` varchar(50) COLLATE utf8_unicode_ci NOT NULL DEFAULT 'INFO',
  `message` text COLLATE utf8_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1074738 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Products`
--

DROP TABLE IF EXISTS `Products`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Products` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `site_id` int(11) unsigned NOT NULL,
  `url_id` int(11) unsigned NOT NULL,
  `category` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `brand` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `title` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `description` text COLLATE utf8_unicode_ci NOT NULL,
  `price` int(11) NOT NULL DEFAULT '0',
  `retail_price` int(11) NOT NULL DEFAULT '0',
  `model` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `item_number` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `product_condition` text COLLATE utf8_unicode_ci NOT NULL,
  `accessories` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `measurements` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `img` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT 'https://via.placeholder.com/300x400',
  `approved` tinyint(1) NOT NULL DEFAULT '0',
  `featured` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `last_updated` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
  `created_at` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
  `color` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `size` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `shoe_size` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `sub_category` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `product_url` text COLLATE utf8_unicode_ci,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_product_url` (`product_url`(767)),
  KEY `site_id` (`site_id`),
  KEY `url_id` (`url_id`),
  KEY `product_site_category_idx` (`site_id`,`category`) USING HASH,
  KEY `product_created_idx` (`created_at`) USING BTREE,
  KEY `products_featured_idx` (`featured`) USING HASH,
  KEY `products_category_idx` (`category`),
  CONSTRAINT `Products_ibfk_1` FOREIGN KEY (`site_id`) REFERENCES `Sites` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
  CONSTRAINT `Products_ibfk_2` FOREIGN KEY (`url_id`) REFERENCES `Urls` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=81855981 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Sites`
--

DROP TABLE IF EXISTS `Sites`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Sites` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(250) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `url` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `sleep` int(11) NOT NULL DEFAULT '0',
  `active` tinyint(1) NOT NULL DEFAULT '0',
  `max_page` int(11) NOT NULL DEFAULT '0',
  `last_scraped` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Urls`
--

DROP TABLE IF EXISTS `Urls`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Urls` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `site_id` int(11) unsigned NOT NULL,
  `url` text COLLATE utf8_unicode_ci,
  `category` varchar(60) COLLATE utf8_unicode_ci NOT NULL,
  `last_updated` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
  `created_at` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
  PRIMARY KEY (`id`),
  KEY `site_id` (`site_id`),
  CONSTRAINT `Urls_ibfk_1` FOREIGN KEY (`site_id`) REFERENCES `Sites` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=1329717 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Users`
--

DROP TABLE IF EXISTS `Users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `guid` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `username` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `first_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `last_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `email` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `password` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `first_time` tinyint(1) NOT NULL DEFAULT '0',
  `fb_id` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
  `acl` int(4) NOT NULL DEFAULT '0',
  `partner_site_id` int(11) unsigned NOT NULL DEFAULT '0',
  `password_reset` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT ' ',
  PRIMARY KEY (`id`),
  UNIQUE KEY `guid` (`guid`)
) ENGINE=InnoDB AUTO_INCREMENT=50 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-08-09 18:21:24
