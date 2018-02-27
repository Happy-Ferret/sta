-- init sql script for sta in a schema called sta
create schema if not exists sta;
set schema 'sta';

-- Context
create table if not exists contexts (
	id          serial primary key,
	name        varchar(64) not null,
	description varchar(1024) not null
);

create table if not exists contexts_container (
	context   int references contexts(id) not null,
	container int references contexts(id) not null,
	primary key (context, container)
);

create table if not exists contexts_links (
	context int references contexts(id) not null,
	name    varchar(64) not null,
	key     varchar(64) not null,
	locked  boolean not null,
	target  int references contexts(id) not null,
	primary key (context, target)
);

create table if not exists contexts_commands (
	context int references contexts(id) not null,
	command varchar(64) not null,
	primary key (context, command)
);

create table if not exists contexts_properties (
	context int references contexts(id) not null,
	name    varchar(64) not null,
	value   varchar(64) not null,
	primary key (context, name)
);
