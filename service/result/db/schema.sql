create type output as enum ('error', 'success');

create table result
(
    uuid       uuid primary key,
    job_uuid   uuid      not null,
    output     text      not null,
    type       output    not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null,
    constraint fk_job foreign key (job_uuid) references job (uuid) on delete cascade
)
