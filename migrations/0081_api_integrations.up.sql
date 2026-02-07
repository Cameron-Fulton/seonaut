CREATE TABLE IF NOT EXISTS `api_credentials` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int unsigned NOT NULL,
  `service` varchar(50) NOT NULL COMMENT 'pagespeed, searchconsole, ga4, ahrefs',
  `credentials` text NOT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `api_credentials_user_service` (`user_id`, `service`),
  CONSTRAINT `api_credentials_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `pagespeed_results` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `pagereport_id` int unsigned NOT NULL,
  `crawl_id` int unsigned NOT NULL,
  `performance_score` float DEFAULT NULL,
  `fcp` float DEFAULT NULL,
  `lcp` float DEFAULT NULL,
  `tbt` float DEFAULT NULL,
  `cls` float DEFAULT NULL,
  `speed_index` float DEFAULT NULL,
  `tti` float DEFAULT NULL,
  `recommendations` text DEFAULT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `pagespeed_results_pagereport` (`pagereport_id`),
  KEY `pagespeed_results_crawl` (`crawl_id`),
  CONSTRAINT `pagespeed_results_crawl` FOREIGN KEY (`crawl_id`) REFERENCES `crawls` (`id`) ON DELETE CASCADE,
  CONSTRAINT `pagespeed_results_pagereport` FOREIGN KEY (`pagereport_id`) REFERENCES `pagereports` (`id`) ON DELETE CASCADE
);
