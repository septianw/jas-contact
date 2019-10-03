-- +migrate Up
INSERT INTO `contact` VALUES ('1','Amara','Vandervort','PhD','0'),
('2','Kallie','Wolf','Jr.','0'),
('3','Theresa','Larkin','V','1'),
('4','Myrl','Steuber','IV','1'),
('5','Maybelle','Hamill','IV','0'),
('6','Rhea','Brown','Sr.','1'),
('7','Deshawn','Metz','Jr.','1'),
('8','Lyda','Metz','V','0'),
('9','Johnathon','Boyle','III','0'),
('10','Misty','Harvey','DDS','1'),
('11','Zaria','Boyle','PhD','1'),
('12','Kellen','Cummings','II','1'),
('13','Eda','Hills','IV','1'),
('14','Gudrun','Dibbert','PhD','0'),
('15','Alia','Luettgen','I','1'),
('16','Adan','Robel','PhD','1'),
('17','Okey','Green','I','0'),
('18','Douglas','Wyman','Jr.','0'),
('19','Jordan','Ziemann','IV','1'),
('20','Bridget','Bruen','I','1'),
('21','Johnny','Johnson','DVM','1'),
('22','Millie','O\'Kon','DDS','1'),
('23','Lolita','Walter','V','1'),
('24','Cora','Bechtelar','MD','1'),
('25','Bettye','Reinger','Sr.','1'),
('26','Laisha','Tillman','Sr.','1'),
('27','Ruthe','Cartwright','IV','0'),
('28','Reuben','Jacobs','Sr.','1'),
('29','Reilly','Beer','V','1'),
('30','Saul','Langworth','III','1'),
('31','Celia','Wolff','Jr.','1'),
('32','Baron','Kerluke','III','0'),
('33','Craig','Carter','IV','1'),
('34','Addie','Cruickshank','I','1'),
('35','Declan','Hills','Sr.','0'),
('36','Giovanni','Cronin','Sr.','0'),
('37','Raymundo','Fay','IV','0'),
('38','Joannie','Emard','DDS','0'),
('39','Amaya','Hand','V','0'),
('40','Mitchell','Aufderhar','Jr.','0'),
('41','Rashad','McLaughlin','MD','1'),
('42','Leda','Rath','DDS','1'),
('43','Lelah','Dicki','MD','0'),
('44','Manuela','Waters','IV','0'),
('45','Elwin','Ernser','DVM','0'),
('46','Shane','Dicki','PhD','1'),
('47','Wallace','Brakus','IV','1'),
('48','Dante','O\'Hara','Jr.','1'),
('49','Breanne','Wisozk','IV','0'),
('50','Stevie','Nitzsche','Sr.','1'),
('51','Nya','Cummings','V','1'),
('52','Schuyler','Mayer','V','1'),
('53','Syble','Hagenes','Jr.','1'),
('54','Rickie','Greenholt','MD','1'),
('55','Melyssa','Senger','DDS','0'),
('56','Kenya','Swaniawski','MD','1'),
('57','Theodore','Mohr','II','0'),
('58','Tierra','Funk','III','0'),
('59','Faustino','Ankunding','Sr.','0'),
('60','Skye','Grant','V','0'),
('61','Alvis','Koelpin','V','1'),
('62','Camden','Conroy','II','0'),
('63','Jamir','McDermott','MD','1'),
('64','Marina','Davis','DDS','0'),
('65','Aiden','Kub','DVM','1'),
('66','Arnoldo','Shields','I','0'),
('67','Shana','Bartoletti','III','0'),
('68','Christophe','Kohler','IV','1'),
('69','Peggie','Luettgen','II','1'),
('70','Kelly','Rogahn','PhD','0'),
('71','Ramon','Feil','III','0'),
('72','Pink','Bogan','Jr.','0'),
('73','Ruthe','Kerluke','DDS','0'),
('74','Kobe','VonRueden','II','1'),
('75','Baby','Conn','I','0'),
('76','Jessika','Hoeger','I','0'),
('77','Melany','Prosacco','V','1'),
('78','Nathanael','Koch','V','1'),
('79','Sean','Stoltenberg','I','0'),
('80','Russel','Bergstrom','III','0'),
('81','Anastacio','Hilll','Jr.','0'),
('82','Bailey','Lehner','DDS','0'),
('83','Rebeca','Oberbrunner','DVM','1'),
('84','Hester','Gulgowski','V','0'),
('85','Guillermo','Reichel','Jr.','0'),
('86','Patrick','Jakubowski','Sr.','1'),
('87','Nikita','Ernser','III','1'),
('88','Haylie','Barrows','MD','0'),
('89','Candelario','Kris','DDS','0'),
('90','Cloyd','Bahringer','DVM','1'),
('91','Shanelle','Kessler','I','0'),
('92','Wendell','Bruen','IV','1'),
('93','Roman','Bins','Sr.','1'),
('94','Ebba','Rodriguez','MD','0'),
('95','Cleveland','Fay','I','1'),
('96','Antone','Greenfelder','Jr.','0'),
('97','Hilario','Robel','I','1'),
('98','Dwight','Trantow','PhD','0'),
('99','Presley','Batz','V','1'),
('100','Dell','Schaden','DVM','0'); 

INSERT INTO `contacttype` VALUES ('1','mitra'),
('2','investor'),
('3','karyawan'),
('4','konsumen'); 

INSERT INTO `contactwtype` VALUES ('1','1'),
('1','2'),
('1','3'),
('2','4'),
('4','2'),
('5','2'),
('7','1'),
('7','3'),
('7','4'),
('8','1'),
('8','2'),
('9','1'),
('10','2'),
('11','4'),
('12','2'),
('12','3'),
('12','4'),
('15','3'),
('17','1'),
('17','2'),
('17','3'),
('19','4'),
('22','1'),
('23','2'),
('25','4'),
('26','2'),
('26','3'),
('27','3'),
('28','1'),
('28','3'),
('29','4'),
('30','1'),
('30','3'),
('31','1'),
('32','2'),
('33','2'),
('33','3'),
('36','2'),
('36','3'),
('37','2'),
('37','3'),
('38','3'),
('41','1'),
('43','4'),
('46','2'),
('47','2'),
('48','4'),
('51','4'),
('53','1'),
('53','2'),
('55','2'),
('57','1'),
('57','4'),
('59','4'),
('60','2'),
('60','3'),
('62','2'),
('63','4'),
('64','1'),
('66','3'),
('67','1'),
('68','1'),
('68','3'),
('68','4'),
('69','2'),
('69','3'),
('70','2'),
('71','1'),
('71','2'),
('72','1'),
('72','2'),
('72','4'),
('73','3'),
('74','3'),
('75','1'),
('75','2'),
('76','2'),
('77','3'),
('78','1'),
('78','4'),
('79','4'),
('81','1'),
('81','2'),
('81','4'),
('83','2'),
('86','2'),
('90','2'),
('91','1'),
('92','3'),
('95','4'),
('98','4'),
('99','2'),
('99','3'); 

-- +migrate Down
TRUNCATE TABLE `contactwtype`;
TRUNCATE TABLE `contacttype`;
TRUNCATE TABLE `contact`;