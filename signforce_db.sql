-- MySQL dump 10.13  Distrib 8.0.34, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: signforceapi
-- ------------------------------------------------------
-- Server version	8.0.34

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
-- Table structure for table `algos`
--

DROP TABLE IF EXISTS `algos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `algos` (
  `idAlgo` int NOT NULL AUTO_INCREMENT,
  `oidAlgo` varchar(30) NOT NULL,
  `asn1Repr` varchar(30) NOT NULL,
  `reprAlt` varchar(30) DEFAULT NULL,
  `useType` int NOT NULL,
  `isDeprecated` tinyint(1) NOT NULL DEFAULT '0',
  `inUse` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`idAlgo`)
) ENGINE=InnoDB AUTO_INCREMENT=76 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `algos`
--

LOCK TABLES `algos` WRITE;
/*!40000 ALTER TABLE `algos` DISABLE KEYS */;
INSERT INTO `algos` VALUES (1,'1.2.840.10040.4.3','dsa-with-sha1','id-dsa-with-sha1',2,0,0),(2,'1.2.840.10045.4.1','ecdsa-with-SHA1','sha1_ecdsa',2,0,0),(3,'1.2.840.10045.4.3.2','ecdsa-with-SHA256','sha256_ecdsa',2,0,1),(4,'1.2.840.10045.4.3.3','ecdsa-with-SHA384','sha384_ecdsa',2,0,0),(5,'1.2.840.10045.4.3.4','ecdsa-with-SHA512','sha512_ecdsa',2,0,0),(6,'1.2.840.113549.1.1.4','md5WithRSAEncryption','md5_rsa',2,1,0),(7,'1.2.840.113549.1.1.5','sha1-with-rsa-signature','sha1_rsa',2,1,0),(8,'1.2.840.113549.1.1.10','rsassa-pss','rsassa_pss',3,0,1),(9,'1.2.840.113549.1.1.11','sha256WithRSAEncryption','sha256_rsa',2,0,1),(10,'1.2.840.113549.1.1.12','sha384WithRSAEncryption','sha384_rsa',2,0,0),(11,'1.2.840.113549.1.1.13','sha512WithRSAEncryption','sha384_rsa',2,0,0),(12,'1.3.14.3.2.2','md4WitRSA','md4_rsa',2,1,0),(13,'1.3.14.3.2.3','md5WithRSA','md5_rsa',2,1,0),(14,'1.3.14.3.2.11','rsaSignature','rsa',2,0,0),(15,'1.3.14.3.2.13','dsaWithSHA','sha_dsa',2,1,0),(16,'1.3.14.3.2.15','shaWithRSASignature','sha_rsa',2,1,0),(17,'1.2.840.113549.2.2','md2','md2',1,1,0),(18,'1.2.840.113549.2.4','md4','md4',1,1,0),(19,'1.3.14.3.2.18','sha','sha',1,1,0),(20,'1.3.14.3.2.26','sha1','sha1',1,1,1),(21,'2.16.840.1.101.3.4.2.1','sha256','sha256',1,0,1),(22,'2.16.840.1.101.3.4.2.2','sha384','sha384',1,0,1),(23,'2.16.840.1.101.3.4.2.3','sha512','sha512',1,0,1),(24,'2.16.840.1.101.3.4.2.4','sha224','sha224',1,0,0),(25,'2.16.840.1.101.3.4.2.5','sha512-224','sha512_224',1,0,0),(26,'2.16.840.1.101.3.4.2.6','sha512-256	','sha512_256',1,0,0),(27,'2.16.840.1.101.3.4.2.7','sha3-224','sha3_224',1,0,0),(28,'2.16.840.1.101.3.4.2.8','sha3-256','sha3_256',1,0,0),(29,'2.16.840.1.101.3.4.2.9','sha3-384','sha3_384',1,0,0),(30,'2.16.840.1.101.3.4.2.10','sha3-512','sha3_512',1,0,0),(31,'2.16.840.1.101.3.4.2.11','shake128','shake128',1,0,0),(32,'2.16.840.1.101.3.4.2.12','shake256','shake256',1,0,0),(33,'2.16.840.1.101.3.4.1.23','aes192-OFB','aes192_OFB',3,0,0),(34,'16.840.1.101.3.4.1.24','aes192-CFB','aes192_CFB',3,0,0),(35,'2.16.840.1.101.3.4.1.25','id-aes192-wrap','id_aes192_wrap',3,0,0),(36,'2.16.840.1.101.3.4.1.26','aes192-GCM','aes192_GCM',3,0,0),(37,'2.16.840.1.101.3.4.1.27','aes192-CCM','aes192_CCM',3,0,0),(38,'2.16.840.1.101.3.4.1.28','aes192-wrap-pad','aes192_wrap_pad',3,0,0),(39,'2.16.840.1.101.3.4.1.41','aes256-ECB','aes256_ECB',3,0,0),(40,'2.16.840.1.101.3.4.1.42','aes256-CBC','aes256_CBC',3,0,0),(41,'2.16.840.1.101.3.4.1.43','aes256-OFB','aes256_OFB',3,0,0),(42,'2.16.840.1.101.3.4.1.44','aes256-CFB','aes256_CFB',3,0,0),(43,'2.16.840.1.101.3.4.1.45','id-aes256-wrap','id_aes256_wrap',3,0,0),(44,'2.16.840.1.101.3.4.1.46','aes256-GCM','aes256_GCM',3,0,0),(45,'2.16.840.1.101.3.4.1.47','aes256-CCM','aes256_CCM',3,0,0),(46,'2.16.840.1.101.3.4.1.48','aes256-wrap-pad','aes256_wrap_pad',3,0,0),(47,'1.3.14.3.2.6','desECB','desECB',3,1,0),(48,'1.3.14.3.2.7','desCBC','desCBC',3,1,0),(49,'1.3.14.3.2.8','desOFB','desOFB',3,1,0),(50,'1.3.14.3.2.9','desCFB','desCFB',3,1,0),(51,'1.3.14.3.2.10','desMAC','desMAC',3,1,0),(52,'1.3.14.3.2.17','desEDE','desEDE',3,0,1),(53,'1.3.6.1.4.1.4929.1.6','3Des','3des',3,0,0),(54,'1.3.6.1.4.1.4929.1.7','3DesECB','3desECB',3,0,0),(55,'1.3.6.1.4.1.4929.1.8','3DesCBC','3DesCBC',3,0,0),(56,'1.3.6.1.4.1.4929.1.9','3DesOFB','3DesOFB',3,0,0),(57,'1.3.6.1.4.1.4929.1.10','3DesCFB','3DesCFB',3,0,0),(58,'1.2.840.113549.1.1.1','rsaEncryption','rsa',4,0,1),(59,'1.2.840.113549.1.1.6','rsaOAEPEncryptionSET','rsaOAEPEncryptionSET',4,0,1),(60,'1.2.840.113549.1.1.7','id-RSAES-OAEP','id-RSAES-OAEP',4,0,1),(61,'1.2.840.113549.1.1.10','rsassa-pss','rsassa-pss',4,0,1),(62,'1.2.840.10045.2.1','ecPublicKey','id-ecPublicKey',4,0,1),(63,'1.3.132.0.6','secp112r1','secp112r1',4,0,1),(64,'1.3.132.0.10','secp256k1','secp256k1',4,0,1),(65,'1.3.132.0.28','secp128r1','secp128r1',4,0,1),(66,'1.3.132.0.31','secp192k1','secp192k1',4,0,1),(67,'1.3.132.0.33','secp224r1','secp224r1',4,0,1),(68,'1.3.132.0.34','secp384r1','secp384r1',4,0,1),(69,'1.3.132.0.35','secp521r1','secp521r1',4,0,1),(70,'1.2.840.10045.3.1.1','secp192r1','prime192v1',4,0,1),(71,'1.2.840.10045.3.1.7','secp256r1','prime256v1',4,0,1),(72,'1.3.101.110','id-X25519','curve25519',4,0,1),(73,'1.3.101.112','ed25519','Ed25519',4,0,1),(74,'1.3.101.113','ed448','Ed448',4,0,1),(75,'1.3.132.1.12','ecdh','ECDH',4,0,1);
/*!40000 ALTER TABLE `algos` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `countries`
--

DROP TABLE IF EXISTS `countries`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `countries` (
  `idCountry` int NOT NULL AUTO_INCREMENT,
  `nameCountry` varchar(80) NOT NULL,
  `aliasCountry` varchar(30) DEFAULT NULL,
  `shrinkName` varchar(5) NOT NULL,
  `phoneCode` int NOT NULL,
  `continent` enum('AFRICA','AMERICA','ASIA','EUROPE','OCEANIA') NOT NULL,
  PRIMARY KEY (`idCountry`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `countries`
--

LOCK TABLES `countries` WRITE;
/*!40000 ALTER TABLE `countries` DISABLE KEYS */;
INSERT INTO `countries` VALUES (1,'Estados Unidos Mexicanos','México','MX',52,'AMERICA'),(2,'United States of America','USA','US',1,'AMERICA'),(3,'Canada','Canada','CA',1,'AMERICA');
/*!40000 ALTER TABLE `countries` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `docstatus`
--

DROP TABLE IF EXISTS `docstatus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `docstatus` (
  `idDocStat` int NOT NULL AUTO_INCREMENT,
  `nameDocStat` varchar(45) DEFAULT NULL,
  `descriptionDocStat` varchar(100) DEFAULT NULL,
  `useDocStat` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`idDocStat`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `docstatus`
--

LOCK TABLES `docstatus` WRITE;
/*!40000 ALTER TABLE `docstatus` DISABLE KEYS */;
/*!40000 ALTER TABLE `docstatus` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `documents`
--

DROP TABLE IF EXISTS `documents`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `documents` (
  `idDocument` bigint NOT NULL AUTO_INCREMENT,
  `hashDoc` varchar(100) NOT NULL,
  `createdAtDoc` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `lastModifiedDoc` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deletedAtDoc` datetime DEFAULT NULL,
  `deletedReasonDoc` text,
  `ownerInstDoc_fk` bigint NOT NULL,
  `ownerTeamDoc_fk` bigint NOT NULL,
  `creatorUserDoc_fk` bigint NOT NULL,
  `nameDoc` varchar(100) NOT NULL,
  `pathDoc` varchar(100) NOT NULL,
  `extDoc` varchar(10) NOT NULL,
  `sizeB` int NOT NULL DEFAULT '0',
  `abstractDoc` text,
  `authUseStatus` int NOT NULL DEFAULT '0',
  `authRoleStatus` int NOT NULL DEFAULT '0',
  `activeDoc` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`idDocument`),
  KEY `fk_documents_owner_inst` (`ownerInstDoc_fk`),
  KEY `fk_documents_owner_team` (`ownerTeamDoc_fk`),
  KEY `fk_documents_creator_user` (`creatorUserDoc_fk`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `documents`
--

LOCK TABLES `documents` WRITE;
/*!40000 ALTER TABLE `documents` DISABLE KEYS */;
INSERT INTO `documents` VALUES (1,'DDUgZ6JF/tfIg1E2CAThpe69QzHt8yHUL/COutRKcAA=','2025-09-04 17:43:14','2025-09-04 17:43:14',NULL,NULL,71,0,26,'CV_MariaBarrera.pdf','././uploaded/documents/71/0/26','pdf',0,NULL,0,0,1),(2,'qIsSgfs8kFYPLMvSYoeH+d9NucWARs8upgqR18FuOxA=','2025-09-04 17:55:22','2025-09-04 17:55:22',NULL,NULL,71,0,26,'CV_2025_JJNM.pdf','./uploaded/documents/71/0/26','pdf',0,NULL,0,0,1),(3,'DDUgZ6JF/tfIg1E2CAThpe69QzHt8yHUL/COutRKcAA=','2025-09-04 17:55:22','2025-09-04 17:55:22',NULL,NULL,71,0,26,'CV_MariaBarrera.pdf','./uploaded/documents/71/0/26','pdf',0,NULL,0,0,1),(4,'thV6h7lh9hIVj8l7JoNdn2eizvhRdoG/AVIRuZX3Qcw=','2025-09-04 17:55:22','2025-09-04 17:55:22',NULL,NULL,71,0,26,'CV2024_eng.pdf','./uploaded/documents/71/0/26','pdf',0,NULL,0,0,1),(5,'1TKmbykynFN/h9JCvJ5+IBSriz+ST9AvAavQ/SdT8Us=','2025-09-04 18:00:56','2025-09-04 18:00:56',NULL,NULL,71,0,26,'configurepkcs11_cloud.png','./uploaded/documents/71/0/26','png',0,NULL,0,0,1),(6,'u1dImPaqo3+qhM4kICFk6tMXPSnms33GOFQFwL7e9AM=','2025-09-04 18:00:56','2025-09-04 18:00:56',NULL,NULL,71,0,26,'Confirmación _ Viva.pdf','./uploaded/documents/71/0/26','pdf',0,NULL,0,0,1),(7,'XbOaERJcgbF8PF4lvpCdTSVqEv1Gk0ic28LuxlSak+w=','2025-09-04 18:00:56','2025-09-04 18:00:56',NULL,NULL,71,0,26,'PRUEBA_SDT_2.pdf','./uploaded/documents/71/0/26','pdf',0,NULL,0,0,1),(8,'1TKmbykynFN/h9JCvJ5+IBSriz+ST9AvAavQ/SdT8Us=','2025-09-04 18:13:47','2025-09-04 18:13:47',NULL,NULL,66,0,20,'configurepkcs11_cloud.png','./uploaded/documents/66/0/20','png',0,NULL,0,0,1),(9,'u1dImPaqo3+qhM4kICFk6tMXPSnms33GOFQFwL7e9AM=','2025-09-04 18:13:47','2025-09-04 18:13:47',NULL,NULL,66,0,20,'Confirmación _ Viva.pdf','./uploaded/documents/66/0/20','pdf',0,NULL,0,0,1),(10,'thV6h7lh9hIVj8l7JoNdn2eizvhRdoG/AVIRuZX3Qcw=','2025-09-04 18:13:47','2025-09-04 18:13:47',NULL,NULL,66,0,20,'CV2024_eng.pdf','./uploaded/documents/66/0/20','pdf',0,NULL,0,0,1),(11,'Ei7mgvSAEbt3UB/m1Nfw5YCtxII4+y822Kl6I81ndh4=','2025-09-04 18:14:43','2025-09-04 18:14:43',NULL,NULL,66,0,20,'Carta de Colaboración Remunerada colaboradores  Julio.docx','./uploaded/documents/66/20/templates','docx',0,NULL,0,0,1);
/*!40000 ALTER TABLE `documents` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `documentstatus`
--

DROP TABLE IF EXISTS `documentstatus`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `documentstatus` (
  `idDstat` bigint NOT NULL AUTO_INCREMENT,
  `idDocument` bigint NOT NULL,
  `lastModifiedDstat` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`idDstat`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `documentstatus`
--

LOCK TABLES `documentstatus` WRITE;
/*!40000 ALTER TABLE `documentstatus` DISABLE KEYS */;
/*!40000 ALTER TABLE `documentstatus` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `emails`
--

DROP TABLE IF EXISTS `emails`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `emails` (
  `idEmail` int NOT NULL AUTO_INCREMENT,
  `idUserSender_fk` int NOT NULL,
  `reason` varchar(1024) NOT NULL,
  `sentAtEmail` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `subjectEmail` varchar(512) DEFAULT NULL,
  `bodyEmail` text NOT NULL,
  `statusEmail` int NOT NULL DEFAULT '1',
  `receivedAtEmail` datetime DEFAULT NULL,
  `idAppSource_fk` int NOT NULL,
  PRIMARY KEY (`idEmail`)
) ENGINE=InnoDB AUTO_INCREMENT=34 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `emails`
--

LOCK TABLES `emails` WRITE;
/*!40000 ALTER TABLE `emails` DISABLE KEYS */;
INSERT INTO `emails` VALUES (5,1,'Esto es una pruebota','2025-07-29 02:07:58','thelegendofmax19@gmail.com','From: thelegendofmax19@gmail.com\r\nTo: thelegendofmax19@gmail.com\r\nSubject: Esto es una pruebota\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<h1>la concha su madre</h1>\n',1,NULL,1),(6,1,'Esto es una pruebota','2025-07-29 02:08:52','thelegendofmax19@gmail.com','Subject: Esto es una pruebota\r\nContent-Type: text/html; charset=\"UTF-8\"\r\nFrom: thelegendofmax19@gmail.com\r\nTo: do.c.001@hotmail.com\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<h1>la concha su madre</h1>\n',1,NULL,1),(7,1,'Esto es una prueba','2025-07-29 02:14:41','thelegendofmax19@gmail.com','Subject: Esto es una prueba\r\nContent-Type: text/html; charset=\"UTF-8\"\r\nFrom: thelegendofmax19@gmail.com\r\nTo: thelegendofmax19@gmail.com\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\nla concha su madre\n',1,NULL,1),(8,1,'Esto es una prueba','2025-07-29 02:15:24','thelegendofmax19@gmail.com','From: thelegendofmax19@gmail.com\r\nTo: thelegendofmax19@gmail.com\r\nSubject: Esto es una prueba\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\nla concha su madre\n',1,NULL,10),(9,1,'hola','2025-08-12 02:26:57','thelegendofmax19@gmail.com','Content-Type: text/html; charset=\"UTF-8\"\r\nFrom: thelegendofmax19@gmail.com\r\nTo: do.c.001@hotmail.com\r\nSubject: hola\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\nabcd\n',1,NULL,3),(10,1,'hola','2025-08-12 16:27:54','thelegendofmax19@gmail.com','Content-Type: text/html; charset=\"UTF-8\"\r\nFrom: thelegendofmax19@gmail.com\r\nTo: do.c.001@hotmail.com\r\nSubject: hola\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\nabcd\n',1,NULL,3),(11,1,'Hola desde Go','2025-08-12 18:05:00','thelegendofmax19@gmail.com','From: thelegendofmax19@gmail.com\r\nTo: contacto.somostec@gmail.com\r\nSubject: Hola desde Go\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\nEste es el contenido del mensaje\n',1,NULL,3),(12,1,'¡Bienvenido a Signforce! Correo de verificación SIG250109PBQ','2025-08-12 18:49:01','thelegendofmax19@gmail.com','Content-Type: text/html; charset=\"UTF-8\"\r\nFrom: thelegendofmax19@gmail.com\r\nTo: contacto.somostec@gmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación SIG250109PBQ\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n        body {\r\n            font-family: \'Segoe UI\', Roboto, Oxygen, Ubuntu, Cantarell, \'Open Sans\', \'Helvetica Neue\', sans-serif;\r\n            line-height: 1.6;\r\n            color: #333;\r\n            max-width: 600px;\r\n            margin: 0 auto;\r\n            padding: 20px;\r\n            background-color: #f9f9f9;\r\n        }\r\n        .container {\r\n            background-color: #ffffff;\r\n            border-radius: 12px;\r\n            padding: 30px;\r\n            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);\r\n        }\r\n        .header {\r\n            text-align: center;\r\n            margin-bottom: 25px;\r\n        }\r\n        .logo {\r\n            max-width: 150px;\r\n            margin-bottom: 20px;\r\n        }\r\n        h1 {\r\n            color: #2c3e50;\r\n            font-size: 24px;\r\n            margin-bottom: 20px;\r\n            font-weight: 600;\r\n        }\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 16px;\r\n            color: #555;\r\n        }\r\n        .cta-button {\r\n            display: inline-block;\r\n            background-color: #4285f4;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 12px 24px;\r\n            border-radius: 8px;\r\n            font-weight: 500;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: background-color 0.3s;\r\n        }\r\n        .cta-button:hover {\r\n            background-color: #3367d6;\r\n        }\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 14px;\r\n            color: #999;\r\n            text-align: center;\r\n        }\r\n        .divider {\r\n            height: 1px;\r\n            background-color: #eee;\r\n            margin: 25px 0;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <p>Este es tu nombre de usuario que podrás cambiar al completar tu registro: contacto.somostec@gmail.com</p>\r\n        \r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p>Para comenzar nuestro viaje juntos, entra a la siguiente liga donde te llevaremos paso a paso por el proceo de registro:</p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"{URL_COMPLETAR_REGISTRO}\" class=\"cta-button\">¡Completa tu registro aqui!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta el 2025-09-11 18:49:00. Si no has solicitado crear una cuenta, puedes ignorar este mensaje.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en responder a este correo o contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #4285f4; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #4285f4; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(13,1,'¡Bienvenido a Signforce! Correo de verificación SIG250109PBR','2025-08-12 18:55:57','thelegendofmax19@gmail.com','Subject: ¡Bienvenido a Signforce! Correo de verificación SIG250109PBR\r\nContent-Type: text/html; charset=\"UTF-8\"\r\nFrom: thelegendofmax19@gmail.com\r\nTo: contacto.somostec@gmail.com\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n        body {\r\n            font-family: \'Segoe UI\', Roboto, Oxygen, Ubuntu, Cantarell, \'Open Sans\', \'Helvetica Neue\', sans-serif;\r\n            line-height: 1.6;\r\n            color: #333;\r\n            max-width: 600px;\r\n            margin: 0 auto;\r\n            padding: 20px;\r\n            background-color: #f9f9f9;\r\n        }\r\n        .container {\r\n            background-color: #ffffff;\r\n            border-radius: 12px;\r\n            padding: 30px;\r\n            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);\r\n        }\r\n        .header {\r\n            text-align: center;\r\n            margin-bottom: 25px;\r\n        }\r\n        .logo {\r\n            max-width: 150px;\r\n            margin-bottom: 20px;\r\n        }\r\n        h1 {\r\n            color: #2c3e50;\r\n            font-size: 24px;\r\n            margin-bottom: 20px;\r\n            font-weight: 600;\r\n        }\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 16px;\r\n            color: #555;\r\n        }\r\n        .cta-button {\r\n            display: inline-block;\r\n            background-color: #4285f4;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 12px 24px;\r\n            border-radius: 8px;\r\n            font-weight: 500;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: background-color 0.3s;\r\n        }\r\n        .cta-button:hover {\r\n            background-color: #3367d6;\r\n        }\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 14px;\r\n            color: #999;\r\n            text-align: center;\r\n        }\r\n        .divider {\r\n            height: 1px;\r\n            background-color: #eee;\r\n            margin: 25px 0;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <p>Este es tu nombre de usuario que podrás cambiar al completar tu registro: contacto.somostec@gmail.com</p>\r\n        \r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p>Para comenzar nuestro viaje juntos, entra a la siguiente liga donde te llevaremos paso a paso por el proceo de registro:</p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"{URL_COMPLETAR_REGISTRO}\" class=\"cta-button\">¡Completa tu registro aqui!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta el 2025-09-11 18:55:56. Si no has solicitado crear una cuenta, puedes ignorar este mensaje.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en responder a este correo o contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #4285f4; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #4285f4; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(14,1,'¡Bienvenido a Signforce! Correo de verificación SIG250109PBS','2025-08-12 18:57:32','thelegendofmax19@gmail.com','From: thelegendofmax19@gmail.com\r\nTo: contacto.somostec@gmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación SIG250109PBS\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n        body {\r\n            font-family: \'Segoe UI\', Roboto, Oxygen, Ubuntu, Cantarell, \'Open Sans\', \'Helvetica Neue\', sans-serif;\r\n            line-height: 1.6;\r\n            color: #333;\r\n            max-width: 600px;\r\n            margin: 0 auto;\r\n            padding: 20px;\r\n            background-color: #f9f9f9;\r\n        }\r\n        .container {\r\n            background-color: #ffffff;\r\n            border-radius: 12px;\r\n            padding: 30px;\r\n            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);\r\n        }\r\n        .header {\r\n            text-align: center;\r\n            margin-bottom: 25px;\r\n        }\r\n        .logo {\r\n            max-width: 150px;\r\n            margin-bottom: 20px;\r\n        }\r\n        h1 {\r\n            color: #2c3e50;\r\n            font-size: 24px;\r\n            margin-bottom: 20px;\r\n            font-weight: 600;\r\n        }\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 16px;\r\n            color: #555;\r\n        }\r\n        .cta-button {\r\n            display: inline-block;\r\n            background-color: #4285f4;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 12px 24px;\r\n            border-radius: 8px;\r\n            font-weight: 500;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: background-color 0.3s;\r\n        }\r\n        .cta-button:hover {\r\n            background-color: #3367d6;\r\n        }\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 14px;\r\n            color: #999;\r\n            text-align: center;\r\n        }\r\n        .divider {\r\n            height: 1px;\r\n            background-color: #eee;\r\n            margin: 25px 0;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <p>Este es tu nombre de usuario que podrás cambiar al completar tu registro: contacto.somostec@gmail.com</p>\r\n        \r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p>Para comenzar nuestro viaje juntos, entra a la siguiente liga donde te llevaremos paso a paso por el proceo de registro:</p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"{URL_COMPLETAR_REGISTRO}\" class=\"cta-button\">¡Completa tu registro aqui!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta el 2025-09-11 18:57:31. Si no has solicitado crear una cuenta, puedes ignorar este mensaje.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en responder a este correo o contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #4285f4; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #4285f4; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(15,1,'¡Bienvenido a Signforce! Correo de verificación SIG250109PBT','2025-08-12 19:06:37','thelegendofmax19@gmail.com','From: thelegendofmax19@gmail.com\r\nTo: contacto.somostec@gmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación SIG250109PBT\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n        body {\r\n            font-family: \'Segoe UI\', Roboto, Oxygen, Ubuntu, Cantarell, \'Open Sans\', \'Helvetica Neue\', sans-serif;\r\n            line-height: 1.6;\r\n            color: #333;\r\n            max-width: 600px;\r\n            margin: 0 auto;\r\n            padding: 20px;\r\n            background-color: #f9f9f9;\r\n        }\r\n        .container {\r\n            background-color: #ffffff;\r\n            border-radius: 12px;\r\n            padding: 30px;\r\n            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);\r\n        }\r\n        .header {\r\n            text-align: center;\r\n            margin-bottom: 25px;\r\n        }\r\n        .logo {\r\n            max-width: 150px;\r\n            margin-bottom: 20px;\r\n        }\r\n        h1 {\r\n            color: #2c3e50;\r\n            font-size: 24px;\r\n            margin-bottom: 20px;\r\n            font-weight: 600;\r\n        }\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 16px;\r\n            color: #555;\r\n        }\r\n        .cta-button {\r\n            display: inline-block;\r\n            background-color: #4285f4;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 12px 24px;\r\n            border-radius: 8px;\r\n            font-weight: 500;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: background-color 0.3s;\r\n        }\r\n        .cta-button:hover {\r\n            background-color: #3367d6;\r\n        }\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 14px;\r\n            color: #999;\r\n            text-align: center;\r\n        }\r\n        .divider {\r\n            height: 1px;\r\n            background-color: #eee;\r\n            margin: 25px 0;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <p>Este es tu nombre de usuario que podrás cambiar al completar tu registro: contacto.somostec@gmail.com</p>\r\n        \r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p>Para comenzar nuestro viaje juntos, entra a la siguiente liga donde te llevaremos paso a paso por el proceo de registro:</p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"{URL_COMPLETAR_REGISTRO}\" class=\"cta-button\">¡Completa tu registro aqui!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta el 2025-09-11 19:06:36. Si no has solicitado crear una cuenta, puedes ignorar este mensaje.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en responder a este correo o contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #4285f4; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #4285f4; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(16,1,'¡Bienvenido a Signforce! Correo de verificación SIG250109PBU','2025-08-12 19:15:56','thelegendofmax19@gmail.com','From: thelegendofmax19@gmail.com\r\nTo: contacto.somostec@gmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación SIG250109PBU\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n        body {\r\n            font-family: \'Segoe UI\', Roboto, Oxygen, Ubuntu, Cantarell, \'Open Sans\', \'Helvetica Neue\', sans-serif;\r\n            line-height: 1.6;\r\n            color: #333;\r\n            max-width: 600px;\r\n            margin: 0 auto;\r\n            padding: 20px;\r\n            background-color: #f9f9f9;\r\n        }\r\n        .container {\r\n            background-color: #ffffff;\r\n            border-radius: 12px;\r\n            padding: 30px;\r\n            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);\r\n        }\r\n        .header {\r\n            text-align: center;\r\n            margin-bottom: 25px;\r\n        }\r\n        .logo {\r\n            max-width: 150px;\r\n            margin-bottom: 20px;\r\n        }\r\n        h1 {\r\n            color: #2c3e50;\r\n            font-size: 24px;\r\n            margin-bottom: 20px;\r\n            font-weight: 600;\r\n        }\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 16px;\r\n            color: #555;\r\n        }\r\n        .cta-button {\r\n            display: inline-block;\r\n            background-color: #4285f4;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 12px 24px;\r\n            border-radius: 8px;\r\n            font-weight: 500;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: background-color 0.3s;\r\n        }\r\n        .cta-button:hover {\r\n            background-color: #3367d6;\r\n        }\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 14px;\r\n            color: #999;\r\n            text-align: center;\r\n        }\r\n        .divider {\r\n            height: 1px;\r\n            background-color: #eee;\r\n            margin: 25px 0;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <p>Este es tu nombre de usuario que podrás cambiar al completar tu registro: contacto.somostec@gmail.com</p>\r\n        \r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p>Para comenzar nuestro viaje juntos, entra a la siguiente liga donde te llevaremos paso a paso por el proceo de registro:</p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"http://192.168.1.68:8000/validation\" class=\"cta-button\" style=\"display:inline-block;background-color:#007bff;color:#fff;padding:10px 20px;text-decoration:none;border-radius:5px;\">¡Completa tu registro aqui!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta el 2025-09-11 19:15:55. Si no has solicitado crear una cuenta, puedes ignorar este mensaje.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en responder a este correo o contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #4285f4; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #4285f4; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(17,1,'¡Bienvenido a Signforce! Correo de verificación SIG250109PBA','2025-08-13 16:17:58','thelegendofmax19@gmail.com','From: thelegendofmax19@gmail.com\r\nTo: do.c.001@hotmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación SIG250109PBA\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n        body {\r\n            font-family: \'Segoe UI\', Roboto, Oxygen, Ubuntu, Cantarell, \'Open Sans\', \'Helvetica Neue\', sans-serif;\r\n            line-height: 1.6;\r\n            color: #333;\r\n            max-width: 600px;\r\n            margin: 0 auto;\r\n            padding: 20px;\r\n            background-color: #f9f9f9;\r\n        }\r\n        .container {\r\n            background-color: #ffffff;\r\n            border-radius: 12px;\r\n            padding: 30px;\r\n            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);\r\n        }\r\n        .header {\r\n            text-align: center;\r\n            margin-bottom: 25px;\r\n        }\r\n        .logo {\r\n            max-width: 150px;\r\n            margin-bottom: 20px;\r\n        }\r\n        h1 {\r\n            color: #2c3e50;\r\n            font-size: 24px;\r\n            margin-bottom: 20px;\r\n            font-weight: 600;\r\n        }\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 16px;\r\n            color: #555;\r\n        }\r\n        .cta-button {\r\n            display: inline-block;\r\n            background-color: #4285f4;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 12px 24px;\r\n            border-radius: 8px;\r\n            font-weight: 500;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: background-color 0.3s;\r\n        }\r\n        .cta-button:hover {\r\n            background-color: #3367d6;\r\n        }\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 14px;\r\n            color: #999;\r\n            text-align: center;\r\n        }\r\n        .divider {\r\n            height: 1px;\r\n            background-color: #eee;\r\n            margin: 25px 0;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <p>Este es tu nombre de usuario que podrás cambiar al completar tu registro: do.c.001@hotmail.com</p>\r\n        \r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p>Para comenzar nuestro viaje juntos, entra a la siguiente liga donde te llevaremos paso a paso por el proceo de registro:</p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"http://192.168.1.68:8000/validation\" class=\"cta-button\">¡Completa tu registro aqui!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta el 2025-09-12 16:17:57. Si no has solicitado crear una cuenta, puedes ignorar este mensaje.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en responder a este correo o contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #4285f4; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #4285f4; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(18,1,'¡Bienvenido a Signforce! Correo de verificación SIG250109PBC','2025-08-13 17:29:13','thelegendofmax19@gmail.com','From: thelegendofmax19@gmail.com\r\nTo: do.c.001@hotmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación SIG250109PBC\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n        body {\r\n            font-family: \'Segoe UI\', Roboto, Oxygen, Ubuntu, Cantarell, \'Open Sans\', \'Helvetica Neue\', sans-serif;\r\n            line-height: 1.6;\r\n            color: #333;\r\n            max-width: 600px;\r\n            margin: 0 auto;\r\n            padding: 20px;\r\n            background-color: #f9f9f9;\r\n        }\r\n        .container {\r\n            background-color: #ffffff;\r\n            border-radius: 12px;\r\n            padding: 30px;\r\n            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);\r\n        }\r\n        .header {\r\n            text-align: center;\r\n            margin-bottom: 25px;\r\n        }\r\n        .logo {\r\n            max-width: 150px;\r\n            margin-bottom: 20px;\r\n        }\r\n        h1 {\r\n            color: #2c3e50;\r\n            font-size: 24px;\r\n            margin-bottom: 20px;\r\n            font-weight: 600;\r\n        }\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 16px;\r\n            color: #555;\r\n        }\r\n        .cta-button {\r\n            display: inline-block;\r\n            background-color: #4285f4;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 12px 24px;\r\n            border-radius: 8px;\r\n            font-weight: 500;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: background-color 0.3s;\r\n        }\r\n        .cta-button:hover {\r\n            background-color: #3367d6;\r\n        }\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 14px;\r\n            color: #999;\r\n            text-align: center;\r\n        }\r\n        .divider {\r\n            height: 1px;\r\n            background-color: #eee;\r\n            margin: 25px 0;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <p>Este es tu nombre de usuario que podrás cambiar al completar tu registro: do.c.001@hotmail.com</p>\r\n        \r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p>Para comenzar nuestro viaje juntos, entra a la siguiente liga donde te llevaremos paso a paso por el proceo de registro:</p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"http://192.168.1.68:8000/login\" class=\"cta-button\">¡Completa tu registro aqui!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta el 2025-09-12 17:29:11. Si no has solicitado crear una cuenta, puedes ignorar este mensaje.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en responder a este correo o contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #4285f4; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #4285f4; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(19,1,'¡Bienvenido a Signforce! Correo de verificación SIG250109PBD','2025-08-13 17:36:03','thelegendofmax19@gmail.com','From: thelegendofmax19@gmail.com\r\nTo: do.c.001@hotmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación SIG250109PBD\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n        body {\r\n            font-family: \'Segoe UI\', Roboto, Oxygen, Ubuntu, Cantarell, \'Open Sans\', \'Helvetica Neue\', sans-serif;\r\n            line-height: 1.6;\r\n            color: #333;\r\n            max-width: 600px;\r\n            margin: 0 auto;\r\n            padding: 20px;\r\n            background-color: #f9f9f9;\r\n        }\r\n        .container {\r\n            background-color: #ffffff;\r\n            border-radius: 12px;\r\n            padding: 30px;\r\n            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);\r\n        }\r\n        .header {\r\n            text-align: center;\r\n            margin-bottom: 25px;\r\n        }\r\n        .logo {\r\n            max-width: 150px;\r\n            margin-bottom: 20px;\r\n        }\r\n        h1 {\r\n            color: #2c3e50;\r\n            font-size: 24px;\r\n            margin-bottom: 20px;\r\n            font-weight: 600;\r\n        }\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 16px;\r\n            color: #555;\r\n        }\r\n        .cta-button {\r\n            display: inline-block;\r\n            background-color: #4285f4;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 12px 24px;\r\n            border-radius: 8px;\r\n            font-weight: 500;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: background-color 0.3s;\r\n        }\r\n        .cta-button:hover {\r\n            background-color: #3367d6;\r\n        }\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 14px;\r\n            color: #999;\r\n            text-align: center;\r\n        }\r\n        .divider {\r\n            height: 1px;\r\n            background-color: #eee;\r\n            margin: 25px 0;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <p>Este es tu nombre de usuario que podrás cambiar al completar tu registro: do.c.001@hotmail.com</p>\r\n        \r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p>Para comenzar nuestro viaje juntos, entra a la siguiente liga donde te llevaremos paso a paso por el proceo de registro:</p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"http://192.168.1.68:8000/login\" class=\"cta-button\">¡Completa tu registro aqui!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta el 2025-09-12 17:36:02. Si no has solicitado crear una cuenta, puedes ignorar este mensaje.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en responder a este correo o contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #4285f4; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #4285f4; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(20,1,'¡Bienvenido a Signforce! Correo de verificación SIG250109PBE','2025-08-13 17:40:16','thelegendofmax19@gmail.com','From: thelegendofmax19@gmail.com\r\nTo: do.c.001@hotmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación SIG250109PBE\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n        body {\r\n            font-family: \'Segoe UI\', Roboto, Oxygen, Ubuntu, Cantarell, \'Open Sans\', \'Helvetica Neue\', sans-serif;\r\n            line-height: 1.6;\r\n            color: #333;\r\n            max-width: 600px;\r\n            margin: 0 auto;\r\n            padding: 20px;\r\n            background-color: #f9f9f9;\r\n        }\r\n        .container {\r\n            background-color: #ffffff;\r\n            border-radius: 12px;\r\n            padding: 30px;\r\n            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);\r\n        }\r\n        .header {\r\n            text-align: center;\r\n            margin-bottom: 25px;\r\n        }\r\n        .logo {\r\n            max-width: 150px;\r\n            margin-bottom: 20px;\r\n        }\r\n        h1 {\r\n            color: #2c3e50;\r\n            font-size: 24px;\r\n            margin-bottom: 20px;\r\n            font-weight: 600;\r\n        }\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 16px;\r\n            color: #555;\r\n        }\r\n        .cta-button {\r\n            display: inline-block;\r\n            background-color: #4285f4;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 12px 24px;\r\n            border-radius: 8px;\r\n            font-weight: 500;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: background-color 0.3s;\r\n        }\r\n        .cta-button:hover {\r\n            background-color: #3367d6;\r\n        }\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 14px;\r\n            color: #999;\r\n            text-align: center;\r\n        }\r\n        .divider {\r\n            height: 1px;\r\n            background-color: #eee;\r\n            margin: 25px 0;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <p>Este es tu nombre de usuario que podrás cambiar al completar tu registro: do.c.001@hotmail.com</p>\r\n        \r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p>Para comenzar nuestro viaje juntos, entra a la siguiente liga donde te llevaremos paso a paso por el proceo de registro:</p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"http://192.168.1.68:8000/login\" class=\"cta-button\">¡Completa tu registro aqui!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta el 2025-09-12 17:40:15. Si no has solicitado crear una cuenta, puedes ignorar este mensaje.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en responder a este correo o contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #4285f4; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #4285f4; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(21,1,'¡Bienvenido a Signforce! Correo de verificación SIG250109PBB','2025-08-14 01:39:02','thelegendofmax19@gmail.com','From: thelegendofmax19@gmail.com\r\nTo: do.c.001@hotmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación SIG250109PBB\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n        body {\r\n            font-family: \'Segoe UI\', Roboto, Oxygen, Ubuntu, Cantarell, \'Open Sans\', \'Helvetica Neue\', sans-serif;\r\n            line-height: 1.6;\r\n            color: #333;\r\n            max-width: 600px;\r\n            margin: 0 auto;\r\n            padding: 20px;\r\n            background-color: #f9f9f9;\r\n        }\r\n        .container {\r\n            background-color: #ffffff;\r\n            border-radius: 12px;\r\n            padding: 30px;\r\n            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);\r\n        }\r\n        .header {\r\n            text-align: center;\r\n            margin-bottom: 25px;\r\n        }\r\n        .logo {\r\n            max-width: 150px;\r\n            margin-bottom: 20px;\r\n        }\r\n        h1 {\r\n            color: #2c3e50;\r\n            font-size: 24px;\r\n            margin-bottom: 20px;\r\n            font-weight: 600;\r\n        }\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 16px;\r\n            color: #555;\r\n        }\r\n        .cta-button {\r\n            display: inline-block;\r\n            background-color: #4285f4;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 12px 24px;\r\n            border-radius: 8px;\r\n            font-weight: 500;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: background-color 0.3s;\r\n        }\r\n        .cta-button:hover {\r\n            background-color: #3367d6;\r\n        }\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 14px;\r\n            color: #999;\r\n            text-align: center;\r\n        }\r\n        .divider {\r\n            height: 1px;\r\n            background-color: #eee;\r\n            margin: 25px 0;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <p>Este es tu nombre de usuario que podrás cambiar al completar tu registro: do.c.001@hotmail.com</p>\r\n        \r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p>Para comenzar nuestro viaje juntos, entra a la siguiente liga donde te llevaremos paso a paso por el proceo de registro:</p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"http://192.168.1.68:8000/login\" class=\"cta-button\">¡Completa tu registro aqui!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta el 2025-09-13 01:39:01. Si no has solicitado crear una cuenta, puedes ignorar este mensaje.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en responder a este correo o contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #4285f4; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #4285f4; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(22,1,'¡Bienvenido a Signforce! Correo de verificación SIG250109PBD','2025-08-19 00:33:57','thelegendofmax19@gmail.com','From: thelegendofmax19@gmail.com\r\nTo: contacto.somostec@gmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación SIG250109PBD\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n        body {\r\n            font-family: \'Segoe UI\', Roboto, Oxygen, Ubuntu, Cantarell, \'Open Sans\', \'Helvetica Neue\', sans-serif;\r\n            line-height: 1.6;\r\n            color: #333;\r\n            max-width: 600px;\r\n            margin: 0 auto;\r\n            padding: 20px;\r\n            background-color: #f9f9f9;\r\n        }\r\n        .container {\r\n            background-color: #ffffff;\r\n            border-radius: 12px;\r\n            padding: 30px;\r\n            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);\r\n        }\r\n        .header {\r\n            text-align: center;\r\n            margin-bottom: 25px;\r\n        }\r\n        .logo {\r\n            max-width: 150px;\r\n            margin-bottom: 20px;\r\n        }\r\n        h1 {\r\n            color: #2c3e50;\r\n            font-size: 24px;\r\n            margin-bottom: 20px;\r\n            font-weight: 600;\r\n        }\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 16px;\r\n            color: #555;\r\n        }\r\n        .cta-button {\r\n            display: inline-block;\r\n            background-color: #4285f4;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 12px 24px;\r\n            border-radius: 8px;\r\n            font-weight: 500;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: background-color 0.3s;\r\n        }\r\n        .cta-button:hover {\r\n            background-color: #3367d6;\r\n        }\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 14px;\r\n            color: #999;\r\n            text-align: center;\r\n        }\r\n        .divider {\r\n            height: 1px;\r\n            background-color: #eee;\r\n            margin: 25px 0;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <p>Este es tu nombre de usuario que podrás cambiar al completar tu registro: contacto.somostec@gmail.com</p>\r\n        \r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p>Para comenzar nuestro viaje juntos, entra a la siguiente liga donde te llevaremos paso a paso por el proceo de registro:</p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"http://192.168.1.68:8000/login\" class=\"cta-button\">¡Completa tu registro aqui!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta el 2025-09-18 00:33:56. Si no has solicitado crear una cuenta, puedes ignorar este mensaje.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en responder a este correo o contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #4285f4; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #4285f4; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(23,1,'¡Bienvenido a Signforce! Correo de verificación SIG250109PBE','2025-08-19 01:12:42','thelegendofmax19@gmail.com','Content-Type: text/html; charset=\"UTF-8\"\r\nFrom: thelegendofmax19@gmail.com\r\nTo: contacto.somostec@gmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación SIG250109PBE\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n        body {\r\n            font-family: \'Segoe UI\', Roboto, Oxygen, Ubuntu, Cantarell, \'Open Sans\', \'Helvetica Neue\', sans-serif;\r\n            line-height: 1.6;\r\n            color: #333;\r\n            max-width: 600px;\r\n            margin: 0 auto;\r\n            padding: 20px;\r\n            background-color: #f9f9f9;\r\n        }\r\n        .container {\r\n            background-color: #ffffff;\r\n            border-radius: 12px;\r\n            padding: 30px;\r\n            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);\r\n        }\r\n        .header {\r\n            text-align: center;\r\n            margin-bottom: 25px;\r\n        }\r\n        .logo {\r\n            max-width: 150px;\r\n            margin-bottom: 20px;\r\n        }\r\n        h1 {\r\n            color: #2c3e50;\r\n            font-size: 24px;\r\n            margin-bottom: 20px;\r\n            font-weight: 600;\r\n        }\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 16px;\r\n            color: #555;\r\n        }\r\n        .cta-button {\r\n            display: inline-block;\r\n            background-color: #4285f4;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 12px 24px;\r\n            border-radius: 8px;\r\n            font-weight: 500;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: background-color 0.3s;\r\n        }\r\n        .cta-button:hover {\r\n            background-color: #3367d6;\r\n        }\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 14px;\r\n            color: #999;\r\n            text-align: center;\r\n        }\r\n        .divider {\r\n            height: 1px;\r\n            background-color: #eee;\r\n            margin: 25px 0;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <h4>Este es tu nombre de usuario que podrás cambiar al completar tu registro: <strong>contacto.somostec@gmail.com</strong></h4>\r\n        <h4>Tu contraseña temporal es: <strong>{TEMPORAL_PASS}</strong></h4>\r\n        \r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p><strong>Para comenzar nuestro viaje, entra a la siguiente liga donde te llevaremos paso a paso por el proceo de registro:</strong></p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"http://192.168.1.68:8000/login\" class=\"cta-button\">¡Completa tu registro aqui!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta la fecha 2025-09-18 01:12:41. Si este email te resulta desconocido, ignoralo.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en responder a este correo o contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #4285f4; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #4285f4; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(24,1,'¡Bienvenido a Signforce! Correo de verificación SIG250109PBG','2025-08-19 01:14:34','thelegendofmax19@gmail.com','From: thelegendofmax19@gmail.com\r\nTo: contacto.somostec@gmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación SIG250109PBG\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n        body {\r\n            font-family: \'Segoe UI\', Roboto, Oxygen, Ubuntu, Cantarell, \'Open Sans\', \'Helvetica Neue\', sans-serif;\r\n            line-height: 1.6;\r\n            color: #333;\r\n            max-width: 600px;\r\n            margin: 0 auto;\r\n            padding: 20px;\r\n            background-color: #f9f9f9;\r\n        }\r\n        .container {\r\n            background-color: #ffffff;\r\n            border-radius: 12px;\r\n            padding: 30px;\r\n            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);\r\n        }\r\n        .header {\r\n            text-align: center;\r\n            margin-bottom: 25px;\r\n        }\r\n        .logo {\r\n            max-width: 150px;\r\n            margin-bottom: 20px;\r\n        }\r\n        h1 {\r\n            color: #2c3e50;\r\n            font-size: 24px;\r\n            margin-bottom: 20px;\r\n            font-weight: 600;\r\n        }\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 16px;\r\n            color: #555;\r\n        }\r\n        .cta-button {\r\n            display: inline-block;\r\n            background-color: #4285f4;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 12px 24px;\r\n            border-radius: 8px;\r\n            font-weight: 500;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: background-color 0.3s;\r\n        }\r\n        .cta-button:hover {\r\n            background-color: #3367d6;\r\n        }\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 14px;\r\n            color: #999;\r\n            text-align: center;\r\n        }\r\n        .divider {\r\n            height: 1px;\r\n            background-color: #eee;\r\n            margin: 25px 0;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <h4>Este es tu nombre de usuario que podrás cambiar al completar tu registro: <strong>contacto.somostec@gmail.com</strong></h4>\r\n        <h4>Tu contraseña temporal es: <strong>rñUMVNBjZOya</strong></h4>\r\n        \r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p><strong>Para comenzar nuestro viaje, entra a la siguiente liga donde te llevaremos paso a paso por el proceo de registro:</strong></p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"http://192.168.1.68:8000/login\" class=\"cta-button\">¡Completa tu registro aqui!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta la fecha 2025-09-18 01:14:33. Si este email te resulta desconocido, ignoralo.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en responder a este correo o contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #4285f4; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #4285f4; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(25,1,'¡Bienvenido a Signforce! Correo de verificación SIG250109PBH','2025-08-19 01:17:19','thelegendofmax19@gmail.com','To: contacto.somostec@gmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación SIG250109PBH\r\nContent-Type: text/html; charset=\"UTF-8\"\r\nFrom: thelegendofmax19@gmail.com\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n        body {\r\n            font-family: \'Segoe UI\', Roboto, Oxygen, Ubuntu, Cantarell, \'Open Sans\', \'Helvetica Neue\', sans-serif;\r\n            line-height: 1.6;\r\n            color: #333;\r\n            max-width: 600px;\r\n            margin: 0 auto;\r\n            padding: 20px;\r\n            background-color: #f9f9f9;\r\n        }\r\n        .container {\r\n            background-color: #ffffff;\r\n            border-radius: 12px;\r\n            padding: 30px;\r\n            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);\r\n        }\r\n        .header {\r\n            text-align: center;\r\n            margin-bottom: 25px;\r\n        }\r\n        .logo {\r\n            max-width: 150px;\r\n            margin-bottom: 20px;\r\n        }\r\n        h1 {\r\n            color: #2c3e50;\r\n            font-size: 24px;\r\n            margin-bottom: 20px;\r\n            font-weight: 600;\r\n        }\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 16px;\r\n            color: #555;\r\n        }\r\n        .cta-button {\r\n            display: inline-block;\r\n            background-color: #4285f4;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 12px 24px;\r\n            border-radius: 8px;\r\n            font-weight: 500;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: background-color 0.3s;\r\n        }\r\n        .cta-button:hover {\r\n            background-color: #3367d6;\r\n        }\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 14px;\r\n            color: #999;\r\n            text-align: center;\r\n        }\r\n        .divider {\r\n            height: 1px;\r\n            background-color: #eee;\r\n            margin: 25px 0;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <h4>Este es tu nombre de usuario que podrás cambiar al completar tu registro: <strong>contacto.somostec@gmail.com</strong></h4>\r\n        <h4>Tu contraseña temporal es: <strong>ZPp:_$NRiFIk</strong></h4>\r\n        \r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p><strong>Para comenzar nuestro viaje, entra a la siguiente liga donde te llevaremos paso a paso por el proceo de registro:</strong></p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"http://192.168.1.68:8000/login\" class=\"cta-button\">¡Completa tu registro aqui!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta la fecha 2025-09-18 01:17:17. Si este email te resulta desconocido, ignoralo.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en responder a este correo o contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #4285f4; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #4285f4; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(26,1,'¡Bienvenido a Signforce! Correo de verificación SIG250109PBI','2025-08-19 01:45:45','thelegendofmax19@gmail.com','To: contacto.somostec@gmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación SIG250109PBI\r\nContent-Type: text/html; charset=\"UTF-8\"\r\nFrom: thelegendofmax19@gmail.com\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n        body {\r\n            font-family: \'Segoe UI\', Roboto, Oxygen, Ubuntu, Cantarell, \'Open Sans\', \'Helvetica Neue\', sans-serif;\r\n            line-height: 1.6;\r\n            color: #333;\r\n            max-width: 600px;\r\n            margin: 0 auto;\r\n            padding: 20px;\r\n            background-color: #f9f9f9;\r\n        }\r\n        .container {\r\n            background-color: #ffffff;\r\n            border-radius: 12px;\r\n            padding: 30px;\r\n            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);\r\n        }\r\n        .header {\r\n            text-align: center;\r\n            margin-bottom: 25px;\r\n        }\r\n        .logo {\r\n            max-width: 150px;\r\n            margin-bottom: 20px;\r\n        }\r\n        h1 {\r\n            color: #2c3e50;\r\n            font-size: 24px;\r\n            margin-bottom: 20px;\r\n            font-weight: 600;\r\n        }\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 16px;\r\n            color: #555;\r\n        }\r\n        .cta-button {\r\n            display: inline-block;\r\n            background-color: #4285f4;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 12px 24px;\r\n            border-radius: 8px;\r\n            font-weight: 500;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: background-color 0.3s;\r\n        }\r\n        .cta-button:hover {\r\n            background-color: #3367d6;\r\n        }\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 14px;\r\n            color: #999;\r\n            text-align: center;\r\n        }\r\n        .divider {\r\n            height: 1px;\r\n            background-color: #eee;\r\n            margin: 25px 0;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <h4>Este es tu nombre de usuario que podrás cambiar al completar tu registro: <strong>contacto.somostec@gmail.com</strong></h4>\r\n        <h4>Tu contraseña es: <strong>)tB_-ewS&$;A</strong> ¡La puedes cambiar cuando quieras!</h4>\r\n        \r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p><strong>Para comenzar nuestro viaje, entra a la siguiente liga donde te llevaremos paso a paso por el proceo de registro:</strong></p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"http://192.168.1.68:8000/login\" class=\"cta-button\">¡Completa tu registro aqui!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta la fecha 2025-09-18 01:45:44. Si este email te resulta desconocido, ignoralo.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en responder a este correo o contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #4285f4; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #4285f4; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(27,1,'¡Bienvenido a Signforce! Correo de verificación SIG250109PBJ','2025-08-19 01:50:06','thelegendofmax19@gmail.com','From: thelegendofmax19@gmail.com\r\nTo: contacto.somostec@gmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación SIG250109PBJ\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n        body {\r\n            font-family: \'Segoe UI\', Roboto, Oxygen, Ubuntu, Cantarell, \'Open Sans\', \'Helvetica Neue\', sans-serif;\r\n            line-height: 1.6;\r\n            color: #333;\r\n            max-width: 600px;\r\n            margin: 0 auto;\r\n            padding: 20px;\r\n            background-color: #f9f9f9;\r\n        }\r\n        .container {\r\n            background-color: #ffffff;\r\n            border-radius: 12px;\r\n            padding: 30px;\r\n            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);\r\n        }\r\n        .header {\r\n            text-align: center;\r\n            margin-bottom: 25px;\r\n        }\r\n        .logo {\r\n            max-width: 150px;\r\n            margin-bottom: 20px;\r\n        }\r\n        h1 {\r\n            color: #2c3e50;\r\n            font-size: 24px;\r\n            margin-bottom: 20px;\r\n            font-weight: 600;\r\n        }\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 16px;\r\n            color: #555;\r\n        }\r\n        .cta-button {\r\n            display: inline-block;\r\n            background-color: #4285f4;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 12px 24px;\r\n            border-radius: 8px;\r\n            font-weight: 500;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: background-color 0.3s;\r\n        }\r\n        .cta-button:hover {\r\n            background-color: #3367d6;\r\n        }\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 14px;\r\n            color: #999;\r\n            text-align: center;\r\n        }\r\n        .divider {\r\n            height: 1px;\r\n            background-color: #eee;\r\n            margin: 25px 0;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <h4>Este es tu nombre de usuario que podrás cambiar al completar tu registro: <strong>contacto.somostec@gmail.com</strong></h4>\r\n        <h4>Tu contraseña es: <strong>JAPz+TUQ;qSo</strong> ¡La puedes cambiar cuando quieras!</h4>\r\n        \r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p><strong>Para comenzar nuestro viaje, entra a la siguiente liga donde te llevaremos paso a paso por el proceo de registro:</strong></p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"http://192.168.1.68:8000/login\" class=\"cta-button\">¡Inicia sesión y completa tu registro!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta la fecha 2025-09-18 01:50:05. Si este email te resulta desconocido, ignoralo.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en responder a este correo o contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #4285f4; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #4285f4; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(28,1,'¡Bienvenido a Signforce! Correo de verificación SIG250109PBK','2025-08-19 01:51:29','thelegendofmax19@gmail.com','To: contacto.somostec@gmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación SIG250109PBK\r\nContent-Type: text/html; charset=\"UTF-8\"\r\nFrom: thelegendofmax19@gmail.com\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n        body {\r\n            font-family: \'Segoe UI\', Roboto, Oxygen, Ubuntu, Cantarell, \'Open Sans\', \'Helvetica Neue\', sans-serif;\r\n            line-height: 1.6;\r\n            color: #333;\r\n            max-width: 600px;\r\n            margin: 0 auto;\r\n            padding: 20px;\r\n            background-color: #f9f9f9;\r\n        }\r\n        .container {\r\n            background-color: #ffffff;\r\n            border-radius: 12px;\r\n            padding: 30px;\r\n            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);\r\n        }\r\n        .header {\r\n            text-align: center;\r\n            margin-bottom: 25px;\r\n        }\r\n        .logo {\r\n            max-width: 150px;\r\n            margin-bottom: 20px;\r\n        }\r\n        h1 {\r\n            color: #2c3e50;\r\n            font-size: 24px;\r\n            margin-bottom: 20px;\r\n            font-weight: 600;\r\n        }\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 16px;\r\n            color: #555;\r\n        }\r\n        .cta-button {\r\n            display: inline-block;\r\n            background-color: #4285f4;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 12px 24px;\r\n            border-radius: 8px;\r\n            font-weight: 500;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: background-color 0.3s;\r\n        }\r\n        .cta-button:hover {\r\n            background-color: #3367d6;\r\n        }\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 14px;\r\n            color: #999;\r\n            text-align: center;\r\n        }\r\n        .divider {\r\n            height: 1px;\r\n            background-color: #eee;\r\n            margin: 25px 0;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <h4>Este es tu nombre de usuario que podrás cambiar al completar tu registro: <strong>contacto.somostec@gmail.com</strong></h4>\r\n        <h4>Tu contraseña es: <strong>W-bjZQm)no:*</strong> ¡La puedes cambiar cuando quieras!</h4>\r\n        \r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p><strong>Para comenzar nuestro viaje, entra a la siguiente liga donde te llevaremos paso a paso por el proceo de registro:</strong></p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"http://192.168.1.68:8000/login\" class=\"cta-button\">¡Inicia sesión y completa tu registro!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta la fecha 2025-09-18 01:51:28. Si este email te resulta desconocido, ignoralo.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en responder a este correo o contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #4285f4; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #4285f4; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(29,15,'¡Bienvenido a Signforce! Correo de verificación SIG250109PBA','2025-08-19 16:16:37','thelegendofmax19@gmail.com','Subject: ¡Bienvenido a Signforce! Correo de verificación SIG250109PBA\r\nContent-Type: text/html; charset=\"UTF-8\"\r\nFrom: thelegendofmax19@gmail.com\r\nTo: do.c.001@hotmail.com\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n        body {\r\n            font-family: \'Segoe UI\', Roboto, Oxygen, Ubuntu, Cantarell, \'Open Sans\', \'Helvetica Neue\', sans-serif;\r\n            line-height: 1.6;\r\n            color: #333;\r\n            max-width: 600px;\r\n            margin: 0 auto;\r\n            padding: 20px;\r\n            background-color: #f9f9f9;\r\n        }\r\n        .container {\r\n            background-color: #ffffff;\r\n            border-radius: 12px;\r\n            padding: 30px;\r\n            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);\r\n        }\r\n        .header {\r\n            text-align: center;\r\n            margin-bottom: 25px;\r\n        }\r\n        .logo {\r\n            max-width: 150px;\r\n            margin-bottom: 20px;\r\n        }\r\n        h1 {\r\n            color: #2c3e50;\r\n            font-size: 24px;\r\n            margin-bottom: 20px;\r\n            font-weight: 600;\r\n        }\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 16px;\r\n            color: #555;\r\n        }\r\n        .cta-button {\r\n            display: inline-block;\r\n            background-color: #4285f4;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 12px 24px;\r\n            border-radius: 8px;\r\n            font-weight: 500;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: background-color 0.3s;\r\n        }\r\n        .cta-button:hover {\r\n            background-color: #3367d6;\r\n        }\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 14px;\r\n            color: #999;\r\n            text-align: center;\r\n        }\r\n        .divider {\r\n            height: 1px;\r\n            background-color: #eee;\r\n            margin: 25px 0;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <h4>Este es tu nombre de usuario que podrás cambiar al completar tu registro: <strong>do.c.001@hotmail.com</strong></h4>\r\n        <h4>Tu contraseña es: <strong>N;*QkWa/de)&</strong> ¡La puedes cambiar cuando quieras!</h4>\r\n        \r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p><strong>Para comenzar nuestro viaje, entra a la siguiente liga donde te llevaremos paso a paso por el proceo de registro:</strong></p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"http://localhost:8000/login\" class=\"cta-button\">¡Inicia sesión y completa tu registro!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta la fecha 2025-09-18 16:16:36. Si este email te resulta desconocido, ignoralo.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en responder a este correo o contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #4285f4; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #4285f4; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(30,16,'¡Bienvenido a Signforce! Correo de verificación JSTFRTS123456','2025-08-19 19:39:29','thelegendofmax19@gmail.com','Content-Type: text/html; charset=\"UTF-8\"\r\nFrom: thelegendofmax19@gmail.com\r\nTo: contacto.somostec@gmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación JSTFRTS123456\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n       body {\r\n            background: #F3E2D4;\r\n            min-height: 100vh;\r\n            font-family: Inter, \"Helvetica\", sans-serif;\r\n            margin: 0;\r\n            display: flex;\r\n            align-items: center;\r\n            justify-content: center;\r\n            padding: 0px;\r\n            color: #eaeaea;\r\n        }\r\n\r\n        .container {\r\n            background: rgb(24, 1, 48);\r\n            border-radius: 16px;\r\n            padding: 20px;\r\n            max-width: 800px;\r\n            width: 100%;\r\n            box-shadow: \r\n                0 8px 32px rgba(0, 0, 0, 0.6),\r\n                0 0 20px rgba(200, 90, 255, 0.25);\r\n            border: 1px solid rgba(255, 255, 255, 0.08);\r\n            text-align: center;\r\n        }\r\n\r\n        .header {\r\n            margin-bottom: 25px;\r\n        }\r\n\r\n        .logo {\r\n            max-width: 120px;\r\n            margin-bottom: 15px;\r\n        }\r\n\r\n        h1, h2, h3, h4 {\r\n            color: #ffffff;\r\n            margin-bottom: 15px;\r\n            line-height: 1;\r\n            padding: 0;\r\n            letter-spacing: 0.5px;\r\n            word-spacing: 1px;\r\n            line-height: 1.8;\r\n        }\r\n        h1 {\r\n            background:  #35004366;\r\n        }\r\n\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 1rem;\r\n            line-height: 1.5;\r\n            color: #dcdcdc;\r\n        }\r\n\r\n        .cta-button {\r\n            display: inline-block;\r\n            background: #410041;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 14px 30px;\r\n            border-radius: 10px;\r\n            font-weight: 600;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: all 0.3s ease;\r\n            letter-spacing: 0.5px;\r\n            box-shadow: 0 4px 12px rgba(0,0,0,0.35);\r\n        }\r\n\r\n        .cta-button:hover {\r\n            transform: translateY(-3px);\r\n            background: #660066;\r\n            box-shadow: 0 8px 25px rgba(0,0,0,0.4);\r\n        }\r\n\r\n        .divider {\r\n            height: 1px;\r\n            background: linear-gradient(90deg, transparent, #afafaf, transparent);\r\n            margin: 30px 0;\r\n        }\r\n\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 0.85rem;\r\n            color: #999;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <h4>Este es tu nombre de usuario que podrás acceder a la plataforma: <strong>contacto.somostec@gmail.com</strong></h4>\r\n        <h4>Tu contraseña es: <strong>Kvt:/iYPa;ñy</strong> ¡La puedes cambiar cuando quieras!</h4>\r\n        <div class=\"divider\"></div>\r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p><strong>Comienza nuestro viaje entrando a la ligadonde <br>iremos paso a paso por el proceso de registro:</strong></p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"http://localhost:8000/login\" class=\"cta-button\">¡Inicia sesión y completa tu registro!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta la fecha <strong>2025-09-18 19:39:27</strong></p>\r\n        <p>Si este email te resulta desconocido, ignoralo.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #fae5ff; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #fae5ff; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(31,20,'¡Bienvenido a Signforce! Correo de verificación SIG250109PBG','2025-08-22 15:56:01','thelegendofmax19@gmail.com','Content-Type: text/html; charset=\"UTF-8\"\r\nFrom: thelegendofmax19@gmail.com\r\nTo: contacto.somostec@gmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación SIG250109PBG\r\n\r\nCorreo enviado de Juan de Jesus Nuñez Mendoza\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n       body {\r\n            background: #F3E2D4;\r\n            min-height: 100vh;\r\n            font-family: Inter, \"Helvetica\", sans-serif;\r\n            margin: 0;\r\n            display: flex;\r\n            align-items: center;\r\n            justify-content: center;\r\n            padding: 0px;\r\n            color: #eaeaea;\r\n        }\r\n\r\n        .container {\r\n            background: rgb(24, 1, 48);\r\n            border-radius: 16px;\r\n            padding: 20px;\r\n            max-width: 800px;\r\n            width: 100%;\r\n            box-shadow: \r\n                0 8px 32px rgba(0, 0, 0, 0.6),\r\n                0 0 20px rgba(200, 90, 255, 0.25);\r\n            border: 1px solid rgba(255, 255, 255, 0.08);\r\n            text-align: center;\r\n        }\r\n\r\n        .header {\r\n            margin-bottom: 25px;\r\n        }\r\n\r\n        .logo {\r\n            max-width: 120px;\r\n            margin-bottom: 15px;\r\n        }\r\n\r\n        h1, h2, h3, h4 {\r\n            color: #ffffff;\r\n            margin-bottom: 15px;\r\n            line-height: 1;\r\n            padding: 0;\r\n            letter-spacing: 0.5px;\r\n            word-spacing: 1px;\r\n            line-height: 1.8;\r\n        }\r\n        h1 {\r\n            background:  #35004366;\r\n        }\r\n\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 1rem;\r\n            line-height: 1.5;\r\n            color: #dcdcdc;\r\n        }\r\n\r\n        .cta-button {\r\n            display: inline-block;\r\n            background: #410041;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 14px 30px;\r\n            border-radius: 10px;\r\n            font-weight: 600;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: all 0.3s ease;\r\n            letter-spacing: 0.5px;\r\n            box-shadow: 0 4px 12px rgba(0,0,0,0.35);\r\n        }\r\n\r\n        .cta-button:hover {\r\n            transform: translateY(-3px);\r\n            background: #660066;\r\n            box-shadow: 0 8px 25px rgba(0,0,0,0.4);\r\n        }\r\n\r\n        .divider {\r\n            height: 1px;\r\n            background: linear-gradient(90deg, transparent, #afafaf, transparent);\r\n            margin: 30px 0;\r\n        }\r\n\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 0.85rem;\r\n            color: #999;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n\r\n        <h4>Este es tu nombre de usuario que podrás acceder a la plataforma: <strong>contacto.somostec@gmail.com</strong></h4>\r\n        <h4>Tu contraseña es: <strong>lñAUBk(NmOJb</strong> ¡La puedes cambiar cuando quieras!</h4>\r\n        <div class=\"divider\"></div>\r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p><strong>Comienza nuestro viaje entrando a la ligadonde <br>iremos paso a paso por el proceso de registro:</strong></p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"http://localhost:8000/login\" class=\"cta-button\">¡Inicia sesión y completa tu registro!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta la fecha <strong>2025-09-21 15:56:00</strong></p>\r\n        <p>Si este email te resulta desconocido, ignoralo.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #fae5ff; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #fae5ff; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(32,25,'¡Bienvenido a Signforce! Correo de verificación JSTFRTS123456','2025-09-02 01:40:54','thelegendofmax19@gmail.com','To: thelegendofmax19@gmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación JSTFRTS123456\r\nContent-Type: text/html; charset=\"UTF-8\"\r\nFrom: thelegendofmax19@gmail.com\r\n\r\nCorreo enviado de Diego Leonardo Reyes Chavez\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n       body {\r\n            background: #F3E2D4;\r\n            min-height: 100vh;\r\n            font-family: Inter, \"Helvetica\", sans-serif;\r\n            margin: 0;\r\n            display: flex;\r\n            align-items: center;\r\n            justify-content: center;\r\n            padding: 0px;\r\n            color: #eaeaea;\r\n        }\r\n\r\n        .container {\r\n            background: rgb(24, 1, 48);\r\n            border-radius: 16px;\r\n            padding: 20px;\r\n            max-width: 800px;\r\n            width: 100%;\r\n            box-shadow: \r\n                0 8px 32px rgba(0, 0, 0, 0.6),\r\n                0 0 20px rgba(200, 90, 255, 0.25);\r\n            border: 1px solid rgba(255, 255, 255, 0.08);\r\n            text-align: center;\r\n        }\r\n\r\n        .header {\r\n            margin-bottom: 25px;\r\n        }\r\n\r\n        .logo {\r\n            max-width: 120px;\r\n            margin-bottom: 15px;\r\n        }\r\n\r\n        h1, h2, h3, h4 {\r\n            color: #ffffff;\r\n            margin-bottom: 15px;\r\n            line-height: 1;\r\n            padding: 0;\r\n            letter-spacing: 0.5px;\r\n            word-spacing: 1px;\r\n            line-height: 1.8;\r\n        }\r\n        h1 {\r\n            background:  #35004366;\r\n        }\r\n\r\n         a {\r\n            color: rgb(255, 0, 255); /* Color de los enlaces por defecto */\r\n        }\r\n\r\n        a:hover {\r\n            color: rgb(0, 255, 200); /* Color al pasar el ratón */\r\n        }\r\n\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 1rem;\r\n            line-height: 1.5;\r\n            color: #dcdcdc;\r\n        }\r\n\r\n        .cta-button {\r\n            display: inline-block;\r\n            background: #410041;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 14px 30px;\r\n            border-radius: 10px;\r\n            font-weight: 600;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: all 0.3s ease;\r\n            letter-spacing: 0.5px;\r\n            box-shadow: 0 4px 12px rgba(0,0,0,0.35);\r\n        }\r\n\r\n        .cta-button:hover {\r\n            transform: translateY(-3px);\r\n            background: #660066;\r\n            box-shadow: 0 8px 25px rgba(0,0,0,0.4);\r\n        }\r\n\r\n        .divider {\r\n            height: 1px;\r\n            background: linear-gradient(90deg, transparent, #afafaf, transparent);\r\n            margin: 30px 0;\r\n        }\r\n\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 0.85rem;\r\n            color: #999;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n        <a href=\"email.com\">enlace</a>\r\n        <h4>Este es tu nombre de usuario que podrás acceder a la plataforma: <strong>thelegendofmax19@gmail.com</strong></h4>\r\n        <h4>Tu contraseña es: <strong>Y/FMwt-keR;v</strong> ¡La puedes cambiar cuando quieras!</h4>\r\n        <div class=\"divider\"></div>\r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p><strong>Comienza nuestro viaje entrando a la ligadonde <br>iremos paso a paso por el proceso de registro:</strong></p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"http://localhost:8000/login\" class=\"cta-button\">¡Inicia sesión y completa tu registro!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta la fecha <strong>2025-10-02 01:40:53</strong></p>\r\n        <p>Si este email te resulta desconocido, ignoralo.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #fae5ff; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #fae5ff; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3),(33,26,'¡Bienvenido a Signforce! Correo de verificación JSTFRTS123456','2025-09-02 01:55:53','thelegendofmax19@gmail.com','To: thelegendofmax19@gmail.com\r\nSubject: ¡Bienvenido a Signforce! Correo de verificación JSTFRTS123456\r\nContent-Type: text/html; charset=\"UTF-8\"\r\nFrom: thelegendofmax19@gmail.com\r\n\r\nCorreo enviado de Diego Leonardo Reyes Chavez\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <title>Bienvenido a Signforce</title>\r\n    <style>\r\n       body {\r\n            background: #F3E2D4;\r\n            min-height: 100vh;\r\n            font-family: Inter, \"Helvetica\", sans-serif;\r\n            margin: 0;\r\n            display: flex;\r\n            align-items: center;\r\n            justify-content: center;\r\n            padding: 0px;\r\n            color: #eaeaea;\r\n        }\r\n\r\n        .container {\r\n            background: rgb(24, 1, 48);\r\n            border-radius: 16px;\r\n            padding: 20px;\r\n            max-width: 800px;\r\n            width: 100%;\r\n            box-shadow: \r\n                0 8px 32px rgba(0, 0, 0, 0.6),\r\n                0 0 20px rgba(200, 90, 255, 0.25);\r\n            border: 1px solid rgba(255, 255, 255, 0.08);\r\n            text-align: center;\r\n        }\r\n\r\n        .header {\r\n            margin-bottom: 25px;\r\n        }\r\n\r\n        .logo {\r\n            max-width: 120px;\r\n            margin-bottom: 15px;\r\n        }\r\n\r\n        h1, h2, h3, h4 {\r\n            color: #ffffff;\r\n            margin-bottom: 15px;\r\n            line-height: 1;\r\n            padding: 0;\r\n            letter-spacing: 0.5px;\r\n            word-spacing: 1px;\r\n            line-height: 1.8;\r\n        }\r\n        h1 {\r\n            background:  #35004366;\r\n        }\r\n\r\n         a {\r\n            color: rgb(255, 0, 255); /* Color de los enlaces por defecto */\r\n        }\r\n\r\n        a:hover {\r\n            color: rgb(0, 255, 200); /* Color al pasar el ratón */\r\n        }\r\n\r\n        p {\r\n            margin-bottom: 20px;\r\n            font-size: 1rem;\r\n            line-height: 1.5;\r\n            color: #dcdcdc;\r\n        }\r\n\r\n        .cta-button {\r\n            display: inline-block;\r\n            background: #410041;\r\n            color: white !important;\r\n            text-decoration: none;\r\n            padding: 14px 30px;\r\n            border-radius: 10px;\r\n            font-weight: 600;\r\n            margin: 20px 0;\r\n            text-align: center;\r\n            transition: all 0.3s ease;\r\n            letter-spacing: 0.5px;\r\n            box-shadow: 0 4px 12px rgba(0,0,0,0.35);\r\n        }\r\n\r\n        .cta-button:hover {\r\n            transform: translateY(-3px);\r\n            background: #660066;\r\n            box-shadow: 0 8px 25px rgba(0,0,0,0.4);\r\n        }\r\n\r\n        .divider {\r\n            height: 1px;\r\n            background: linear-gradient(90deg, transparent, #afafaf, transparent);\r\n            margin: 30px 0;\r\n        }\r\n\r\n        .footer {\r\n            margin-top: 30px;\r\n            font-size: 0.85rem;\r\n            color: #999;\r\n        }\r\n    </style>\r\n</head>\r\n<body>\r\n    <div class=\"container\">\r\n        <div class=\"header\">\r\n            <!-- Reemplaza con tu logo\r\n            <img src=\"[URL_LOGO]\" alt=\"[Nombre de la Plataforma]\" class=\"logo\">\r\n            -->\r\n            <h1>¡Bienvenido a Signforce!</h1>\r\n        </div>\r\n        <a href=\"email.com\">enlace</a>\r\n        <h4>Este es tu nombre de usuario que podrás acceder a la plataforma: <strong>thelegendofmax19@gmail.com</strong></h4>\r\n        <h4>Tu contraseña es: <strong>cQZpEv*+ñsId</strong> ¡La puedes cambiar cuando quieras!</h4>\r\n        <div class=\"divider\"></div>\r\n        <p>Nos emociona formar un nuevo vinculo contigo y tus colaboradores. En Signforce podrás firmar tus documentos facilmente y con toda la seguridad y certeza legal.</p>\r\n        \r\n        <p><strong>Comienza nuestro viaje entrando a la ligadonde <br>iremos paso a paso por el proceso de registro:</strong></p>\r\n        \r\n        <div style=\"text-align: center;\">\r\n            <a href=\"http://localhost:8000/login\" class=\"cta-button\">¡Inicia sesión y completa tu registro!</a>\r\n        </div>\r\n        \r\n        <p>¡Tomate tu tiempo! tienes hasta la fecha <strong>2025-10-02 01:55:52</strong></p>\r\n        <p>Si este email te resulta desconocido, ignoralo.</p>\r\n        \r\n        <div class=\"divider\"></div>\r\n        \r\n        <p>Si tienes alguna pregunta, no dudes en contactarnos en preguntas@signforce.com</p>\r\n        \r\n        <p>¡Nos vemos pronto!</p>\r\n        \r\n        <p>Atte: Equipo Signforce</p>\r\n        \r\n        <div class=\"footer\">\r\n            <p>© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n            <p>CDMX, México</p>\r\n            <p>\r\n                <a href=\"https://www.signforce.com/privacy\" style=\"color: #fae5ff; text-decoration: none;\">Política de Privacidad</a> | \r\n                <a href=\"https://www.signforce.com/terms\" style=\"color: #fae5ff; text-decoration: none;\">Términos de Servicio</a>\r\n            </p>\r\n        </div>\r\n    </div>\r\n</body>\r\n</html>\n',1,NULL,3);
/*!40000 ALTER TABLE `emails` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `folders`
--

DROP TABLE IF EXISTS `folders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `folders` (
  `idFolder` bigint NOT NULL AUTO_INCREMENT,
  `ownerInsFoldert_fk` bigint NOT NULL,
  `ownerTeamFolder_fk` bigint NOT NULL,
  `creatorUserFolder_fk` bigint NOT NULL,
  `creationAtFolder` datetime NOT NULL,
  `lastModifiedFolder` datetime NOT NULL,
  `deletedAtFolder` datetime DEFAULT NULL,
  `deletedReasonFolder` text,
  `purposeFolder` varchar(20) NOT NULL,
  `descriptionFolder` text,
  `closedAtFolder` datetime DEFAULT NULL,
  `secuentialSign` tinyint(1) NOT NULL DEFAULT '0',
  `pathSerializedFolder` varchar(100) NOT NULL,
  `expirationDateFolder` datetime NOT NULL,
  `numDocsFolder` int NOT NULL DEFAULT '0',
  `numDocsSignFolder` int NOT NULL DEFAULT '0',
  `numSignersFolder` int NOT NULL DEFAULT '0',
  `numReceiversFolder` int NOT NULL DEFAULT '0',
  `completedAtFolder` datetime DEFAULT NULL,
  PRIMARY KEY (`idFolder`),
  KEY `fk_folders_owner_inst` (`ownerInsFoldert_fk`),
  KEY `fk_folders_owner_team` (`ownerTeamFolder_fk`),
  KEY `fk_folders_creator_user` (`creatorUserFolder_fk`),
  CONSTRAINT `fk_folders_owner_inst` FOREIGN KEY (`ownerInsFoldert_fk`) REFERENCES `institutions` (`idInstitution`) ON DELETE CASCADE,
  CONSTRAINT `fk_folders_owner_team` FOREIGN KEY (`ownerTeamFolder_fk`) REFERENCES `teams` (`idTeam`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `folders`
--

LOCK TABLES `folders` WRITE;
/*!40000 ALTER TABLE `folders` DISABLE KEYS */;
/*!40000 ALTER TABLE `folders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `institutions`
--

DROP TABLE IF EXISTS `institutions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `institutions` (
  `idInstitution` bigint NOT NULL AUTO_INCREMENT,
  `legalNameInst` varchar(100) NOT NULL,
  `aliasNameInst` varchar(50) NOT NULL,
  `taxNumInst` varchar(45) NOT NULL,
  `streetAddress` text,
  `addressLine` text,
  `postalCode` varchar(10) DEFAULT NULL,
  `neighborhood` varchar(256) DEFAULT NULL,
  `locality` varchar(256) DEFAULT NULL,
  `stateCodeInst_fk` int DEFAULT '1',
  `countryCodeInst_fk` int DEFAULT '1',
  `formattedAddress` text,
  `contactPhoneInst` varchar(15) NOT NULL,
  `contactEmailInst` varchar(45) NOT NULL,
  `createdAtInst` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `lastModifiedInst` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deletedAtInst` datetime DEFAULT CURRENT_TIMESTAMP,
  `statusInst_fk` int NOT NULL DEFAULT '1',
  `activeInst` tinyint(1) NOT NULL DEFAULT '0',
  `typeContractInst` int NOT NULL DEFAULT '0',
  `paymentDataInst_fk` int DEFAULT NULL,
  `logoUrlInst` varchar(100) DEFAULT NULL,
  `legalSignupName` varchar(100) CHARACTER SET armscii8 COLLATE armscii8_general_ci NOT NULL,
  `legalSignupLastname` varchar(100) NOT NULL,
  `registerSignupName` varchar(100) NOT NULL,
  `registerSignupLastname` varchar(100) NOT NULL,
  `rootUser_fk` bigint DEFAULT NULL,
  PRIMARY KEY (`idInstitution`)
) ENGINE=InnoDB AUTO_INCREMENT=72 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `institutions`
--

LOCK TABLES `institutions` WRITE;
/*!40000 ALTER TABLE `institutions` DISABLE KEYS */;
INSERT INTO `institutions` VALUES (66,'SIGNFORCE S.A. de C.V.','SIGNFORCE','SIG250109PBA','Laguna de la mancha 62',NULL,'11520','Granada','Miguel Hidalgo',1,1,NULL,'5529996911','contacto.somostec@gmail.com','2025-08-22 15:56:00','2025-08-22 15:56:00','2025-08-22 15:56:00',6,0,0,1,NULL,'Juan de Jesus','Nuñez Mendoza','Juan de Jesus','Nuñez Mendoza',20),(71,'JUSTFRUITS Ltd.','JUSTFRUITS','JSTFRTS123456','Laguna de la mancha 64',NULL,'11520','Granada','Miguel Hidalgo',1,1,NULL,'5555555555','thelegendofmax19@gmail.com','2025-09-02 01:55:52','2025-09-02 01:55:52','2025-09-02 01:55:52',3,0,0,NULL,NULL,'Diego Leonardo','Reyes Chavez','Diego Leonardo','Reyes Chavez',26);
/*!40000 ALTER TABLE `institutions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invites`
--

DROP TABLE IF EXISTS `invites`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `invites` (
  `idInvite` varchar(50) NOT NULL,
  `idFolderInvite_fk` bigint NOT NULL,
  `idDocumentInvite_fk` bigint NOT NULL,
  `idUserIssuerInvite_fk` bigint NOT NULL,
  `idUserInvite_fk` bigint NOT NULL,
  `idTeamInvite_fk` bigint NOT NULL,
  `IdSignatureInvite_fk` bigint DEFAULT NULL,
  `requireAliveProof` tinyint(1) NOT NULL,
  `limitAliveProofTrys` int NOT NULL DEFAULT '3',
  `okAliveProof` tinyint(1) DEFAULT NULL,
  `expirationDateInvite` datetime DEFAULT NULL,
  `isSignInvite` tinyint(1) NOT NULL,
  `sentAtInvite` datetime DEFAULT NULL,
  `openedAtInvite` datetime DEFAULT NULL,
  `closedAtInvite` datetime DEFAULT NULL,
  `descriptionTextInvite` text,
  PRIMARY KEY (`idInvite`),
  KEY `fk_invite_folder` (`idFolderInvite_fk`),
  KEY `fk_invite_document` (`idDocumentInvite_fk`),
  KEY `fk_invite_issuer_user` (`idUserIssuerInvite_fk`),
  KEY `fk_invite_target_user` (`idUserInvite_fk`),
  KEY `fk_invite_team` (`idTeamInvite_fk`),
  KEY `fk_invite_signature` (`IdSignatureInvite_fk`),
  CONSTRAINT `fk_invite_document` FOREIGN KEY (`idDocumentInvite_fk`) REFERENCES `documents` (`idDocument`) ON DELETE CASCADE,
  CONSTRAINT `fk_invite_folder` FOREIGN KEY (`idFolderInvite_fk`) REFERENCES `folders` (`idFolder`) ON DELETE CASCADE,
  CONSTRAINT `fk_invite_signature` FOREIGN KEY (`IdSignatureInvite_fk`) REFERENCES `signatures` (`idSignature`) ON DELETE SET NULL,
  CONSTRAINT `fk_invite_team` FOREIGN KEY (`idTeamInvite_fk`) REFERENCES `teams` (`idTeam`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invites`
--

LOCK TABLES `invites` WRITE;
/*!40000 ALTER TABLE `invites` DISABLE KEYS */;
/*!40000 ALTER TABLE `invites` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `kyc`
--

DROP TABLE IF EXISTS `kyc`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `kyc` (
  `idInsttitution_fk` bigint NOT NULL,
  `idUser_fk` bigint NOT NULL,
  `documentHash` varchar(100) NOT NULL,
  `documentType` varchar(30) DEFAULT NULL,
  `documentName` varchar(50) DEFAULT NULL,
  `documentClass` varchar(50) DEFAULT NULL,
  `documentPath` varchar(255) DEFAULT NULL,
  `expirationDate` datetime DEFAULT NULL,
  PRIMARY KEY (`idInsttitution_fk`,`idUser_fk`,`documentHash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `kyc`
--

LOCK TABLES `kyc` WRITE;
/*!40000 ALTER TABLE `kyc` DISABLE KEYS */;
INSERT INTO `kyc` VALUES (66,20,'DDUgZ6JF/tfIg1E2CAThpe69QzHt8yHUL/COutRKcAA=','pdf','1756794810466092800_CV_MariaBarrera','PoderRepresentante','temp/66/20/',NULL),(66,20,'qIsSgfs8kFYPLMvSYoeH+d9NucWARs8upgqR18FuOxA=','pdf','1756794847496311300_CV_2025_JJNM','PruebaResidencia','temp/66/20/',NULL),(66,20,'thV6h7lh9hIVj8l7JoNdn2eizvhRdoG/AVIRuZX3Qcw=','pdf','CV2024_eng','IdentidadOficial','C:/Users/USER/Desktop/',NULL),(66,20,'WBmwaNYeqZ+Hr10vIqgq2Kvr3j23u/g3sbMCs98joq0=','pdf','Diagrama_AWS','ActaConstitutiva','C:/Users/USER/Desktop/',NULL),(71,26,'DDUgZ6JF/tfIg1E2CAThpe69QzHt8yHUL/COutRKcAA=','pdf','1756799873497022300_CV_MariaBarrera','PoderRepresentante','temp/71/26/',NULL),(71,26,'qIsSgfs8kFYPLMvSYoeH+d9NucWARs8upgqR18FuOxA=','pdf','1756799859709880700_CV_2025_JJNM','ActaConstitutiva','temp/71/26/',NULL),(71,26,'thV6h7lh9hIVj8l7JoNdn2eizvhRdoG/AVIRuZX3Qcw=','pdf','1756799889950012700_CV2024_eng','IdentidadOficial','temp/71/26/',NULL),(71,26,'WBmwaNYeqZ+Hr10vIqgq2Kvr3j23u/g3sbMCs98joq0=','pdf','1756799900331305800_Diagrama_AWS','PruebaResidencia','temp/71/26/',NULL);
/*!40000 ALTER TABLE `kyc` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `microapps`
--

DROP TABLE IF EXISTS `microapps`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `microapps` (
  `idapp` int NOT NULL AUTO_INCREMENT,
  `domainApp` varchar(200) NOT NULL DEFAULT 'localhost',
  `portApp` int NOT NULL DEFAULT '8000',
  `nameApp` varchar(45) NOT NULL DEFAULT 'SignforceApp',
  `description` varchar(300) DEFAULT NULL,
  `isActive` tinyint(1) NOT NULL DEFAULT '1',
  `versionApp` int NOT NULL DEFAULT '1',
  `currPathApp` varchar(500) DEFAULT 'C:/',
  PRIMARY KEY (`idapp`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `microapps`
--

LOCK TABLES `microapps` WRITE;
/*!40000 ALTER TABLE `microapps` DISABLE KEYS */;
INSERT INTO `microapps` VALUES (3,'localhost',8000,'sfmiddle','Middleware para comunicar con app y pagina principal',1,1,'C:/Users/USER/Desktop/signForce/sfmiddle'),(4,'localhost',8001,'emailServ','App de envío de correo electrónico y mensajería',1,1,'C:/Users/USER/Desktop/signForce/emailServ');
/*!40000 ALTER TABLE `microapps` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `payment`
--

DROP TABLE IF EXISTS `payment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `payment` (
  `idInstitution` bigint NOT NULL AUTO_INCREMENT,
  `statusPayment_fk` smallint DEFAULT NULL,
  `planId_fk` int DEFAULT NULL,
  `expirationPlan` datetime DEFAULT NULL,
  `cardNumber` varchar(45) DEFAULT NULL,
  `expirationDate` varchar(6) DEFAULT NULL,
  `nameOwner` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`idInstitution`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `payment`
--

LOCK TABLES `payment` WRITE;
/*!40000 ALTER TABLE `payment` DISABLE KEYS */;
INSERT INTO `payment` VALUES (1,0,3,'2025-10-12 17:35:53','','','Juan Nuñez');
/*!40000 ALTER TABLE `payment` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `platformrolepermissions`
--

DROP TABLE IF EXISTS `platformrolepermissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `platformrolepermissions` (
  `role_id` smallint NOT NULL,
  `permission_id` smallint NOT NULL,
  `description` text NOT NULL,
  PRIMARY KEY (`role_id`,`permission_id`),
  KEY `fk_permissions_permission` (`permission_id`),
  CONSTRAINT `fk_permissions_role` FOREIGN KEY (`role_id`) REFERENCES `userroleplatform` (`id`) ON DELETE CASCADE,
  CONSTRAINT `platformrolepermissions_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `userroleplatform` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `platformrolepermissions`
--

LOCK TABLES `platformrolepermissions` WRITE;
/*!40000 ALTER TABLE `platformrolepermissions` DISABLE KEYS */;
/*!40000 ALTER TABLE `platformrolepermissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `signatures`
--

DROP TABLE IF EXISTS `signatures`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `signatures` (
  `idSignature` bigint NOT NULL AUTO_INCREMENT,
  `digestValueSign` varchar(100) NOT NULL,
  `digestAlgoSign_fk` int NOT NULL,
  `signatureAlgoSign_fk` int DEFAULT NULL,
  `signatureValueSign` text,
  `genTimeSign` datetime DEFAULT NULL,
  `idInstSign_fk` bigint NOT NULL,
  `idUserSign_fk` bigint NOT NULL,
  `idTeamSign_fk` bigint NOT NULL,
  `idFolderSign_fk` bigint NOT NULL,
  `idUserKeysSign_fk` bigint DEFAULT NULL,
  `pathSign` varchar(100) DEFAULT NULL,
  `typeSign_fk` int DEFAULT NULL,
  `nonceSign` bigint DEFAULT NULL,
  `xSign` float DEFAULT NULL,
  `ySign` float DEFAULT NULL,
  `wSign` float DEFAULT NULL,
  `hSign` float DEFAULT NULL,
  `pageSign` int DEFAULT NULL,
  `ipSignerSign` varchar(36) DEFAULT NULL,
  `inviteNumSign` varchar(50) DEFAULT NULL,
  `notifyCreatorSign` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`idSignature`),
  KEY `fk_sign_digest_algo` (`digestAlgoSign_fk`),
  KEY `fk_sign_signature_algo` (`signatureAlgoSign_fk`),
  KEY `fk_sign_institution` (`idInstSign_fk`),
  KEY `fk_sign_user` (`idUserSign_fk`),
  KEY `fk_sign_team` (`idTeamSign_fk`),
  KEY `fk_sign_folder` (`idFolderSign_fk`),
  KEY `fk_sign_user_keys` (`idUserKeysSign_fk`),
  CONSTRAINT `fk_sign_digest_algo` FOREIGN KEY (`digestAlgoSign_fk`) REFERENCES `algos` (`idAlgo`) ON DELETE CASCADE,
  CONSTRAINT `fk_sign_folder` FOREIGN KEY (`idFolderSign_fk`) REFERENCES `folders` (`idFolder`) ON DELETE CASCADE,
  CONSTRAINT `fk_sign_institution` FOREIGN KEY (`idInstSign_fk`) REFERENCES `institutions` (`idInstitution`) ON DELETE CASCADE,
  CONSTRAINT `fk_sign_signature_algo` FOREIGN KEY (`signatureAlgoSign_fk`) REFERENCES `algos` (`idAlgo`) ON DELETE CASCADE,
  CONSTRAINT `fk_sign_team` FOREIGN KEY (`idTeamSign_fk`) REFERENCES `teams` (`idTeam`) ON DELETE CASCADE,
  CONSTRAINT `fk_sign_user_keys` FOREIGN KEY (`idUserKeysSign_fk`) REFERENCES `userkeys` (`idUserKeys`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `signatures`
--

LOCK TABLES `signatures` WRITE;
/*!40000 ALTER TABLE `signatures` DISABLE KEYS */;
/*!40000 ALTER TABLE `signatures` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `states`
--

DROP TABLE IF EXISTS `states`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `states` (
  `idState` int NOT NULL AUTO_INCREMENT,
  `nameState` varchar(100) NOT NULL,
  `shrinkName` varchar(5) NOT NULL,
  `localUTCTime` int NOT NULL,
  `DSTTime` int NOT NULL,
  `country_fk` int NOT NULL,
  PRIMARY KEY (`idState`),
  KEY `fk_state_country` (`country_fk`),
  CONSTRAINT `fk_state_country` FOREIGN KEY (`country_fk`) REFERENCES `countries` (`idCountry`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `states`
--

LOCK TABLES `states` WRITE;
/*!40000 ALTER TABLE `states` DISABLE KEYS */;
INSERT INTO `states` VALUES (1,'Ciudad de México','CDMX',-6,-6,1),(2,'Aguascalientes','AGS',-6,-6,1),(3,'Baja California','BCN',-7,-7,1),(4,'Baja California Sur','BCS',-7,-7,1),(5,'Campeche','CAMP',-6,-6,1),(6,'Coahuila','COAH',-6,-6,1),(7,'Colima','COL',-6,-6,1),(8,'Chiapas','CHH',-6,-6,1),(9,'Chihuahua','CHIH',-7,-7,1),(10,'Durango','DGO',-6,-6,1),(11,'Estado de México','MEX',-6,-6,1),(12,'Guanajuato','GTO',-6,-6,1),(13,'Guerrero','GRO',-6,-6,1),(14,'Hidalgo','HGO',-6,-6,1),(15,'Jalisco','JAL',-6,-6,1),(16,'Michoacán','MICH',-6,-6,1),(17,'Morelos','MOR',-6,-6,1),(18,'Nayarit','NAY',-7,-7,1),(19,'Nuevo León','NL',-6,-6,1),(20,'Oaxaca','OAX',-6,-6,1),(21,'Puebla','PUE',-6,-6,1),(22,'Queretaro','QRO',-6,-6,1),(23,'Quintana Roo','Q.ROO',-5,-5,1),(24,'San Luis Potosi','SLP',-6,-6,1),(25,'Sinaloa','SIN',-7,-7,1),(26,'Sonora','SON',-7,-7,1),(27,'Tabasco','TAB',-6,-6,1),(28,'Tamaulipas','TAMPS',-6,-6,1),(29,'Tlaxcala','TLAX',-6,-6,1),(30,'Veracruz','VER',-6,-6,1),(31,'Yucatán','YUC',-6,-6,1),(32,'Zacatecas','ZAC',-6,-6,1);
/*!40000 ALTER TABLE `states` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `statusinstitution`
--

DROP TABLE IF EXISTS `statusinstitution`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `statusinstitution` (
  `idStatusInst` int NOT NULL AUTO_INCREMENT,
  `nameStatus` varchar(30) NOT NULL,
  `permissionStatusInst_fk` varchar(10) NOT NULL,
  `descriptionStatusInst` text,
  PRIMARY KEY (`idStatusInst`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `statusinstitution`
--

LOCK TABLES `statusinstitution` WRITE;
/*!40000 ALTER TABLE `statusinstitution` DISABLE KEYS */;
INSERT INTO `statusinstitution` VALUES (1,'Nulo','1','Sin estatus o no requiere'),(2,'PENDIENTE_REGISTRO','1','El usuario inició el registro pero no lo completó (por ejemplo, aún no verificó el correo o no llenó todos los datos).'),(3,'REVISION_REGISTRO','1','El usuario ya completo sus datos y subió sus documentos de identidad, comprobante de domicilio, etc., pero están en revisión.'),(4,'RECHAZO_DOCUMENTOS','1','Algún documento fue rechazado por ser ilegible, inválido o por no coincidir con los datos del usuario. Se debe permitir que el usuario los reenvíe.'),(5,'CONTRATO_PENDIENTE','1','El usuario está verificado pero aún no ha firmado el contrato de servicio (por ejemplo, falta aceptación digital o firma electrónica).'),(6,'PLAN_PAGO','1','Falta plan de pago. El contrato fue firmado y registrado correctamente (puede incluir sello de tiempo y almacenamiento del documento).'),(7,'ACTIVO','1','Usuario completamente habilitado para operar dentro del sistema (por ejemplo, puede firmar documentos, generar sellos de tiempo, etc.).'),(8,'SUSPENDIDO','1','Usuario que fue suspendido temporalmente por problemas administrativos, seguridad o requerimientos legales.'),(9,'INACTIVO','1','Usuario que no presenta actividad durante un periodo de tiempo (puede segir con autorización para usar ciertos elementos de la plataforma).'),(10,'BAJA','1','Usuario que fue dado de baja o canceló su cuenta voluntariamente.'),(11,'REVOCADO','1','Usuario que fue dado de baja o suspendido por falta de pago, incumplimiento de politicas o mal uso de la plataforma.');
/*!40000 ALTER TABLE `statusinstitution` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `teammembers`
--

DROP TABLE IF EXISTS `teammembers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `teammembers` (
  `idTeam` bigint NOT NULL,
  `idUser` bigint NOT NULL,
  `idRole` smallint NOT NULL,
  `assignedBy` bigint NOT NULL,
  `assignedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`idTeam`,`idUser`),
  KEY `fk_team_user` (`idUser`),
  KEY `fk_team_role` (`idRole`),
  KEY `fk_team_assignedBy` (`assignedBy`),
  CONSTRAINT `fk_team_role` FOREIGN KEY (`idRole`) REFERENCES `userroleplatform` (`id`),
  CONSTRAINT `fk_team_team` FOREIGN KEY (`idTeam`) REFERENCES `teams` (`idTeam`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `teammembers`
--

LOCK TABLES `teammembers` WRITE;
/*!40000 ALTER TABLE `teammembers` DISABLE KEYS */;
/*!40000 ALTER TABLE `teammembers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `teams`
--

DROP TABLE IF EXISTS `teams`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `teams` (
  `idTeam` bigint NOT NULL,
  `idInstitutionTeam_fk` int NOT NULL,
  `creatorUserTeam_fk` int NOT NULL,
  `nameTeam` varchar(40) NOT NULL,
  `createdAtTeam` datetime NOT NULL,
  `lastModifiedTeam` datetime NOT NULL,
  `deletedAtTeam` datetime DEFAULT NULL,
  `logoUrlTeam` varchar(100) DEFAULT NULL,
  `limitSignersTeam` int NOT NULL DEFAULT '10',
  `limitUsersTeam` int NOT NULL DEFAULT '20',
  PRIMARY KEY (`idTeam`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `teams`
--

LOCK TABLES `teams` WRITE;
/*!40000 ALTER TABLE `teams` DISABLE KEYS */;
/*!40000 ALTER TABLE `teams` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `userkeys`
--

DROP TABLE IF EXISTS `userkeys`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `userkeys` (
  `idUserKeys` bigint NOT NULL AUTO_INCREMENT,
  `keyFilePath` varchar(200) DEFAULT NULL,
  `certFilePath` varchar(200) NOT NULL,
  `serialNumberX509` varchar(100) NOT NULL,
  `certVersion` tinyint DEFAULT NULL,
  `signAlgoIssuer_fk` tinyint NOT NULL,
  `digestAlgoIssuer_fk` tinyint NOT NULL,
  `issuerRFC4514` varchar(100) NOT NULL,
  `notValidAfter` datetime NOT NULL,
  `notValidBefore` datetime NOT NULL,
  `subjectRFC4514` varchar(100) NOT NULL,
  `ocspUrl` varchar(100) DEFAULT NULL,
  `crlsPath` varchar(100) DEFAULT NULL,
  `signatureIssuer` text,
  `createdAtKey` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `revokedAtKey` datetime DEFAULT NULL,
  `isRevoked` tinyint(1) NOT NULL DEFAULT '0',
  `keyLenKey` int NOT NULL,
  `algoSignType` varchar(10) NOT NULL,
  `hashKey` varchar(50) NOT NULL,
  `hashCer` varchar(50) NOT NULL,
  PRIMARY KEY (`idUserKeys`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `userkeys`
--

LOCK TABLES `userkeys` WRITE;
/*!40000 ALTER TABLE `userkeys` DISABLE KEYS */;
INSERT INTO `userkeys` VALUES (1,'C:/Users/USER/Desktop/certs/Claveprivada_FIEL_NUMJ900112T99_20220707_112849.key','C:/Users/USER/Desktop/certs/numj900112t99.cer','NUJ900112HDFXNN09',3,9,21,'2.5.4.45 = SAT970701NN3','2026-07-07 10:44:44','2022-07-07 10:44:44','SERIALNUMBER = NUMJ900112HDFXNN09','www.ocsp.com','C:/a/b/c.crls','abcde','2025-05-19 15:43:12',NULL,0,2048,'9','','');
/*!40000 ALTER TABLE `userkeys` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `userplatformpermissions`
--

DROP TABLE IF EXISTS `userplatformpermissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `userplatformpermissions` (
  `id` smallint NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `userplatformpermissions`
--

LOCK TABLES `userplatformpermissions` WRITE;
/*!40000 ALTER TABLE `userplatformpermissions` DISABLE KEYS */;
/*!40000 ALTER TABLE `userplatformpermissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `userroleplatform`
--

DROP TABLE IF EXISTS `userroleplatform`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `userroleplatform` (
  `id` smallint NOT NULL AUTO_INCREMENT,
  `name` varchar(30) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `userroleplatform`
--

LOCK TABLES `userroleplatform` WRITE;
/*!40000 ALTER TABLE `userroleplatform` DISABLE KEYS */;
INSERT INTO `userroleplatform` VALUES (1,'ROOT','Primer usuario de la plataforma, por defecto es el representante legal'),(2,'MASTER','Comparte mismos privilegios que root excepto darse de baja del sistema o contratar servicios nuevos'),(3,'TEAMOWNER','Permisos de configuracion de su equipo, entorno, documentos y miembros locales'),(4,'TEAMMEMBER','Permiso sobre su configuracion de usuario y sus documentos'),(5,'VIEWER','Acceso a dashboard global de equipos y metricas globales. Solo puede consultar y descargar documentos, sin posibilidad de editarlos.'),(6,'INVITADO','Acceso temporal y limitado a documentos acorde al Id de su invitacion');
/*!40000 ALTER TABLE `userroleplatform` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `userroles`
--

DROP TABLE IF EXISTS `userroles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `userroles` (
  `idUser` bigint NOT NULL,
  `idRole` smallint NOT NULL,
  `grantedBy` bigint NOT NULL,
  `grantedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`idUser`,`idRole`),
  KEY `fk_userroles_role` (`idRole`),
  KEY `fk_userroles_grantedby` (`grantedBy`),
  CONSTRAINT `fk_user_role_role` FOREIGN KEY (`idRole`) REFERENCES `userroleplatform` (`id`),
  CONSTRAINT `fk_userroles_role` FOREIGN KEY (`idRole`) REFERENCES `userroleplatform` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `userroles`
--

LOCK TABLES `userroles` WRITE;
/*!40000 ALTER TABLE `userroles` DISABLE KEYS */;
INSERT INTO `userroles` VALUES (20,1,20,'2025-08-22 15:56:00'),(21,1,21,'2025-09-02 01:31:46'),(25,1,25,'2025-09-02 01:40:53'),(26,1,26,'2025-09-02 01:55:52');
/*!40000 ALTER TABLE `userroles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `idUser` bigint NOT NULL AUTO_INCREMENT,
  `taxNumUser` varchar(20) DEFAULT NULL,
  `pobUidUser` varchar(20) DEFAULT NULL,
  `nameUser` varchar(60) NOT NULL,
  `lastNameUser` varchar(60) DEFAULT NULL,
  `aliasUser` varchar(80) DEFAULT NULL,
  `emailUser` varchar(45) DEFAULT NULL,
  `phoneUser` varchar(15) DEFAULT NULL,
  `countryPhoneCode` int DEFAULT '52',
  `activeUser` tinyint(1) NOT NULL DEFAULT '0',
  `roleAppUser_fk` int NOT NULL DEFAULT '0',
  `idKeysUser_fk` int DEFAULT NULL,
  `appPassHash` blob NOT NULL,
  `kycUser_fk` int DEFAULT NULL,
  `isAliveUser` tinyint(1) NOT NULL DEFAULT '0',
  `idTeam_fk` int DEFAULT '0',
  `idInstitution_fk` int NOT NULL DEFAULT '0',
  `createdAtUser` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `lastModifiedUser` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deletedAtUser` datetime DEFAULT NULL,
  PRIMARY KEY (`idUser`),
  UNIQUE KEY `idUser_UNIQUE` (`idUser`),
  UNIQUE KEY `emailUser_UNIQUE` (`emailUser`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (20,NULL,NULL,'Juan de Jesus','Nuñez Mendoza','root_OtqS8otM7aOFYuvwR8b/BUANTFcjUqEULu3+9n0h5mI=','contacto.somostec@gmail.com','5529996911',52,1,1,NULL,_binary 'Yn/N6OgYq9dCXdTIMml4MSyJLK7YXjOtfylVE3XIfU8=',NULL,0,0,66,'2025-08-22 15:56:00','2025-08-22 15:56:00',NULL),(26,NULL,NULL,'Diego Leonardo','Reyes Chavez','root_fyJT1+IosioIvaHwnFFvb+rYHfZTbrAvqZGjS7ONm+g=','thelegendofmax19@gmail.com','5555555555',52,1,1,NULL,_binary 'HLMlnH4n4MVQiYPi71FmiOR7/a2/w94SI28kI6ppWUY=',NULL,0,0,71,'2025-09-02 01:55:52','2025-09-02 01:55:52',NULL);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `userstatushistory`
--

DROP TABLE IF EXISTS `userstatushistory`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `userstatushistory` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `idUser` bigint NOT NULL,
  `prevStatus` enum('PENDING','ACTIVE','SUSPENDED','REVOKED') NOT NULL,
  `newStatus` enum('PENDING','ACTIVE','SUSPENDED','REVOKED') NOT NULL,
  `changedBy` bigint NOT NULL,
  `changedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `reason` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_userstatus_user` (`idUser`),
  KEY `fk_userstatus_changedby` (`changedBy`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `userstatushistory`
--

LOCK TABLES `userstatushistory` WRITE;
/*!40000 ALTER TABLE `userstatushistory` DISABLE KEYS */;
/*!40000 ALTER TABLE `userstatushistory` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `usetypealgos`
--

DROP TABLE IF EXISTS `usetypealgos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `usetypealgos` (
  `idusetype` int NOT NULL AUTO_INCREMENT,
  `useAlgo` varchar(15) NOT NULL,
  `descriptionAlgo` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`idusetype`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `usetypealgos`
--

LOCK TABLES `usetypealgos` WRITE;
/*!40000 ALTER TABLE `usetypealgos` DISABLE KEYS */;
INSERT INTO `usetypealgos` VALUES (1,'hash','Validate messages integrity'),(2,'sign','Validate integrity, autenticity, non repudiation'),(3,'simencrypt','Allow confidentiality'),(4,'asymencrypt','Allow confidentiality and non repudiation'),(5,'checksum','Validate message integrity'),(6,'errorcorrecting','Allow integrity');
/*!40000 ALTER TABLE `usetypealgos` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-09-12 17:49:59
