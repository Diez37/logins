CREATE TABLE IF NOT EXISTS logins
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid       CHAR(36)    NOT NULL,
    login      VARCHAR(56) NOT NULL,
    banned     BIT         NOT NULL DEFAULT 0,
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_at  TIMESTAMP   NULL
);

CREATE UNIQUE INDEX logins_uuid ON logins (uuid);
CREATE UNIQUE INDEX logins_login ON logins (login);
