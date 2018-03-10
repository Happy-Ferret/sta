-- init sql script for sta in a schema called sta
drop schema contexts cascade;
create schema if not exists contexts;
set schema 'contexts';

-- Context
create table if not exists contexts (
	id          serial primary key,
	name        varchar(64) not null,
	description varchar(1024) not null
);

-- Context container
create table if not exists container (
	context   int references contexts(id) not null,
	container int references contexts(id) not null,
	primary key (context, container)
);

-- Links
create table if not exists links (
	context int references contexts(id) not null,
	name    varchar(64) not null,
	key     varchar(64) not null,
	locked  boolean not null,
	target  int references contexts(id) not null,
	primary key (context, target)
);

-- links slaves
create table if not exists links_slaves (
	masterct int not null,
	mastertg int not null,
	slavect  int not null,
	slavetg  int not null,
	constraint links_slaves_pk primary key (masterct, mastertg, slavect, slavetg),
	constraint links_slaves_fg_master foreign key (masterct, mastertg) references links(context, target),
	constraint links_slaves_fg_slave foreign key (slavect, slavetg) references links(context, target)
);

-- Context commands
create table if not exists commands (
	context int references contexts(id) not null,
	command varchar(64) not null,
	primary key (context, command)
);

-- Context properties
create table if not exists properties (
	context int references contexts(id) not null,
	name    varchar(64) not null,
	value   varchar(64) not null,
	primary key (context, name)
);

-- World
create table if not exists worlds (
	id serial primary key,
	name  varchar(64) not null,
	spawn int references contexts(id) not null
);

-- sample data
insert into contexts (id, name, description) values (1, 'The First Room', 'You are here.');
insert into contexts (id, name, description) values (2, 'The Second Room', 'You are there.');

insert into links (context, name, key, locked, target) values (1, 'The Door', 'key#0', true, 2);
insert into links (context, name, key, locked, target) values (2, 'The Door', 'key#0', true, 1);
insert into links_slaves (masterct, mastertg, slavect, slavetg) values (1, 2, 2, 1);
insert into links_slaves (masterct, mastertg, slavect, slavetg) values (2, 1, 1, 2);

insert into contexts (id, name, description) values (3, 'The Key', 'This is it.');
insert into container (context, container) values (3, 1);
insert into properties (context, name, value) values (3, 'key', 'key#0');

insert into worlds (name, spawn) values ('The One True World', 1);
