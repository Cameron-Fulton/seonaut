CREATE TABLE IF NOT EXISTS `custom_extractors` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `project_id` int unsigned NOT NULL,
  `name` varchar(256) NOT NULL,
  `type` varchar(20) NOT NULL COMMENT 'regex, css, or xpath',
  `pattern` text NOT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `custom_extractors_project` (`project_id`),
  CONSTRAINT `custom_extractors_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `extraction_results` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `pagereport_id` int unsigned NOT NULL,
  `crawl_id` int unsigned NOT NULL,
  `extractor_id` int unsigned NOT NULL,
  `value` text DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `extraction_results_pagereport` (`pagereport_id`),
  KEY `extraction_results_crawl` (`crawl_id`),
  KEY `extraction_results_extractor` (`extractor_id`),
  CONSTRAINT `extraction_results_crawl` FOREIGN KEY (`crawl_id`) REFERENCES `crawls` (`id`) ON DELETE CASCADE,
  CONSTRAINT `extraction_results_pagereport` FOREIGN KEY (`pagereport_id`) REFERENCES `pagereports` (`id`) ON DELETE CASCADE,
  CONSTRAINT `extraction_results_extractor` FOREIGN KEY (`extractor_id`) REFERENCES `custom_extractors` (`id`) ON DELETE CASCADE
);
