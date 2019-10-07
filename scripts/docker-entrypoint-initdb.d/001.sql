CREATE TABLE events (
    id serial primary key,
    name	text not null,
    type	int not null,
    start_time  timestamp,
    end_time    timestamp
);
CREATE index start_idx ON events (start_time);
