CREATE TABLE `sequence` (
                            `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
                            `stub` VARCHAR(1) NOT NULL,
                            `timestamp` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                            PRIMARY KEY (`id`),
                            UNIQUE KEY `idx_uniq_stub` (`stub`)
) ENGINE=MYISAM DEFAULT CHARSET=utf8 COMMENT = '序号表';

DESCRIBE sequence;
REPLACE INTO sequence (stub) VALUES ('a');
SELECT LAST_INSERT_ID();
SELECT * FROM sequence;