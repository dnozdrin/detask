begin;
create table boards
(
    id          serial primary key,
    created_at  timestamp not null default now(),
    updated_at  timestamp not null default now(),

    name        varchar(500),
    description varchar(1000) not null default ''
);

create table columns
(
    id         serial primary key,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),

    name       varchar(255),
    board      int       not null,
    position   serial    not null,

    unique (name, board),
    unique (position, board),
    foreign key (board) references boards (id) on delete cascade
);

create table tasks
(
    id          serial primary key,
    created_at  timestamp not null default now(),
    updated_at  timestamp not null default now(),

    name        varchar(500),
    description varchar(5000) not null default '',
    "column"    int       not null,
    position    serial    not null,

    unique (position, "column"),
    foreign key ("column") references columns (id) on delete cascade
);

create table comments
(
    id         serial primary key,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),

    text       varchar(5000),
    task       int       not null,
    foreign key (task) references tasks (id) on delete cascade
);
commit;