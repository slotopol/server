
/*
Tables of databases creates in code if them absent.
Statements below with inserting executes if `club` table is empty.
Users access levels (`gal` or `access` fields) are sum of followed ints:
   1 - *member*, user have access to club
   2 - *dealer*, can change club game settings
   4 - *booker*, can change user balance and move user money to/from club deposit
   8 - *master*, can change club bank, fund, deposit
  16 - *admin*, can change same access levels to other users
  31 - all rights
*/

INSERT INTO `club` (`cid`,`name`,`bank`,`fund`,`lock`,`rate`,`mrtp`) VALUES
(1,'virtual',10000,1000000,0,2.5,95),
(2,'global',0,0,0,2.5,0);

INSERT INTO `user` (`uid`,`email`,`secret`,`name`,`status`,`gal`) VALUES
(1,'admin@example.org','0YBoaT','admin',1,31),
(2,'dealer@example.org','LtpkAr','dealer',1,3),
(3,'player@example.org','iVI05M','player',1,1);

INSERT INTO `props` (`cid`,`uid`,`wallet`,`access`,`mrtp`) VALUES
(1,1,10000,1,0),
(2,1,0,0,0),
(1,2,10000,1+4+8,0),
(2,2,0,0,0),
(1,3,1000,1,98),
(2,3,0,0,98);
