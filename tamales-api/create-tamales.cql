drop keyspace if exists tamales;

create keyspace tamales with replication = {'class':'SimpleStrategy', 'replication_factor' : 1};

/*
 *  Vendor Tables
 */

create table tamales.vendor_invite
(
    user_id uuid,
    email   text,
    primary key (user_id)
) with default_time_to_live = 3600; // 1 hour

create table tamales.vendor
(
    user_id uuid,
    email   text,
    primary key (user_id)
);

/*
 *  Vendor Location Tables
 */

create table tamales.vendor_location
(
    user_id       uuid,
    location_time timestamp,
    location      text,
    primary key (user_id)
);

create table tamales.vendor_location_history
(
    user_id       uuid,
    location_time timestamp,
    location      text,
    primary key ((user_id), location_time)
);

/*
 *  User Tables
 */

create table tamales.user_lookup_by_email
(
    email   text,
    user_id uuid,
    type    text,
    primary key (email)
);

create table tamales.login_token
(
    login_token uuid,
    user_id     uuid,
    primary key (login_token)
) with default_time_to_live = 600; // 10 minutes

create table tamales.auth_token
(
    auth_token uuid,
    user_id    uuid,
    primary key (auth_token)
) with default_time_to_live = 86400; // 1 day
