create table users(
    id serial primary key,
    f_name varchar(256) not null,
    l_name varchar(256) not null,

    username varchar(256) not null unique,
    password_hash varchar(256) not null,
    -- tg_id int,

    rating int,
    num_people_rated int
);

create table comments(
    from_id int not null,
    to_id int not null,

    primary key(from_id, to_id),
    foreign key(from_id) references users(id),
    foreign key(to_id) references users(id)
);

create table trips(
    id serial primary key,
    
    admin_id int not null,
	tg_alias            varchar(256), 
	is_passanger        boolean         not null,

	places_max          int             not null,
    places_taken        int             not null,

	chosen_date_time    varchar(256)    not null,

	from_point          int             not null,
	to_point            int             not null,

	description         varchar(256)
);

create table users_trips(
    user_id int not null,
    trip_id int not null,

    primary key(user_id, trip_id), 
    foreign key(user_id) references users(id),
    foreign key(trip_id) references trips(id)
);