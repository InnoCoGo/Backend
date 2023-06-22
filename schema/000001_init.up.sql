create table users(
    id bigserial primary key,
    first_name varchar(256) not null,
    last_name varchar(256) not null,

    username varchar(256) not null unique,
    password_hash varchar(256) not null,
    rating int,
    num_people_rated int,

    tg_id int unique
);

create table trips(
    id serial primary key,
    
    admin_id int not null,
	is_driver        boolean         not null,

	places_max          int             not null,
    places_taken        int             not null,

	chosen_date_time    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

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

create table comments(
    from_id int not null,
    to_id int not null,

    primary key(from_id, to_id),
    foreign key(from_id) references users(id),
    foreign key(to_id) references users(id)
);