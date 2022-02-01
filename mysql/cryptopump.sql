CREATE DATABASE  IF NOT EXISTS `cryptopump` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `cryptopump`;
-- MySQL dump 10.13  Distrib 8.0.28, for macos11 (x86_64)
--
-- Host: 127.0.0.1    Database: cryptopump
-- ------------------------------------------------------
-- Server version	8.0.27

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
-- Table structure for table `global`
--

DROP TABLE IF EXISTS `global`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `global` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `Profit` float NOT NULL,
  `ProfitNet` float NOT NULL,
  `ProfitPct` float NOT NULL,
  `TransactTime` varchar(45) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci MAX_ROWS=1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `orders`
--

DROP TABLE IF EXISTS `orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `orders` (
  `ClientOrderId` varchar(45) NOT NULL,
  `CummulativeQuoteQty` float NOT NULL,
  `ExecutedQuantity` float NOT NULL,
  `OrderID` bigint NOT NULL,
  `OrderIDSource` bigint NOT NULL,
  `Price` float NOT NULL,
  `Side` varchar(45) NOT NULL,
  `Status` varchar(45) NOT NULL,
  `Symbol` varchar(45) NOT NULL,
  `TransactTime` bigint NOT NULL,
  `ThreadID` varchar(45) NOT NULL,
  `ThreadIDSession` varchar(45) NOT NULL,
  PRIMARY KEY (`OrderID`),
  UNIQUE KEY `OrderID_UNIQUE` (`OrderID`),
  KEY `orders_idx_side_status` (`Side`,`Status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `session`
--

