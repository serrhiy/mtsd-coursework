create table "users" (
  "id"        bigint generated always as identity,
  "token"     char(32),
  "username"  varchar(64) not null check (length("username") >= 3),
  "createdAt" timestamptz default current_timestamp
);

alter table "users" add constraint "pkUsers" primary key ("id");
create unique index "akUserUsername" on "users" ("username");
create unique index "akUserToken" on "users" ("token");

create table "rooms" (
  "id"        bigint generated always as identity,
  "title"     varchar(64) not null check (length("title") >= 3),
  "createdAt" timestamptz default current_timestamp,
  "creatorId" bigint not null
);

alter table "rooms" add constraint "pkRooms" primary key ("id");
alter table "rooms" add constraint "fkRoomCreatorId" 
  foreign key ("creatorId") references "users" ("id");
create unique index "akRoomsTitle" on "rooms" ("title");

create table "messages" (
  "message"   text not null check (length("message") >= 1),
  "userId"    bigint not null,
  "roomId"    bigint not null,
  "createdAt" timestamptz default current_timestamp
);

alter table "messages" add constraint "pkMessages" primary key ("id");
alter table "messages" add constraint "fkMessagesUserIdUsers"
  foreign key ("userId") references "users" ("id");
alter table "messages" add constraint "fkMessagesRoomIdRooms"
  foreign key ("roomId") references "rooms" ("id");