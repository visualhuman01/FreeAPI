-- MySQL dump 10.13  Distrib 8.0.12, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: freeapi_db
-- ------------------------------------------------------
-- Server version	8.0.12

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
 SET NAMES utf8 ;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `dbconfig_datatype`
--

DROP TABLE IF EXISTS `dbconfig_datatype`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `dbconfig_datatype` (
  `datatype_id` int(11) NOT NULL AUTO_INCREMENT,
  `datatype_type` int(11) DEFAULT NULL,
  `datatype_name` varchar(16) DEFAULT NULL,
  `datatype_len` int(11) DEFAULT NULL,
  `datatype_is_fixed` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`datatype_id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `dbconfig_datatype`
--

LOCK TABLES `dbconfig_datatype` WRITE;
/*!40000 ALTER TABLE `dbconfig_datatype` DISABLE KEYS */;
INSERT INTO `dbconfig_datatype` VALUES (1,1,'TINYINT',1,1),(2,1,'SMALLINT',2,1),(3,1,'MEDIUMINT',3,1),(4,1,'INT',4,1),(5,1,'BIGINT',8,1),(6,1,'FLOAT',4,1),(7,1,'DOUBLE',8,1),(8,1,'DATE',3,1),(9,1,'TIME',3,1),(10,1,'YEAR',1,1),(11,1,'DATETIME',8,1),(12,1,'TIMESTAMP',8,1),(13,1,'CHAR',NULL,0),(14,1,'VARCHAR',NULL,0),(15,1,'TEXT',NULL,0);
/*!40000 ALTER TABLE `dbconfig_datatype` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2018-09-04  1:32:40
