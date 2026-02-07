CREATE TABLE IF NOT EXISTS `ai_analyses` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `pagereport_id` int unsigned NOT NULL,
  `crawl_id` int unsigned NOT NULL,
  `provider` varchar(50) NOT NULL COMMENT 'claude or lmstudio',
  `analysis_type` varchar(100) NOT NULL,
  `result` text DEFAULT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `ai_analyses_pagereport` (`pagereport_id`),
  KEY `ai_analyses_crawl` (`crawl_id`),
  CONSTRAINT `ai_analyses_crawl` FOREIGN KEY (`crawl_id`) REFERENCES `crawls` (`id`) ON DELETE CASCADE,
  CONSTRAINT `ai_analyses_pagereport` FOREIGN KEY (`pagereport_id`) REFERENCES `pagereports` (`id`) ON DELETE CASCADE
);
