CREATE TABLE entry (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    origin TEXT NOT NULL,
    desc TEXT NOT NULL
);

CREATE TABLE asset (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    entry_id INTEGER,
    "location" TEXT,
    FOREIGN KEY (entry_id) REFERENCES entry(id)
);