create table job
(
    uuid       uuid primary key,
    name       text      not null,
    command    text      not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null
)
