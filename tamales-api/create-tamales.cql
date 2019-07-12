drop keyspace if exists tamales;

create keyspace tamales with replication = {'class':'SimpleStrategy', 'replication_factor' : 1};

/*
 *  Vendor Tables
 */
create table tamales.vendor_invite
(
    vendor_id uuid,
    email     text,
    primary key (vendor_id)
) with default_time_to_live = 3600;

create table tamales.vendor
(
    vendor_id uuid,
    email     text,
    primary key (vendor_id)
);

/*
 *  Vendor Location Tables
 */

create table tamales.vendor_location
(
    vendor_id     uuid,
    location_time timestamp,
    location      text,
    primary key (vendor_id)
);

create table tamales.vendor_location_history
(
    vendor_id     uuid,
    location_time timestamp,
    location      text,
    primary key ((vendor_id), location_time)
);

/*
 *  User Tables
 */
create table tamales.user_lookup_by_email
(
    email   text,
    user_id uuid,
    type text,
    primary key (email)
);

create table tamales.login_token
(
    login_token uuid,
    user_id     uuid,
    primary key (login_token)
) with default_time_to_live = 600;
