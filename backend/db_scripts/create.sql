CREATE TABLE contents (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    data TEXT UNIQUE NOT NULL
);

CREATE TABLE nodes (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    method INTEGER,
    content_id INTEGER,

    FOREIGN KEY (content_id) REFERENCES contents (id)
);

CREATE TABLE edges (
    src_id INTEGER,
    tgt_id INTEGER,

    CONSTRAINT pk_edges PRIMARY KEY (src_id, tgt_id),
    FOREIGN KEY (src_id) REFERENCES nodes (id),
    FOREIGN KEY (tgt_id) REFERENCES nodes (id)
);

CREATE TABLE users (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name TEXT,
    pass TEXT
);

CREATE TABLE sessions (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,

    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE dialogs (
    session_id INTEGER,
    content_id INTEGER,

    CONSTRAINT pk_dialogs PRIMARY KEY (session_id, content_id),
    FOREIGN KEY (session_id) REFERENCES sessions (id),
    FOREIGN KEY (content_id) REFERENCES contents (id)
);