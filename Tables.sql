CREATE TABLE `Users` (
 `id` int unsigned primary key NOT NULL AUTO_INCREMENT,
  `Telegram_Id` bigint UNSIGNED UNIQUE NOT NULL,
  `First_Name` varchar(150) DEFAULT NULL,
  `Last_Name` varchar(150) DEFAULT NULL,
  `Chat_Id` int NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL
  );
  
CREATE TABLE `Tasks` (
`id` int unsigned primary key NOT NULL AUTO_INCREMENT,
 `Title` varchar(150) NOT NULL,
 `Description` varchar(300) DEFAULT NULL,
 `End_Date` timestamp NULL DEFAULT NULL,
 `User_Id` bigint unsigned NOT NULL,
 `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
 `updated_at` datetime DEFAULT NULL,
 `deleted_at` datetime DEFAULT NULL,
 FOREIGN KEY (`User_Id`) 
 REFERENCES `Users` (`Telegram_Id`) 
 ON DELETE CASCADE 
 ON UPDATE CASCADE
 );