# drop database entrytask;
# create database entrytask;
use entrytask;
create table `t_user` (
    `id` bigint unsigned not null auto_increment,
    `username` varchar(64) collate utf8mb4_unicode_ci not null comment 'username',
    `nickname` varchar(64) collate utf8mb4_unicode_ci not null comment 'nickname',
    `profile_pic` varchar(1024) collate utf8mb4_unicode_ci not null comment 'user avatar url',
    `password` varchar(1024) character set utf8mb4 collate utf8mb4_unicode_ci not null comment 'user password (encrypted)',
    `ctime` TIMESTAMP(6) not null comment 'create timestamp',
    `mtime` TIMESTAMP(6) not null comment 'modify timestamp',
    primary key (`id`),
    unique key `uniq_name` (`username`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_unicode_ci;