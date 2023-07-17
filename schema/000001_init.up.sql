create table users(
    id bigserial primary key,
    first_name varchar(256),
    last_name varchar(256),

    username varchar(256) not null unique,
    password_hash varchar(256) not null,
    rating int,
    num_people_rated int,

    tg_id bigint
);

create table trips(
    id serial primary key,
    
    admin_id bigint not null,
    admin_username varchar(256),
    admin_tg_id bigint,
    
	is_driver        boolean         not null,

	places_max          int             not null,
    places_taken        int             not null,

	chosen_timestamp    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

	from_point          int             not null,
	to_point            int             not null,

	description         varchar(256)
);

create table users_trips(
    user_id bigint not null,
    trip_id bigint not null,

    primary key(user_id, trip_id), 
    foreign key(user_id) references users(id),
    foreign key(trip_id) references trips(id)
);

create table comments(
    from_id bigint not null,
    to_id bigint not null,

    primary key(from_id, to_id),
    foreign key(from_id) references users(id),
    foreign key(to_id) references users(id)
);

create table messages(
    id serial primary key,
    user_id bigint not null,
    room_id bigint not null,

    content varchar(256) not null,
    content_type int not null,
    url varchar(256),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    foreign key(user_id) references users(id),
    foreign key(room_id) references users(id),

);