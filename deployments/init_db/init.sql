create table if not exists segment
(
    id          serial primary key,
    slug        varchar(255) not null unique,
    description text
);

create table if not exists client
(
    id          integer primary key check (id > 0)
);

create table if not exists user_segment
(
    user_id     integer references client (id) on delete cascade,
    segment_id  integer references segment (id) on delete cascade,
    deletion_time timestamp,
    primary key (user_id, segment_id)
);

create table if not exists user_segment_audit
(
    user_id         integer references client (id) on delete cascade,
    segment_id      integer not null,
    segment_slug    varchar(255),
    operation       char(1) not null,
    stamp           timestamp not null
);

create or replace function process_user_segment_audit() returns trigger as $$
begin
        if ((select id from client where id in (old.user_id, new.user_id)) is null) then
            return null;
        elsif (tg_op = 'DELETE') then
            insert into user_segment_audit select old.user_id, old.segment_id, (select slug from segment where id = old.segment_id) as segment_slug, 'D', now();
            return old;
        elsif (tg_op = 'INSERT') then
            insert into user_segment_audit select new.user_id, new.segment_id, (select slug from segment where id = new.segment_id) as segment_slug, 'I', now();
            return new;
        end if;
        return null;
end;
$$ language plpgsql;

create trigger process_user_segment_audit after insert or delete on user_segment for each row execute procedure process_user_segment_audit();


-- insert into segment(slug, description)
-- values ('AVITO_VOICE_MESSAGES', 'Голосовые сообщения в чатах'),
--        ('AVITO_PERFORMANCE_VAS', 'Новые услуги продвижения'),
--        ('AVITO_DISCOUNT_30', 'Скидка 30% на услуги продвижения'),
--        ('AVITO_DISCOUNT_50', 'Скидка 50% на услуги продвижения');
--
-- do $$begin
--     for i in 1..50 loop
--             insert into client(id) values (i);
--             if i&1 > 0 then
--                 insert into user_segment(user_id, segment_id) values (i, 1);
--             end if;
--             if i&2 > 0 then
--                 insert into user_segment(user_id, segment_id) values (i, 2);
--             end if;
--             if i&4 > 0 then
--                 insert into user_segment(user_id, segment_id) values (i, 3);
--             end if;
--             if i&8 > 0 then
--                 insert into user_segment(user_id, segment_id) values (i, 4);
--             end if;
-- end loop;
-- end;$$