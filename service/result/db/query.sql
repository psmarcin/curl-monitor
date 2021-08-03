-- name: GetResult :one
select *
from result
where uuid = $1
limit 1;

-- name: ListResult :many
select *
from result
where job_uuid = $1
  and created_at > $2
order by created_at
limit $2;

-- name: CreateResult :one
insert into result (uuid, job_uuid, output, type, created_at, updated_at)
values (@uuid, @job_uuid, @output, @type, @created_at, @updated_at)
returning *;
