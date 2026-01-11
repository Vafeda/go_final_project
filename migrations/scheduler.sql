CREATE TABLE scheduler
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    date      CHAR(8) NOT NULL DEFAULT "",
    title     VARCHAR(100) NOT NULL,
    comment   TEXT,
    `repeat`  VARCHAR(128)
);

CREATE INDEX idx_scheduler_date ON scheduler(date);