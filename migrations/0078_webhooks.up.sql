CREATE TABLE IF NOT EXISTS `webhooks` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `project_id` int unsigned NOT NULL,
  `event_type` varchar(100) NOT NULL,
  `url` varchar(2048) NOT NULL,
  `secret` varchar(256) DEFAULT NULL,
  `active` tinyint NOT NULL DEFAULT 1,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `webhooks_project` (`project_id`),
  CONSTRAINT `webhooks_project` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `webhook_logs` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `webhook_id` int unsigned NOT NULL,
  `event_type` varchar(100) NOT NULL,
  `payload` text NOT NULL,
  `response_code` int DEFAULT NULL,
  `error` text DEFAULT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `webhook_logs_webhook` (`webhook_id`),
  CONSTRAINT `webhook_logs_webhook` FOREIGN KEY (`webhook_id`) REFERENCES `webhooks` (`id`) ON DELETE CASCADE
);
