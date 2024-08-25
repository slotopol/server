
/*
Tables of databases creates in code if them absent.
Statements below with inserting executes if `club` table is empty.
Users access levels (`gal` or `access` fields) are sum of followed ints:
  1 - *ban*, user have no access to club
  2 - *game*, can change club game settings
  4 - *user*, can change user balance and move user money to/from club deposit
  8 - *club*, can change club bank, fund, deposit
  16 - *admin*, can change same access levels to other users
  30 - all rights
*/

INSERT INTO `club` (`cid`,`ctime`,`utime`,`name`,`bank`,`fund`,`lock`,`jptrate`,`mrtp`) VALUES
(1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,'virtual',10000,1000000,0,0.015,0);

INSERT INTO `user` (`uid`,`ctime`,`utime`,`email`,`secret`,`name`,`gal`) VALUES
(1,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,'admin@example.org','0YBoaT','admin',30),
(2,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,'dealer@example.org','LtpkAr','dealer',2),
(3,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,'player@example.org','iVI05M','player',0);

INSERT INTO `props` (`cid`,`uid`,`ctime`,`utime`,`wallet`,`access`) VALUES
(1,2,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,12000,4+8),
(1,3,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP,1000,0);