DROP TABLE IF EXISTS `session`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `session` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `ThreadID` varchar(45) NOT NULL,
  `ThreadIDSession` varchar(45) NOT NULL,
  `Exchange` varchar(45) NOT NULL,
  `FiatSymbol` varchar(45) NOT NULL,
  `FiatFunds` float NOT NULL,
  `DiffTotal` float NOT NULL,
  `Status` tinyint(1) NOT NULL,
  PRIMARY KEY (`ID`),
  UNIQUE KEY `ThreadID_UNIQUE` (`ThreadID`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `thread`
--

DROP TABLE IF EXISTS `thread`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `thread` (
  `ID` int NOT NULL AUTO_INCREMENT,
  `ThreadID` varchar(45) NOT NULL,
  `ThreadIDSession` varchar(45) NOT NULL,
  `OrderID` bigint DEFAULT NULL,
  `CummulativeQuoteQty` float NOT NULL,
  `Price` float NOT NULL,
  `ExecutedQuantity` float NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping routines for database 'cryptopump'
--
/*!50003 DROP PROCEDURE IF EXISTS `DeleteSession` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `DeleteSession`(IN in_ThreadID varchar(45))
BEGIN
	DECLARE ThreadID varchar(45);
	SET SQL_SAFE_UPDATES = 0;
	SET ThreadID = in_ThreadID;
	DELETE FROM session ft
	WHERE ft.ThreadID = in_ThreadID;
	SET SQL_SAFE_UPDATES = 1;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `DeleteThreadTransactionAll` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `DeleteThreadTransactionAll`()
BEGIN
    SET SQL_SAFE_UPDATES = 0;
    DELETE FROM thread ft;
    SET SQL_SAFE_UPDATES = 1;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `DeleteThreadTransactionByOrderID` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `DeleteThreadTransactionByOrderID`(IN in_param_OrderID bigint)
BEGIN
	DECLARE declared_in_param_OrderID bigint;
    SET SQL_SAFE_UPDATES = 0;
    SET declared_in_param_OrderID = in_param_OrderID;
    DELETE FROM thread ft
    WHERE ft.OrderID = in_param_OrderID;
    SET SQL_SAFE_UPDATES = 1;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetGlobal` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetGlobal`()
BEGIN
SELECT 
    `global`.`Profit` AS `Profit`,
    `global`.`ProfitNet` AS `ProfitNet`,
    `global`.`ProfitPct` AS `ProfitPct`,
    `global`.`TransactTime` AS `TransactTime`
FROM
    `global`
WHERE
    `global`.`ID` = 1
LIMIT 1;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetLastOrderTransactionPrice` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetLastOrderTransactionPrice`(IN in_param_ThreadID varchar(45), IN in_param_Side varchar(45))
BEGIN
	DECLARE declared_in_param_ThreadID CHAR(45);
    DECLARE declared_in_param_Side CHAR(45);
    SET declared_in_param_ThreadID = in_param_ThreadID;
    SET declared_in_param_Side = in_param_Side;
    SELECT `orders`.`Price` AS `Price`
	FROM `orders`
	WHERE (`orders`.`ThreadID` = declared_in_param_ThreadID
		AND `orders`.`Side` = declared_in_param_Side AND (`orders`.`Status` <> 'CANCELED'
		OR `orders`.`Status` IS NULL))
	ORDER BY from_unixtime((`orders`.`TransactTime` / 1000)) DESC
	LIMIT 1;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetLastOrderTransactionSide` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetLastOrderTransactionSide`(IN in_param_ThreadID varchar(45))
BEGIN
	DECLARE declared_in_param_ThreadID CHAR(45);
    SET declared_in_param_ThreadID = in_param_ThreadID;
	SELECT `orders`.`Side` AS `Side`
	FROM `orders`
	WHERE (`orders`.`ThreadID` = declared_in_param_ThreadID
	   AND `orders`.`Status` = 'FILLED')
	ORDER BY from_unixtime((`orders`.`TransactTime` / 1000)) DESC
	LIMIT 1;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetOrderByOrderID` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetOrderByOrderID`(IN in_param_OrderID bigint, IN in_param_ThreadID varchar(45))
BEGIN
DECLARE declared_in_param_orderid BIGINT; 
DECLARE declared_in_param_threadid CHAR(50); 
SET declared_in_param_orderid = in_param_orderid; 
SET declared_in_param_threadid = in_param_threadid;
SELECT 
    `orders`.`orderid` AS `OrderID`,
    `orders`.`price` AS `Price`,
    `orders`.`executedquantity` AS `ExecutedQuantity`,
    `orders`.`cummulativequoteqty` AS `CummulativeQuoteQty`,
    `orders`.`transacttime` AS `TransactTime`
FROM
    `orders`
WHERE
    (`orders`.`orderid` = declared_in_param_orderid
        AND `orders`.`threadid` = declared_in_param_threadid)
LIMIT 1; 
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetOrderSymbol` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetOrderSymbol`(IN in_param varchar(45))
BEGIN
	DECLARE declared_in_param CHAR(45);
    SET declared_in_param = in_param;
	SELECT Symbol from orders ft 
    WHERE ft.ThreadID = declared_in_param
    ORDER BY TransactTime DESC LIMIT 1;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetOrderTransactionCount` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetOrderTransactionCount`(IN in_param_ThreadID varchar(45), IN in_param_Side varchar(45), IN in_param_Minutes int)
BEGIN
	DECLARE declared_in_param_ThreadID CHAR(45);
    DECLARE declared_in_param_Side CHAR(45);
    DECLARE declared_in_param_Minutes int;
    SET declared_in_param_ThreadID = in_param_ThreadID;
    SET declared_in_param_Side = in_param_Side;
    SET declared_in_param_Minutes = in_param_Minutes;
SELECT COALESCE(count(*),0) AS `count`
FROM `orders`
WHERE (`orders`.`Side` = declared_in_param_Side
   AND `orders`.`Status` = 'FILLED' AND str_to_date(date_format(CAST(from_unixtime((`orders`.`TransactTime` / 1000)) AS DATETIME), '%Y-%m-%d %H:%i'), '%Y-%m-%d %H:%i') BETWEEN str_to_date(date_format(CAST(date_add(now(6), INTERVAL declared_in_param_Minutes minute) AS DATETIME), '%Y-%m-%d %H:%i'), '%Y-%m-%d %H:%i') AND str_to_date(date_format(CAST(now(6) AS DATETIME), '%Y-%m-%d %H:%i'), '%Y-%m-%d %H:%i') AND `orders`.`ThreadID` = declared_in_param_ThreadID);
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetOrderTransactionPending` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetOrderTransactionPending`(IN in_param_ThreadID varchar(45))
BEGIN
	DECLARE declared_in_param_ThreadID CHAR(45);
    SET declared_in_param_ThreadID = in_param_ThreadID;
SELECT `orders`.`OrderID` AS `OrderID`, `orders`.`Symbol` AS `Symbol`
FROM `orders`
WHERE (`orders`.`ThreadID` = declared_in_param_ThreadID
   AND (`orders`.`Status` <> 'FILLED'
    OR `orders`.`Status` IS NULL) AND (`orders`.`Status` <> 'CANCELED' OR `orders`.`Status` IS NULL) AND `orders`.`Status` IS NOT NULL AND (`orders`.`Status` <> '' OR `orders`.`Status` IS NULL))
ORDER BY from_unixtime((`orders`.`TransactTime` / 1000)) ASC
LIMIT 1;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetOrderTransactionSideLastTwo` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetOrderTransactionSideLastTwo`(IN in_param_ThreadID varchar(45))
BEGIN
	DECLARE declared_in_param_ThreadID CHAR(45);
    SET declared_in_param_ThreadID = in_param_ThreadID;
	SELECT `A`.`Side` AS `Last`, `B`.`Side` AS `SecondLast` FROM (
	(SELECT `orders`.`Side` AS `Side`
	FROM `orders`
	WHERE (`orders`.`ThreadID` = declared_in_param_ThreadID
	   AND (`orders`.`Status` <> 'CANCELED'
		OR `orders`.`Status` IS NULL))
	ORDER BY from_unixtime((`orders`.`TransactTime` / 1000)) DESC
	LIMIT 1) A
	INNER JOIN
	(SELECT `orders`.`Side` AS `Side`
	FROM `orders`
	WHERE (`orders`.`ThreadID` = declared_in_param_ThreadID
	   AND (`orders`.`Status` <> 'CANCELED'
		OR `orders`.`Status` IS NULL))
	ORDER BY from_unixtime((`orders`.`TransactTime` / 1000)) DESC
	LIMIT 1,1) B
	);
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetOrderTransactionTimeByOrderID` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetOrderTransactionTimeByOrderID`(IN in_param_OrderID bigint)
BEGIN
	DECLARE declared_in_param_OrderID bigint;
    SET declared_in_param_OrderID = in_param_OrderID;
	SELECT `orders`.`TransactTime` AS `TransactTime`
	FROM `orders`
	WHERE `orders`.`OrderID` = declared_in_param_OrderID
	LIMIT 1;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetProfit` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetProfit`()
BEGIN
SELECT
        SUM(`source`.`Profit`) AS `profit`,
        SUM(`source`.`Profit`) + (`source`.`Diff`) AS `netprofit`,
        AVG(`source`.`Percentage`) AS `avg` 
    FROM
        (SELECT
            `orders`.`Side` AS `Side`,
            `Orders`.`Side` AS `Orders__Side`,
            `orders`.`Status` AS `Status`,
            `Orders`.`Status` AS `Orders__Status`,
            `orders`.`ThreadID` AS `ThreadID`,
            `Orders`.`CummulativeQuoteQty` AS `Orders__CummulativeQuoteQty`,
            `orders`.`CummulativeQuoteQty` AS `CummulativeQuoteQty`,
            (`Orders`.`CummulativeQuoteQty` - `orders`.`CummulativeQuoteQty`) AS `Profit`,
            ((`Orders`.`CummulativeQuoteQty` - `orders`.`CummulativeQuoteQty`) / CASE 
                WHEN `Orders`.`CummulativeQuoteQty` = 0 THEN NULL 
                ELSE `Orders`.`CummulativeQuoteQty` END) AS `Percentage`,
(SELECT
    sum(`session`.`DiffTotal`) AS `sum` 
FROM
    `session`) AS `Diff` 
FROM
`orders` 
INNER JOIN
`orders` `Orders` 
    ON `orders`.`OrderID` = `Orders`.`OrderIDSource` 
WHERE
(
    `orders`.`Side` = 'BUY'
) 
AND (
    `orders`.`Status` = 'FILLED'
)
) `source` 
WHERE
(
1 = 1 
AND `source`.`Orders__Side` = 'SELL' 
AND 1 = 1 
AND `source`.`Orders__Status` = 'FILLED'
);
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetProfitByThreadID` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetProfitByThreadID`(IN in_param_ThreadID varchar(45))
BEGIN
DECLARE declared_in_param_ThreadID CHAR(50);
    SET declared_in_param_ThreadID = in_param_ThreadID;
SELECT 
    SUM(`source`.`Profit`) + (`source`.`Diff`) AS `sum`,
    AVG(`source`.`Percentage`) AS `avg`
FROM
    (SELECT 
        `orders`.`Side` AS `Side`,
            `Orders`.`Side` AS `Orders__Side`,
            `orders`.`Status` AS `Status`,
            `Orders`.`Status` AS `Orders__Status`,
            `orders`.`ThreadID` AS `ThreadID`,
            `Orders`.`CummulativeQuoteQty` AS `Orders__CummulativeQuoteQty`,
            `orders`.`CummulativeQuoteQty` AS `CummulativeQuoteQty`,
            (`Orders`.`CummulativeQuoteQty` - `orders`.`CummulativeQuoteQty`) AS `Profit`,
            ((`Orders`.`CummulativeQuoteQty` - `orders`.`CummulativeQuoteQty`) / CASE
                WHEN `Orders`.`CummulativeQuoteQty` = 0 THEN NULL
                ELSE `Orders`.`CummulativeQuoteQty`
            END) AS `Percentage`,
            (SELECT 
                    SUM(`session`.`DiffTotal`) AS `sum`
                FROM
                    `session`
                WHERE
                    `session`.`ThreadID` = declared_in_param_ThreadID) AS `Diff`
    FROM
        `orders`
    INNER JOIN `orders` `Orders` ON `orders`.`OrderID` = `Orders`.`OrderIDSource`) `source`
WHERE
    (`source`.`Side` = 'BUY'
        AND `source`.`Orders__Side` = 'SELL'
        AND `source`.`Status` = 'FILLED'
        AND `source`.`Orders__Status` = 'FILLED'
        AND `source`.`ThreadID` = declared_in_param_ThreadID);
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetSessionStatus` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetSessionStatus`()
BEGIN
SELECT `session`.`ThreadID` AS `ThreadID`, `session`.`Status` AS `Status`
FROM cryptopump.session
WHERE `session`.`Status` = 1;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetThreadCount` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetThreadCount`()
BEGIN
SELECT 
    COUNT(DISTINCT `session`.`ThreadID`) AS `count`
FROM
    `cryptopump`.`session`;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetThreadLastTransaction` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetThreadLastTransaction`(IN in_param_ThreadID varchar(45))
BEGIN
	DECLARE declared_in_param_ThreadID CHAR(50);
    SET declared_in_param_ThreadID = in_param_ThreadID;
	SELECT `thread`.`CummulativeQuoteQty` AS `CummulativeQuoteQty`, `thread`.`OrderID` AS `OrderID`, `thread`.`Price` AS `Price`, `thread`.`ExecutedQuantity` AS `ExecutedQuantity`, `Orders`.`TransactTime` AS `TransactTime`
	FROM `thread`
	LEFT JOIN `orders` `Orders` ON `thread`.`OrderID` = `Orders`.`OrderID`
	WHERE (`thread`.`ThreadID` = declared_in_param_ThreadID)
	ORDER BY `thread`.`Price` ASC
	LIMIT 1;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetThreadTransactionAmount` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetThreadTransactionAmount`()
BEGIN
SELECT 
    SUM(`thread`.`CummulativeQuoteQty`) AS `sum`
FROM
    `cryptopump`.`thread`;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetThreadTransactionByPrice` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetThreadTransactionByPrice`(IN in_param_ThreadID varchar(45), IN in_param_Price float)
BEGIN
	DECLARE declared_in_param_ThreadID CHAR(50);
	DECLARE declared_in_param_Price FLOAT;
    SET declared_in_param_ThreadID = in_param_ThreadID;
    SET declared_in_param_Price = in_param_Price;
	SELECT `thread`.`CummulativeQuoteQty` AS `CummulativeQuoteQty`, `thread`.`OrderID` AS `OrderID`, `thread`.`Price` AS `Price`, `thread`.`ExecutedQuantity` AS `ExecutedQuantity`, `Orders`.`TransactTime` AS `TransactTime`
	FROM `thread`
	LEFT JOIN `orders` `Orders` ON `thread`.`OrderID` = `Orders`.`OrderID`
	WHERE (`thread`.`ThreadID` = declared_in_param_ThreadID
	   AND `thread`.`Price` < declared_in_param_Price)
	ORDER BY `thread`.`Price` ASC
	LIMIT 1;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetThreadTransactionByPriceHigher` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetThreadTransactionByPriceHigher`(IN in_param_ThreadID varchar(45), IN in_param_Price float)
BEGIN
	DECLARE declared_in_param_ThreadID CHAR(50);
	DECLARE declared_in_param_Price FLOAT;
    SET declared_in_param_ThreadID = in_param_ThreadID;
    SET declared_in_param_Price = in_param_Price;
	SELECT 
    `thread`.`CummulativeQuoteQty` AS `CummulativeQuoteQty`,
    `thread`.`OrderID` AS `OrderID`,
    `thread`.`Price` AS `Price`,
    `thread`.`ExecutedQuantity` AS `ExecutedQuantity`,
    `Orders`.`TransactTime` AS `TransactTime`
FROM
    `thread`
        LEFT JOIN
    `orders` `Orders` ON `thread`.`OrderID` = `Orders`.`OrderID`
WHERE
    (`thread`.`ThreadID` = declared_in_param_ThreadID
        AND `thread`.`Price` > declared_in_param_Price)
ORDER BY `thread`.`Price` DESC
LIMIT 1;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetThreadTransactionByThreadID` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetThreadTransactionByThreadID`(IN in_param_ThreadID varchar(45))
BEGIN
	DECLARE declared_in_param_ThreadID CHAR(50);
    SET declared_in_param_ThreadID = in_param_ThreadID;
SELECT 
    `thread`.`OrderID` AS `OrderID`,
    `thread`.`CummulativeQuoteQty` AS `CummulativeQuoteQty`,
    `thread`.`Price` AS `Price`,
    `thread`.`ExecutedQuantity` AS `ExecutedQuantity`
FROM
    `thread`
        LEFT JOIN
    `orders` `Orders` ON `thread`.`OrderID` = `Orders`.`OrderID`
WHERE
    `thread`.`ThreadID` = declared_in_param_ThreadID
ORDER BY `thread`.`Price` ASC;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetThreadTransactionCount` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetThreadTransactionCount`(IN in_param varchar(45))
BEGIN
	DECLARE declared_in_param CHAR(50);
    SET declared_in_param = in_param;
    SELECT count(*) AS count FROM thread ft
    WHERE ft.ThreadID = declared_in_param;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetThreadTransactionDistinct` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetThreadTransactionDistinct`()
BEGIN
	SELECT DISTINCT ThreadID, ThreadIDSession FROM thread;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetThreadTransactiontUpmarketPriceCount` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `GetThreadTransactiontUpmarketPriceCount`(IN in_param_ThreadID varchar(45), IN in_param_Price float)
BEGIN
	DECLARE declared_in_param_ThreadID CHAR(45);
	DECLARE declared_in_param_Price float;
    SET declared_in_param_ThreadID = in_param_ThreadID;
    SET declared_in_param_Price = in_param_Price;
	SELECT count(*) AS `count`
	FROM `thread`
	WHERE (`thread`.`Price` < declared_in_param_Price
	   AND `thread`.`ThreadID` = declared_in_param_ThreadID);
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `SaveGlobal` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `SaveGlobal`(in_Profit float, in_ProfitNet float, in_ProfitPct float, in_TransactTime bigint)
BEGIN
INSERT INTO global (Profit, ProfitNet, ProfitPct, TransactTime)
VALUES (in_Profit, in_ProfitNet, in_ProfitPct, in_TransactTime);
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `SaveOrder` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `SaveOrder`(ClientOrderId varchar(45), CummulativeQuoteQty float, ExecutedQuantity float, OrderID bigint, OrderIDSource bigint, Price float, Side varchar(45), Status varchar(45), Symbol varchar(45), TransactTime bigint, ThreadID varchar(45), ThreadIDSession varchar(45))
BEGIN
INSERT INTO orders (ClientOrderId, CummulativeQuoteQty, ExecutedQuantity, OrderID, OrderIDSource, Price, Side, Status, Symbol, TransactTime, ThreadID, ThreadIDSession)
VALUES (ClientOrderId, CummulativeQuoteQty, ExecutedQuantity, OrderID, OrderIDSource, Price, Side, Status, Symbol, TransactTime, ThreadID, ThreadIDSession);
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `SaveSession` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `SaveSession`(in_ThreadID varchar(45), in_ThreadIDSession varchar(45), in_Exchange varchar(45), in_FiatSymbol varchar(45), in_FiatFunds float, in_DiffTotal float, in_Status tinyint(1))
BEGIN
INSERT INTO session (ThreadID, ThreadIDSession, Exchange, FiatSymbol, FiatFunds, DiffTotal, Status)
VALUES (in_ThreadID, in_ThreadIDSession, in_Exchange, in_FiatSymbol, in_FiatFunds, in_DiffTotal, in_Status);
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `SaveThreadTransaction` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `SaveThreadTransaction`(ThreadID varchar(45), ThreadIDSession varchar(45), OrderID bigint, CummulativeQuoteQty float, Price float, ExecutedQuantity float)
BEGIN
INSERT INTO thread (ThreadID, ThreadIDSession, OrderID, CummulativeQuoteQty, Price, ExecutedQuantity)
VALUES (ThreadID, ThreadIDSession, OrderID, CummulativeQuoteQty, Price, ExecutedQuantity);
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `UpdateGlobal` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `UpdateGlobal`(in_Profit float, in_ProfitNet float, in_ProfitPct float, in_TransactTime bigint)
BEGIN
SET SQL_SAFE_UPDATES = 0;
UPDATE global 
SET 
    Profit = in_Profit,
    ProfitNet = in_ProfitNet,
    ProfitPct = in_ProfitPct,
    TransactTime = in_TransactTime
WHERE
    ID = 1;
SET SQL_SAFE_UPDATES = 1;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `UpdateOrder` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `UpdateOrder`(in_OrderID bigint, CummulativeQuoteQty float, ExecutedQuantity float, Price float, Status varchar(45))
BEGIN
SET SQL_SAFE_UPDATES = 0;
UPDATE orders
SET  CummulativeQuoteQty = CummulativeQuoteQty,
	ExecutedQuantity = ExecutedQuantity,
    Price = Price,
    Status = Status
WHERE OrderID = in_OrderID;
SET SQL_SAFE_UPDATES = 1;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `UpdateSession` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`root`@`%` PROCEDURE `UpdateSession`(in_ThreadID varchar(45), in_ThreadIDSession varchar(45), in_Exchange varchar(45), in_FiatSymbol varchar(45), in_FiatFunds float, in_DiffTotal float, in_Status tinyint(1))
BEGIN
SET SQL_SAFE_UPDATES = 0;
	UPDATE `session` 
SET 
    `session`.`FiatFunds` = in_FiatFunds,
    `session`.`DiffTotal` = in_DiffTotal,
    `session`.`Status` = in_Status
WHERE
    `session`.`ThreadID` = in_ThreadID;
SET SQL_SAFE_UPDATES = 1;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-02-01 18:48:08
