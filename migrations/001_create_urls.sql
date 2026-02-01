create table if not exists urls (
    id bigserial primary key,
    long_url text not null,
    code varchar(8) not null unique,
    created_at timestamptz not null default now()
);
