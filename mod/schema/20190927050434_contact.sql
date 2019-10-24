-- +migrate Up
CREATE TABLE `contact` (
  `contactid` int(11) NOT NULL AUTO_INCREMENT,
  `fname` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `lname` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `prefix` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  `deleted` TINYINT(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`contactid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE `contacttype` (
  `ctypeid` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(45) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`ctypeid`),
  UNIQUE INDEX `name_UNIQUE` (`name` ASC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE `contactwtype` (
  `contact_contactid` int(11) NOT NULL,
  `contacttype_ctypeid` int(11) NOT NULL,
  PRIMARY KEY (`contact_contactid`,`contacttype_ctypeid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- +migrate Down
DROP TABLE IF EXISTS `contactwtype`;
DROP TABLE IF EXISTS `contacttype`;
DROP TABLE IF EXISTS `contact`;