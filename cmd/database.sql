SET NAMES utf8;
SET time_zone = '+00:00';
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP TABLE IF EXISTS `posts`;
CREATE TABLE `posts` (
    `id` int NOT NULL AUTO_INCREMENT,
    `title` varchar(255) NOT NULL,
    `author` varchar(255) DEFAULT NULL,
    `text` text NOT NULL,
    `updated` varchar(255) DEFAULT NULL,   
    PRIMARY KEY (`id`)   
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `posts` (`id`, `title`, `author`, `text`, `updated`) VALUES (0, "MY CREATION", NULL, "Lets talk about...", NULL);


DROP TABLE IF EXISTS `news`;
CREATE TABLE `news` (
    `id` int NOT NULL AUTO_INCREMENT,
    `title` varchar(255) NOT NULL,
    `date` varchar(255) NOT NULL,
    `author` text NOT NULL,
    `text` varchar(255) DEFAULT NULL,   
    PRIMARY KEY (`id`)   
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `news` (`id`, `title`, `date`, `author`, `text`) VALUES (0, "MY CREATION", "21.10.2024", "admin", "hello fucking world");