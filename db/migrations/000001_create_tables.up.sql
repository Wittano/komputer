create table jokes
(
    id       integer primary key autoincrement,
    question text,
    answer   text not null,
    type     text not null,
    category text not null,
    userID   text,
    guildID  text
);

create index jokes_index on jokes (id, guildID);

create table jokes_audit
(
    id primary key,
    status   text default 'CREATED',
    joke_id  integer not null,
    question text,
    answer   text    not null,
    type     text    not null,
    category text    not null,
    userID   text,
    guildID  text
);

create trigger joke_add_event_trigger
    after insert
    on jokes
begin
    insert
    into jokes_audit(joke_id, question, answer, type, category, userID, guildID)
    values (new.id, new.question, new.answer, new.type, new.category, new.userID, new.guildID);
end;

create trigger joke_update_event_trigger
    after update
    on jokes
begin
    insert
    into jokes_audit(status, joke_id, question, answer, type, category, userID, guildID)
    values ('UPDATED', new.id, new.question, new.answer, new.type,
            new.category,
            new.userID, new.guildID);
end;

create trigger joke_delete_event_trigger
    before delete
    on jokes
begin
    insert
    into jokes_audit(status, joke_id, question, answer, type, category, userID, guildID)
    values ('DELETE', old.id, old.question, old.answer, old.type,
            old.category,
            old.userID, old.guildID);
end;

create table admins
(
    id primary key,
    name text unique not null
);