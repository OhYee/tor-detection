CREATE TABLE `domain` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `domain` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `domain_domain_IDX` (`domain`,`time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=194202 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

CREATE TABLE `http` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `host` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `request` text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `http_host_IDX` (`host`,`time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=215035 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci