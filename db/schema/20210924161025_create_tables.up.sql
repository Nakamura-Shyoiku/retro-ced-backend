CREATE TABLE `Clicks` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `link` mediumtext COLLATE utf8_unicode_ci NOT NULL,
    `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=693 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE `Users` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `guid` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `username` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `first_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `last_name` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `email` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `password` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
    `first_time` tinyint(1) NOT NULL DEFAULT '0',
    `fb_id` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `created_at` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
    `acl` int(4) NOT NULL DEFAULT '0',
    `partner_site_id` int(11) unsigned NOT NULL DEFAULT '0',
    `password_reset` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT ' ',
    PRIMARY KEY (`id`),
    UNIQUE KEY `guid` (`guid`)
    ) ENGINE=InnoDB AUTO_INCREMENT=50 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


CREATE TABLE `Sites` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name` varchar(250) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `url` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `sleep` int(11) NOT NULL DEFAULT '0',
    `active` tinyint(1) NOT NULL DEFAULT '0',
    `max_page` int(11) NOT NULL DEFAULT '0',
    `last_scraped` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


CREATE TABLE `Urls` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `site_id` int(11) unsigned NOT NULL,
    `url` text COLLATE utf8_unicode_ci,
    `category` varchar(60) COLLATE utf8_unicode_ci NOT NULL,
    `last_updated` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
    `created_at` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
    PRIMARY KEY (`id`),
    KEY `site_id` (`site_id`),
    CONSTRAINT `Urls_ibfk_1` FOREIGN KEY (`site_id`) REFERENCES `Sites` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
    ) ENGINE=InnoDB AUTO_INCREMENT=1329717 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE `Products` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `site_id` int(11) unsigned NOT NULL,
    `url_id` int(11) unsigned NOT NULL,
    `category` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `brand` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `title` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `description` text COLLATE utf8_unicode_ci NOT NULL,
    `price` int(11) NOT NULL DEFAULT '0',
    `retail_price` int(11) NOT NULL DEFAULT '0',
    `model` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `item_number` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `product_condition` text COLLATE utf8_unicode_ci NOT NULL,
    `accessories` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `measurements` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `img` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT 'https://via.placeholder.com/300x400',
    `approved` tinyint(1) NOT NULL DEFAULT '0',
    `featured` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
    `last_updated` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
    `created_at` datetime NOT NULL DEFAULT '1000-01-01 00:00:00',
    `color` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `size` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `shoe_size` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `sub_category` varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    `url` text,
    `is_favourited` bool,
    `product_url` text COLLATE utf8_unicode_ci,
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_product_url` (`product_url`(767)),
    KEY `site_id` (`site_id`),
    KEY `url_id` (`url_id`),
    KEY `product_site_category_idx` (`site_id`,`category`) USING HASH,
    KEY `product_created_idx` (`created_at`) USING BTREE,
    KEY `products_featured_idx` (`featured`) USING HASH,
    KEY `products_category_idx` (`category`)
    ) ENGINE=InnoDB AUTO_INCREMENT=81855981 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE `Favourites` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `product_id` int(11) unsigned NOT NULL,
    `user_id` int(11) unsigned NOT NULL,
    PRIMARY KEY (`id`),
    KEY `user_id` (`user_id`),
    KEY `product_id` (`product_id`),
    CONSTRAINT `Favourites_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `Users` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT `Favourites_ibfk_2` FOREIGN KEY (`product_id`) REFERENCES `Products` (`id`) ON DELETE CASCADE ON UPDATE NO ACTION
    ) ENGINE=InnoDB AUTO_INCREMENT=290 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;


CREATE TABLE `Logs` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `site_id` int(11) unsigned NOT NULL,
    `url_id` int(11) unsigned NOT NULL,
    `level` varchar(50) COLLATE utf8_unicode_ci NOT NULL DEFAULT 'INFO',
    `message` text COLLATE utf8_unicode_ci NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
    ) ENGINE=InnoDB AUTO_INCREMENT=1074738 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;



