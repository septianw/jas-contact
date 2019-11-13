-- +migrate Up
ALTER TABLE `contact` 
ADD COLUMN `phone` VARCHAR(60) NULL AFTER `prefix`,
ADD COLUMN `email` VARCHAR(45) NULL AFTER `phone`;

-- +migrate Down
ALTER TABLE `contact` 
DROP COLUMN `email`,
DROP COLUMN `phone`;