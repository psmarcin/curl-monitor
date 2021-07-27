-- name: GetJob :one
select * from job where uuid = $1 limit 1;

-- name: ListJob :many
select * from job where created_at > $1 order by created_at limit $2;

-- name: CreateJob :one
insert into job (uuid, name, command, created_at, updated_at) values (@uuid, @name, @command, @created_at, @updated_at) returning *;

-- name: UpdateJob :one
update job
    set name = @name, command = @command, updated_at = @updated_at
where uuid = @uuid
returning *;
