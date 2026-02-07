CREATE TABLE IF NOT EXISTS `schema_markup` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `pagereport_id` int unsigned NOT NULL,
  `crawl_id` int unsigned NOT NULL,
  `schema_type` varchar(256) NOT NULL,
  `json_ld` text DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `schema_markup_pagereport` (`pagereport_id`),
  KEY `schema_markup_crawl` (`crawl_id`),
  CONSTRAINT `schema_markup_crawl` FOREIGN KEY (`crawl_id`) REFERENCES `crawls` (`id`) ON DELETE CASCADE,
  CONSTRAINT `schema_markup_pagereport` FOREIGN KEY (`pagereport_id`) REFERENCES `pagereports` (`id`) ON DELETE CASCADE
);
