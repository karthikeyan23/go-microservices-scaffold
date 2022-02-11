create table if not exists entity
(
    id         uuid      not null
        constraint entity_pk
            primary key,
    name       varchar   not null,
    created_at timestamp not null
);

alter table entity
    owner to beezlabs;

create unique index if not exists entity_id_uindex
    on entity (id);

