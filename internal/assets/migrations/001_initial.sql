-- +migrate Up

create table blobs (
    id bigserial primary key,
    data jsonb not null,
    owner_id character varying(56) REFERENCES accounts(account_id)
);

-- +migrate Down

drop table blobs;

