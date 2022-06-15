CREATE TABLE `user_account` (
     `id` bigint unsigned NOT NULL AUTO_INCREMENT,
     `user_id` int(11) not NULL UNIQUE ,
     `balance` decimal(10,2) NOT NULL DEFAULT '0.00',
     `create_time` datetime DEFAULT now(),
     `update_time` datetime DEFAULT now(),
     PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4;