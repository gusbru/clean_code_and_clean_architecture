create schema ccca;

create table ccca.users (
    account_id UUID primary key,
    name varchar(255) not null,
    email varchar(255) unique not null,
    document varchar(255) not null,
    password varchar(255) not null
);

create table ccca.assets (
    asset_id varchar(10) not null,
    account_id UUID not null references ccca.users(account_id),
    amount decimal(10, 2) not null default 0,
    primary key (asset_id, account_id)
);
