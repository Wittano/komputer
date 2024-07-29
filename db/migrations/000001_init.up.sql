-- create admins
create table admins
(
    id   integer primary key autoincrement,
    name text unique not null
);

-- create default categories
create table joke_categories
(
    id   integer primary key autoincrement,
    name text unique not null
);

insert into joke_categories(name)
values ('Any');
insert into joke_categories(name)
values ('Programming');
insert into joke_categories(name)
values ('Misc');
insert into joke_categories(name)
values ('Dark');
insert into joke_categories(name)
values ('YoMama');

-- create default types

create table joke_types
(
    id   integer primary key autoincrement,
    name text unique not null
);

create trigger joke_types_lower_case_name
    after insert
    on joke_types
begin
    update joke_types set name = lower(new.name) where id = new.id;
end;

insert into joke_types(name)
values ('single');
insert into joke_types(name)
values ('twopart');

-- create jokes

create table jokes
(
    id          integer primary key autoincrement,
    question    text,
    answer      text    not null,
    type_id     integer not null references joke_types      default 1,
    category_id integer not null references joke_categories default 1,
    userID      text,
    guildID     text
);

create index jokes_index on jokes (id, guildID);

create table jokes_audit
(
    id          integer primary key autoincrement,
    status      text default 'CREATED',
    joke_id     integer not null,
    question    text,
    answer      text    not null,
    type_id     integer not null references joke_types,
    category_id integer not null references joke_categories,
    userID      text,
    guildID     text
);

create trigger joke_add_event_trigger
    after insert
    on jokes
begin
    insert
    into jokes_audit(joke_id, question, answer, type_id, category_id, userID, guildID)
    values (new.id, new.question, new.answer, new.type_id, new.category_id, new.userID, new.guildID);
end;

create trigger joke_update_event_trigger
    after update
    on jokes
begin
    insert
    into jokes_audit(status, joke_id, question, answer, type_id, category_id, userID, guildID)
    values ('UPDATED', new.id, new.question, new.answer, new.type_id,
            new.category_id,
            new.userID, new.guildID);
end;

create trigger joke_delete_event_trigger
    before delete
    on jokes
begin
    insert
    into jokes_audit(status, joke_id, question, answer, type_id, category_id, userID, guildID)
    values ('DELETE', old.id, old.question, old.answer, old.type_id,
            old.category_id,
            old.userID, old.guildID);
end;