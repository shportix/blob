-- +migrate Up

create table blobs (
    id bigserial primary key,
    data jsonb not null,
    owner_id character varying(56) REFERENCES accounts(account_id)
);

create table blobs_requests (
    d bigserial primary key,
    sign bpchar not null,
    real_request_target bpchar not null,
    date timestamp with time zone not null
);

-- +migrate Down

drop table blobs;
drop table blobs_requests;

