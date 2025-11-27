CREATE DATABASE  IF NOT EXISTS `signforceapi` /*!40100 DEFAULT CHARACTER SET utf8mb3 */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `signforceapi`;
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
INSERT INTO `algos` VALUES (1,'1.2.840.10040.4.3','dsa-with-sha1','id-dsa-with-sha1',2,0,0),(2,'1.2.840.10045.4.1','ecdsa-with-SHA1','sha1_ecdsa',2,0,0),(3,'1.2.840.10045.4.3.2','ecdsa-with-SHA256','sha256_ecdsa',2,0,1),(4,'1.2.840.10045.4.3.3','ecdsa-with-SHA384','sha384_ecdsa',2,0,0),(5,'1.2.840.10045.4.3.4','ecdsa-with-SHA512','sha512_ecdsa',2,0,0),(6,'1.2.840.113549.1.1.4','md5WithRSAEncryption','md5_rsa',2,1,0),(7,'1.2.840.113549.1.1.5','sha1-with-rsa-signature','sha1_rsa',2,1,0),(8,'1.2.840.113549.1.1.10','rsassa-pss','rsassa_pss',3,0,1),(9,'1.2.840.113549.1.1.11','sha256WithRSAEncryption','sha256_rsa',2,0,1),(10,'1.2.840.113549.1.1.12','sha384WithRSAEncryption','sha384_rsa',2,0,0),(11,'1.2.840.113549.1.1.13','sha512WithRSAEncryption','sha384_rsa',2,0,0),(12,'1.3.14.3.2.2','md4WitRSA','md4_rsa',2,1,0),(13,'1.3.14.3.2.3','md5WithRSA','md5_rsa',2,1,0),(14,'1.3.14.3.2.11','rsaSignature','rsa',2,0,0),(15,'1.3.14.3.2.13','dsaWithSHA','sha_dsa',2,1,0),(16,'1.3.14.3.2.15','shaWithRSASignature','sha_rsa',2,1,0),(17,'1.2.840.113549.2.2','md2','md2',1,1,0),(18,'1.2.840.113549.2.4','md4','md4',1,1,0),(19,'1.3.14.3.2.18','sha','sha',1,1,0),(20,'1.3.14.3.2.26','sha1','sha1',1,1,1),(21,'2.16.840.1.101.3.4.2.1','sha256','sha256',1,0,1),(22,'2.16.840.1.101.3.4.2.2','sha384','sha384',1,0,1),(23,'2.16.840.1.101.3.4.2.3','sha512','sha512',1,0,1),(24,'2.16.840.1.101.3.4.2.4','sha224','sha224',1,0,0),(25,'2.16.840.1.101.3.4.2.5','sha512-224','sha512_224',1,0,0),(26,'2.16.840.1.101.3.4.2.6','sha512-256	','sha512_256',1,0,0),(27,'2.16.840.1.101.3.4.2.7','sha3-224','sha3_224',1,0,0),(28,'2.16.840.1.101.3.4.2.8','sha3-256','sha3_256',1,0,0),(29,'2.16.840.1.101.3.4.2.9','sha3-384','sha3_384',1,0,0),(30,'2.16.840.1.101.3.4.2.10','sha3-512','sha3_512',1,0,0),(31,'2.16.840.1.101.3.4.2.11','shake128','shake128',1,0,0),(32,'2.16.840.1.101.3.4.2.12','shake256','shake256',1,0,0),(33,'2.16.840.1.101.3.4.1.23','aes192-OFB','aes192_ofb',3,0,0),(34,'16.840.1.101.3.4.1.24','aes192-CFB','aes192_cfb',3,0,0),(35,'2.16.840.1.101.3.4.1.25','id-aes192-wrap','id_aes192_wrap',3,0,0),(36,'2.16.840.1.101.3.4.1.26','aes192-GCM','aes192_gcm',3,0,0),(37,'2.16.840.1.101.3.4.1.27','aes192-CCM','aes192_ccm',3,0,0),(38,'2.16.840.1.101.3.4.1.28','aes192-wrap-pad','aes192_wrap_pad',3,0,0),(39,'2.16.840.1.101.3.4.1.41','aes256-ECB','aes256_ecb',3,0,0),(40,'2.16.840.1.101.3.4.1.42','aes256-CBC','aes256_cbc',3,0,0),(41,'2.16.840.1.101.3.4.1.43','aes256-OFB','aes256_ofb',3,0,0),(42,'2.16.840.1.101.3.4.1.44','aes256-CFB','aes256_cbf',3,0,0),(43,'2.16.840.1.101.3.4.1.45','id-aes256-wrap','id_aes256_wrap',3,0,0),(44,'2.16.840.1.101.3.4.1.46','aes256-GCM','aes256_gcm',3,0,0),(45,'2.16.840.1.101.3.4.1.47','aes256-CCM','aes256_ccm',3,0,0),(46,'2.16.840.1.101.3.4.1.48','aes256-wrap-pad','aes256_wrap_pad',3,0,0),(47,'1.3.14.3.2.6','desECB','desecb',3,1,0),(48,'1.3.14.3.2.7','desCBC','descbc',3,1,0),(49,'1.3.14.3.2.8','desOFB','desofb',3,1,0),(50,'1.3.14.3.2.9','desCFB','descfb',3,1,0),(51,'1.3.14.3.2.10','desMAC','desmac',3,1,0),(52,'1.3.14.3.2.17','desEDE','desede',3,0,1),(53,'1.3.6.1.4.1.4929.1.6','3Des','3des',3,0,0),(54,'1.3.6.1.4.1.4929.1.7','3DesECB','3desecb',3,0,0),(55,'1.3.6.1.4.1.4929.1.8','3DesCBC','3Descbc',3,0,0),(56,'1.3.6.1.4.1.4929.1.9','3DesOFB','3Desofb',3,0,0),(57,'1.3.6.1.4.1.4929.1.10','3DesCFB','3Descfb',3,0,0),(58,'1.2.840.113549.1.1.1','rsaEncryption','rsa',4,0,1),(59,'1.2.840.113549.1.1.6','rsaOAEPEncryptionSET','rsaOAEPEncryptionSET',4,0,1),(60,'1.2.840.113549.1.1.7','id-RSAES-OAEP','id-RSAES-OAEP',4,0,1),(61,'1.2.840.113549.1.1.10','rsassa-pss','rsassa-pss',4,0,1),(62,'1.2.840.10045.2.1','ecPublicKey','id-ecPublicKey',4,0,1),(63,'1.3.132.0.6','secp112r1','secp112r1',4,0,1),(64,'1.3.132.0.10','secp256k1','secp256k1',4,0,1),(65,'1.3.132.0.28','secp128r1','secp128r1',4,0,1),(66,'1.3.132.0.31','secp192k1','secp192k1',4,0,1),(67,'1.3.132.0.33','secp224r1','secp224r1',4,0,1),(68,'1.3.132.0.34','secp384r1','secp384r1',4,0,1),(69,'1.3.132.0.35','secp521r1','secp521r1',4,0,1),(70,'1.2.840.10045.3.1.1','secp192r1','prime192v1',4,0,1),(71,'1.2.840.10045.3.1.7','secp256r1','prime256v1',4,0,1),(72,'1.3.101.110','id-X25519','curve25519',4,0,1),(73,'1.3.101.112','ed25519','Ed25519',4,0,1),(74,'1.3.101.113','ed448','Ed448',4,0,1),(75,'1.3.132.1.12','ecdh','ECDH',4,0,1);
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
  `nameCountry` varchar(80) COLLATE utf8mb4_unicode_ci NOT NULL,
  `aliasCountry` varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `shrinkName` varchar(5) COLLATE utf8mb4_unicode_ci NOT NULL,
  `phoneCode` int NOT NULL,
  `continent` enum('AFRICA','AMERICA','ASIA','EUROPE','OCEANIA') COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`idCountry`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
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
  `nameDocStat` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `descriptionDocStat` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `useDocStat` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`idDocStat`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
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
  `documentHash` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `createdAtDoc` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `lastModifiedDoc` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deletedAtDoc` datetime DEFAULT NULL,
  `deletedReasonDoc` text COLLATE utf8mb4_unicode_ci,
  `ownerInstDoc_fk` bigint NOT NULL,
  `ownerTeamDoc_fk` bigint NOT NULL,
  `creatorUserDoc_fk` bigint NOT NULL,
  `documentPath` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `documentName` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `documentExt` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL,
  `sizeB` int NOT NULL DEFAULT '0',
  `abstractDoc` text COLLATE utf8mb4_unicode_ci,
  `authUseStatus` int NOT NULL DEFAULT '0',
  `authRoleStatus` int NOT NULL DEFAULT '0',
  `activeDoc` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`idDocument`),
  KEY `fk_documents_owner_inst` (`ownerInstDoc_fk`),
  KEY `fk_documents_owner_team` (`ownerTeamDoc_fk`),
  KEY `fk_documents_creator_user` (`creatorUserDoc_fk`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `documents`
--

LOCK TABLES `documents` WRITE;
/*!40000 ALTER TABLE `documents` DISABLE KEYS */;
INSERT INTO `documents` VALUES (1,'U6r7VsslAuhUO4G5TFaifr3nn51caDdIUD0LU3UElfA=','2025-11-26 14:40:57','2025-11-26 14:40:57',NULL,NULL,2,2,2,'uploaded/documents/2/2/2/','CV_2025_JJNM_ENG','pdf',82721,' Resumen de Juan De Jesús Nuñez Mendoza Es un ingeniero mecánico con experiencia en desarrollo de software, análisis de datos y seguridad cibernética. Ha trabajado en empresas como Notaría 230, Coppel y Intermedia México. Tiene experiencia en tecnologías como AWS, Azure, Google Sheets, AppScripts, MongoDB, Redis, ReactJS, Fast API y Flask. Experiencia Laboral Mechanical Engineer en Notaría 230 (Jun/2024 - Jul/2025) Refactorio de middleware multichannel y análisis de transacciones en la base de datos. Análisis y corrección de vulnerabilidades de código estático. Python Developer Sr. en Coppel (Jan/2024 - Jun/2024) Desarrollo de APIs para web scraping de sitios digitales. Creación de un sistema de extracción multimedia para radiofrecuencias y streaming. Developer en Intermedia México (Aug/2022 - Jan/2024) Diseño e implementación de sistemas de recolección de información y informes automatizados. Desarrollo de un sistema de mensajería automática para WhatsApp. Habilidades Técnicas Lenguajes de programación: C, C , Python, Golang, Javascript Bases de datos: SQL, MongoDB, Redis Frameworks web: ReactJS, Fast API, Flask Tecnologías: Linux, Docker, Arduino, Git, Pandas, TensorFlow - PyTorch Educación Ingeniero mecánico en IPN - ESIME Azcapotzalco (2011-2017) Idiomas Español (nativo) Inglés (B2)',1,1,1),(2,'2JeEA34CnaGNei+KraRf/EDmMhN+ILGUHZxPrXHYomM=','2025-11-26 14:40:57','2025-11-26 14:40:57',NULL,NULL,2,2,2,'uploaded/documents/2/2/2/','discurso','pdf',50467,' Resumen La autora de la carta recibe una invitación para dar un discurso sobre el trabajo del instituto en la atención a las víctimas de inundaciones recientes en la Sierra de Hidalgo. Expresa su gratitud por poder compartir sus experiencias y agradecer al equipo de trabajo que la apoyó durante el evento. Hechos El autor tuvo la oportunidad de acudir a una zona de desastre después de las inundaciones recientes en la Sierra de Hidalgo. El trabajo consistió en brindar atención médica y consuelo a las víctimas. La autora expresó su orgullo por pertenecer a un instituto que prioriza la compasión y el compromiso con la salud y el bienestar de sus miembros. Oportunidades La carta ofrece una oportunidad para compartir experiencias y reflexionar sobre el trabajo en situaciones de emergencia. La autora invita al público a conocer más sobre las actividades del instituto y su compromiso con la atención médica y social.',1,1,1),(3,'XbOaERJcgbF8PF4lvpCdTSVqEv1Gk0ic28LuxlSak+w=','2025-11-26 16:26:18','2025-11-26 16:26:18',NULL,NULL,2,2,2,'uploaded/documents/2/2/2/','PRUEBA_SDT_2','pdf',34509,'El archivo electrónico identificado como SDT de DDSF 1180 - FIRMACLIENTEOK - ACTA-CIRCUNSTANCIADA ha sido certificado con un Sello Digital de Tiempo (SELDO DIGITAL DE TIEMPO) que acredita que el dato existía antes de la estampa del Sello, en consonancia con lo establecido por el artículo 89 del Código de Comercio y la Norma Oficial Mexicana NOM-151-SCFI-2016. La información detallada del certificado es la siguiente: Fecha UTC: 29/01/2024 16:34:19 p.m. Hash: Kc/rJN5x7swlysTuTgW7tICEk/ke4BhxU1ew1sbIIsE Algoritmo: 2.16.840.1.101.3.4.2.1 Número de serie: -14548616704993510685347173171720569092 ID: 1247 Policy: 2.16.484.101.10.316.100.8.1.3.1.3 Fecha de emisión: 26-feb-2024 07:01:01 El certificado es emitido por la Notaría 230 del Distrito Federal, en representación de Alfredo Bazúa Witte, quien es el titular de la notaría y prestador de servicios de certificación.',1,1,1),(4,'DDUgZ6JF/tfIg1E2CAThpe69QzHt8yHUL/COutRKcAA=','2025-11-26 16:26:18','2025-11-26 16:26:18',NULL,NULL,2,2,2,'uploaded/documents/2/2/2/','CV_MariaBarrera','pdf',80806,'Resumen de María Barrera: Información Personal - Nacimiento: No disponible - Nacionalidad: Mexicana - Lugar de residencia: Ciudad de México, México - Correo electrónico: mbaarrera415@gmail.com - Perfil en LinkedIn: linkedin.com/in/MariaBarrera415 Resumen Profesional - Cargo: Ejecutiva de Ventas - Experiencia laboral: 5 años en empresas líderes del sector tecnológico como Salesforce, SAP y Microsoft. - Logros: Gestión exitosa de ciclos de ventas completos Compromiso con la venta de soluciones complejas a ejecutivos de alto nivel Construcción de relaciones estratégicas con clientes y superación constante de los objetivos de ingresos - Habilidades lingüísticas: español nativo, inglés fluente y japonés básico Experiencia Laboral 1. Ejecutiva de Ventas - Salesforce (2024-Present) : Gestiona ciclos de ventas completos para cuentas de alto nivel Desarrolla soluciones de Marketing Cloud Utiliza recurso de ingenieros de ventas, servicios profesionales y marketing 2. Ejecutiva de Ventas - SAP (2021-2024) : Gestiona ciclos de ventas completos para cuentas de alto nivel en industrias como manufactura, retail y turismo Cumple con la meta del 100% de cuota laboral anual durante tres años consecutivos Dirige un equipo de soluciones adicionales para aumentar las ventas y el acoplamiento en más de cinco líneas de negocio 3. Gerente de Éxito de los Clientes - Microsoft (2019-2020) : Fomenta relaciones estratégicas con stakeholders clave para garantizar la entrega exitosa de soluciones Establece un 100% de tasa de renovación en cuentas clave y impulsa iniciativas de consumo de nube Monitorea el uso de incidentes y asegura la eficiencia operativa Educación - MBA - Tecnológico de Monterrey (2024-Present) - B.S. Ingeniería Industrial - Tecnológico de Monterrey (2016-2020) - Gestión de Proyectos e Operaciones - Instituto Nacional Polytechnique de Toulouse (2019) Habilidades y Certificaciones - Salesforce - SAP - Power BI - Excel - Planificación de Cuentas - Gerencia del flujo de la cadena de proveedores - Venta estratégica - Comunicación - SCRUM Fundamentos - Gestión de Materiales en SAP - Contabilidad financiera en SAP Idiomas - Español: hablado nativo - Inglés: hablado fluído - Japonés: básico',1,1,1),(5,'C6wfK7XKQRtQvHuUd0CWi9sMy2WXqwSNt5EX/ZKIon4=','2025-11-26 18:28:44','2025-11-26 18:28:44',NULL,NULL,3,3,3,'uploaded/documents/3/3/3/','Ensamblaje1','pdf',106517,'Este es un documento de diseño de hardware (ASSEMBLY 1) de una empresa privada. Fecha: No se proporciona Escala: 1:20 Tolerancias: Fractions y dos y tres decimales para dimensiones y ángulos, respectivamente Material: No se especifica Propietario: La empresa mencionada en el documento (que debe ser insertada)',1,1,1),(6,'wCVeNEyl6+LS6wbc01ZadmcLdkcGooaTa4pru7Qa29U=','2025-11-26 18:28:44','2025-11-26 18:28:44',NULL,NULL,3,3,3,'uploaded/documents/3/3/3/','Lista de materiales','pdf',710321,' Resumen del ensamble Evita El ensamble Evita cuenta con los siguientes elementos: Motor Siemens 0.5 HP 220V 2 polos mod.1RA3 254-2YK34 Perforación de tolva cedula 14 (2mm) Tornillos B18.6.7M - M3.5 x 0.6 x 13 Chumacera con unidad UCP207-20 y perfil aluminio 50x50 2mm Angulo soporte invertido perfil aluminio 50x50 2mm Resistencia de 127 V 1000 W ValvVar AISI 304 Tornillo2 AISI 304 Cañón acero galvanizado 1.25\" Boquilla laton rosca NPT Perno acero 1045 Piñon1 AISI 4140 modulo 2 16 dientes Engrane3 AISI 4140 modulo 2 56 dientes Chumacera2 con unidad UCP203 Piñon3 AISI 4140 modulo 3 16 dientes Engrane4 AISI 4140 modulo 3 56 dientes Bancada soporte rigido metálico soldable Bancada1 barra rigida metálica soldable Variador HITACHI WJ200-007SF (própuesta en dibujo) Pirometro analógico 0-400 °C',1,1,1),(7,'oR/5ro0qPrNs6wgbQRRfNk3RonWOs/FvtGiTxeCo3Nk=','2025-11-26 18:28:44','2025-11-26 18:28:44',NULL,NULL,3,3,3,'uploaded/documents/3/3/3/','tolva croquis','pdf',343527,' Resumen de Eva Perón Eva Perón (1919-1952) fue la segunda primera dama de Argentina durante el gobierno de su esposo, Juan Domingo Perón. A continuación, se presentan algunos hechos clave sobre su vida y legado: Infancia y juventud : Eva Perón nació en Los Toldos, provincia de Río Negro, Argentina. Se mudó a Buenos Aires con su familia a una edad temprana. Carrera política : Entró en la política en la década de 1940, al apoyar a su esposo en sus campañas electorales. Primer dame : Se convirtió en primera dama de Argentina en 1946 y mantuvo este cargo hasta 1952, cuando falleció. Legado : Eva Perón es recordada como una defensora de los derechos humanos y la educación. Fundó la Asociación Mutualista de Trabajadores (AMT), que proporcionaba servicios médicos y sociales a trabajadores pobres. Es importante destacar que, aunque Eva Perón fue una figura influyente en la política argentina, su legado es objeto de debate entre historiadores y críticos. Algunos la ven como un símbolo de la revolución popular, mientras que otros la critican por su relación con el régimen autoritario del peronismo.',1,1,1),(8,'frDMWhwmOBvEJvbXwB7SKM9VEbJq7yI97JgNbKY3/ww=','2025-11-26 18:28:44','2025-11-26 18:28:44',NULL,NULL,3,3,3,'uploaded/documents/3/3/3/','Tren','pdf',71082,'Se presentan datos sobre un molde de fundición metalurgica: Temperatura: 15.88°C, 4.76°C y 9.96°C Presión de aire: 36 l/min, 27.33 l/min y 20° (presión atmosférica) Espesor del molde: 32 mm, 54 mm y 2 cm (con diferentes espesores en distintas secciones) Se mencionan también dos tipos de datos adicionales: Dientes presentes en el molde: 56 dientes Ángulo de presión: 20° Espesor del molde con un error de medición modulo 2 mm y 3 mm',1,1,1);
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
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
  `reason` varchar(1024) COLLATE utf8mb4_unicode_ci NOT NULL,
  `sentAtEmail` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `subjectEmail` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `bodyEmail` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `statusEmail` int NOT NULL DEFAULT '1',
  `receivedAtEmail` datetime DEFAULT NULL,
  `idAppSource_fk` int NOT NULL,
  PRIMARY KEY (`idEmail`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `emails`
--

LOCK TABLES `emails` WRITE;
/*!40000 ALTER TABLE `emails` DISABLE KEYS */;
INSERT INTO `emails` VALUES (1,2,'¡Bienvenido a Signforce! Correo de verificación JSTFRTS200592','2025-11-26 14:36:29','thelegendofmax19@gmail.com','Email bienvenida',1,NULL,3),(2,2,'¡Te invitaron a firmar! Diego Reyes','2025-11-26 14:42:44','thelegendofmax19@gmail.com','From: thelegendofmax19@gmail.com\r\nTo: prof.jjnm.ockham@gmail.com\r\nSubject: ¡Te invitaron a firmar! Diego Reyes\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\nCorreo enviado de Diego Reyes\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n  <head>\r\n    <meta charset=\"UTF-8\" />\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />\r\n    <title>Invitación para Firmar Documento - Signforce</title>\r\n  </head>\r\n  <body style=\"margin:0; padding:0; background-color:#e4e9ec; font-family:\'Segoe UI\',system-ui,sans-serif; color:#ffffff;\">\r\n    <table align=\"center\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"max-width:820px; background-color:#2a2e3566; padding:20px;\">\r\n      <tr>\r\n        <td align=\"center\" style=\"padding:40px; background-color:#33333366; border-radius:16px; border:1px solid #444;\">\r\n          \r\n          <!-- HEADER -->\r\n          <h1 style=\"font-size:22px; color:#ffffff; margin-bottom:10px;\">¡Tienes un documento pendiente por firmar!</h1>\r\n          <p style=\"font-size:16px; color:#dddddd;\">\r\n            <strong>Diego Reyes</strong> te ha enviado un documento importante para su firma.\r\n          </p>\r\n\r\n          <!-- DOCUMENT INFO -->\r\n          <table width=\"100%\" cellpadding=\"0\" cellspacing=\"0\" style=\"margin-top:30px; background-color:#3a3a3a; border-radius:12px; padding:15px; border:1px solid #555;\">\r\n            <tr>\r\n              <td style=\"text-align:center; padding-bottom:10px;\">\r\n                <h3 style=\"color:#ffffff; font-size:18px;\">Información del Documento</h3>\r\n              </td>\r\n            </tr>\r\n            <tr>\r\n              <td>\r\n                <table width=\"100%\" cellpadding=\"4\" cellspacing=\"0\" style=\"color:#ffffff;\">\r\n                  <tr>\r\n                    <td width=\"50%\" valign=\"top\" style=\"font-size:14px;\">\r\n                      <p><strong>Remitente:</strong> Diego Reyes</p>\r\n                      <p><strong>Fecha de envío:</strong> 2025-11-26 14:42:43</p>\r\n                    </td>\r\n                    <td width=\"50%\" valign=\"top\" style=\"font-size:14px;\">\r\n                      <div style=\"background-color:#444; border-radius:8px; padding:10px; margin-top:8px;\">\r\n                        <strong>Mensaje del remitente:</strong>\r\n                        <p style=\"font-style:italic; margin:5px 0;\">firme ya por favor</p>\r\n                      </div>\r\n                    </td>\r\n                  </tr>\r\n                </table>\r\n              </td>\r\n            </tr>\r\n          </table>\r\n\r\n          <!-- CTA BUTTON -->\r\n          <div style=\"margin-top:25px; margin-bottom:25px;\">\r\n            <a href=\"http://192.168.1.66:8000/viewinvite?idInvite=44e2ede0-6d9f-4150-be28-625f769d98b6\" \r\n               style=\"background:linear-gradient(135deg,#2210cc,#320064); color:#ffffff; text-decoration:none; padding:14px 28px; border-radius:40px; font-weight:bold; display:inline-block;\">\r\n              Revisar y Firmar Documento\r\n            </a>\r\n          </div>\r\n\r\n          <!-- DIVIDER -->\r\n          <hr style=\"border:none; height:1px; background:linear-gradient(90deg,transparent,#999,transparent); margin:25px 0;\">\r\n\r\n          <!-- SECURITY NOTICE -->\r\n          <p style=\"color:#ff66cc; font-size:14px;\"><strong>Importante:</strong> Este enlace es personal e intransferible. Por seguridad, no lo compartas con nadie.</p>\r\n          <p style=\"font-size:14px; color:#cccccc;\">\r\n            Si tienes problemas para acceder, contacta a <strong>Diego Reyes</strong> al correo:\r\n            <a href=\"mailto:prof.jjnm.ockham@gmail.com\" style=\"color:#66aaff;\">prof.jjnm.ockham@gmail.com</a>\r\n          </p>\r\n          <p style=\"font-size:14px; color:#cccccc;\">\r\n            ¿Dudas sobre el proceso de firma? Escríbenos a \r\n            <a href=\"mailto:soporte@signforce.com\" style=\"color:#66aaff;\">soporte@signforce.com</a>\r\n          </p>\r\n\r\n          <!-- FOOTER -->\r\n          <hr style=\"border:none; height:1px; background:linear-gradient(90deg,transparent,#999,transparent); margin:25px 0;\">\r\n          <p style=\"font-size:12px; color:#aaaaaa;\">© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n          <p style=\"font-size:12px;\">\r\n            <a href=\"https://www.signforce.com/privacy\" style=\"color:#66aaff;\">Política de Privacidad</a> |\r\n            <a href=\"https://www.signforce.com/terms\" style=\"color:#66aaff;\">Términos de Servicio</a>\r\n          </p>\r\n\r\n        </td>\r\n      </tr>\r\n    </table>\r\n  </body>\r\n</html>\n',1,NULL,3),(3,2,'¡Te invitaron a firmar! Diego Reyes','2025-11-26 16:27:52','thelegendofmax19@gmail.com','From: thelegendofmax19@gmail.com\r\nTo: prof.jjnm.ockham@gmail.com\r\nSubject: ¡Te invitaron a firmar! Diego Reyes\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\nCorreo enviado de Diego Reyes\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n  <head>\r\n    <meta charset=\"UTF-8\" />\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />\r\n    <title>Invitación para Firmar Documento - Signforce</title>\r\n  </head>\r\n  <body style=\"margin:0; padding:0; background-color:#e4e9ec; font-family:\'Segoe UI\',system-ui,sans-serif; color:#ffffff;\">\r\n    <table align=\"center\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"max-width:820px; background-color:#2a2e3566; padding:20px;\">\r\n      <tr>\r\n        <td align=\"center\" style=\"padding:40px; background-color:#33333366; border-radius:16px; border:1px solid #444;\">\r\n          \r\n          <!-- HEADER -->\r\n          <h1 style=\"font-size:22px; color:#ffffff; margin-bottom:10px;\">¡Tienes un documento pendiente por firmar!</h1>\r\n          <p style=\"font-size:16px; color:#dddddd;\">\r\n            <strong>Diego Reyes</strong> te ha enviado un documento importante para su firma.\r\n          </p>\r\n\r\n          <!-- DOCUMENT INFO -->\r\n          <table width=\"100%\" cellpadding=\"0\" cellspacing=\"0\" style=\"margin-top:30px; background-color:#3a3a3a; border-radius:12px; padding:15px; border:1px solid #555;\">\r\n            <tr>\r\n              <td style=\"text-align:center; padding-bottom:10px;\">\r\n                <h3 style=\"color:#ffffff; font-size:18px;\">Información del Documento</h3>\r\n              </td>\r\n            </tr>\r\n            <tr>\r\n              <td>\r\n                <table width=\"100%\" cellpadding=\"4\" cellspacing=\"0\" style=\"color:#ffffff;\">\r\n                  <tr>\r\n                    <td width=\"50%\" valign=\"top\" style=\"font-size:14px;\">\r\n                      <p><strong>Remitente:</strong> Diego Reyes</p>\r\n                      <p><strong>Fecha de envío:</strong> 2025-11-26 16:27:51</p>\r\n                    </td>\r\n                    <td width=\"50%\" valign=\"top\" style=\"font-size:14px;\">\r\n                      <div style=\"background-color:#444; border-radius:8px; padding:10px; margin-top:8px;\">\r\n                        <strong>Mensaje del remitente:</strong>\r\n                        <p style=\"font-style:italic; margin:5px 0;\">firme ya</p>\r\n                      </div>\r\n                    </td>\r\n                  </tr>\r\n                </table>\r\n              </td>\r\n            </tr>\r\n          </table>\r\n\r\n          <!-- CTA BUTTON -->\r\n          <div style=\"margin-top:25px; margin-bottom:25px;\">\r\n            <a href=\"http://192.168.1.66:8000/viewinvite?idInvite=7d8cd08f-f0f5-4ffa-82e6-5bdcae314156\" \r\n               style=\"background:linear-gradient(135deg,#2210cc,#320064); color:#ffffff; text-decoration:none; padding:14px 28px; border-radius:40px; font-weight:bold; display:inline-block;\">\r\n              Revisar y Firmar Documento\r\n            </a>\r\n          </div>\r\n\r\n          <!-- DIVIDER -->\r\n          <hr style=\"border:none; height:1px; background:linear-gradient(90deg,transparent,#999,transparent); margin:25px 0;\">\r\n\r\n          <!-- SECURITY NOTICE -->\r\n          <p style=\"color:#ff66cc; font-size:14px;\"><strong>Importante:</strong> Este enlace es personal e intransferible. Por seguridad, no lo compartas con nadie.</p>\r\n          <p style=\"font-size:14px; color:#cccccc;\">\r\n            Si tienes problemas para acceder, contacta a <strong>Diego Reyes</strong> al correo:\r\n            <a href=\"mailto:prof.jjnm.ockham@gmail.com\" style=\"color:#66aaff;\">prof.jjnm.ockham@gmail.com</a>\r\n          </p>\r\n          <p style=\"font-size:14px; color:#cccccc;\">\r\n            ¿Dudas sobre el proceso de firma? Escríbenos a \r\n            <a href=\"mailto:soporte@signforce.com\" style=\"color:#66aaff;\">soporte@signforce.com</a>\r\n          </p>\r\n\r\n          <!-- FOOTER -->\r\n          <hr style=\"border:none; height:1px; background:linear-gradient(90deg,transparent,#999,transparent); margin:25px 0;\">\r\n          <p style=\"font-size:12px; color:#aaaaaa;\">© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n          <p style=\"font-size:12px;\">\r\n            <a href=\"https://www.signforce.com/privacy\" style=\"color:#66aaff;\">Política de Privacidad</a> |\r\n            <a href=\"https://www.signforce.com/terms\" style=\"color:#66aaff;\">Términos de Servicio</a>\r\n          </p>\r\n\r\n        </td>\r\n      </tr>\r\n    </table>\r\n  </body>\r\n</html>\n',1,NULL,3),(4,2,'¡Te invitaron a firmar! Diego Reyes','2025-11-26 16:30:55','thelegendofmax19@gmail.com','From: thelegendofmax19@gmail.com\r\nTo: prof.jjnm.ockham@gmail.com\r\nSubject: ¡Te invitaron a firmar! Diego Reyes\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\nCorreo enviado de Diego Reyes\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n  <head>\r\n    <meta charset=\"UTF-8\" />\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />\r\n    <title>Invitación para Firmar Documento - Signforce</title>\r\n  </head>\r\n  <body style=\"margin:0; padding:0; background-color:#e4e9ec; font-family:\'Segoe UI\',system-ui,sans-serif; color:#ffffff;\">\r\n    <table align=\"center\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"max-width:820px; background-color:#2a2e3566; padding:20px;\">\r\n      <tr>\r\n        <td align=\"center\" style=\"padding:40px; background-color:#33333366; border-radius:16px; border:1px solid #444;\">\r\n          \r\n          <!-- HEADER -->\r\n          <h1 style=\"font-size:22px; color:#ffffff; margin-bottom:10px;\">¡Tienes un documento pendiente por firmar!</h1>\r\n          <p style=\"font-size:16px; color:#dddddd;\">\r\n            <strong>Diego Reyes</strong> te ha enviado un documento importante para su firma.\r\n          </p>\r\n\r\n          <!-- DOCUMENT INFO -->\r\n          <table width=\"100%\" cellpadding=\"0\" cellspacing=\"0\" style=\"margin-top:30px; background-color:#3a3a3a; border-radius:12px; padding:15px; border:1px solid #555;\">\r\n            <tr>\r\n              <td style=\"text-align:center; padding-bottom:10px;\">\r\n                <h3 style=\"color:#ffffff; font-size:18px;\">Información del Documento</h3>\r\n              </td>\r\n            </tr>\r\n            <tr>\r\n              <td>\r\n                <table width=\"100%\" cellpadding=\"4\" cellspacing=\"0\" style=\"color:#ffffff;\">\r\n                  <tr>\r\n                    <td width=\"50%\" valign=\"top\" style=\"font-size:14px;\">\r\n                      <p><strong>Remitente:</strong> Diego Reyes</p>\r\n                      <p><strong>Fecha de envío:</strong> 2025-11-26 16:30:54</p>\r\n                    </td>\r\n                    <td width=\"50%\" valign=\"top\" style=\"font-size:14px;\">\r\n                      <div style=\"background-color:#444; border-radius:8px; padding:10px; margin-top:8px;\">\r\n                        <strong>Mensaje del remitente:</strong>\r\n                        <p style=\"font-style:italic; margin:5px 0;\">firme ya</p>\r\n                      </div>\r\n                    </td>\r\n                  </tr>\r\n                </table>\r\n              </td>\r\n            </tr>\r\n          </table>\r\n\r\n          <!-- CTA BUTTON -->\r\n          <div style=\"margin-top:25px; margin-bottom:25px;\">\r\n            <a href=\"http://192.168.1.66:8000/viewinvite?idInvite=261acc76-8ccd-4af4-8773-95f0d8e388fc\" \r\n               style=\"background:linear-gradient(135deg,#2210cc,#320064); color:#ffffff; text-decoration:none; padding:14px 28px; border-radius:40px; font-weight:bold; display:inline-block;\">\r\n              Revisar y Firmar Documento\r\n            </a>\r\n          </div>\r\n\r\n          <!-- DIVIDER -->\r\n          <hr style=\"border:none; height:1px; background:linear-gradient(90deg,transparent,#999,transparent); margin:25px 0;\">\r\n\r\n          <!-- SECURITY NOTICE -->\r\n          <p style=\"color:#ff66cc; font-size:14px;\"><strong>Importante:</strong> Este enlace es personal e intransferible. Por seguridad, no lo compartas con nadie.</p>\r\n          <p style=\"font-size:14px; color:#cccccc;\">\r\n            Si tienes problemas para acceder, contacta a <strong>Diego Reyes</strong> al correo:\r\n            <a href=\"mailto:prof.jjnm.ockham@gmail.com\" style=\"color:#66aaff;\">prof.jjnm.ockham@gmail.com</a>\r\n          </p>\r\n          <p style=\"font-size:14px; color:#cccccc;\">\r\n            ¿Dudas sobre el proceso de firma? Escríbenos a \r\n            <a href=\"mailto:soporte@signforce.com\" style=\"color:#66aaff;\">soporte@signforce.com</a>\r\n          </p>\r\n\r\n          <!-- FOOTER -->\r\n          <hr style=\"border:none; height:1px; background:linear-gradient(90deg,transparent,#999,transparent); margin:25px 0;\">\r\n          <p style=\"font-size:12px; color:#aaaaaa;\">© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n          <p style=\"font-size:12px;\">\r\n            <a href=\"https://www.signforce.com/privacy\" style=\"color:#66aaff;\">Política de Privacidad</a> |\r\n            <a href=\"https://www.signforce.com/terms\" style=\"color:#66aaff;\">Términos de Servicio</a>\r\n          </p>\r\n\r\n        </td>\r\n      </tr>\r\n    </table>\r\n  </body>\r\n</html>\n',1,NULL,3),(5,2,'¡Te invitaron a firmar! Diego Reyes','2025-11-26 16:39:38','thelegendofmax19@gmail.com','From: thelegendofmax19@gmail.com\r\nTo: prof.jjnm.ockham@gmail.com\r\nSubject: ¡Te invitaron a firmar! Diego Reyes\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\nCorreo enviado de Diego Reyes\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n  <head>\r\n    <meta charset=\"UTF-8\" />\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />\r\n    <title>Invitación para Firmar Documento - Signforce</title>\r\n  </head>\r\n  <body style=\"margin:0; padding:0; background-color:#e4e9ec; font-family:\'Segoe UI\',system-ui,sans-serif; color:#ffffff;\">\r\n    <table align=\"center\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"max-width:820px; background-color:#2a2e3566; padding:20px;\">\r\n      <tr>\r\n        <td align=\"center\" style=\"padding:40px; background-color:#33333366; border-radius:16px; border:1px solid #444;\">\r\n          \r\n          <!-- HEADER -->\r\n          <h1 style=\"font-size:22px; color:#ffffff; margin-bottom:10px;\">¡Tienes un documento pendiente por firmar!</h1>\r\n          <p style=\"font-size:16px; color:#dddddd;\">\r\n            <strong>Diego Reyes</strong> te ha enviado un documento importante para su firma.\r\n          </p>\r\n\r\n          <!-- DOCUMENT INFO -->\r\n          <table width=\"100%\" cellpadding=\"0\" cellspacing=\"0\" style=\"margin-top:30px; background-color:#3a3a3a; border-radius:12px; padding:15px; border:1px solid #555;\">\r\n            <tr>\r\n              <td style=\"text-align:center; padding-bottom:10px;\">\r\n                <h3 style=\"color:#ffffff; font-size:18px;\">Información del Documento</h3>\r\n              </td>\r\n            </tr>\r\n            <tr>\r\n              <td>\r\n                <table width=\"100%\" cellpadding=\"4\" cellspacing=\"0\" style=\"color:#ffffff;\">\r\n                  <tr>\r\n                    <td width=\"50%\" valign=\"top\" style=\"font-size:14px;\">\r\n                      <p><strong>Remitente:</strong> Diego Reyes</p>\r\n                      <p><strong>Fecha de envío:</strong> 2025-11-26 16:39:37</p>\r\n                    </td>\r\n                    <td width=\"50%\" valign=\"top\" style=\"font-size:14px;\">\r\n                      <div style=\"background-color:#444; border-radius:8px; padding:10px; margin-top:8px;\">\r\n                        <strong>Mensaje del remitente:</strong>\r\n                        <p style=\"font-style:italic; margin:5px 0;\">firme ya</p>\r\n                      </div>\r\n                    </td>\r\n                  </tr>\r\n                </table>\r\n              </td>\r\n            </tr>\r\n          </table>\r\n\r\n          <!-- CTA BUTTON -->\r\n          <div style=\"margin-top:25px; margin-bottom:25px;\">\r\n            <a href=\"http://192.168.1.66:8000/viewinvite?idInvite=e76a3868-47c5-44b9-88ed-b8e946ccfa25\" \r\n               style=\"background:linear-gradient(135deg,#2210cc,#320064); color:#ffffff; text-decoration:none; padding:14px 28px; border-radius:40px; font-weight:bold; display:inline-block;\">\r\n              Revisar y Firmar Documento\r\n            </a>\r\n          </div>\r\n\r\n          <!-- DIVIDER -->\r\n          <hr style=\"border:none; height:1px; background:linear-gradient(90deg,transparent,#999,transparent); margin:25px 0;\">\r\n\r\n          <!-- SECURITY NOTICE -->\r\n          <p style=\"color:#ff66cc; font-size:14px;\"><strong>Importante:</strong> Este enlace es personal e intransferible. Por seguridad, no lo compartas con nadie.</p>\r\n          <p style=\"font-size:14px; color:#cccccc;\">\r\n            Si tienes problemas para acceder, contacta a <strong>Diego Reyes</strong> al correo:\r\n            <a href=\"mailto:prof.jjnm.ockham@gmail.com\" style=\"color:#66aaff;\">prof.jjnm.ockham@gmail.com</a>\r\n          </p>\r\n          <p style=\"font-size:14px; color:#cccccc;\">\r\n            ¿Dudas sobre el proceso de firma? Escríbenos a \r\n            <a href=\"mailto:soporte@signforce.com\" style=\"color:#66aaff;\">soporte@signforce.com</a>\r\n          </p>\r\n\r\n          <!-- FOOTER -->\r\n          <hr style=\"border:none; height:1px; background:linear-gradient(90deg,transparent,#999,transparent); margin:25px 0;\">\r\n          <p style=\"font-size:12px; color:#aaaaaa;\">© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n          <p style=\"font-size:12px;\">\r\n            <a href=\"https://www.signforce.com/privacy\" style=\"color:#66aaff;\">Política de Privacidad</a> |\r\n            <a href=\"https://www.signforce.com/terms\" style=\"color:#66aaff;\">Términos de Servicio</a>\r\n          </p>\r\n\r\n        </td>\r\n      </tr>\r\n    </table>\r\n  </body>\r\n</html>\n',1,NULL,3),(6,2,'¡Te invitaron a firmar! Diego Reyes','2025-11-26 16:39:38','thelegendofmax19@gmail.com','From: thelegendofmax19@gmail.com\r\nTo: prof.jjnm.ockham@gmail.com\r\nSubject: ¡Te invitaron a firmar! Diego Reyes\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\nCorreo enviado de Diego Reyes\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n  <head>\r\n    <meta charset=\"UTF-8\" />\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />\r\n    <title>Invitación para Firmar Documento - Signforce</title>\r\n  </head>\r\n  <body style=\"margin:0; padding:0; background-color:#e4e9ec; font-family:\'Segoe UI\',system-ui,sans-serif; color:#ffffff;\">\r\n    <table align=\"center\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"max-width:820px; background-color:#2a2e3566; padding:20px;\">\r\n      <tr>\r\n        <td align=\"center\" style=\"padding:40px; background-color:#33333366; border-radius:16px; border:1px solid #444;\">\r\n          \r\n          <!-- HEADER -->\r\n          <h1 style=\"font-size:22px; color:#ffffff; margin-bottom:10px;\">¡Tienes un documento pendiente por firmar!</h1>\r\n          <p style=\"font-size:16px; color:#dddddd;\">\r\n            <strong>Diego Reyes</strong> te ha enviado un documento importante para su firma.\r\n          </p>\r\n\r\n          <!-- DOCUMENT INFO -->\r\n          <table width=\"100%\" cellpadding=\"0\" cellspacing=\"0\" style=\"margin-top:30px; background-color:#3a3a3a; border-radius:12px; padding:15px; border:1px solid #555;\">\r\n            <tr>\r\n              <td style=\"text-align:center; padding-bottom:10px;\">\r\n                <h3 style=\"color:#ffffff; font-size:18px;\">Información del Documento</h3>\r\n              </td>\r\n            </tr>\r\n            <tr>\r\n              <td>\r\n                <table width=\"100%\" cellpadding=\"4\" cellspacing=\"0\" style=\"color:#ffffff;\">\r\n                  <tr>\r\n                    <td width=\"50%\" valign=\"top\" style=\"font-size:14px;\">\r\n                      <p><strong>Remitente:</strong> Diego Reyes</p>\r\n                      <p><strong>Fecha de envío:</strong> 2025-11-26 16:39:37</p>\r\n                    </td>\r\n                    <td width=\"50%\" valign=\"top\" style=\"font-size:14px;\">\r\n                      <div style=\"background-color:#444; border-radius:8px; padding:10px; margin-top:8px;\">\r\n                        <strong>Mensaje del remitente:</strong>\r\n                        <p style=\"font-style:italic; margin:5px 0;\">firme ya</p>\r\n                      </div>\r\n                    </td>\r\n                  </tr>\r\n                </table>\r\n              </td>\r\n            </tr>\r\n          </table>\r\n\r\n          <!-- CTA BUTTON -->\r\n          <div style=\"margin-top:25px; margin-bottom:25px;\">\r\n            <a href=\"http://192.168.1.66:8000/viewinvite?idInvite=e76a3868-47c5-44b9-88ed-b8e946ccfa25\" \r\n               style=\"background:linear-gradient(135deg,#2210cc,#320064); color:#ffffff; text-decoration:none; padding:14px 28px; border-radius:40px; font-weight:bold; display:inline-block;\">\r\n              Revisar y Firmar Documento\r\n            </a>\r\n          </div>\r\n\r\n          <!-- DIVIDER -->\r\n          <hr style=\"border:none; height:1px; background:linear-gradient(90deg,transparent,#999,transparent); margin:25px 0;\">\r\n\r\n          <!-- SECURITY NOTICE -->\r\n          <p style=\"color:#ff66cc; font-size:14px;\"><strong>Importante:</strong> Este enlace es personal e intransferible. Por seguridad, no lo compartas con nadie.</p>\r\n          <p style=\"font-size:14px; color:#cccccc;\">\r\n            Si tienes problemas para acceder, contacta a <strong>Diego Reyes</strong> al correo:\r\n            <a href=\"mailto:prof.jjnm.ockham@gmail.com\" style=\"color:#66aaff;\">prof.jjnm.ockham@gmail.com</a>\r\n          </p>\r\n          <p style=\"font-size:14px; color:#cccccc;\">\r\n            ¿Dudas sobre el proceso de firma? Escríbenos a \r\n            <a href=\"mailto:soporte@signforce.com\" style=\"color:#66aaff;\">soporte@signforce.com</a>\r\n          </p>\r\n\r\n          <!-- FOOTER -->\r\n          <hr style=\"border:none; height:1px; background:linear-gradient(90deg,transparent,#999,transparent); margin:25px 0;\">\r\n          <p style=\"font-size:12px; color:#aaaaaa;\">© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n          <p style=\"font-size:12px;\">\r\n            <a href=\"https://www.signforce.com/privacy\" style=\"color:#66aaff;\">Política de Privacidad</a> |\r\n            <a href=\"https://www.signforce.com/terms\" style=\"color:#66aaff;\">Términos de Servicio</a>\r\n          </p>\r\n\r\n        </td>\r\n      </tr>\r\n    </table>\r\n  </body>\r\n</html>\n',1,NULL,3),(7,2,'¡Te invitaron a firmar! Diego Reyes','2025-11-26 16:39:39','thelegendofmax19@gmail.com','Content-Type: text/html; charset=\"UTF-8\"\r\nFrom: thelegendofmax19@gmail.com\r\nTo: prof.jjnm.ockham@gmail.com\r\nSubject: ¡Te invitaron a firmar! Diego Reyes\r\n\r\nCorreo enviado de Diego Reyes\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n  <head>\r\n    <meta charset=\"UTF-8\" />\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />\r\n    <title>Invitación para Firmar Documento - Signforce</title>\r\n  </head>\r\n  <body style=\"margin:0; padding:0; background-color:#e4e9ec; font-family:\'Segoe UI\',system-ui,sans-serif; color:#ffffff;\">\r\n    <table align=\"center\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"max-width:820px; background-color:#2a2e3566; padding:20px;\">\r\n      <tr>\r\n        <td align=\"center\" style=\"padding:40px; background-color:#33333366; border-radius:16px; border:1px solid #444;\">\r\n          \r\n          <!-- HEADER -->\r\n          <h1 style=\"font-size:22px; color:#ffffff; margin-bottom:10px;\">¡Tienes un documento pendiente por firmar!</h1>\r\n          <p style=\"font-size:16px; color:#dddddd;\">\r\n            <strong>Diego Reyes</strong> te ha enviado un documento importante para su firma.\r\n          </p>\r\n\r\n          <!-- DOCUMENT INFO -->\r\n          <table width=\"100%\" cellpadding=\"0\" cellspacing=\"0\" style=\"margin-top:30px; background-color:#3a3a3a; border-radius:12px; padding:15px; border:1px solid #555;\">\r\n            <tr>\r\n              <td style=\"text-align:center; padding-bottom:10px;\">\r\n                <h3 style=\"color:#ffffff; font-size:18px;\">Información del Documento</h3>\r\n              </td>\r\n            </tr>\r\n            <tr>\r\n              <td>\r\n                <table width=\"100%\" cellpadding=\"4\" cellspacing=\"0\" style=\"color:#ffffff;\">\r\n                  <tr>\r\n                    <td width=\"50%\" valign=\"top\" style=\"font-size:14px;\">\r\n                      <p><strong>Remitente:</strong> Diego Reyes</p>\r\n                      <p><strong>Fecha de envío:</strong> 2025-11-26 16:39:37</p>\r\n                    </td>\r\n                    <td width=\"50%\" valign=\"top\" style=\"font-size:14px;\">\r\n                      <div style=\"background-color:#444; border-radius:8px; padding:10px; margin-top:8px;\">\r\n                        <strong>Mensaje del remitente:</strong>\r\n                        <p style=\"font-style:italic; margin:5px 0;\">firme ya</p>\r\n                      </div>\r\n                    </td>\r\n                  </tr>\r\n                </table>\r\n              </td>\r\n            </tr>\r\n          </table>\r\n\r\n          <!-- CTA BUTTON -->\r\n          <div style=\"margin-top:25px; margin-bottom:25px;\">\r\n            <a href=\"http://192.168.1.66:8000/viewinvite?idInvite=e76a3868-47c5-44b9-88ed-b8e946ccfa25\" \r\n               style=\"background:linear-gradient(135deg,#2210cc,#320064); color:#ffffff; text-decoration:none; padding:14px 28px; border-radius:40px; font-weight:bold; display:inline-block;\">\r\n              Revisar y Firmar Documento\r\n            </a>\r\n          </div>\r\n\r\n          <!-- DIVIDER -->\r\n          <hr style=\"border:none; height:1px; background:linear-gradient(90deg,transparent,#999,transparent); margin:25px 0;\">\r\n\r\n          <!-- SECURITY NOTICE -->\r\n          <p style=\"color:#ff66cc; font-size:14px;\"><strong>Importante:</strong> Este enlace es personal e intransferible. Por seguridad, no lo compartas con nadie.</p>\r\n          <p style=\"font-size:14px; color:#cccccc;\">\r\n            Si tienes problemas para acceder, contacta a <strong>Diego Reyes</strong> al correo:\r\n            <a href=\"mailto:prof.jjnm.ockham@gmail.com\" style=\"color:#66aaff;\">prof.jjnm.ockham@gmail.com</a>\r\n          </p>\r\n          <p style=\"font-size:14px; color:#cccccc;\">\r\n            ¿Dudas sobre el proceso de firma? Escríbenos a \r\n            <a href=\"mailto:soporte@signforce.com\" style=\"color:#66aaff;\">soporte@signforce.com</a>\r\n          </p>\r\n\r\n          <!-- FOOTER -->\r\n          <hr style=\"border:none; height:1px; background:linear-gradient(90deg,transparent,#999,transparent); margin:25px 0;\">\r\n          <p style=\"font-size:12px; color:#aaaaaa;\">© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n          <p style=\"font-size:12px;\">\r\n            <a href=\"https://www.signforce.com/privacy\" style=\"color:#66aaff;\">Política de Privacidad</a> |\r\n            <a href=\"https://www.signforce.com/terms\" style=\"color:#66aaff;\">Términos de Servicio</a>\r\n          </p>\r\n\r\n        </td>\r\n      </tr>\r\n    </table>\r\n  </body>\r\n</html>\n',1,NULL,3),(8,2,'¡Te invitaron a firmar! Diego Reyes','2025-11-26 17:18:22','thelegendofmax19@gmail.com','Content-Type: text/html; charset=\"UTF-8\"\r\nFrom: thelegendofmax19@gmail.com\r\nTo: prof.jjnm.ockham@gmail.com\r\nSubject: ¡Te invitaron a firmar! Diego Reyes\r\n\r\nCorreo enviado de Diego Reyes\n\n<!DOCTYPE html>\r\n<html lang=\"es\">\r\n  <head>\r\n    <meta charset=\"UTF-8\" />\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />\r\n    <title>Invitación para Firmar Documento - Signforce</title>\r\n  </head>\r\n  <body style=\"margin:0; padding:0; background-color:#e4e9ec; font-family:\'Segoe UI\',system-ui,sans-serif; color:#ffffff;\">\r\n    <table align=\"center\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"max-width:820px; background-color:#2a2e3566; padding:20px;\">\r\n      <tr>\r\n        <td align=\"center\" style=\"padding:40px; background-color:#33333366; border-radius:16px; border:1px solid #444;\">\r\n          \r\n          <!-- HEADER -->\r\n          <h1 style=\"font-size:22px; color:#ffffff; margin-bottom:10px;\">¡Tienes un documento pendiente por firmar!</h1>\r\n          <p style=\"font-size:16px; color:#dddddd;\">\r\n            <strong>Diego Reyes</strong> te ha enviado un documento importante para su firma.\r\n          </p>\r\n\r\n          <!-- DOCUMENT INFO -->\r\n          <table width=\"100%\" cellpadding=\"0\" cellspacing=\"0\" style=\"margin-top:30px; background-color:#3a3a3a; border-radius:12px; padding:15px; border:1px solid #555;\">\r\n            <tr>\r\n              <td style=\"text-align:center; padding-bottom:10px;\">\r\n                <h3 style=\"color:#ffffff; font-size:18px;\">Información del Documento</h3>\r\n              </td>\r\n            </tr>\r\n            <tr>\r\n              <td>\r\n                <table width=\"100%\" cellpadding=\"4\" cellspacing=\"0\" style=\"color:#ffffff;\">\r\n                  <tr>\r\n                    <td width=\"50%\" valign=\"top\" style=\"font-size:14px;\">\r\n                      <p><strong>Remitente:</strong> Diego Reyes</p>\r\n                      <p><strong>Fecha de envío:</strong> 2025-11-26 17:18:20</p>\r\n                    </td>\r\n                    <td width=\"50%\" valign=\"top\" style=\"font-size:14px;\">\r\n                      <div style=\"background-color:#444; border-radius:8px; padding:10px; margin-top:8px;\">\r\n                        <strong>Mensaje del remitente:</strong>\r\n                        <p style=\"font-style:italic; margin:5px 0;\">firme ya</p>\r\n                      </div>\r\n                    </td>\r\n                  </tr>\r\n                </table>\r\n              </td>\r\n            </tr>\r\n          </table>\r\n\r\n          <!-- CTA BUTTON -->\r\n          <div style=\"margin-top:25px; margin-bottom:25px;\">\r\n            <a href=\"http://192.168.1.66:8000/viewinvite?idInvite=e489641a-7850-425c-a7c4-1c0fff6a0dd1\" \r\n               style=\"background:linear-gradient(135deg,#2210cc,#320064); color:#ffffff; text-decoration:none; padding:14px 28px; border-radius:40px; font-weight:bold; display:inline-block;\">\r\n              Revisar y Firmar Documento\r\n            </a>\r\n          </div>\r\n\r\n          <!-- DIVIDER -->\r\n          <hr style=\"border:none; height:1px; background:linear-gradient(90deg,transparent,#999,transparent); margin:25px 0;\">\r\n\r\n          <!-- SECURITY NOTICE -->\r\n          <p style=\"color:#ff66cc; font-size:14px;\"><strong>Importante:</strong> Este enlace es personal e intransferible. Por seguridad, no lo compartas con nadie.</p>\r\n          <p style=\"font-size:14px; color:#cccccc;\">\r\n            Si tienes problemas para acceder, contacta a <strong>Diego Reyes</strong> al correo:\r\n            <a href=\"mailto:prof.jjnm.ockham@gmail.com\" style=\"color:#66aaff;\">prof.jjnm.ockham@gmail.com</a>\r\n          </p>\r\n          <p style=\"font-size:14px; color:#cccccc;\">\r\n            ¿Dudas sobre el proceso de firma? Escríbenos a \r\n            <a href=\"mailto:soporte@signforce.com\" style=\"color:#66aaff;\">soporte@signforce.com</a>\r\n          </p>\r\n\r\n          <!-- FOOTER -->\r\n          <hr style=\"border:none; height:1px; background:linear-gradient(90deg,transparent,#999,transparent); margin:25px 0;\">\r\n          <p style=\"font-size:12px; color:#aaaaaa;\">© 2025 SIGNFORCE S.A. DE C.V. Todos los derechos reservados.</p>\r\n          <p style=\"font-size:12px;\">\r\n            <a href=\"https://www.signforce.com/privacy\" style=\"color:#66aaff;\">Política de Privacidad</a> |\r\n            <a href=\"https://www.signforce.com/terms\" style=\"color:#66aaff;\">Términos de Servicio</a>\r\n          </p>\r\n\r\n        </td>\r\n      </tr>\r\n    </table>\r\n  </body>\r\n</html>\n',1,NULL,3),(9,3,'¡Bienvenido a Signforce! Correo de verificación AVGSBV65SDHJB','2025-11-26 18:04:54','thelegendofmax19@gmail.com','Email bienvenida',1,NULL,3);
/*!40000 ALTER TABLE `emails` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `folderdocuments`
--

DROP TABLE IF EXISTS `folderdocuments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `folderdocuments` (
  `idfolderdocument` varchar(50) NOT NULL,
  `idDocument` bigint NOT NULL,
  `idFolder` bigint NOT NULL,
  PRIMARY KEY (`idfolderdocument`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `folderdocuments`
--

LOCK TABLES `folderdocuments` WRITE;
/*!40000 ALTER TABLE `folderdocuments` DISABLE KEYS */;
INSERT INTO `folderdocuments` VALUES ('125334d1-bcb9-42e5-86ed-0b47e2448bfe',3,2),('759c325e-bf9a-4994-a6aa-e1c6546e2140',6,5),('82d6c8a0-c813-42d8-8412-f9a59da62c4e',1,1),('881ec49c-c3cf-4318-b53d-093271ec9eb0',3,4),('97313e3e-8656-4677-88e2-e068dba69532',4,3),('a9437e58-ac54-43fa-91f1-1a4254462ae3',2,3),('c73d3105-37c2-4a86-a70b-7e657c7efb8a',3,3),('c9425ac2-2ff3-4e93-a57c-0ff1099f02e7',5,5),('e67e08aa-fb12-472d-8d10-52a87c076dd9',4,2),('f42e18b2-9248-4dac-a827-870f4a7b9333',2,1);
/*!40000 ALTER TABLE `folderdocuments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `folders`
--

DROP TABLE IF EXISTS `folders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `folders` (
  `idFolder` bigint NOT NULL AUTO_INCREMENT,
  `ownerInst_fk` bigint NOT NULL,
  `ownerTeam_fk` bigint NOT NULL,
  `creatorUser_fk` bigint NOT NULL,
  `creationAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `lastModified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deletedAt` datetime DEFAULT NULL,
  `deletedReason` text COLLATE utf8mb4_unicode_ci,
  `purpose` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `closedAt` datetime DEFAULT NULL,
  `secuentialSign` tinyint(1) NOT NULL DEFAULT '0',
  `pathSerialized` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `expirationDate` datetime NOT NULL,
  `numDocs` int NOT NULL,
  `numDocsSign` int NOT NULL DEFAULT '0',
  `numSigners` int NOT NULL DEFAULT '0',
  `numReceivers` int NOT NULL DEFAULT '0',
  `completedAt` datetime DEFAULT NULL,
  PRIMARY KEY (`idFolder`),
  KEY `fk_folders_creator_user` (`creatorUser_fk`),
  KEY `fk_folders_owner_inst` (`ownerInst_fk`),
  KEY `fk_folders_owner_team` (`ownerTeam_fk`),
  CONSTRAINT `fk_folders_owner_inst` FOREIGN KEY (`ownerInst_fk`) REFERENCES `institutions` (`idInstitution`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `folders`
--

LOCK TABLES `folders` WRITE;
/*!40000 ALTER TABLE `folders` DISABLE KEYS */;
INSERT INTO `folders` VALUES (1,2,2,2,'2025-11-26 14:41:49','2025-11-26 14:41:49',NULL,NULL,NULL,NULL,NULL,0,'./temp/folders/2/2/2/','2025-11-29 14:41:49',2,0,0,0,NULL),(2,2,2,2,'2025-11-26 16:27:05','2025-11-26 16:27:05',NULL,NULL,NULL,NULL,NULL,0,'./temp/folders/2/2/2/','2025-11-29 16:27:05',2,0,0,0,NULL),(3,2,2,2,'2025-11-26 16:38:57','2025-11-26 16:38:57',NULL,NULL,NULL,NULL,NULL,0,'./temp/folders/2/2/2/','2025-11-29 16:38:57',3,0,0,0,NULL),(4,2,2,2,'2025-11-26 17:16:22','2025-11-26 17:16:22',NULL,NULL,NULL,NULL,NULL,0,'./temp/folders/2/2/2/','2025-11-29 17:16:22',1,0,0,0,NULL),(5,3,3,3,'2025-11-26 18:34:00','2025-11-26 18:34:00',NULL,NULL,NULL,NULL,NULL,0,'./temp/folders/3/3/3/','2025-11-29 18:34:00',2,0,0,0,NULL);
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
  `legalNameInst` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `aliasNameInst` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `taxNumInst` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL,
  `streetAddress` text COLLATE utf8mb4_unicode_ci,
  `addressLine` text COLLATE utf8mb4_unicode_ci,
  `postalCode` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `neighborhood` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `locality` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `stateCodeInst_fk` int DEFAULT '1',
  `countryCodeInst_fk` int DEFAULT '1',
  `formattedAddress` text COLLATE utf8mb4_unicode_ci,
  `contactPhoneInst` varchar(15) COLLATE utf8mb4_unicode_ci NOT NULL,
  `contactEmailInst` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL,
  `createdAtInst` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `lastModifiedInst` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deletedAtInst` datetime DEFAULT CURRENT_TIMESTAMP,
  `statusInst_fk` int NOT NULL DEFAULT '1',
  `activeInst` tinyint(1) NOT NULL DEFAULT '0',
  `typeContractInst` int NOT NULL DEFAULT '0',
  `paymentDataInst_fk` int DEFAULT NULL,
  `logoUrlInst` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `legalSignupName` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `legalSignupLastname` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `registerSignupName` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `registerSignupLastname` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `rootUser_fk` bigint DEFAULT NULL,
  PRIMARY KEY (`idInstitution`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `institutions`
--

LOCK TABLES `institutions` WRITE;
/*!40000 ALTER TABLE `institutions` DISABLE KEYS */;
INSERT INTO `institutions` VALUES (1,'SIGNFORCE S.A. DE C.V.','SIGNFORCE','SIGPBA901201','Laguna de la mancha 62',NULL,'11520','Granada','Miguel Hidalgo',1,1,NULL,'5529996911','do.c.001@hotmail.com','2025-11-26 13:56:15','2025-11-26 13:56:15','2025-11-26 13:56:15',8,1,0,NULL,NULL,'Juan','Nuñez','Juan','Nuñez',1),(2,'JUSTFRUITS S.A. DE C.V.','JUSTFRUITS','JSTFRTS200592','Laguna de la mancha 64',NULL,'11520','Granada','Miguel Hidalgo',1,1,NULL,'5511223344','prof.jjnm.ockham@gmail.com','2025-11-26 14:36:28','2025-11-26 14:36:28','2025-11-26 14:36:28',6,0,0,2,NULL,'Diego','Reyes','Diego','Reyes',2),(3,'AIDAVAGA LTD.','AIDAVAGA','AVGSBV65SDHJB','Niños Heroes 125',NULL,'11225','Niños Heroes','Benito Juarez',1,1,NULL,'5510317898','contacto.somostec@gmail.com','2025-11-26 18:04:53','2025-11-26 18:04:53','2025-11-26 18:04:53',7,0,0,3,NULL,'Aida','Valencia','Aida','Valencia',3);
/*!40000 ALTER TABLE `institutions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invites`
--

DROP TABLE IF EXISTS `invites`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `invites` (
  `idInvite` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `idFolder` bigint NOT NULL,
  `idUserOwnner_fk` bigint NOT NULL,
  `idUserDest_fk` bigint NOT NULL,
  `requireAliveProof` tinyint(1) NOT NULL DEFAULT '0',
  `expirationDate` datetime DEFAULT NULL,
  `sentAt` datetime DEFAULT NULL,
  `openedAt` datetime DEFAULT NULL,
  `closedAt` datetime DEFAULT NULL,
  `descriptionText` text COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`idInvite`),
  KEY `fk_invite_target_user` (`idUserDest_fk`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invites`
--

LOCK TABLES `invites` WRITE;
/*!40000 ALTER TABLE `invites` DISABLE KEYS */;
INSERT INTO `invites` VALUES ('261acc76-8ccd-4af4-8773-95f0d8e388fc',2,2,2,0,'2025-11-26 00:00:00','2025-11-26 16:30:54',NULL,NULL,'firme ya'),('44e2ede0-6d9f-4150-be28-625f769d98b6',1,2,2,0,'2025-11-28 00:00:00','2025-11-26 14:42:43',NULL,NULL,'firme ya por favor'),('7d8cd08f-f0f5-4ffa-82e6-5bdcae314156',2,2,2,0,'2025-11-26 00:00:00','2025-11-26 16:27:51',NULL,NULL,'firme ya'),('e489641a-7850-425c-a7c4-1c0fff6a0dd1',4,2,2,0,'2025-11-26 00:00:00','2025-11-26 17:18:20',NULL,NULL,'firme ya'),('e76a3868-47c5-44b9-88ed-b8e946ccfa25',3,2,2,0,'2025-11-26 00:00:00','2025-11-26 16:39:37',NULL,NULL,'firme ya');
/*!40000 ALTER TABLE `invites` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invitesdetail`
--

DROP TABLE IF EXISTS `invitesdetail`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `invitesdetail` (
  `idInviteDetail` bigint NOT NULL AUTO_INCREMENT,
  `idInvite` varchar(50) NOT NULL,
  `idfolderdocument` varchar(50) NOT NULL,
  PRIMARY KEY (`idInviteDetail`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invitesdetail`
--

LOCK TABLES `invitesdetail` WRITE;
/*!40000 ALTER TABLE `invitesdetail` DISABLE KEYS */;
INSERT INTO `invitesdetail` VALUES (1,'44e2ede0-6d9f-4150-be28-625f769d98b6','82d6c8a0-c813-42d8-8412-f9a59da62c4e'),(2,'44e2ede0-6d9f-4150-be28-625f769d98b6','f42e18b2-9248-4dac-a827-870f4a7b9333'),(3,'7d8cd08f-f0f5-4ffa-82e6-5bdcae314156','125334d1-bcb9-42e5-86ed-0b47e2448bfe'),(4,'7d8cd08f-f0f5-4ffa-82e6-5bdcae314156','e67e08aa-fb12-472d-8d10-52a87c076dd9'),(5,'261acc76-8ccd-4af4-8773-95f0d8e388fc','125334d1-bcb9-42e5-86ed-0b47e2448bfe'),(6,'261acc76-8ccd-4af4-8773-95f0d8e388fc','e67e08aa-fb12-472d-8d10-52a87c076dd9'),(7,'e76a3868-47c5-44b9-88ed-b8e946ccfa25','97313e3e-8656-4677-88e2-e068dba69532'),(8,'e76a3868-47c5-44b9-88ed-b8e946ccfa25','a9437e58-ac54-43fa-91f1-1a4254462ae3'),(9,'e76a3868-47c5-44b9-88ed-b8e946ccfa25','c73d3105-37c2-4a86-a70b-7e657c7efb8a'),(10,'e76a3868-47c5-44b9-88ed-b8e946ccfa25','a9437e58-ac54-43fa-91f1-1a4254462ae3'),(11,'e76a3868-47c5-44b9-88ed-b8e946ccfa25','c73d3105-37c2-4a86-a70b-7e657c7efb8a'),(12,'e76a3868-47c5-44b9-88ed-b8e946ccfa25','97313e3e-8656-4677-88e2-e068dba69532'),(13,'e76a3868-47c5-44b9-88ed-b8e946ccfa25','97313e3e-8656-4677-88e2-e068dba69532'),(14,'e76a3868-47c5-44b9-88ed-b8e946ccfa25','a9437e58-ac54-43fa-91f1-1a4254462ae3'),(15,'e76a3868-47c5-44b9-88ed-b8e946ccfa25','c73d3105-37c2-4a86-a70b-7e657c7efb8a'),(16,'e489641a-7850-425c-a7c4-1c0fff6a0dd1','881ec49c-c3cf-4318-b53d-093271ec9eb0');
/*!40000 ALTER TABLE `invitesdetail` ENABLE KEYS */;
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
  `documentHash` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `documentExt` varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `documentName` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `documentClass` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `documentPath` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `expirationDate` datetime DEFAULT NULL,
  `createdAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`idInsttitution_fk`,`idUser_fk`,`documentHash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `kyc`
--

LOCK TABLES `kyc` WRITE;
/*!40000 ALTER TABLE `kyc` DISABLE KEYS */;
INSERT INTO `kyc` VALUES (2,2,'MgSgDxvqT0C2jGPqrli/danEJZkdIpLPGplpm0YGtJk=','png','3dxmodeling','IdentidadOficial','temp/2/2/',NULL,'2025-11-26 14:37:22'),(2,2,'OXUcBI8Keum4KcUd1Mqbg6rTNLe6fINCoqhS3aTcKTg=','pdf','CV_2025_JJNM','PoderRepresentante','temp/2/2/',NULL,'2025-11-26 14:37:18'),(2,2,'u1dImPaqo3+qhM4kICFk6tMXPSnms33GOFQFwL7e9AM=','pdf','Confirmación _ Viva','ActaConstitutiva','temp/2/2/',NULL,'2025-11-26 14:37:12'),(2,2,'WBmwaNYeqZ+Hr10vIqgq2Kvr3j23u/g3sbMCs98joq0=','pdf','Diagrama_AWS','PruebaResidencia','temp/2/2/',NULL,'2025-11-26 14:37:26'),(3,3,'C6wfK7XKQRtQvHuUd0CWi9sMy2WXqwSNt5EX/ZKIon4=','PDF','Ensamblaje1','ActaConstitutiva','temp/3/3/',NULL,'2025-11-26 18:21:00'),(3,3,'frDMWhwmOBvEJvbXwB7SKM9VEbJq7yI97JgNbKY3/ww=','PDF','Tren','PruebaResidencia','temp/3/3/',NULL,'2025-11-26 18:21:15'),(3,3,'oR/5ro0qPrNs6wgbQRRfNk3RonWOs/FvtGiTxeCo3Nk=','PDF','tolva croquis','IdentidadOficial','temp/3/3/',NULL,'2025-11-26 18:21:12'),(3,3,'wCVeNEyl6+LS6wbc01ZadmcLdkcGooaTa4pru7Qa29U=','PDF','Lista de materiales','PoderRepresentante','temp/3/3/',NULL,'2025-11-26 18:21:09');
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
  `domainApp` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'localhost',
  `portApp` int NOT NULL DEFAULT '8000',
  `nameApp` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'SignforceApp',
  `description` varchar(300) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `isActive` tinyint(1) NOT NULL DEFAULT '1',
  `versionApp` int NOT NULL DEFAULT '1',
  `currPathApp` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT 'C:/',
  PRIMARY KEY (`idapp`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `microapps`
--

LOCK TABLES `microapps` WRITE;
/*!40000 ALTER TABLE `microapps` DISABLE KEYS */;
INSERT INTO `microapps` VALUES (3,'localhost',8000,'sfmiddle','Middleware para comunicar con app y pagina principal',1,1,'C:/Users/USER/Desktop/signForce/sfmiddle'),(4,'localhost',8002,'emailServ','App de envío de correo electrónico y mensajería',1,1,'C:/Users/USER/Desktop/signForce/emailServ'),(5,'localhost',8001,'sfback','Backend de firma electrónica y operaciones critpograficas',1,1,'C:/Users/USER/Desktop/signForce/sfback'),(6,'192.168.1.66',4999,'llmServ','Api para consumir modelo de LLM',1,1,'C:/Users/USER/Desktop/signForce/sfia');
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
  `cardNumber` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `expirationDate` varchar(6) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `nameOwner` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`idInstitution`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `payment`
--

LOCK TABLES `payment` WRITE;
/*!40000 ALTER TABLE `payment` DISABLE KEYS */;
INSERT INTO `payment` VALUES (1,1,3,'2025-12-26 14:38:54','1234567890987654','01/35','Diego Leonardo Reyes Chavez'),(2,1,3,'2025-12-26 17:48:45','1234123412412341','12/54','Juan de Jesus Nuñez Mendoza'),(3,1,3,'2025-12-26 18:25:16','2342135245345623','12/41','Aida Valencia Galicia');
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
  `description` text COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`role_id`,`permission_id`),
  KEY `fk_permissions_permission` (`permission_id`),
  CONSTRAINT `fk_permissions_role` FOREIGN KEY (`role_id`) REFERENCES `userroleplatform` (`id`) ON DELETE CASCADE,
  CONSTRAINT `platformrolepermissions_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `userroleplatform` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
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
  `idUser_fk` bigint NOT NULL,
  `idInvite_fk` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `idUserKeys_fk` bigint DEFAULT NULL,
  `digestValueSign` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `digestAlgoSign_fk` int NOT NULL DEFAULT '21',
  `signatureAlgoSign_fk` int NOT NULL DEFAULT '9',
  `signatureValueSign` text COLLATE utf8mb4_unicode_ci,
  `genTimeSign` datetime DEFAULT NULL,
  `pathSign` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `typeSign_fk` int DEFAULT NULL,
  `nonceSign` bigint DEFAULT NULL,
  `ipSignerSign` varchar(36) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `notifyCreatorSign` tinyint(1) DEFAULT NULL,
  `readOnly` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`idSignature`),
  KEY `fk_sign_user_keys` (`idUserKeys_fk`),
  KEY `fk_sign_digest_algo` (`digestAlgoSign_fk`),
  KEY `fk_sign_signature_algo` (`signatureAlgoSign_fk`),
  CONSTRAINT `fk_sign_digest_algo` FOREIGN KEY (`digestAlgoSign_fk`) REFERENCES `algos` (`idAlgo`) ON DELETE CASCADE,
  CONSTRAINT `fk_sign_signature_algo` FOREIGN KEY (`signatureAlgoSign_fk`) REFERENCES `algos` (`idAlgo`) ON DELETE CASCADE,
  CONSTRAINT `fk_sign_user_keys` FOREIGN KEY (`idUserKeys_fk`) REFERENCES `userkeys` (`idUserKeys`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `signatures`
--

LOCK TABLES `signatures` WRITE;
/*!40000 ALTER TABLE `signatures` DISABLE KEYS */;
INSERT INTO `signatures` VALUES (1,2,'44e2ede0-6d9f-4150-be28-625f769d98b6',NULL,'U6r7VsslAuhUO4G5TFaifr3nn51caDdIUD0LU3UElfA=',21,9,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(2,2,'44e2ede0-6d9f-4150-be28-625f769d98b6',NULL,'2JeEA34CnaGNei+KraRf/EDmMhN+ILGUHZxPrXHYomM=',21,9,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(3,2,'7d8cd08f-f0f5-4ffa-82e6-5bdcae314156',NULL,'XbOaERJcgbF8PF4lvpCdTSVqEv1Gk0ic28LuxlSak+w=',21,9,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(4,2,'7d8cd08f-f0f5-4ffa-82e6-5bdcae314156',NULL,'DDUgZ6JF/tfIg1E2CAThpe69QzHt8yHUL/COutRKcAA=',21,9,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(5,2,'261acc76-8ccd-4af4-8773-95f0d8e388fc',NULL,'XbOaERJcgbF8PF4lvpCdTSVqEv1Gk0ic28LuxlSak+w=',21,9,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(6,2,'261acc76-8ccd-4af4-8773-95f0d8e388fc',NULL,'DDUgZ6JF/tfIg1E2CAThpe69QzHt8yHUL/COutRKcAA=',21,9,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(7,2,'e76a3868-47c5-44b9-88ed-b8e946ccfa25',NULL,'2JeEA34CnaGNei+KraRf/EDmMhN+ILGUHZxPrXHYomM=',21,9,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(8,2,'e76a3868-47c5-44b9-88ed-b8e946ccfa25',NULL,'XbOaERJcgbF8PF4lvpCdTSVqEv1Gk0ic28LuxlSak+w=',21,9,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(9,2,'e76a3868-47c5-44b9-88ed-b8e946ccfa25',NULL,'DDUgZ6JF/tfIg1E2CAThpe69QzHt8yHUL/COutRKcAA=',21,9,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL),(10,2,'e489641a-7850-425c-a7c4-1c0fff6a0dd1',NULL,'XbOaERJcgbF8PF4lvpCdTSVqEv1Gk0ic28LuxlSak+w=',21,9,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL);
/*!40000 ALTER TABLE `signatures` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `signstamps`
--

DROP TABLE IF EXISTS `signstamps`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `signstamps` (
  `idStamp` bigint NOT NULL AUTO_INCREMENT,
  `idSignature_fk` bigint NOT NULL,
  `xSign` float NOT NULL,
  `ySign` float NOT NULL,
  `wSign` float NOT NULL,
  `hSign` float NOT NULL,
  `pageSign` int NOT NULL,
  `pathImg` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`idStamp`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `signstamps`
--

LOCK TABLES `signstamps` WRITE;
/*!40000 ALTER TABLE `signstamps` DISABLE KEYS */;
INSERT INTO `signstamps` VALUES (1,1,368.581,497.98,154.987,98.0067,1,NULL),(2,2,414.021,447.753,154.988,85.8896,1,NULL),(3,3,81.0502,739.214,70.4976,61.9532,1,NULL),(4,4,512.97,711.211,66.2338,56.9498,1,NULL),(5,5,144.908,705.027,219.557,65.953,1,NULL),(6,6,205.506,464.575,225.83,67.8857,1,NULL),(7,7,106.863,260.485,140.898,42.2408,1,NULL),(8,8,302.413,219.652,140.995,42.2408,1,NULL),(9,9,278.503,221.177,132.468,39.7324,1,NULL),(10,10,408.159,765.967,140.995,42.2408,1,NULL);
/*!40000 ALTER TABLE `signstamps` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `states`
--

DROP TABLE IF EXISTS `states`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `states` (
  `idState` int NOT NULL AUTO_INCREMENT,
  `nameState` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `shrinkName` varchar(5) COLLATE utf8mb4_unicode_ci NOT NULL,
  `localUTCTime` int NOT NULL,
  `DSTTime` int NOT NULL,
  `country_fk` int NOT NULL,
  PRIMARY KEY (`idState`),
  KEY `fk_state_country` (`country_fk`),
  CONSTRAINT `fk_state_country` FOREIGN KEY (`country_fk`) REFERENCES `countries` (`idCountry`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
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
  `nameStatus` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL,
  `permissionStatusInst_fk` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL,
  `descriptionStatusInst` text COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`idStatusInst`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `statusinstitution`
--

LOCK TABLES `statusinstitution` WRITE;
/*!40000 ALTER TABLE `statusinstitution` DISABLE KEYS */;
INSERT INTO `statusinstitution` VALUES (1,'Nulo','1','Sin estatus o no requiere'),(2,'PENDIENTE_REGISTRO','1','El usuario inició el registro pero no lo completó (por ejemplo, aún no verificó el correo o no llenó todos los datos).'),(3,'REVISION_REGISTRO','1','El usuario ya completo sus datos y subió sus documentos de identidad, comprobante de domicilio, etc., pero están en revisión.'),(4,'RECHAZO_DOCUMENTOS','0','Algún documento fue rechazado por ser ilegible, inválido o por no coincidir con los datos del usuario. Se debe permitir que el usuario los reenvíe.'),(5,'PLAN_PAGO','1','Falta plan de pago. El contrato fue firmado y registrado correctamente (puede incluir sello de tiempo y almacenamiento del documento).'),(6,'LLAVES','1','El usuario está verificado pero aún no ha subido la llave del usuario root.'),(7,'CONTRATO','1','El usuario ya puede firmar y tiene llaves pero aún no ha firmado el contrato de servicio (por ejemplo plan basico).'),(8,'ACTIVO','1','Usuario completamente habilitado para operar dentro del sistema (por ejemplo, puede firmar documentos, generar sellos de tiempo, etc.).'),(9,'SUSPENDIDO','0','Usuario que fue suspendido temporalmente por problemas administrativos, seguridad o requerimientos legales.'),(10,'BAJA','1','Usuario que fue dado de baja o canceló su cuenta voluntariamente.'),(11,'REVOCADO','0','Usuario que se le retiró el servicio por alguna causa especificada.'),(12,'RECHAZADO','1','Usuario rechazado por motivos de pago');
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
  CONSTRAINT `fk_team_role` FOREIGN KEY (`idRole`) REFERENCES `userroleplatform` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
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
  `idTeam` bigint NOT NULL AUTO_INCREMENT,
  `idInstitution_fk` int NOT NULL,
  `creatorUser_fk` int NOT NULL,
  `nameTeam` varchar(40) COLLATE utf8mb4_unicode_ci NOT NULL,
  `createdAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `lastModified` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deletedAt` datetime DEFAULT NULL,
  `logoUrlTeam` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `limitSigners` int NOT NULL DEFAULT '-1',
  `limitUsers` int NOT NULL DEFAULT '-1',
  `description` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`idTeam`),
  UNIQUE KEY `unique_team_per_client` (`idInstitution_fk`,`nameTeam`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `teams`
--

LOCK TABLES `teams` WRITE;
/*!40000 ALTER TABLE `teams` DISABLE KEYS */;
INSERT INTO `teams` VALUES (1,1,1,'mainteam_1','2025-11-26 14:16:43','2025-11-26 14:16:43',NULL,NULL,-1,-1,'first team'),(2,2,1,'mainteam_2','2025-11-26 14:36:28','2025-11-26 14:36:28',NULL,NULL,-1,-1,'Primer equipo de 2'),(3,3,1,'mainteam_3','2025-11-26 18:04:53','2025-11-26 18:04:53',NULL,NULL,-1,-1,'Primer equipo de 3');
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
  `idUser_fk` bigint NOT NULL,
  `keyFilePath` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `certFilePath` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `serialNumber` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `certVersion` tinyint NOT NULL,
  `issuerRFC4514` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL,
  `notValidAfter` datetime NOT NULL,
  `notValidBefore` datetime NOT NULL,
  `subjectRFC4514` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ocspUrl` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `crlsUrl` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `signature` varchar(1024) COLLATE utf8mb4_unicode_ci NOT NULL,
  `signAlgo` varchar(40) COLLATE utf8mb4_unicode_ci NOT NULL,
  `createdAtKey` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `revokedAtKey` datetime DEFAULT NULL,
  `validKeys` tinyint(1) NOT NULL DEFAULT '0',
  `keyLenKey` int NOT NULL,
  `hashKey` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `hashCer` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `subjectUniqueId` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `subjectSerialNumber` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`idUserKeys`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `userkeys`
--

LOCK TABLES `userkeys` WRITE;
/*!40000 ALTER TABLE `userkeys` DISABLE KEYS */;
INSERT INTO `userkeys` VALUES (1,2,'C:/Users/USER/Desktop/signForce/keys/2/2/funk671228ph6/Claveprivada_FIEL_FUNK671228PH6_20230509_114807.key','C:/Users/USER/Desktop/signForce/keys/2/2/funk671228ph6/funk671228ph6.cer','292233162870206001759766198462772978647764711477',3,'CN=AC UAT,O=SERVICIO DE ADMINISTRACION TRIBUTARIA,OU=SAT-IES Authority,E=oscar.martinez@sat.gob.mx,STREET=3ra cerrada de caliz,postalCode=06370,C=MX,ST=CIUDAD DE MEXICO,L=COYOACAN,OID.2.5.4.45=2.5.4.45,OID.1.2.840.113549.1.9.2=responsable: ACDMA-SAT,','2027-05-08 18:07:00','2023-05-09 18:07:00','CN=KARLA FUENTE NOLASCO,OID.2.5.4.41=KARLA FUENTE NOLASCO,O=KARLA FUENTE NOLASCO,C=MX,E=pruebas@pruebas.gob.mx,OID.2.5.4.45=FUNK671228PH6,SERIALNUMBER=FUNK671228MCLNLR05,','','','U2ymIrVSaT5vdmUNj/dG87uaf2Pwf2ChxdJJ33Kxt8EZ1ZVbCmsqjJQ51xRUo8wOa+ALxpEfxSr7YBxmXPsZHRAsYwEtt11xm5fjkx02Yie/QxqAr9VuLK3WtCOPo1eZDf9KLhyq+zAsHgO1sPknz16TY+l7EMtt/FKXe0TdROuZ8DXAWZy0lxvbIfzUPjV67+a8GDoyQpCSjGMpV8CCTwTeVgS7NpRnLr5eSU4sasouBoWot4FAA5Eky5YR6HY8xIalV4zAbaZx/1XL30tacQ6B42lpQVCb4Vobw4c3B0YxbNjybkXOgmCnXdrmxz7QG90650Bv+cbqJ3yuMvBy4oxQ6EfD+ZW/kiSWbZ/PQM2iWXuCtQY7Ifa9ARGyRcBJ0RFRl9ts7d56knkDavJ8Nc58Drs27leTXGvVWiLpffMTKvNGpFDygTE0hcXlEqDMY4K1F6aPG7RUwgqn7Z8sw648EfpGIsKrRbGzLETBbsZ7+nqzHOTUMG2AhAYJ3IVYDJoNYeEvjM9jW29+MO0xjZ2CnroZR+CCK2t3YUneVHg+H86P+UIEbSHPujaLjOgmDoW4WWhugA/a2QIsWVQEJZndbNQVcOhrBf+C53aqmQmuasFpH8ZnEvcv2C+BI94/MK8b+nFV60DTKOOl3xUT7K5FJH+/8nq9C9R+0rturYg=','sha256-rsa','2025-11-26 14:39:51',NULL,1,2048,'16ptiC3tSNTtF57VDVXsO15dVwwW7UBykXmGjPHbaIw=','Y0d+kTUiQ+L40IHdFrzMJ5F8qbjV8+/Viy6Ep077mYo=','FUNK671228PH6','FUNK671228MCLNLR05'),(2,2,'C:/Users/USER/Desktop/signForce/keys/2/2/iañl750210963/Claveprivada_FIEL_IAÑL750210963_20230509_114916.key','C:/Users/USER/Desktop/signForce/keys/2/2/iañl750210963/iañl750210963.cer','292233162870206001759766198462772978647764711478',3,'CN=AC UAT,O=SERVICIO DE ADMINISTRACION TRIBUTARIA,OU=SAT-IES Authority,E=oscar.martinez@sat.gob.mx,STREET=3ra cerrada de caliz,postalCode=06370,C=MX,ST=CIUDAD DE MEXICO,L=COYOACAN,OID.2.5.4.45=2.5.4.45,OID.1.2.840.113549.1.9.2=responsable: ACDMA-SAT,','2027-05-08 18:07:33','2023-05-09 18:07:33','CN=LUIS IAN ÑUZCO,OID.2.5.4.41=LUIS IAN ÑUZCO,O=LUIS IAN ÑUZCO,C=MX,E=pruebas.sat@pruebas.gob.mx,OID.2.5.4.45=IAÑL750210963,SERIALNUMBER=IANL750210HSRNZS09,','','','n7pA6WzxSIAddZ0IQzF9hEs7WR4eq3twFNqH27g3wZGYndeNjYFNWikJAvU4EP9DC3EXCKexLpxtcJdSDwwcSsogOcNwTl9z8cny6vhs0Fho6hY0fm7hwxqFMePXa4LN46Bc9Vm9CHusXUBS4yK188vSUBe0c0S4c2XPwT8eVk7f9gONJmmCUBkxyg3DCsuHnXTGaVLexZM8J8iTCcp9/le/hENmQY7cuJQYh55fmoUZa4a+22QZpFuRWRsRIgC6eSTdC+SyoeEkA26D6bOTFg5g04hlKLHk2KN5EHutoxIz3HaZly6mHk3nfJRvgzMRbdzCmYUpPDkID9K1csfiSnPwzsGXoLx9KkjKfcH1OtD8Ha33L3ncGEkThkY9gDsQRhKp1bCfhmVQ3rJN/3CsJ66vaHkx6baAwr+79VfsDi22v4XmPx2YvDrFwdO7WmOVbuEr093+wF0fRrF44UiHSp1QFLo4cXiyBBUDRH3giA+N9HdK4SjzMOnUTZlAbRhX4HOE6W55i99lc+RU5Wo9mdrI7eEyHllSFxSWeXRfY30o6Fh6Mg13LGaiiRXp34xI09ME0RKylZS6TzQ5AUkVmlDnBZPj7U+aTD5NZ/iEaeF9IOISd8QsyA7mZHAyUoYYrSzjoGfHd5rkZT7iDPqa8Kdbvk8iu37YCaeGr5qzxQA=','sha256-rsa','2025-11-26 14:40:20',NULL,1,2048,'c75MySlEiuR+EAQAq0Ok2JFETMIG1v/sIM/OBnTEf38=','ckqhmn1fzDxyyxDzCdduvExyjakIldSFmxIxjBeAwLE=','IAÑL750210963','IANL750210HSRNZS09'),(3,3,'C:/Users/USER/Desktop/signForce/keys/3/3/ewe1709045u0/Claveprivada_FIEL_EWE1709045U0_20230518_055331.key','C:/Users/USER/Desktop/signForce/keys/3/3/ewe1709045u0/ewe1709045u0.cer','292233162870206001759766198462772978647764840759',3,'CN=AC UAT,O=SERVICIO DE ADMINISTRACION TRIBUTARIA,OU=SAT-IES Authority,E=oscar.martinez@sat.gob.mx,STREET=3ra cerrada de caliz,postalCode=06370,C=MX,ST=CIUDAD DE MEXICO,L=COYOACAN,OID.2.5.4.45=2.5.4.45,OID.1.2.840.113549.1.9.2=responsable: ACDMA-SAT,','2027-05-17 12:10:57','2023-05-18 12:10:57','CN=ESCUELA WILSON ESQUIVEL S DE CV,OID.2.5.4.41=ESCUELA WILSON ESQUIVEL S DE CV,O=ESCUELA WILSON ESQUIVEL S DE CV,C=MX,E=SATpruebas@pruebas.gob.mx,OID.2.5.4.45=EWE1709045U0 / VADA800927DJ3,SERIALNUMBER= / VADA800927HSRSRL05,','','','mZaUbguH9e0KhGn6b+4nRMJTEzH/x1L+2oZB3DqVdtYH3JqY2b+Rmg64VCf5eZTiEW2gaslW3iN5OipRjjx+6WzU5P7NxzxIButINziBVUdeClYC3AdUwWnf+hNs9xM5iK+WtgODI+1Apwz0zS6Nu/osApooTad8vVZZONRoIJ6TdeV9b6oNc1qz5RPdXM+pp3JnkPrAFEGVQQQ6fyNFS6lMNSlWIkYAyA2s6rKvfy4OmLSokyXwf8kop1IQN6g0ogHZu4Ivsfsc3RO6QEdO6Oa9hB3z6rIGowidxExF7WYRT/gNMynGrCWy97Ug2Pr4QRzGTui2uaJhRw0gN4PsXJrFyqhPDBhXRmk1o4yUb4Id94M0Tp7WapVcBxtazniChUDMaH5ZWd23jgqHHZ9GZPh/5CNqUK+P5D+3h0/zfUotf6FH80eNwvAmGWfUXPKomAnqmGVC5RlXUeVuF6NysVUuK0+dhf1dzlE1sKAhO1Tb8NlMcT9FdSJRlAyG4IwYyDuXcEiq0LNHHprs3exgQ4jAMUCd/hieH2qi1tOe9kYxvNvO/5u2yGJ0Xu8L/oz/cm1grYa2lumZeb1Ik/x+fRT0K1qO9NoXeVJvPG/mfjCfksd5baaJ4uyLInnCFkZDJJoD3I7aUVTdd7tdcsbqbbV8NFBKSmOcDHw67Qb2HBM=','sha256-rsa','2025-11-26 18:25:59',NULL,1,2048,'IbMefuazIcEjAr4uarTd/Ogm4bWiDt71bWYu+cn8Re4=','5sNuNkRuU0Aeu2jYfS8fag64HiAHL5gKywkEiZZi4ow=','EWE1709045U0 / VADA800927DJ3','VADA800927HSRSRL05');
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
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
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
  `name` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `userroles`
--

LOCK TABLES `userroles` WRITE;
/*!40000 ALTER TABLE `userroles` DISABLE KEYS */;
INSERT INTO `userroles` VALUES (1,1,1,'2025-11-26 13:56:15'),(2,1,2,'2025-11-26 14:36:28'),(3,1,3,'2025-11-26 18:04:53');
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
  `taxNumUser` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pobUidUser` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `nameUser` varchar(60) COLLATE utf8mb4_unicode_ci NOT NULL,
  `lastNameUser` varchar(60) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `aliasUser` varchar(80) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `emailUser` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `phoneUser` varchar(15) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
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
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,NULL,NULL,'Juan','Nuñez','root_a4ayc/80/OGda4BO/1o/V0etpOqiLx1JwB5S3beHW0s=','do.c.001@hotmail.com','5529996911',52,1,1,NULL,_binary 'cmTbLk5FpBxZIvPcp09xJw8OPgo9F+9ymE2K7FaVKLQ=',NULL,0,1,1,'2025-11-26 13:56:15','2025-11-26 13:56:15',NULL),(2,NULL,NULL,'Diego','Reyes','root_1HNeOiZeFu7gP1lxi5tdAwGcB9i2xR+Q2jpmbuwTqzU=','prof.jjnm.ockham@gmail.com','5511223344',52,1,1,2,_binary 'U8JlQCLCmLz0Hi8BzT6SDv47lk7E5+pkAQrPO5nWFIo=',NULL,0,2,2,'2025-11-26 14:36:28','2025-11-26 14:36:28',NULL),(3,NULL,NULL,'Aida','Valencia','root_TgdAhWK+24tgzgXB3s/jrRa3IjCWfeAfZAt+Rym0n84=','contacto.somostec@gmail.com','5510317898',52,1,1,3,_binary 'aDETUIJftgmtpF+RIyDxm5U850Adw6VNdPzgdHn00u8=',NULL,0,3,3,'2025-11-26 18:04:53','2025-11-26 18:04:53',NULL);
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
  `prevStatus` enum('PENDING','ACTIVE','SUSPENDED','REVOKED') COLLATE utf8mb4_unicode_ci NOT NULL,
  `newStatus` enum('PENDING','ACTIVE','SUSPENDED','REVOKED') COLLATE utf8mb4_unicode_ci NOT NULL,
  `changedBy` bigint NOT NULL,
  `changedAt` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `reason` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_userstatus_user` (`idUser`),
  KEY `fk_userstatus_changedby` (`changedBy`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
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
  `useAlgo` varchar(15) COLLATE utf8mb4_unicode_ci NOT NULL,
  `descriptionAlgo` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`idusetype`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
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

-- Dump completed on 2025-11-26 18:42:03
