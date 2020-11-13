CREATE TABLE `tor`.`domain` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `domain` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `domain_domain_IDX` (`domain`,`time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=194202 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

CREATE TABLE `tor`.`http` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `host` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `request` text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `http_host_IDX` (`host`,`time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=215035 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

CREATE TABLE `tor`.`ip` (
  `ip` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `count` bigint(20) NOT NULL DEFAULT 0,
  `domain` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tor` tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`ip`),
  KEY `ip_count_IDX` (`count`,`tor`) USING BTREE,
  KEY `ip_domain_IDX` (`domain`,`ip`) USING BTREE,
  KEY `ip_ip_IDX` (`ip`) USING BTREE,
  KEY `count_IDX` (`count`) USING BTREE,
  KEY `domain_IDX` (`domain`) USING BTREE,
  KEY `ip_tor_IDX` (`tor`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci