-- name: GetJokeById :one
select *
from jokes
where id == ?
limit 1;

-- name: GetTypes :many
select *
from joke_types;

-- name: GetCategories :many
select *
from joke_categories c;

-- name: AddJoke :exec
insert
into jokes(question, answer, type_id, category_id, userID, guildID)
values (?, ?, ?, ?, ?, ?);

-- name: RemoveJoke :exec
delete
from jokes
where id = ?;

-- name: UpdateJoke :exec
update jokes
set question    = ?,
    answer      = ?,
    type_id     = ?,
    category_id = ?,
    userID      = ?,
    guildID     = ?
where id = ?;