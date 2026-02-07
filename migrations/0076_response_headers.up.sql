CREATE TABLE IF NOT EXISTS `response_headers` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `pagereport_id` int unsigned NOT NULL,
  `crawl_id` int unsigned NOT NULL,
  `name` varchar(256) NOT NULL,
  `value` text NOT NULL,
  PRIMARY KEY (`id`),
  KEY `response_headers_pagereport` (`pagereport_id`),
  KEY `response_headers_crawl` (`crawl_id`),
  CONSTRAINT `response_headers_crawl` FOREIGN KEY (`crawl_id`) REFERENCES `crawls` (`id`) ON DELETE CASCADE,
  CONSTRAINT `response_headers_pagereport` FOREIGN KEY (`pagereport_id`) REFERENCES `pagereports` (`id`) ON DELETE CASCADE
);
